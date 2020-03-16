package v1

import (
	"fmt"
	"gitlab.com/daystram/cast/cast-be/datatransfers"
	"gitlab.com/daystram/cast/cast-be/handlers"

	"github.com/astaxie/beego"
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
	videos, err := c.Handler.GetVideo(variant, count, offset)
	if err != nil {
		fmt.Printf("[VideoController::GetList] failed retrieving fresh videos. %+v\n", err)
		return datatransfers.Response{Error: "Failed retrieving fresh videos", Code: 500}
	}
	return datatransfers.Response{Data: videos, Code: 200}
}

// @Title Search
// @Success 200 {object} models.Object
// @router /search [get]
func (c *VideoController) Search() {

}

// @Title Get Details
// @Success 200 {object} models.Object
// @router /details [get]
func (c *VideoController) GetDetails() {

}
