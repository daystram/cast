package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["gitlab.com/daystram/cast/cast-be/controller/v1:LiveController"] = append(beego.GlobalControllerRouter["gitlab.com/daystram/cast/cast-be/controller/v1:LiveController"],
        beego.ControllerComments{
            Method: "PlayLive",
            Router: `/:username`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(
				param.New("username", param.InPath),
			),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["gitlab.com/daystram/cast/cast-be/controller/v1:PingController"] = append(beego.GlobalControllerRouter["gitlab.com/daystram/cast/cast-be/controller/v1:PingController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["gitlab.com/daystram/cast/cast-be/controller/v1:VideoController"] = append(beego.GlobalControllerRouter["gitlab.com/daystram/cast/cast-be/controller/v1:VideoController"],
        beego.ControllerComments{
            Method: "GetDetails",
            Router: `/details`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["gitlab.com/daystram/cast/cast-be/controller/v1:VideoController"] = append(beego.GlobalControllerRouter["gitlab.com/daystram/cast/cast-be/controller/v1:VideoController"],
        beego.ControllerComments{
            Method: "GetList",
            Router: `/fresh`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(
				param.New("variant", param.IsRequired),
				param.New("count", param.Default("false")),
				param.New("offset", param.Default("false")),
			),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["gitlab.com/daystram/cast/cast-be/controller/v1:VideoController"] = append(beego.GlobalControllerRouter["gitlab.com/daystram/cast/cast-be/controller/v1:VideoController"],
        beego.ControllerComments{
            Method: "Search",
            Router: `/search`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
