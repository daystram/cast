package models

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/daystram/cast/cast-be/config"
	"github.com/daystram/cast/cast-be/constants"
	"github.com/daystram/cast/cast-be/datatransfers"
)

type UserOrmer interface {
	GetOneByID(ID primitive.ObjectID) (user datatransfers.User, err error)
	GetOneByEmail(email string) (user datatransfers.User, err error)
	GetOneByUsername(username string) (user datatransfers.User, err error)
	CheckUnique(field, value string) (err error)
	InsertUser(user datatransfers.User) (ID primitive.ObjectID, err error)
	EditUser(user datatransfers.User) (err error)
	SetVerified(ID primitive.ObjectID) (err error)
	DeleteOneByID(ID primitive.ObjectID) (err error)
}

type userOrm struct {
	collection *mongo.Collection
}

func NewUserOrmer(db *mongo.Client) UserOrmer {
	return &userOrm{db.Database(config.AppConfig.MongoDBName).Collection(constants.DBCollectionUser)}
}

func (o *userOrm) GetOneByID(ID primitive.ObjectID) (user datatransfers.User, err error) {
	err = o.collection.FindOne(context.Background(), bson.M{"_id": ID}).Decode(&user)
	return
}

func (o *userOrm) GetOneByEmail(email string) (user datatransfers.User, err error) {
	err = o.collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
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

func (o *userOrm) InsertUser(user datatransfers.User) (ID primitive.ObjectID, err error) {
	result := &mongo.InsertOneResult{}
	if user.ID.IsZero() {
		user.ID = primitive.NewObjectID()
	}
	if result, err = o.collection.InsertOne(context.Background(), user); err != nil {
		return
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

func (o *userOrm) EditUser(user datatransfers.User) (err error) {
	return o.collection.FindOneAndUpdate(context.Background(),
		bson.M{"_id": user.ID},
		bson.D{{"$set", bson.D{
			{"name", user.Name},
			{"email", user.Email},
			{"password", user.Password},
		}}},
	).Err()
}

func (o *userOrm) SetVerified(ID primitive.ObjectID) (err error) {
	return o.collection.FindOneAndUpdate(context.Background(),
		bson.M{"_id": ID},
		bson.D{{"$set", bson.D{{"verified", true}}}},
	).Err()
}

func (o *userOrm) DeleteOneByID(ID primitive.ObjectID) (err error) {
	_, err = o.collection.DeleteOne(context.Background(), bson.M{"_id": ID})
	return
}
