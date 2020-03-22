package v1

import (
	"fmt"
	"net/http"

	"gitlab.com/daystram/cast/cast-be/constants"
	"gitlab.com/daystram/cast/cast-be/datatransfers"
	"gitlab.com/daystram/cast/cast-be/handlers"

	"github.com/astaxie/beego"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserControllerAuth struct {
	beego.Controller
	Handler handlers.Handler
	userID  primitive.ObjectID
}

func (c *UserControllerAuth) Prepare() {
	c.userID, _ = primitive.ObjectIDFromHex(c.Ctx.Input.Param(constants.ContextParamUserID))
}

// @Title Get User Details
// @Success 200 {object} models.Object
// @router /info [get]
func (c *UserControllerAuth) ProfileInfo() datatransfers.Response {
	user, err := c.Handler.UserDetails(c.userID)
	if err != nil {
		fmt.Printf("[UserControllerAuth::ProfileInfo] user cannot be found. %+v\n", err)
		return datatransfers.Response{Error: "failed setting stream window", Code: http.StatusInternalServerError}
	}
	return datatransfers.Response{Data: user, Code: http.StatusOK}
}

// @Title Update User Details
// @Success 200 {object} models.Object
// @Param   info    body	{datatransfers.UserEdit}	true	"user field info"
// @router /edit [put]
func (c *UserControllerAuth) UpdateProfile(info datatransfers.UserEdit) datatransfers.Response {
	err := c.Handler.UpdateUser(info, c.userID)
	if err != nil {
		fmt.Printf("[UserControllerAuth::ProfileInfo] user cannot be found. %+v\n", err)
		return datatransfers.Response{Error: "failed setting stream window", Code: http.StatusInternalServerError}
	}
	return datatransfers.Response{Code: http.StatusOK}
}
