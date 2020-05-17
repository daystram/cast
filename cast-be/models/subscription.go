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

type SubscriptionOrmer interface {
	GetSubscriptionsByAuthorID(authorID primitive.ObjectID) (subscribers []datatransfers.Subscription, err error)
	GetCountByAuthorID(authorID primitive.ObjectID) (count int, err error)
	GetOneByAuthorIDUserID(authorID, userID primitive.ObjectID) (subscription datatransfers.Subscription, err error)
	InsertSubscription(subscription datatransfers.Subscription) (ID primitive.ObjectID, err error)
	RemoveSubscriptionByAuthorIDUserID(authorID, userID primitive.ObjectID) (err error)
}

type subscriptionOrmer struct {
	collection *mongo.Collection
}

func NewSubscriptionOrmer(db *mongo.Client) SubscriptionOrmer {
	return &subscriptionOrmer{db.Database(config.AppConfig.MongoDBName).Collection(constants.DBCollectionSubscription)}
}

func (o *subscriptionOrmer) GetSubscriptionsByAuthorID(authorID primitive.ObjectID) (result []datatransfers.Subscription, err error) {
	query := &mongo.Cursor{}
	if query, err = o.collection.Aggregate(context.Background(), mongo.Pipeline{
		{{"$match", bson.D{{"author", authorID}}}},
	}); err != nil {
		return
	}
	for query.Next(context.Background()) {
		var subscription datatransfers.Subscription
		if err = query.Decode(&subscription); err != nil {
			return
		}
		result = append(result, subscription)
	}
	return
}

func (o *subscriptionOrmer) GetCountByAuthorID(authorID primitive.ObjectID) (count int, err error) {
	var count64 int64
	count64, err = o.collection.CountDocuments(context.Background(), bson.M{"author": authorID})
	count = int(count64)
	return
}

func (o *subscriptionOrmer) GetOneByAuthorIDUserID(authorID, userID primitive.ObjectID) (subscription datatransfers.Subscription, err error) {
	err = o.collection.FindOne(context.Background(), bson.M{"author": authorID, "user": userID}).Decode(&subscription)
	return
}

func (o *subscriptionOrmer) InsertSubscription(subscription datatransfers.Subscription) (ID primitive.ObjectID, err error) {
	result := &mongo.InsertOneResult{}
	subscription.ID = primitive.NewObjectID()
	if result, err = o.collection.InsertOne(context.Background(), subscription); err != nil {
		return
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

func (o *subscriptionOrmer) RemoveSubscriptionByAuthorIDUserID(authorID, userID primitive.ObjectID) (err error) {
	_, err = o.collection.DeleteOne(context.Background(), bson.M{"author": authorID, "user": userID})
	return
}
