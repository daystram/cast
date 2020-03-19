package models

import (
	"context"
	"gitlab.com/daystram/cast/cast-be/config"
	"gitlab.com/daystram/cast/cast-be/constants"
	"gitlab.com/daystram/cast/cast-be/datatransfers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type VideoOrmer interface {
	GetRecent(variant string, count int, offset int) ([]datatransfers.Video, error)
	GetOneByHash(hash string) (datatransfers.Video, error)
	IncrementViews(hash string) error
	SetLive(authorID primitive.ObjectID, live bool) (err error)
	InsertVideo(video datatransfers.VideoInsert) (ID primitive.ObjectID, err error)
}

type videoOrm struct {
	collection *mongo.Collection
}

func NewVideoOrmer(db *mongo.Client) VideoOrmer {
	return &videoOrm{db.Database(config.AppConfig.MongoDBName).Collection(constants.DBCollectionVideo)}
}

func (o *videoOrm) GetRecent(variant string, count int, offset int) (result []datatransfers.Video, err error) {
	query := &mongo.Cursor{}
	if query, err = o.collection.Aggregate(context.TODO(), mongo.Pipeline{
		{{"$match", bson.D{{"type", variant}}}},
		{{"$match", bson.D{{"is_live", variant == constants.VideoTypeLive}}}},
		{{"$sort", bson.D{{"created_at", -1}}}},
		{{"$skip", offset}},
		{{"$limit", count}},
		{{"$lookup", bson.D{
			{"from", constants.DBCollectionUser},
			{"localField", "author"},
			{"foreignField", "_id"},
			{"as", "author"},
		}}},
		{{"$unwind", "$author"}}}); err != nil {
		return
	}
	for query.Next(context.TODO()) {
		var video datatransfers.Video
		if err = query.Decode(&video); err != nil {
			return
		}
		result = append(result, video)
	}
	return
}

func (o *videoOrm) GetOneByHash(hash string) (video datatransfers.Video, err error) {
	query := &mongo.Cursor{}
	if query, err = o.collection.Aggregate(context.TODO(), mongo.Pipeline{
		{{"$match", bson.D{{"hash", hash}}}},
		{{"$limit", 1}},
		{{"$lookup", bson.D{
			{"from", constants.DBCollectionUser},
			{"localField", "author"},
			{"foreignField", "_id"},
			{"as", "author"},
		}}},
		{{"$unwind", "$author"}}}); err != nil {
		return
	}
	if exists := query.Next(context.TODO()); exists {
		err = query.Decode(&video)
		return
	}
	return datatransfers.Video{}, mongo.ErrNoDocuments
}

func (o *videoOrm) IncrementViews(hash string) error {
	return o.collection.FindOneAndUpdate(context.TODO(), bson.M{"hash": hash}, bson.M{"$inc": bson.M{"views": 1}}).Err()
}

func (o *videoOrm) SetLive(authorID primitive.ObjectID, live bool) (err error) {
	return o.collection.FindOneAndUpdate(context.TODO(), bson.M{"author": authorID, "type": constants.VideoTypeLive}, bson.M{"is_live": live}).Err()
}

func (o *videoOrm) InsertVideo(video datatransfers.VideoInsert) (ID primitive.ObjectID, err error) {
	result := &mongo.InsertOneResult{}
	if video.Hash == "" {
		video.ID = primitive.NewObjectID()
		video.Hash = video.ID.Hex()
	}
	if video.ID.IsZero() {
		video.ID = primitive.NewObjectID()
	}
	if result, err = o.collection.InsertOne(context.TODO(), video); err != nil {
		return
	}
	return result.InsertedID.(primitive.ObjectID), nil
}
