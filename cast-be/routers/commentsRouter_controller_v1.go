package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["gitlab.com/daystram/cast/cast-be/controller/v1:AuthController"] = append(beego.GlobalControllerRouter["gitlab.com/daystram/cast/cast-be/controller/v1:AuthController"],
        beego.ControllerComments{
            Method: "PostCheckUnique",
            Router: `/check`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(
				param.New("info", param.IsRequired, param.InBody),
			),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["gitlab.com/daystram/cast/cast-be/controller/v1:AuthController"] = append(beego.GlobalControllerRouter["gitlab.com/daystram/cast/cast-be/controller/v1:AuthController"],
        beego.ControllerComments{
            Method: "PostAuthenticate",
            Router: `/login`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(
				param.New("info", param.IsRequired, param.InBody),
			),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["gitlab.com/daystram/cast/cast-be/controller/v1:AuthController"] = append(beego.GlobalControllerRouter["gitlab.com/daystram/cast/cast-be/controller/v1:AuthController"],
        beego.ControllerComments{
            Method: "PostDeAuthenticate",
            Router: `/logout`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["gitlab.com/daystram/cast/cast-be/controller/v1:AuthController"] = append(beego.GlobalControllerRouter["gitlab.com/daystram/cast/cast-be/controller/v1:AuthController"],
        beego.ControllerComments{
            Method: "PostRegister",
            Router: `/signup`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(
				param.New("info", param.IsRequired, param.InBody),
			),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["gitlab.com/daystram/cast/cast-be/controller/v1:LiveController"] = append(beego.GlobalControllerRouter["gitlab.com/daystram/cast/cast-be/controller/v1:LiveController"],
        beego.ControllerComments{
            Method: "PlayLive",
            Router: `/stream/:username`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(
				param.New("username", param.InPath),
			),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["gitlab.com/daystram/cast/cast-be/controller/v1:LiveControllerAuth"] = append(beego.GlobalControllerRouter["gitlab.com/daystram/cast/cast-be/controller/v1:LiveControllerAuth"],
        beego.ControllerComments{
            Method: "ControlWindow",
            Router: `/window`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(
				param.New("open", param.IsRequired),
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
            MethodParams: param.Make(
				param.New("hash", param.IsRequired),
			),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["gitlab.com/daystram/cast/cast-be/controller/v1:VideoController"] = append(beego.GlobalControllerRouter["gitlab.com/daystram/cast/cast-be/controller/v1:VideoController"],
        beego.ControllerComments{
            Method: "GetList",
            Router: `/list`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(
				param.New("variant"),
				param.New("author"),
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
            MethodParams: param.Make(
				param.New("query", param.IsRequired),
			),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["gitlab.com/daystram/cast/cast-be/controller/v1:VideoControllerAuth"] = append(beego.GlobalControllerRouter["gitlab.com/daystram/cast/cast-be/controller/v1:VideoControllerAuth"],
        beego.ControllerComments{
            Method: "GetCheckUnique",
            Router: `/check`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(
				param.New("title", param.IsRequired),
			),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["gitlab.com/daystram/cast/cast-be/controller/v1:VideoControllerAuth"] = append(beego.GlobalControllerRouter["gitlab.com/daystram/cast/cast-be/controller/v1:VideoControllerAuth"],
        beego.ControllerComments{
            Method: "DeleteVideo",
            Router: `/delete`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(
				param.New("hash", param.IsRequired),
			),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["gitlab.com/daystram/cast/cast-be/controller/v1:VideoControllerAuth"] = append(beego.GlobalControllerRouter["gitlab.com/daystram/cast/cast-be/controller/v1:VideoControllerAuth"],
        beego.ControllerComments{
            Method: "EditVideo",
            Router: `/edit`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(
				param.New("video", param.IsRequired, param.InBody),
			),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["gitlab.com/daystram/cast/cast-be/controller/v1:VideoControllerAuth"] = append(beego.GlobalControllerRouter["gitlab.com/daystram/cast/cast-be/controller/v1:VideoControllerAuth"],
        beego.ControllerComments{
            Method: "UploadVideo",
            Router: `/upload`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
