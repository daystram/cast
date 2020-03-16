package v1

import (
	"github.com/astaxie/beego"

	"gitlab.com/daystram/cast/cast-be/handlers"
)

// Live Stream Upstream Controller
type LiveController struct {
	beego.Controller
	Handler handlers.Handler
}
