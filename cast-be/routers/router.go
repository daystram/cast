package routers

import (
	"context"
	"fmt"
	"log"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
	"github.com/nareix/joy4/format"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	conf "github.com/daystram/cast/cast-be/config"
	"github.com/daystram/cast/cast-be/controller/middleware"
	v1 "github.com/daystram/cast/cast-be/controller/v1"
	"github.com/daystram/cast/cast-be/handlers"
)

func init() {
	conf.InitializeAppConfig()

	var err error
	// Init MongoDB
	var db *mongo.Client
	if db, err = mongo.Connect(context.Background(), options.Client().ApplyURI(conf.AppConfig.MongoDBURI)); err != nil {
		log.Fatalf("[Initialization] Failed connecting to MongoDB at %s. %+v\n", conf.AppConfig.MongoDBURI, err)
	}
	log.Printf("[Initialization] Successfully connected to MongoDB!\n")

	// Init RTMP Formats
	format.RegisterAll()

	// TODO: initialize S3

	// Init RabbitMQ
	var mqConn *amqp.Connection
	if mqConn, err = amqp.Dial(conf.AppConfig.RabbitMQURI); err != nil {
		log.Fatalf("[Initialization] Failed connecting to RabbitMQ at %s. %+v\n", conf.AppConfig.RabbitMQURI, err)
	}
	var mq *amqp.Channel
	if mq, err = mqConn.Channel(); err != nil {
		log.Fatalf("[Initialization] Failed opening RabbitMQ channel. %+v\n", err)
	}
	if _, err = mq.QueueDeclare(conf.AppConfig.RabbitMQQueueTask, true, false, false, false, nil); err != nil {
		log.Fatalf("[Initialization] Failed declaring RabbitMQ queue %s. %+v\n", conf.AppConfig.RabbitMQQueueTask, err)
	}
	if _, err = mq.QueueDeclare(conf.AppConfig.RabbitMQQueueProgress, true, false, false, false, nil); err != nil {
		log.Fatalf("[Initialization] Failed declaring RabbitMQ queue %s. %+v\n", conf.AppConfig.RabbitMQQueueProgress, err)
	}
	log.Printf("[Initialization] Successfully connected to RabbitMQ!\n")

	// Init Handler
	h := handlers.NewHandler(handlers.Component{DB: db, MQClient: nil})

	// Init RTMP Uplink
	h.CreateRTMPUpLink()

	// Init Transcoder Listener
	//go h.TranscodeListenerWorker()
	fmt.Printf("[Initialization] Initialization complete!\n")

	apiV1 := beego.NewNamespace("api/v1",
		beego.NSNamespace("/ping",
			beego.NSInclude(
				&v1.PingController{},
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
			beego.NSBefore(middleware.AuthenticateAccessToken),
			beego.NSNamespace("/auth",
				beego.NSInclude(
					&v1.AuthControllerAuth{Handler: h},
				),
			),
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
	beego.AddNamespace(apiV1)
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins:     []string{"http://localhost:3000", "https://dev.cast.daystram.com", "https://cast.daystram.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin"},
		AllowCredentials: true,
	}))
}
