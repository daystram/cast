package main

import (
	_ "gitlab.com/daystram/cast/cast-be/routers"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
)

func main() {
	//_ = logs.SetLogger(logs.AdapterFile, `{"filename":"logs/test.log", "separate":["error"], "level": 3}`)

	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		ExposeHeaders:    []string{"*", "Authorization"},
		AllowCredentials: true,
	}))

	beego.Run()
}
