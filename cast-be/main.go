package main

import (
	_ "gitlab.com/daystram/cast/cast-be/routers"

	"github.com/astaxie/beego"
)

func main() {
	beego.Run()
}
