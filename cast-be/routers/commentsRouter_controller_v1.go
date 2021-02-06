package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/daystram/cast/cast-be/controller/v1:AuthControllerAuth"] = append(beego.GlobalControllerRouter["github.com/daystram/cast/cast-be/controller/v1:AuthControllerAuth"],
        beego.ControllerComments{
            Method: "PostCheckRegister",
            Router: "/check",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(
				param.New("info", param.IsRequired, param.InBody),
			),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/daystram/cast/cast-be/controller/v1:LiveController"] = append(beego.GlobalControllerRouter["github.com/daystram/cast/cast-be/controller/v1:LiveController"],
        beego.ControllerComments{
            Method: "PlayLive",
            Router: "/stream/:username",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(
				param.New("username", param.InPath),
			),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/daystram/cast/cast-be/controller/v1:LiveControllerAuth"] = append(beego.GlobalControllerRouter["github.com/daystram/cast/cast-be/controller/v1:LiveControllerAuth"],
        beego.ControllerComments{
            Method: "GetWindow",
            Router: "/window",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(
				param.New("_"),
			),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/daystram/cast/cast-be/controller/v1:LiveControllerAuth"] = append(beego.GlobalControllerRouter["github.com/daystram/cast/cast-be/controller/v1:LiveControllerAuth"],
        beego.ControllerComments{
            Method: "ControlWindow",
            Router: "/window",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(
				param.New("open", param.IsRequired),
			),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/daystram/cast/cast-be/controller/v1:PingController"] = append(beego.GlobalControllerRouter["github.com/daystram/cast/cast-be/controller/v1:PingController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: "/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/daystram/cast/cast-be/controller/v1:UserControllerAuth"] = append(beego.GlobalControllerRouter["github.com/daystram/cast/cast-be/controller/v1:UserControllerAuth"],
        beego.ControllerComments{
            Method: "ProfileDetails",
            Router: "/info",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(
				param.New("_"),
			),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/daystram/cast/cast-be/controller/v1:VideoController"] = append(beego.GlobalControllerRouter["github.com/daystram/cast/cast-be/controller/v1:VideoController"],
        beego.ControllerComments{
            Method: "GetDetails",
            Router: "/details",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(
				param.New("hash", param.IsRequired),
				param.New("username"),
			),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/daystram/cast/cast-be/controller/v1:VideoController"] = append(beego.GlobalControllerRouter["github.com/daystram/cast/cast-be/controller/v1:VideoController"],
        beego.ControllerComments{
            Method: "GetList",
            Router: "/list",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(
				param.New("variant"),
				param.New("author"),
				param.New("count", param.Default("false")),
				param.New("offset", param.Default("false")),
			),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/daystram/cast/cast-be/controller/v1:VideoController"] = append(beego.GlobalControllerRouter["github.com/daystram/cast/cast-be/controller/v1:VideoController"],
        beego.ControllerComments{
            Method: "Search",
            Router: "/search",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(
				param.New("query", param.IsRequired),
				param.New("count", param.Default("false")),
				param.New("offset", param.Default("false")),
			),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/daystram/cast/cast-be/controller/v1:VideoControllerAuth"] = append(beego.GlobalControllerRouter["github.com/daystram/cast/cast-be/controller/v1:VideoControllerAuth"],
        beego.ControllerComments{
            Method: "GetCheckUnique",
            Router: "/check",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(
				param.New("title", param.IsRequired),
			),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/daystram/cast/cast-be/controller/v1:VideoControllerAuth"] = append(beego.GlobalControllerRouter["github.com/daystram/cast/cast-be/controller/v1:VideoControllerAuth"],
        beego.ControllerComments{
            Method: "CommentVideo",
            Router: "/comment",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(
				param.New("info", param.IsRequired, param.InBody),
			),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/daystram/cast/cast-be/controller/v1:VideoControllerAuth"] = append(beego.GlobalControllerRouter["github.com/daystram/cast/cast-be/controller/v1:VideoControllerAuth"],
        beego.ControllerComments{
            Method: "DeleteVideo",
            Router: "/delete",
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(
				param.New("hash", param.IsRequired),
			),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/daystram/cast/cast-be/controller/v1:VideoControllerAuth"] = append(beego.GlobalControllerRouter["github.com/daystram/cast/cast-be/controller/v1:VideoControllerAuth"],
        beego.ControllerComments{
            Method: "EditVideo",
            Router: "/edit",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(
				param.New("_"),
			),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/daystram/cast/cast-be/controller/v1:VideoControllerAuth"] = append(beego.GlobalControllerRouter["github.com/daystram/cast/cast-be/controller/v1:VideoControllerAuth"],
        beego.ControllerComments{
            Method: "LikeVideo",
            Router: "/like",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(
				param.New("info", param.IsRequired, param.InBody),
			),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/daystram/cast/cast-be/controller/v1:VideoControllerAuth"] = append(beego.GlobalControllerRouter["github.com/daystram/cast/cast-be/controller/v1:VideoControllerAuth"],
        beego.ControllerComments{
            Method: "GetList",
            Router: "/list",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(
				param.New("variant"),
				param.New("author"),
				param.New("count", param.Default("false")),
				param.New("offset", param.Default("false")),
			),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/daystram/cast/cast-be/controller/v1:VideoControllerAuth"] = append(beego.GlobalControllerRouter["github.com/daystram/cast/cast-be/controller/v1:VideoControllerAuth"],
        beego.ControllerComments{
            Method: "SubscribeAuthor",
            Router: "/subscribe",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(
				param.New("info", param.IsRequired, param.InBody),
			),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/daystram/cast/cast-be/controller/v1:VideoControllerAuth"] = append(beego.GlobalControllerRouter["github.com/daystram/cast/cast-be/controller/v1:VideoControllerAuth"],
        beego.ControllerComments{
            Method: "UploadVideo",
            Router: "/upload",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(
				param.New("_"),
			),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/daystram/cast/cast-be/controller/v1:WebSocketController"] = append(beego.GlobalControllerRouter["github.com/daystram/cast/cast-be/controller/v1:WebSocketController"],
        beego.ControllerComments{
            Method: "Connect",
            Router: "/chat/:hash",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(
				param.New("hash", param.InPath),
				param.New("_"),
			),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/daystram/cast/cast-be/controller/v1:WebSocketControllerAuth"] = append(beego.GlobalControllerRouter["github.com/daystram/cast/cast-be/controller/v1:WebSocketControllerAuth"],
        beego.ControllerComments{
            Method: "ConnectChat",
            Router: "/chat/:hash",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(
				param.New("hash", param.InPath),
				param.New("_"),
			),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/daystram/cast/cast-be/controller/v1:WebSocketControllerAuth"] = append(beego.GlobalControllerRouter["github.com/daystram/cast/cast-be/controller/v1:WebSocketControllerAuth"],
        beego.ControllerComments{
            Method: "ConnectNotification",
            Router: "/notification",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(
				param.New("_"),
			),
            Filters: nil,
            Params: nil})

}
