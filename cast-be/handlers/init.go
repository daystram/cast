package handlers

import (
	"net/http"
	"sync"

	"gitlab.com/daystram/cast/cast-be/config"
	data "gitlab.com/daystram/cast/cast-be/datatransfers"
	"gitlab.com/daystram/cast/cast-be/models"

	googlePS "cloud.google.com/go/pubsub"
	"github.com/nareix/joy4/av/pubsub"
	"github.com/nareix/joy4/format/rtmp"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type module struct {
	db   func() *Entity
	mq   func() *MQ
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
	DB       *mongo.Client
	MQClient *googlePS.Client
}

type Entity struct {
	videoOrm models.VideoOrmer
	userOrm  models.UserOrmer
}

type MQ struct {
	transcodeTopic       *googlePS.Topic
	completeSubscription *googlePS.Subscription
}

type Handler interface {
	CreateRTMPUpLink()
	ControlUpLinkWindow(userID primitive.ObjectID, open bool) (err error)
	StreamLive(username string, w http.ResponseWriter, r *http.Request) (err error)

	CheckUniqueUserField(field, value string) (err error)
	Register(info data.UserRegister) (err error)
	Authenticate(info data.UserLogin) (token string, err error)

	UserDetails(userID primitive.ObjectID) (user data.User, err error)
	UpdateUser(user data.UserEdit, ID primitive.ObjectID) (err error)

	FreshList(variant string, count, offset int) (videos []data.Video, err error)
	AuthorList(author string, count, offset int) (videos []data.Video, err error)
	Search(query string, tags []string) (videos []data.Video, err error)
	VideoDetails(hash string) (video data.Video, err error)
	CreateVOD(upload data.VideoUpload, userID primitive.ObjectID) (ID primitive.ObjectID, err error)
	DeleteVideo(ID, userID primitive.ObjectID) (err error)
	UpdateVideo(video data.VideoEdit, userID primitive.ObjectID) (err error)
	CheckUniqueVideoTitle(title string) (err error)
	NormalizeThumbnail(ID primitive.ObjectID) (err error)

	TranscodeListenerWorker()
	StartTranscode(hash string)
}

func NewHandler(component Component) Handler {
	return &module{
		db: func() (e *Entity) {
			return &Entity{
				videoOrm: models.NewVideoOrmer(component.DB),
				userOrm:  models.NewUserOrmer(component.DB),
			}
		},
		mq: func() (m *MQ) {
			return &MQ{
				transcodeTopic:       component.MQClient.Topic(config.AppConfig.TopicNameTranscode),
				completeSubscription: component.MQClient.Subscription(config.AppConfig.SubscriptionNameComplete),
			}
		},
	}
}
