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

// @Title Register user from OID, skip if already exists
// @Param   idToken    body	{datatransfers.UserRegister}	true	"registration info"
// @Success 200 success
// @router /check [post]
func (c *AuthControllerAuth) PostCheckRegister(idToken datatransfers.UserRegister) datatransfers.Response {
	err := c.Handler.Register(idToken)
	if err != nil {
		log.Printf("[AuthControllerAuth::PostCheckRegister] failed registering. %+v\n", err)
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		return datatransfers.Response{Error: "failed registering", Code: http.StatusInternalServerError}
	}
	return datatransfers.Response{Code: http.StatusOK}
}
