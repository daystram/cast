package handlers

import (
	data "gitlab.com/daystram/cast/cast-be/datatransfers"
	"gitlab.com/daystram/cast/cast-be/models"

	"go.mongodb.org/mongo-driver/mongo"
)

type module struct {
	db func() *Entity
}

type Entity struct {
	videoOrm models.VideoOrmer
	userOrm  models.UserOrmer
}

type Component struct {
	DB *mongo.Client
}

type Handler interface {
	GetVideo(variant string, count, offset int) (videos []data.Video, err error)
	Search(query string, tags []string) (videos []data.Video, err error)
	VODDetails(hash string) (videos data.Video, err error)
	LiveDetails(username string) (videos data.Video, err error)
}

func NewHandler(component Component) Handler {
	return &module{
		db: func() (e *Entity) {
			return &Entity{
				videoOrm: models.NewVideoOrmer(component.DB),
				userOrm:  models.NewUserOrmer(component.DB),
			}
		},
	}
}
