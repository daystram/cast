package models

import (
	"gitlab.com/daystram/cast/cast-be/config"
	"gitlab.com/daystram/cast/cast-be/constants"
	"gitlab.com/daystram/cast/cast-be/datatransfers"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserOrmer interface {
	GetOneByID(email string) (datatransfers.User, error)
}

type userOrm struct {
	collection *mongo.Collection
}

func NewUserOrmer(db *mongo.Client) UserOrmer {
	return &userOrm{db.Database(config.AppConfig.MongoDBName).Collection(constants.DBCollectionUser)}
}

func (o *userOrm) GetOneByID(ID string) (datatransfers.User, error) {
	return datatransfers.User{}, nil
}


