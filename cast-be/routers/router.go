package routers

import (
	conf "gitlab.com/daystram/cast/cast-be/config"
	v1 "gitlab.com/daystram/cast/cast-be/controller/v1"

	"github.com/astaxie/beego"
)

func init() {
	conf.InitializeAppConfig()

	nsPublic := beego.NewNamespace("api/v1",
		beego.NSNamespace("/ping",
			beego.NSInclude(
				&v1.PingController{},
			),
		),
	)

	beego.AddNamespace(nsPublic)
}
