package v1

import (
	"fmt"
	"net/http"

	"github.com/astaxie/beego"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/daystram/cast/cast-be/constants"
	"github.com/daystram/cast/cast-be/datatransfers"
	"github.com/daystram/cast/cast-be/handlers"
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
	return datatransfers.Response{Code: 200}
}

type LiveControllerAuth struct {
	beego.Controller
	Handler handlers.Handler
	userID  primitive.ObjectID
}

func (c *LiveControllerAuth) Prepare() {
	c.userID, _ = primitive.ObjectIDFromHex(c.Ctx.Input.Param(constants.ContextParamUserID))
}

// @Title Get RTMP UpLink Status
// @Success 200 {object} models.Object
// @Param   stub		query	string	false	"stub"
// @router /window  [get]
func (c *LiveControllerAuth) GetWindow(_ string) datatransfers.Response {
	var err error
	var user datatransfers.UserDetail
	if user, err = c.Handler.UserDetails(c.userID); err != nil {
		fmt.Printf("[LiveControllerAuth::GetWindow] failed retrieving user %s info. %+v\n", c.userID.Hex(), err)
		return datatransfers.Response{Error: "failed retrieving user info", Code: http.StatusInternalServerError}
	}
	var live datatransfers.Video
	if live, err = c.Handler.VideoDetails(user.Username); err != nil {
		fmt.Printf("[LiveControllerAuth::GetWindow] failed retrieving stream info. %+v\n", err)
		return datatransfers.Response{Error: "failed retrieving stream info", Code: http.StatusInternalServerError}
	}
	return datatransfers.Response{Data: live.IsLive, Code: 200}
}

// @Title Set RTMP UpLink Window
// @Success 200 {object} models.Object
// @Param   open		query	bool	true	"open"
// @router /window  [put]
func (c *LiveControllerAuth) ControlWindow(open bool) datatransfers.Response {
	if err := c.Handler.ControlUpLinkWindow(c.userID, open); err != nil {
		fmt.Printf("[LiveControllerAuth::ControlWindow] failed setting stream window for %s. %+v\n", c.userID.Hex(), err)
		return datatransfers.Response{Error: "failed setting stream window", Code: http.StatusInternalServerError}
	}
	return datatransfers.Response{Code: 200}
}
