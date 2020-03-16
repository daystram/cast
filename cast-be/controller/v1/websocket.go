package v1

import (
	"github.com/astaxie/beego"

	"gitlab.com/daystram/cast/cast-be/handlers"
)

// WebSocket Feed Controller
type WebSocketController struct {
	beego.Controller
	Handler handlers.Handler
}
