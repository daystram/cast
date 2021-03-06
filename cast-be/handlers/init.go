package handlers

import (
	"net/http"
	"sync"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gorilla/websocket"
	"github.com/nareix/joy4/av/pubsub"
	"github.com/nareix/joy4/format/rtmp"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	data "github.com/daystram/cast/cast-be/datatransfers"
	"github.com/daystram/cast/cast-be/models"
)

type module struct {
	db           *Entity
	mq           *amqp.Channel
	s3           *s3.S3
	chat         *Chat
	notification *Notification
	live         Live
}

type Live struct {
	streams map[string]*Stream
	uplink  *rtmp.Server
	mutex   *sync.RWMutex
}

type Chat struct {
	sockets  map[string][]*websocket.Conn
	upgrader websocket.Upgrader
}

type Notification struct {
	sockets  map[string]*websocket.Conn
	upgrader websocket.Upgrader
}

type Stream struct {
	queue *pubsub.Queue
}

type Component struct {
	DB *mongo.Client
	MQ *amqp.Channel
	S3 *s3.S3
}

type Entity struct {
	videoOrm        models.VideoOrmer
	userOrm         models.UserOrmer
	likeOrm         models.LikeOrmer
	subscriptionOrm models.SubscriptionOrmer
	commentOrm      models.CommentOrmer
}

type Handler interface {
	CreateRTMPUpLink()
	ControlUpLinkWindow(userID string, open bool) (err error)
	StreamLive(username string, w http.ResponseWriter, r *http.Request) (err error)

	Register(idToken data.UserRegister) (err error)

	UserDetails(userID string) (detail data.UserDetail, err error)
	UserGetOneByID(userID string) (user data.User, err error)

	CastList(variant string, count, offset int, userID ...string) (videos []data.Video, err error)
	AuthorList(author string, withUnlisted bool, count, offset int) (videos []data.Video, err error)
	SearchVideo(query string, tags []string, count, offset int) (videos []data.Video, err error)
	VideoDetails(hash string) (video data.Video, err error)
	CreateVOD(upload data.VideoUpload, controller beego.Controller, userID string) (ID primitive.ObjectID, err error)
	DeleteVideo(ID primitive.ObjectID, userID string) (err error)
	UpdateVideo(video data.VideoEdit, controller beego.Controller, userID string) (err error)
	CheckUniqueVideoTitle(title string) (err error)
	LikeVideo(userID string, hash string, like bool) (err error)
	Subscribe(userID string, username string, subscribe bool) (err error)
	CheckUserLikes(hash, username string) (liked bool, err error)
	CheckUserSubscribes(authorID string, username string) (subscribed bool, err error)
	CommentVideo(userID string, hash, content string) (comment data.Comment, err error)

	TranscodeListenerWorker()
	StartTranscode(hash string)

	ConnectNotificationWS(ctx *context.Context, userID string) (err error)
	ConnectChatWS(ctx *context.Context, hash string, userID ...string) (err error)
	NotificationPingWorker(conn *websocket.Conn)
	ChatReaderWorker(conn *websocket.Conn, hash string, user data.User, live bool)
	PushNotification(userID string, message data.NotificationOutgoing)
	BroadcastNotificationSubscriber(authorID string, message data.NotificationOutgoing)
}

func NewHandler(component Component) Handler {
	return &module{
		db: &Entity{
			videoOrm:        models.NewVideoOrmer(component.DB),
			userOrm:         models.NewUserOrmer(component.DB),
			likeOrm:         models.NewLikeOrmer(component.DB),
			subscriptionOrm: models.NewSubscriptionOrmer(component.DB),
			commentOrm:      models.NewCommentOrmer(component.DB),
		},
		mq: component.MQ,
		s3: component.S3,
		chat: &Chat{
			sockets:  make(map[string][]*websocket.Conn),
			upgrader: websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }},
		},
		notification: &Notification{
			sockets:  make(map[string]*websocket.Conn),
			upgrader: websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }},
		},
	}
}
