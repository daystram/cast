package v1

import (
	"fmt"
	"github.com/astaxie/beego"
	"gitlab.com/daystram/cast/cast-be/datatransfers"

	"gitlab.com/daystram/cast/cast-be/handlers"
)

// Live Stream Upstream Controller
type LiveController struct {
	beego.Controller
	Handler handlers.Handler
}

// @Title Receive RTMP UpLink
// @Success 200 {object} models.Object
// @router /:username  [get]
func (c *LiveController) PlayLive(username string) datatransfers.Response{
	if err := c.Handler.StreamLive(username, c.Ctx.ResponseWriter, c.Ctx.Request); err != nil {
		fmt.Printf("[LiveController::PlayLive] failed playing live stream from %s. %+v\n", username, err)
		return datatransfers.Response{Error: "failed playing live stream", Code:  500}
	}
	return datatransfers.Response{}
}
