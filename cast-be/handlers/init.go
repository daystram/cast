package handlers

import (
	"net/http"
	"sync"

	"github.com/daystram/cast/cast-be/config"
	data "github.com/daystram/cast/cast-be/datatransfers"
	"github.com/daystram/cast/cast-be/models"

	googlePS "cloud.google.com/go/pubsub"
	"github.com/astaxie/beego/context"
	"github.com/gorilla/websocket"
	"github.com/mailgun/mailgun-go"
	"github.com/nareix/joy4/av/pubsub"
	"github.com/nareix/joy4/format/rtmp"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type module struct {
	db           *Entity
	mq           *MQ
	chat         *Chat
	notification *Notification
	mailer       *mailgun.MailgunImpl
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
	DB       *mongo.Client
	MQClient *googlePS.Client
	Mailer   *mailgun.MailgunImpl
}

type Entity struct {
	videoOrm        models.VideoOrmer
	userOrm         models.UserOrmer
	likeOrm         models.LikeOrmer
	subscriptionOrm models.SubscriptionOrmer
	commentOrm      models.CommentOrmer
	tokenOrm        models.TokenOrmer
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
	SendResetToken(user data.User) (err error)
	CheckResetToken(key string) (err error)
	UpdatePassword(info data.UserUpdatePassword) (err error)
	Verify(key string) (err error)
	Authenticate(info data.UserLogin) (user data.User, token string, err error)

	SendSingleEmail(subject, recipient, template string, variable map[string]string)

	UserDetails(userID primitive.ObjectID) (detail data.UserDetail, err error)
	GetUserByEmail(email string) (user data.User, err error)
	UpdateUser(user data.UserEditForm, ID primitive.ObjectID) (err error)
	NormalizeProfile(username string) (err error)

	CastList(variant string, count, offset int, userID ...primitive.ObjectID) (videos []data.Video, err error)
	AuthorList(author string, count, offset int) (videos []data.Video, err error)
	SearchVideo(query string, tags []string, count, offset int) (videos []data.Video, err error)
	VideoDetails(hash string) (video data.Video, err error)
	CreateVOD(upload data.VideoUpload, userID primitive.ObjectID) (ID primitive.ObjectID, err error)
	DeleteVideo(ID, userID primitive.ObjectID) (err error)
	UpdateVideo(video data.VideoEdit, userID primitive.ObjectID) (err error)
	CheckUniqueVideoTitle(title string) (err error)
	NormalizeThumbnail(hash string) (err error)
	LikeVideo(userID primitive.ObjectID, hash string, like bool) (err error)
	Subscribe(userID primitive.ObjectID, username string, subscribe bool) (err error)
	CheckUserLikes(hash, username string) (liked bool, err error)
	CheckUserSubscribes(authorID primitive.ObjectID, username string) (subscribed bool, err error)
	CommentVideo(userID primitive.ObjectID, hash, content string) (comment data.Comment, err error)

	TranscodeListenerWorker()
	StartTranscode(hash string)

	ConnectNotificationWS(ctx *context.Context, userID primitive.ObjectID) (err error)
	ConnectChatWS(ctx *context.Context, hash string, userID ...primitive.ObjectID) (err error)
	NotificationPingWorker(conn *websocket.Conn)
	ChatReaderWorker(conn *websocket.Conn, hash string, user data.User, live bool)
	PushNotification(userID primitive.ObjectID, message data.NotificationOutgoing)
	BroadcastNotificationSubscriber(authorID primitive.ObjectID, message data.NotificationOutgoing)
}

func NewHandler(component Component) Handler {
	return &module{
		db: &Entity{
			videoOrm:        models.NewVideoOrmer(component.DB),
			userOrm:         models.NewUserOrmer(component.DB),
			likeOrm:         models.NewLikeOrmer(component.DB),
			subscriptionOrm: models.NewSubscriptionOrmer(component.DB),
			commentOrm:      models.NewCommentOrmer(component.DB),
			tokenOrm:        models.NewTokenOrmer(component.DB),
		},
		mq: &MQ{
			transcodeTopic:       component.MQClient.Topic(config.AppConfig.TopicNameTranscode),
			completeSubscription: component.MQClient.Subscription(config.AppConfig.SubscriptionNameComplete),
		},
		chat: &Chat{
			sockets:  make(map[string][]*websocket.Conn),
			upgrader: websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }},
		},
		notification: &Notification{
			sockets:  make(map[string]*websocket.Conn),
			upgrader: websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }},
		},
		mailer: component.Mailer,
	}
}
