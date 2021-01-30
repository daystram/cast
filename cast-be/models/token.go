package models

import (
	"context"

	"github.com/daystram/cast/cast-be/config"
	"github.com/daystram/cast/cast-be/constants"
	"github.com/daystram/cast/cast-be/datatransfers"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TokenOrmer interface {
	GetOneByInvokerPurpose(invoker primitive.ObjectID, purpose string) (token datatransfers.Token, err error)
	GetOneByHash(hash string) (token datatransfers.Token, err error)
	InsertToken(token datatransfers.Token) (ID primitive.ObjectID, err error)
	DeleteOneByInvokerPurpose(invoker primitive.ObjectID, purpose string) (err error)
	DeleteOneByHash(hash string) (err error)
}

type tokenOrm struct {
	collection *mongo.Collection
}

func NewTokenOrmer(db *mongo.Client) TokenOrmer {
	return &tokenOrm{db.Database(config.AppConfig.MongoDBName).Collection(constants.DBCollectionToken)}
}

func (o *tokenOrm) GetOneByInvokerPurpose(invoker primitive.ObjectID, purpose string) (token datatransfers.Token, err error) {
	err = o.collection.FindOne(context.Background(), bson.M{"invoker": invoker, "purpose": purpose}).Decode(&token)
	return
}

func (o *tokenOrm) GetOneByHash(hash string) (token datatransfers.Token, err error) {
	err = o.collection.FindOne(context.Background(), bson.M{"hash": hash}).Decode(&token)
	return
}

func (o *tokenOrm) InsertToken(token datatransfers.Token) (ID primitive.ObjectID, err error) {
	result := &mongo.InsertOneResult{}
	if token.ID.IsZero() {
		token.ID = primitive.NewObjectID()
	}
	if result, err = o.collection.InsertOne(context.Background(), token); err != nil {
		return
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

func (o *tokenOrm) DeleteOneByInvokerPurpose(invoker primitive.ObjectID, purpose string) (err error) {
	_, err = o.collection.DeleteOne(context.Background(), bson.M{"invoker": invoker, "purpose": purpose})
	return
}

func (o *tokenOrm) DeleteOneByHash(hash string) (err error) {
	_, err = o.collection.DeleteOne(context.Background(), bson.M{"hash": hash})
	return
}
