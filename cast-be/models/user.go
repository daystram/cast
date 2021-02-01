package models

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/daystram/cast/cast-be/config"
	"github.com/daystram/cast/cast-be/constants"
	"github.com/daystram/cast/cast-be/datatransfers"
)

type UserOrmer interface {
	GetOneByID(ID string) (user datatransfers.User, err error)
	GetOneByUsername(username string) (user datatransfers.User, err error)
	CheckUnique(field, value string) (err error)
	InsertUser(user datatransfers.User) (err error)
	DeleteOneByID(ID string) (err error)
}

type userOrm struct {
	collection *mongo.Collection
}

func NewUserOrmer(db *mongo.Client) UserOrmer {
	return &userOrm{db.Database(config.AppConfig.MongoDBName).Collection(constants.DBCollectionUser)}
}

func (o *userOrm) GetOneByID(ID string) (user datatransfers.User, err error) {
	err = o.collection.FindOne(context.Background(), bson.M{"_id": ID}).Decode(&user)
	return
}

func (o *userOrm) GetOneByUsername(username string) (user datatransfers.User, err error) {
	err = o.collection.FindOne(context.Background(), bson.M{"username": username}).Decode(&user)
	return
}

func (o *userOrm) CheckUnique(field, value string) (err error) {
	uniqueOptions := &options.FindOneOptions{Collation: &options.Collation{Locale: "en", Strength: 2}}
	if err = o.collection.FindOne(context.Background(), bson.M{field: value}, uniqueOptions).Err(); err == mongo.ErrNoDocuments {
		return nil
	}
	if err == nil {
		return errors.New("[CheckUnique] duplicate entry found")
	}
	return
}

func (o *userOrm) InsertUser(user datatransfers.User) (err error) {
	_, err = o.collection.InsertOne(context.Background(), user)
	return
}

func (o *userOrm) DeleteOneByID(ID string) (err error) {
	_, err = o.collection.DeleteOne(context.Background(), bson.M{"_id": ID})
	return
}
