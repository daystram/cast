package routers

import (
	"context"
	"fmt"
	"log"
	"time"

	conf "gitlab.com/daystram/cast/cast-be/config"
	"gitlab.com/daystram/cast/cast-be/controller/middleware"
	v1 "gitlab.com/daystram/cast/cast-be/controller/v1"
	"gitlab.com/daystram/cast/cast-be/handlers"

	"github.com/astaxie/beego"
	"github.com/nareix/joy4/format"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	conf.InitializeAppConfig()

	// Init MongoDB
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	db, err := mongo.Connect(ctx, options.Client().ApplyURI(conf.AppConfig.MongoDBURI))
	if err != nil {
		log.Fatalf("Failed connecting to mongoDB at %s\n", conf.AppConfig.MongoDBURI)
	}
	fmt.Printf("[Initialization] MongoDB connected\n")

	// Init RTMP UpLink
	format.RegisterAll()

	h := handlers.NewHandler(handlers.Component{DB: db})
	h.CreateRTMPUpLink()
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
			beego.NSNamespace("/video",
				beego.NSInclude(
					&v1.VideoControllerAuth{Handler: h},
				),
			), beego.NSNamespace("/live",
				beego.NSInclude(
					&v1.LiveControllerAuth{Handler: h},
				),
			),
		),

	)

	beego.AddNamespace(nsPublic)
}
