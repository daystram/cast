package routers

import (
	"context"
	"log"
	"time"

	conf "gitlab.com/daystram/cast/cast-be/config"
	v1 "gitlab.com/daystram/cast/cast-be/controller/v1"
	"gitlab.com/daystram/cast/cast-be/handlers"

	"github.com/astaxie/beego"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	conf.InitializeAppConfig()

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	db, err := mongo.Connect(ctx, options.Client().ApplyURI(conf.AppConfig.MongoDBURI))
	if err != nil {
		log.Fatalf("Failed connecting to mongoDB at %s\n", conf.AppConfig.MongoDBURI)
	}

	h := handlers.NewHandler(handlers.Component{DB: db})

	nsPublic := beego.NewNamespace("api/v1",
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
		beego.NSNamespace("/ws",
			beego.NSInclude(
				&v1.WebSocketController{Handler: h},
			),
		),
	)

	beego.AddNamespace(nsPublic)
}
