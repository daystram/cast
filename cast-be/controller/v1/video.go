package v1

import (
	"fmt"
	"gitlab.com/daystram/cast/cast-be/config"
	"net/http"
	"os"

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
// @Param   variant		query	string	false	"variant"
// @Param   author		query	string	false	"author"
// @Param   count		query   int 	false 8	"count"
// @Param   offset		query   int 	false 0	"offset"
// @router /list [get]
func (c *VideoController) GetList(variant, author string, count, offset int) datatransfers.Response {
	var videos []datatransfers.Video
	var err error
	if author == "" {
		videos, err = c.Handler.FreshList(variant, count, offset)
	} else {
		videos, err = c.Handler.AuthorList(author, count, offset)
	}
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
// @router /upload [post]
func (c *VideoControllerAuth) UploadVideo() datatransfers.Response {
	upload := datatransfers.VideoUpload{}
	err := c.ParseForm(&upload)
	if err != nil {
		fmt.Printf("[VideoController::UploadVideo] failed parsing video details. %+v\n", err)
		return datatransfers.Response{Error: "Failed parsing video detail", Code: http.StatusInternalServerError}
	}
	var videoID primitive.ObjectID
	videoID, err = c.Handler.CreateVOD(upload, c.userID)
	if err != nil {
		fmt.Printf("[VideoController::UploadVideo] failed creating video. %+v\n", err)
		return datatransfers.Response{Error: "Failed creating video", Code: http.StatusInternalServerError}
	}
	_ = os.Mkdir(fmt.Sprintf("cast-uploaded-videos/%s", videoID.Hex()), 755)
	err = c.SaveToFile("video", fmt.Sprintf("%s/%s/video_original.mp4", config.AppConfig.UploadsDirectory, videoID.Hex()))
	if err != nil {
		_ = c.Handler.DeleteVideo(videoID)
		fmt.Printf("[VideoController::UploadVideo] failed saving video file. %+v\n", err)
		return datatransfers.Response{Error: "Failed saving video file", Code: http.StatusInternalServerError}
	}
	err = c.SaveToFile("thumbnail", fmt.Sprintf("%s/thumbnail/%s.jpg", config.AppConfig.UploadsDirectory, videoID.Hex()))
	if err != nil {
		_ = c.Handler.DeleteVideo(videoID)
		fmt.Printf("[VideoController::UploadVideo] failed saving thumbnail file. %+v\n", err)
		return datatransfers.Response{Error: "Failed saving thumbnail file", Code: http.StatusInternalServerError}
	}
	// TODO: push for transcoding by cast-is
	return datatransfers.Response{Code: http.StatusOK}
}
