package handlers

import (
	"net/http"
	"sync"

	"gitlab.com/daystram/cast/cast-be/config"
	data "gitlab.com/daystram/cast/cast-be/datatransfers"
	"gitlab.com/daystram/cast/cast-be/models"

	googlePS "cloud.google.com/go/pubsub"
	"github.com/mailgun/mailgun-go"
	"github.com/nareix/joy4/av/pubsub"
	"github.com/nareix/joy4/format/rtmp"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type module struct {
	db     func() *Entity
	mq     func() *MQ
	mailer *mailgun.MailgunImpl
	live   Live
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
	Mailer   *mailgun.MailgunImpl
}

type Entity struct {
	videoOrm   models.VideoOrmer
	userOrm    models.UserOrmer
	likeOrm    models.LikeOrmer
	commentOrm models.CommentOrmer
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
	SendVerification(user data.User) (err error)
	Verify(key string) (err error)
	Authenticate(info data.UserLogin) (token string, err error)

	SendSingleEmail(subject, content string, user data.User)

	UserDetails(userID primitive.ObjectID) (detail data.UserDetail, err error)
	GetUserByEmail(email string) (user data.User, err error)
	UpdateUser(user data.UserEditForm, ID primitive.ObjectID) (err error)
	NormalizeProfile(username string) (err error)

	FreshList(variant string, count, offset int) (videos []data.Video, err error)
	AuthorList(author string, count, offset int) (videos []data.Video, err error)
	Search(query string, tags []string) (videos []data.Video, err error)
	VideoDetails(hash string) (video data.Video, err error)
	CreateVOD(upload data.VideoUpload, userID primitive.ObjectID) (ID primitive.ObjectID, err error)
	DeleteVideo(ID, userID primitive.ObjectID) (err error)
	UpdateVideo(video data.VideoEdit, userID primitive.ObjectID) (err error)
	CheckUniqueVideoTitle(title string) (err error)
	NormalizeThumbnail(ID primitive.ObjectID) (err error)
	LikeVideo(userID primitive.ObjectID, hash string, like bool) (err error)
	CheckUserLikes(hash, username string) (liked bool, err error)
	CommentVideo(userID primitive.ObjectID, hash, content string) (comment data.Comment, err error)

	TranscodeListenerWorker()
	StartTranscode(hash string)
}

func NewHandler(component Component) Handler {
	return &module{
		db: func() (e *Entity) {
			return &Entity{
				videoOrm:   models.NewVideoOrmer(component.DB),
				userOrm:    models.NewUserOrmer(component.DB),
				likeOrm:    models.NewLikeOrmer(component.DB),
				commentOrm: models.NewCommentOrmer(component.DB),
			}
		},
		mq: func() (m *MQ) {
			return &MQ{
				transcodeTopic:       component.MQClient.Topic(config.AppConfig.TopicNameTranscode),
				completeSubscription: component.MQClient.Subscription(config.AppConfig.SubscriptionNameComplete),
			}
		},
		mailer: component.Mailer,
	}
}
