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

// Video Content Controller
type VideoController struct {
	beego.Controller
	Handler handlers.Handler
}

// @Title Get List
// @Success 200 {object} models.Object
// @Param   variant		query	string	true	"variant"
// @Param   count		query   int 	false 8	"count"
// @Param   offset		query   int 	false 0	"offset"
// @router /fresh [get]
func (c *VideoController) GetList(variant string, count, offset int) datatransfers.Response {
	videos, err := c.Handler.VideoList(variant, count, offset)
	if err != nil {
		fmt.Printf("[VideoController::GetList] failed retrieving fresh videos. %+v\n", err)
		return datatransfers.Response{Error: "Failed retrieving fresh videos", Code: http.StatusInternalServerError}
	}
	return datatransfers.Response{Data: videos, Code: http.StatusOK}
}

// @Title Search
// @Success 200 {object} models.Object
// @Param   query		query	string	true	"query"
// @router /search [get]
func (c *VideoController) Search(query string) {

}

// @Title Get Details
// @Success 200 {object} models.Object
// @Param   hash		query	string	true	"hash"
// @router /details [get]
func (c *VideoController) GetDetails(hash string) datatransfers.Response {
	// TODO: add view count
	video, err := c.Handler.VideoDetails(hash)
	if err != nil {
		fmt.Printf("[VideoController::GetDetails] failed retrieving video detail. %+v\n", err)
		return datatransfers.Response{Error: "Failed retrieving video detail", Code: http.StatusInternalServerError}
	}
	return datatransfers.Response{Data: video, Code: http.StatusOK}
}

// Video Content Controller
type VideoControllerAuth struct {
	beego.Controller
	Handler handlers.Handler
	userID  primitive.ObjectID
}

func (c *VideoControllerAuth) Prepare() {
	c.userID, _ = primitive.ObjectIDFromHex(c.Ctx.Input.Param(constants.ContextParamUserID))
}

// @Title Upload Video
// @Success 200 {object} models.Object
// @Param   variant		query	string	true	"variant"
// @Param   count		query   int 	false 8	"count"
// @Param   offset		query   int 	false 0	"offset"
// @router /fresh [get]
func (c *VideoControllerAuth) UploadVideo(variant string, count, offset int) datatransfers.Response {
	return datatransfers.Response{}
}
