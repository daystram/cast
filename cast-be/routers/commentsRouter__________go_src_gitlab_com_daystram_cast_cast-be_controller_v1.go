package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["gitlab.com/daystram/cast/cast-be/controller/v1:PingController"] = append(beego.GlobalControllerRouter["gitlab.com/daystram/cast/cast-be/controller/v1:PingController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
