package models

import (
	"context"
	"gitlab.com/daystram/cast/cast-be/config"
	"gitlab.com/daystram/cast/cast-be/constants"
	"gitlab.com/daystram/cast/cast-be/datatransfers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type VideoOrmer interface {
	GetRecent(variant string, count int, offset int) ([]datatransfers.Video, error)
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
		{{"$match", bson.D{{"Type", variant}}}},
		{{"$sort", bson.D{{"CreatedAt", -1}}}},
		{{"$skip", offset}},
		{{"$limit", count}},
		{{"$lookup", bson.D{
			{"from", constants.DBCollectionUser},
			{"localField", "Author"},
			{"foreignField", "_id"},
			{"as", "Author"},
		}}},
		{{"$unwind", "$Author"}}}); err != nil {
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
