package v1

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"gitlab.com/daystram/cast/cast-be/config"
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
// @Param   username	query	string	false	"username"
// @router /details [get]
func (c *VideoController) GetDetails(hash, username string) datatransfers.Response {
	video, err := c.Handler.VideoDetails(hash)
	if err != nil {
		fmt.Printf("[VideoController::GetDetails] failed retrieving video detail. %+v\n", err)
		return datatransfers.Response{Error: "Failed retrieving video detail", Code: http.StatusInternalServerError}
	}
	if username != "" {
		video.Liked, _ = c.Handler.CheckUserLikes(hash, username)
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

// @Title Check Title
// @Param   title    query	string	true	"title"
// @Success 200 success
// @router /check [get]
func (c *VideoControllerAuth) GetCheckUnique(title string) datatransfers.Response {
	err := c.Handler.CheckUniqueVideoTitle(title)
	if err != nil {
		log.Printf("[VideoControllerAuth::GetCheckUnique] title already used. %+v\n", err)
		return datatransfers.Response{Error: "Title already used", Code: http.StatusConflict}
	}
	return datatransfers.Response{Code: http.StatusOK}
}

// @Title Edit Video
// @Success 200 {object} models.Object
// @Param   video    body	{datatransfers.VideoEditForm}	true	"video"
// @router /edit [put]
func (c *VideoControllerAuth) EditVideo(video datatransfers.VideoEditForm) datatransfers.Response {
	fmt.Println(video.Tags)
	err := c.Handler.UpdateVideo(datatransfers.VideoEdit{
		Hash:        video.Hash,
		Title:       video.Title,
		Description: video.Description,
		Tags:        strings.Split(video.Tags, ","),
	}, c.userID)
	if err != nil {
		fmt.Printf("[VideoController::EditVideo] failed editing video. %+v\n", err)
		return datatransfers.Response{Error: "Failed editing video", Code: http.StatusInternalServerError}
	}
	return datatransfers.Response{Code: http.StatusOK}
}

// @Title Delete Video
// @Success 200 {object} models.Object
// @Param   hash		query	string	true	"hash"
// @router /delete [delete]
func (c *VideoControllerAuth) DeleteVideo(hash string) datatransfers.Response {
	videoID, err := primitive.ObjectIDFromHex(hash)
	if err != nil {
		fmt.Printf("[VideoController::DeleteVideo] invalid video hash. %+v\n", err)
		return datatransfers.Response{Error: "Invalid video hash", Code: http.StatusInternalServerError}
	}
	err = c.Handler.DeleteVideo(videoID, c.userID)
	if err != nil {
		fmt.Printf("[VideoController::DeleteVideo] failed deleting video. %+v\n", err)
		return datatransfers.Response{Error: "Failed deleting video", Code: http.StatusInternalServerError}
	}
	return datatransfers.Response{Code: http.StatusOK}
}

// @Title Upload Video
// @Success 200 {object} models.Object
// @router /upload [post]
func (c *VideoControllerAuth) UploadVideo() datatransfers.Response {
	upload := datatransfers.VideoUploadForm{}
	err := c.ParseForm(&upload)
	if err != nil {
		fmt.Printf("[VideoController::UploadVideo] failed parsing video details. %+v\n", err)
		return datatransfers.Response{Error: "Failed parsing video detail", Code: http.StatusInternalServerError}
	}
	var videoID primitive.ObjectID
	videoID, err = c.Handler.CreateVOD(datatransfers.VideoUpload{
		Title:       upload.Title,
		Description: upload.Description,
		Tags:        strings.Split(upload.Tags, ","),
	}, c.userID)
	if err != nil {
		fmt.Printf("[VideoController::UploadVideo] failed creating video. %+v\n", err)
		return datatransfers.Response{Error: "Failed creating video", Code: http.StatusInternalServerError}
	}
	_ = os.Mkdir(fmt.Sprintf("%s/%s", config.AppConfig.UploadsDirectory, videoID.Hex()), 755)
	err = c.SaveToFile("video", fmt.Sprintf("%s/%s/video.mp4", config.AppConfig.UploadsDirectory, videoID.Hex()))
	if err != nil {
		_ = c.Handler.DeleteVideo(videoID, c.userID)
		fmt.Printf("[VideoController::UploadVideo] failed saving video file. %+v\n", err)
		return datatransfers.Response{Error: "Failed saving video file", Code: http.StatusInternalServerError}
	}
	err = c.SaveToFile("thumbnail", fmt.Sprintf("%s/thumbnail/%s.ori", config.AppConfig.UploadsDirectory, videoID.Hex()))
	if err != nil {
		_ = c.Handler.DeleteVideo(videoID, c.userID)
		fmt.Printf("[VideoController::UploadVideo] failed saving thumbnail file. %+v\n", err)
		return datatransfers.Response{Error: "Failed saving thumbnail file", Code: http.StatusInternalServerError}
	}
	err = c.Handler.NormalizeThumbnail(videoID)
	if err != nil {
		_ = c.Handler.DeleteVideo(videoID, c.userID)
		fmt.Printf("[VideoController::UploadVideo] failed normalizing thumbnail image. %+v\n", err)
		return datatransfers.Response{Error: "Failed normalizing thumbnail image", Code: http.StatusInternalServerError}
	}
	c.Handler.StartTranscode(videoID.Hex())
	return datatransfers.Response{Code: http.StatusOK}
}

// @Title Like Video
// @Success 200 {object} models.Object
// @Param   hash		query	string	true	"hash"
// @Param   like		query	bool	true	"like"
// @router /like [get]
func (c *VideoControllerAuth) LikeVideo(hash string, like bool) datatransfers.Response {
	err := c.Handler.LikeVideo(c.userID, hash, like)
	if err != nil {
		fmt.Printf("[VideoController::LikeVideo] failed liking video. %+v\n", err)
		return datatransfers.Response{Error: "Already liked", Code: http.StatusConflict}
	}
	return datatransfers.Response{Code: http.StatusOK}
}

// @Title Content Video
// @Success 200 {object} models.Object
// @Param   hash		query	string	true	"hash"
// @Param   content		query	string	true	"content"
// @router /comment [get]
func (c *VideoControllerAuth) CommentVideo(hash, content string) datatransfers.Response {
	comment, err := c.Handler.CommentVideo(c.userID, hash, content)
	if err != nil {
		fmt.Printf("[VideoController::CommentVideo] failed liking video. %+v\n", err)
		return datatransfers.Response{Error: "Failed commenting video", Code: http.StatusInternalServerError}
	}
	return datatransfers.Response{Data: comment, Code: http.StatusOK}
}
