package handlers

import (
	"github.com/nareix/joy4/av/pubsub"
	"github.com/nareix/joy4/format/rtmp"
	data "gitlab.com/daystram/cast/cast-be/datatransfers"
	"gitlab.com/daystram/cast/cast-be/models"
	"net/http"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
)

type module struct {
	db   func() *Entity
	live Live
}

type Live struct {
	streams map[string]*Stream
	uplink  *rtmp.Server
	mutex   *sync.RWMutex
}

type Stream struct {
	queue *pubsub.Queue
}

type Component struct {
	DB *mongo.Client
}

type Entity struct {
	videoOrm models.VideoOrmer
	userOrm  models.UserOrmer
}

type Handler interface {
	CreateRTMPUpLink()
	StreamLive(username string, w http.ResponseWriter, r *http.Request) (err error)

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
