package models

import (
	"context"
	"errors"
	"gitlab.com/daystram/cast/cast-be/config"
	"gitlab.com/daystram/cast/cast-be/constants"
	"gitlab.com/daystram/cast/cast-be/datatransfers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type VideoOrmer interface {
	GetRecent(variant string, count int, offset int) (videos []datatransfers.Video, err error)
	GetTrending(count int, offset int) (videos []datatransfers.Video, err error)
	GetLiked(userID primitive.ObjectID, count int, offset int) (videos []datatransfers.Video, err error)
	GetSubscribed(userID primitive.ObjectID, count int, offset int) (videos []datatransfers.Video, err error)
	GetAllVODByAuthor(author primitive.ObjectID) (videos []datatransfers.Video, err error)
	GetAllVODByAuthorPaginated(author primitive.ObjectID, count int, offset int) (videos []datatransfers.Video, err error)
	Search(query string, count, offset int) (videos []datatransfers.Video, err error)
	GetLiveByAuthor(userID primitive.ObjectID) (datatransfers.Video, error)
	GetOneByHash(hash string) (datatransfers.Video, error)
	IncrementViews(hash string, decrement ...bool) (err error)
	SetLive(authorID primitive.ObjectID, pending, live bool) (err error)
	SetResolution(hash string, resolution int) (err error)
	InsertVideo(video datatransfers.VideoInsert) (ID primitive.ObjectID, err error)
	EditVideo(video datatransfers.VideoInsert) (err error)
	DeleteOneByID(ID primitive.ObjectID) (err error)
	CheckUnique(title string) (err error)
}

type videoOrm struct {
	collection *mongo.Collection
}

func NewVideoOrmer(db *mongo.Client) VideoOrmer {
	return &videoOrm{db.Database(config.AppConfig.MongoDBName).Collection(constants.DBCollectionVideo)}
}

func (o *videoOrm) GetRecent(variant string, count int, offset int) (result []datatransfers.Video, err error) {
	query := &mongo.Cursor{}
	if query, err = o.collection.Aggregate(context.Background(), mongo.Pipeline{
		{{"$match", bson.D{{"type", variant}}}},
		{{"$match", bson.D{{"resolutions", bson.D{{"$ne", 0}}}}}},
		{{"$match", bson.D{{"is_live", true}}}},
		{{"$sort", bson.D{{"created_at", -1}, {"_id", 1}}}},
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
	for query.Next(context.Background()) {
		var video datatransfers.Video
		if err = query.Decode(&video); err != nil {
			return
		}
		result = append(result, video)
	}
	return
}

func (o *videoOrm) GetLiked(userID primitive.ObjectID, count int, offset int) (result []datatransfers.Video, err error) {
	query := &mongo.Cursor{}
	if query, err = o.collection.Aggregate(context.Background(), mongo.Pipeline{
		{{"$match", bson.D{{"resolutions", bson.D{{"$ne", 0}}}}}},
		{{"$match", bson.D{{"is_live", true}}}},
		{{"$sort", bson.D{{"created_at", -1}, {"_id", 1}}}},
		{{"$lookup", bson.D{
			{"from", constants.DBCollectionLike},
			{"localField", "hash"},
			{"foreignField", "hash"},
			{"as", "likes"},
		}}},
		{{"$match", bson.D{{"$expr", bson.D{{"$in", bson.A{userID, "$likes.author"}}}}}}},
		{{"$project", bson.D{{"likes", 0}}}},
		{{"$skip", offset}},
		{{"$limit", count}},
		{{"$lookup", bson.D{
			{"from", constants.DBCollectionUser},
			{"localField", "author"},
			{"foreignField", "_id"},
			{"as", "author"},
		}}},
		{{"$unwind", "$author"}},
	}); err != nil {
		return
	}
	for query.Next(context.Background()) {
		var video datatransfers.Video
		if err = query.Decode(&video); err != nil {
			return
		}
		result = append(result, video)
	}
	return
}

func (o *videoOrm) GetSubscribed(userID primitive.ObjectID, count int, offset int) (result []datatransfers.Video, err error) {
	query := &mongo.Cursor{}
	if query, err = o.collection.Aggregate(context.Background(), mongo.Pipeline{
		{{"$match", bson.D{{"resolutions", bson.D{{"$ne", 0}}}}}},
		{{"$match", bson.D{{"is_live", true}}}},
		{{"$sort", bson.D{{"created_at", -1}, {"_id", 1}}}},
		{{"$lookup", bson.D{
			{"from", constants.DBCollectionSubscription},
			{"localField", "author"},
			{"foreignField", "author"},
			{"as", "subscribers"},
		}}},
		{{"$match", bson.D{{"$expr", bson.D{{"$in", bson.A{userID, "$subscribers.user"}}}}}}},
		{{"$project", bson.D{{"subscribers", 0}}}},
		{{"$skip", offset}},
		{{"$limit", count}},
		{{"$lookup", bson.D{
			{"from", constants.DBCollectionUser},
			{"localField", "author"},
			{"foreignField", "_id"},
			{"as", "author"},
		}}},
		{{"$unwind", "$author"}},
	}); err != nil {
		return
	}
	for query.Next(context.Background()) {
		var video datatransfers.Video
		if err = query.Decode(&video); err != nil {
			return
		}
		result = append(result, video)
	}
	return
}

func (o *videoOrm) GetTrending(count int, offset int) (result []datatransfers.Video, err error) {
	query := &mongo.Cursor{}
	if query, err = o.collection.Aggregate(context.Background(), mongo.Pipeline{
		{{"$match", bson.D{{"resolutions", bson.D{{"$ne", 0}}}}}},
		{{"$match", bson.D{{"is_live", true}}}},
		{{"$lookup", bson.D{
			{"from", constants.DBCollectionLike},
			{"localField", "hash"},
			{"foreignField", "hash"},
			{"as", "likes"},
		}}},
		{{"$addFields", bson.D{
			{"weight", bson.D{{
				"$add", bson.A{"$views",
					bson.D{{"$multiply",
						bson.A{bson.D{{"$size", "$likes"}}, 5}}}}}}},
		}}},
		{{"$sort", bson.D{{"weight", -1,}, {"_id", 1}}}},
		{{"$project", bson.D{{"weight", 0}, {"likes", 0}}}},
		{{"$skip", offset}},
		{{"$limit", count}},
		{{"$lookup", bson.D{
			{"from", constants.DBCollectionUser},
			{"localField", "author"},
			{"foreignField", "_id"},
			{"as", "author"},
		}}},
		{{"$unwind", "$author"}},
	}); err != nil {
		return
	}
	for query.Next(context.Background()) {
		var video datatransfers.Video
		if err = query.Decode(&video); err != nil {
			return
		}
		result = append(result, video)
	}
	return
}

func (o *videoOrm) GetAllVODByAuthor(author primitive.ObjectID) (videos []datatransfers.Video, err error) {
	query := &mongo.Cursor{}
	if query, err = o.collection.Aggregate(context.Background(), mongo.Pipeline{
		{{"$match", bson.D{{"author", author}}}},
		{{"$match", bson.D{{"type", constants.VideoTypeVOD}}}},
		{{"$sort", bson.D{{"created_at", -1}, {"_id", 1}}}},
		{{"$lookup", bson.D{
			{"from", constants.DBCollectionUser},
			{"localField", "author"},
			{"foreignField", "_id"},
			{"as", "author"},
		}}},
		{{"$unwind", "$author"}}}); err != nil {
		return
	}
	for query.Next(context.Background()) {
		var video datatransfers.Video
		if err = query.Decode(&video); err != nil {
			return
		}
		videos = append(videos, video)
	}
	return
}

func (o *videoOrm) GetAllVODByAuthorPaginated(author primitive.ObjectID, count int, offset int) (videos []datatransfers.Video, err error) {
	query := &mongo.Cursor{}
	if query, err = o.collection.Aggregate(context.Background(), mongo.Pipeline{
		{{"$match", bson.D{{"author", author}}}},
		{{"$match", bson.D{{"type", constants.VideoTypeVOD}}}},
		{{"$sort", bson.D{{"created_at", -1}, {"_id", 1}}}},
		{{"$skip", offset}},
		{{"$limit", count}},
		{{"$lookup", bson.D{
			{"from", constants.DBCollectionUser},
			{"localField", "author"},
			{"foreignField", "_id"},
			{"as", "author"},
		}}},
		{{"$unwind", "$author"}},
		{{"$lookup", bson.D{
			{"from", constants.DBCollectionLike},
			{"localField", "hash"},
			{"foreignField", "hash"},
			{"as", "likes"},
		}}},
		{{"$addFields", bson.D{
			{"likes", bson.D{{"$size", "$likes"}}},
		}}},
	}); err != nil {
		return
	}
	for query.Next(context.Background()) {
		var video datatransfers.Video
		if err = query.Decode(&video); err != nil {
			return
		}
		videos = append(videos, video)
	}
	return
}

func (o *videoOrm) Search(queryString string, count int, offset int) (result []datatransfers.Video, err error) {
	query := &mongo.Cursor{}
	if query, err = o.collection.Aggregate(context.Background(), mongo.Pipeline{
		{{"$match", bson.D{{"$text", bson.D{{"$search", queryString}}}}}},
		{{"$match", bson.D{{"resolutions", bson.D{{"$ne", 0}}}}}},
		{{"$match", bson.D{{"is_live", true}}}},
		{{"$sort", bson.D{{"views", -1}, {"created_at", -1}, {"_id", 1}}}},
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
	for query.Next(context.Background()) {
		var video datatransfers.Video
		if err = query.Decode(&video); err != nil {
			return
		}
		result = append(result, video)
	}
	return
}

func (o *videoOrm) GetLiveByAuthor(userID primitive.ObjectID) (video datatransfers.Video, err error) {
	query := &mongo.Cursor{}
	if query, err = o.collection.Aggregate(context.Background(), mongo.Pipeline{
		{{"$match", bson.D{{"author", userID}}}},
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
	if exists := query.Next(context.Background()); exists {
		err = query.Decode(&video)
		return
	}
	return datatransfers.Video{}, mongo.ErrNoDocuments
}

func (o *videoOrm) GetOneByHash(hash string) (video datatransfers.Video, err error) {
	query := &mongo.Cursor{}
	if query, err = o.collection.Aggregate(context.Background(), mongo.Pipeline{
		{{"$match", bson.D{{"hash", hash}}}},
		{{"$limit", 1}},
		{{"$lookup", bson.D{
			{"from", constants.DBCollectionUser},
			{"localField", "author"},
			{"foreignField", "_id"},
			{"as", "author"},
		}}},
		{{"$unwind", "$author"}},
		{{"$lookup", bson.D{
			{"from", constants.DBCollectionLike},
			{"localField", "hash"},
			{"foreignField", "hash"},
			{"as", "likes"},
		}}},
		{{"$addFields", bson.D{
			{"likes", bson.D{{"$size", "$likes"}}},
		}}},
	}); err != nil {
		return
	}
	if exists := query.Next(context.Background()); exists {
		err = query.Decode(&video)
		return
	}
	return datatransfers.Video{}, mongo.ErrNoDocuments
}

func (o *videoOrm) IncrementViews(hash string, decrement ...bool) error {
	delta := 1
	if len(decrement) > 0 && decrement[0] {
		delta = -1
	}
	return o.collection.FindOneAndUpdate(context.Background(), bson.M{"hash": hash}, bson.M{"$inc": bson.M{"views": delta}}).Err()
}

func (o *videoOrm) SetLive(authorID primitive.ObjectID, pending, live bool) (err error) {
	var stream datatransfers.VideoInsert
	if err = o.collection.FindOneAndUpdate(context.Background(),
		bson.M{"author": authorID, "type": constants.VideoTypeLive},
		bson.D{{"$set", bson.D{
			{"is_live", live},
			{"pending", pending},
		}}},
	).Decode(&stream); err != nil {
		return
	}
	if stream.IsLive != live {
		if err = o.collection.FindOneAndUpdate(context.Background(),
			bson.M{"author": authorID, "type": constants.VideoTypeLive},
			bson.D{{"$set", bson.D{
				{"created_at", time.Now()},
			}}},
		).Err(); err != nil {
			return
		}
	}
	if !pending && live {
		err = o.collection.FindOneAndUpdate(context.Background(),
			bson.M{"author": authorID, "type": constants.VideoTypeLive},
			bson.D{{"$set", bson.D{
				{"views", 0},
			}}},
		).Err()
	}
	return
}

func (o *videoOrm) SetResolution(hash string, resolution int) (err error) {
	return o.collection.FindOneAndUpdate(context.Background(),
		bson.M{"hash": hash, "type": constants.VideoTypeVOD},
		bson.D{{"$set", bson.D{{"resolutions", resolution}}}},
	).Err()
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
	if result, err = o.collection.InsertOne(context.Background(), video); err != nil {
		return
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

func (o *videoOrm) EditVideo(video datatransfers.VideoInsert) (err error) {
	return o.collection.FindOneAndUpdate(context.Background(),
		bson.M{"hash": video.Hash, "author": video.Author},
		bson.D{{"$set", bson.D{
			{"title", video.Title},
			{"description", video.Description},
			{"tags", video.Tags},
		}}},
	).Err()
}

func (o *videoOrm) DeleteOneByID(ID primitive.ObjectID) (err error) {
	_, err = o.collection.DeleteOne(context.Background(), bson.M{"_id": ID})
	return
}

func (o *videoOrm) CheckUnique(title string) (err error) {
	uniqueOptions := &options.FindOneOptions{Collation: &options.Collation{Locale: "en", Strength: 2}}
	if err = o.collection.FindOne(context.Background(), bson.M{"title": title}, uniqueOptions).Err(); err == mongo.ErrNoDocuments {
		return nil
	}
	if err == nil {
		return errors.New("[CheckUnique] duplicate entry found")
	}
	return
}
