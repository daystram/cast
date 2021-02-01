package v1

import (
	"log"
	"net/http"

	"github.com/astaxie/beego"

	"github.com/daystram/cast/cast-be/datatransfers"
	"github.com/daystram/cast/cast-be/handlers"
)

type AuthControllerAuth struct {
	beego.Controller
	Handler handlers.Handler
}
