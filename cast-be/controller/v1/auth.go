package v1

import (
	"github.com/astaxie/beego"

	"gitlab.com/daystram/cast/cast-be/handlers"
)

// Authentication Controller
type AuthController struct {
	beego.Controller
	Handler handlers.Handler
}
