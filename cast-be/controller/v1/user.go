package v1

import (
	"fmt"
	"net/http"

	"github.com/astaxie/beego"

	"github.com/daystram/cast/cast-be/constants"
	"github.com/daystram/cast/cast-be/datatransfers"
	"github.com/daystram/cast/cast-be/handlers"
)

type UserControllerAuth struct {
	beego.Controller
	Handler handlers.Handler
	userID  string
}

func (c *UserControllerAuth) Prepare() {
	c.userID = c.Ctx.Input.Param(constants.ContextParamUserID)
}

// @Title Get User Details
// @Success 200 {object} models.Object
// @Param   stub		query	string	false	"stub"
// @router /info [get]
func (c *UserControllerAuth) ProfileDetails(_ string) datatransfers.Response {
	user, err := c.Handler.UserDetails(c.userID)
	if err != nil {
		fmt.Printf("[UserControllerAuth::ProfileInfo] user cannot be found. %+v\n", err)
		return datatransfers.Response{Error: "failed setting stream window", Code: http.StatusInternalServerError}
	}
	return datatransfers.Response{Data: user, Code: http.StatusOK}
}
