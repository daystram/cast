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

type CommentOrmer interface {
	GetAllByHash(hash string) (comments []datatransfers.Comment, err error)
	InsertComment(comment datatransfers.CommentInsert) (ID primitive.ObjectID, err error)
}

type commentOrm struct {
	collection *mongo.Collection
}

func NewCommentOrmer(db *mongo.Client) CommentOrmer {
	return &commentOrm{db.Database(config.AppConfig.MongoDBName).Collection(constants.DBCollectionComment)}
}

func (o *commentOrm) GetAllByHash(hash string) (comments []datatransfers.Comment, err error) {
	query := &mongo.Cursor{}
	if query, err = o.collection.Aggregate(context.TODO(), mongo.Pipeline{
		{{"$match", bson.D{{"hash", hash}}}},
		{{"$sort", bson.D{{"created_at", -1}}}},
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
		var comment datatransfers.Comment
		if err = query.Decode(&comment); err != nil {
			return
		}
		comments = append(comments, comment)
	}
	return
}

func (o *commentOrm) InsertComment(comment datatransfers.CommentInsert) (ID primitive.ObjectID, err error) {
	result := &mongo.InsertOneResult{}
	comment.ID = primitive.NewObjectID()
	if result, err = o.collection.InsertOne(context.TODO(), comment); err != nil {
		return
	}
	return result.InsertedID.(primitive.ObjectID), nil
}
