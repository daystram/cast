package main

import (
	_ "gitlab.com/daystram/cast/cast-be/routers"

	"github.com/astaxie/beego"
)

func main() {
	//_ = logs.SetLogger(logs.AdapterFile, `{"filename":"logs/test.log", "separate":["error"], "level": 3}`)
	beego.Run()
}
