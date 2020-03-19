package v1

import (
	"fmt"
	"gitlab.com/daystram/cast/cast-be/constants"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"

	"gitlab.com/daystram/cast/cast-be/datatransfers"
	"gitlab.com/daystram/cast/cast-be/handlers"

	"github.com/astaxie/beego"
)

// Live Stream Upstream Controller
type LiveController struct {
	beego.Controller
	Handler handlers.Handler
}

// @Title Receive RTMP UpLink
// @Success 200 {object} models.Object
// @router /stream/:username  [get]
func (c *LiveController) PlayLive(username string) datatransfers.Response {
	if err := c.Handler.StreamLive(username, c.Ctx.ResponseWriter, c.Ctx.Request); err != nil {
		fmt.Printf("[LiveController::PlayLive] failed playing live stream from %s. %+v\n", username, err)
		return datatransfers.Response{Error: "failed playing live stream", Code: http.StatusInternalServerError}
	}
	return datatransfers.Response{}
}

type LiveControllerAuth struct {
	beego.Controller
	Handler handlers.Handler
	userID  primitive.ObjectID
}

func (c *LiveControllerAuth) Prepare() {
	c.userID, _ = primitive.ObjectIDFromHex(c.Ctx.Input.Param(constants.ContextParamUserID))
}

// @Title Receive RTMP UpLink
// @Success 200 {object} models.Object
// @Param   open		query	bool	true	"open"
// @router /window  [get]
func (c *LiveControllerAuth) ControlWindow(open bool) datatransfers.Response {
	if err := c.Handler.ControlUpLinkWindow(c.userID, open); err != nil {
		fmt.Printf("[LiveController::ControlWindow] failed setting stream window for %s. %+v\n", c.userID, err)
		return datatransfers.Response{Error: "failed setting stream window", Code: http.StatusInternalServerError}
	}
	return datatransfers.Response{}
}
