package models

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/daystram/cast/cast-be/config"
	"github.com/daystram/cast/cast-be/constants"
	"github.com/daystram/cast/cast-be/datatransfers"
)

type LikeOrmer interface {
	GetOneByUserIDHash(userID string, hash string) (like datatransfers.Like, err error)
	GetCountByHash(hash string) (count int, err error)
	InsertLike(like datatransfers.Like) (ID primitive.ObjectID, err error)
	RemoveLikeByUserIDHash(userID string, hash string) (err error)
}

type likeOrm struct {
	collection *mongo.Collection
}

func NewLikeOrmer(db *mongo.Client) LikeOrmer {
	return &likeOrm{db.Database(config.AppConfig.MongoDBName).Collection(constants.DBCollectionLike)}
}

func (o *likeOrm) GetOneByUserIDHash(userID string, hash string) (like datatransfers.Like, err error) {
	err = o.collection.FindOne(context.Background(), bson.M{"author": userID, "hash": hash}).Decode(&like)
	return
}

func (o *likeOrm) GetCountByHash(hash string) (count int, err error) {
	var count64 int64
	count64, err = o.collection.CountDocuments(context.Background(), bson.M{"hash": hash})
	count = int(count64)
	return
}

func (o *likeOrm) InsertLike(like datatransfers.Like) (ID primitive.ObjectID, err error) {
	result := &mongo.InsertOneResult{}
	like.ID = primitive.NewObjectID()
	if result, err = o.collection.InsertOne(context.Background(), like); err != nil {
		return
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

func (o *likeOrm) RemoveLikeByUserIDHash(userID string, hash string) (err error) {
	_, err = o.collection.DeleteOne(context.Background(), bson.M{"author": userID, "hash": hash})
	return
}
