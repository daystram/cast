package routers

import (
	"context"
	"fmt"
	"log"

	conf "github.com/daystram/cast/cast-be/config"
	"github.com/daystram/cast/cast-be/controller/middleware"
	v1 "github.com/daystram/cast/cast-be/controller/v1"
	"github.com/daystram/cast/cast-be/handlers"

	"cloud.google.com/go/pubsub"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
	"github.com/mailgun/mailgun-go"
	"github.com/nareix/joy4/format"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/api/option"
)

func init() {
	conf.InitializeAppConfig()

	// Init MongoDB
	db, err := mongo.Connect(context.Background(), options.Client().ApplyURI(conf.AppConfig.MongoDBURI))
	if err != nil {
		log.Fatalf("Failed connecting to mongoDB at %s. %+v\n", conf.AppConfig.MongoDBURI, err)
	}
	fmt.Printf("[Initialization] MongoDB connected\n")

	// Init RTMP UpLink
	format.RegisterAll()

	// Init Google PubSub
	var pubsubClient *pubsub.Client
	pubsubClient, err = pubsub.NewClient(context.Background(), conf.AppConfig.GoogleProjectID, option.WithCredentialsFile(conf.AppConfig.JSONKey))
	if err != nil {
		log.Fatalf("Failed connecting to Google PubSub. %+v\n", err)
	}
	fmt.Printf("[Initialization] Google PubSub connected\n")

	// Init mailgun Client
	mailer := mailgun.NewMailgun(conf.AppConfig.MailgunDomain, conf.AppConfig.MailgunAPIKey)
	fmt.Printf("[Initialization] mailgun connected\n")

	h := handlers.NewHandler(handlers.Component{DB: db, MQClient: pubsubClient, Mailer: mailer})
	h.CreateRTMPUpLink()
	go h.TranscodeListenerWorker()
	fmt.Printf("[Initialization] Initialization completed\n")

	nsPublic := beego.NewNamespace("api/v1",
		beego.NSNamespace("/ping",
			beego.NSInclude(
				&v1.PingController{},
			),
		),
		beego.NSNamespace("/auth",
			beego.NSInclude(
				&v1.AuthController{Handler: h},
			),
		),
		beego.NSNamespace("/video",
			beego.NSInclude(
				&v1.VideoController{Handler: h},
			),
		),
		beego.NSNamespace("/live",
			beego.NSInclude(
				&v1.LiveController{Handler: h},
			),
		),
		beego.NSNamespace("/ws",
			beego.NSInclude(
				&v1.WebSocketController{Handler: h},
			),
		),
		beego.NSNamespace("/p",
			beego.NSBefore(middleware.AuthenticateJWT),
			beego.NSNamespace("/user",
				beego.NSInclude(
					&v1.UserControllerAuth{Handler: h},
				),
			),
			beego.NSNamespace("/video",
				beego.NSInclude(
					&v1.VideoControllerAuth{Handler: h},
				),
			),
			beego.NSNamespace("/live",
				beego.NSInclude(
					&v1.LiveControllerAuth{Handler: h},
				),
			),
			beego.NSNamespace("/ws",
				beego.NSInclude(
					&v1.WebSocketControllerAuth{Handler: h},
				),
			),
		),
	)
	beego.AddNamespace(nsPublic)
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins:     []string{"http://localhost:3000", "https://dev.cast.daystram.com", "https://cast.daystram.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin"},
		AllowCredentials: true,
	}))
}
