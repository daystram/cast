package v1

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/astaxie/beego"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/daystram/cast/cast-be/constants"
	"github.com/daystram/cast/cast-be/datatransfers"
	"github.com/daystram/cast/cast-be/handlers"
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
		videos, err = c.Handler.CastList(variant, count, offset)
	} else {
		videos, err = c.Handler.AuthorList(author, count, offset)
	}
	if err != nil {
		fmt.Printf("[VideoController::GetList] failed retrieving video list. %+v\n", err)
		return datatransfers.Response{Error: "Failed retrieving video list", Code: http.StatusInternalServerError}
	}
	return datatransfers.Response{Data: videos, Code: http.StatusOK}
}

// @Title Search
// @Success 200 {object} models.Object
// @Param   query		query	string	true	"query"
// @Param   count		query   int 	false 8	"count"
// @Param   offset		query   int 	false 0	"offset"
// @router /search [get]
func (c *VideoController) Search(query string, count, offset int) datatransfers.Response {
	var videos []datatransfers.Video
	var err error
	videos, err = c.Handler.SearchVideo(query, []string{}, count, offset)
	if err != nil {
		fmt.Printf("[VideoController::Search] failed searching videos. %+v\n", err)
		return datatransfers.Response{Error: "Failed searching videos", Code: http.StatusInternalServerError}
	}
	return datatransfers.Response{Data: videos, Code: http.StatusOK}
}

// @Title Get Details
// @Success 200 {object} models.Object
// @Param   hash		query	string	true	"hash"
// @Param   username	query	string	false	"username"
// @router /details [get]
func (c *VideoController) GetDetails(hash, username string) datatransfers.Response {
	video, err := c.Handler.VideoDetails(hash)
	if err != nil {
		fmt.Printf("[VideoController::GetDetails] video not found. %+v\n", err)
		return datatransfers.Response{Code: http.StatusNotFound}
	}
	if username != "" {
		video.Liked, _ = c.Handler.CheckUserLikes(video.Hash, username)
		video.Subscribed, _ = c.Handler.CheckUserSubscribes(video.Author.ID, username)
	}
	return datatransfers.Response{Data: video, Code: http.StatusOK}
}

// Video Content Controller
type VideoControllerAuth struct {
	beego.Controller
	Handler handlers.Handler
	userID  string
}

func (c *VideoControllerAuth) Prepare() {
	c.userID = c.Ctx.Input.Param(constants.ContextParamUserID)
}

// @Title Get List
// @Success 200 {object} models.Object
// @Param   variant		query	string	false	"variant"
// @Param   count		query   int 	false 8	"count"
// @Param   offset		query   int 	false 0	"offset"
// @router /list [get]
func (c *VideoControllerAuth) GetList(variant string, count, offset int) datatransfers.Response {
	var videos []datatransfers.Video
	var err error
	videos, err = c.Handler.CastList(variant, count, offset, c.userID)
	if err != nil {
		fmt.Printf("[VideoControllerAuth::GetList] failed retrieving video list. %+v\n", err)
		return datatransfers.Response{Error: "Failed retrieving video list", Code: http.StatusInternalServerError}
	}
	return datatransfers.Response{Data: videos, Code: http.StatusOK}
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
// @Param   stub		query	string	false	"stub"
// @router /edit [put]
func (c *VideoControllerAuth) EditVideo(_ string) datatransfers.Response {
	video := datatransfers.VideoEditForm{}
	err := c.ParseForm(&video)
	if err != nil {
		fmt.Printf("[VideoControllerAuth::EditVideo] failed parsing video details. %+v\n", err)
		return datatransfers.Response{Error: "Failed parsing video detail", Code: http.StatusInternalServerError}
	}
	if matchVideo, err := c.Handler.VideoDetails(video.Hash); err != nil {
		log.Printf("[VideoControllerAuth::EditVideo] failed retrieving video. %+v\n", err)
		return datatransfers.Response{Error: "Failed editing video", Code: http.StatusInternalServerError}
	} else if matchVideo.Title != video.Title {
		if err = c.Handler.CheckUniqueVideoTitle(video.Title); err != nil {
			log.Printf("[VideoControllerAuth::EditVideo] title already used. %+v\n", err)
			return datatransfers.Response{Error: "Title already used", Code: http.StatusConflict}
		}
	}
	err = c.Handler.UpdateVideo(datatransfers.VideoEdit{
		Hash:        video.Hash,
		Title:       video.Title,
		Description: video.Description,
		Tags:        strings.Split(video.Tags, ","),
	}, c.Controller, c.userID)
	if err != nil {
		fmt.Printf("[VideoControllerAuth::EditVideo] failed editing video. %+v\n", err)
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
		fmt.Printf("[VideoControllerAuth::DeleteVideo] invalid video hash. %+v\n", err)
		return datatransfers.Response{Error: "Invalid video hash", Code: http.StatusInternalServerError}
	}
	err = c.Handler.DeleteVideo(videoID, c.userID)
	if err != nil {
		fmt.Printf("[VideoControllerAuth::DeleteVideo] failed deleting video. %+v\n", err)
		return datatransfers.Response{Error: "Failed deleting video", Code: http.StatusInternalServerError}
	}
	return datatransfers.Response{Code: http.StatusOK}
}

// @Title Upload Video
// @Success 200 {object} models.Object
// @Param   stub		query	string	false	"stub"
// @router /upload [post]
func (c *VideoControllerAuth) UploadVideo(_ string) datatransfers.Response {
	upload := datatransfers.VideoUploadForm{}
	err := c.ParseForm(&upload)
	if err != nil {
		fmt.Printf("[VideoControllerAuth::UploadVideo] failed parsing video details. %+v\n", err)
		return datatransfers.Response{Error: "Failed parsing video detail", Code: http.StatusInternalServerError}
	}
	err = c.Handler.CheckUniqueVideoTitle(upload.Title)
	if err != nil {
		log.Printf("[VideoControllerAuth::UploadVideo] title already used. %+v\n", err)
		return datatransfers.Response{Error: "Title already used", Code: http.StatusConflict}
	}
	var videoID primitive.ObjectID
	videoID, err = c.Handler.CreateVOD(datatransfers.VideoUpload{
		Title:       upload.Title,
		Description: upload.Description,
		Tags:        strings.Split(upload.Tags, ","),
	}, c.Controller, c.userID)
	if err != nil {
		fmt.Printf("[VideoControllerAuth::UploadVideo] failed creating video. %+v\n", err)
		return datatransfers.Response{Error: "Failed creating video", Code: http.StatusInternalServerError}
	}
	// Trigger transcode sequence by cast-is
	c.Handler.StartTranscode(videoID.Hex())
	return datatransfers.Response{Code: http.StatusOK}
}

// @Title Like Video
// @Success 200 {object} models.Object
// @Param   info    body	{datatransfers.LikeBody}	true	"body"
// @router /like [post]
func (c *VideoControllerAuth) LikeVideo(info datatransfers.LikeBody) datatransfers.Response {
	err := c.Handler.LikeVideo(c.userID, info.Hash, info.Like)
	if err != nil {
		fmt.Printf("[VideoController::LikeVideo] failed liking video. %+v\n", err)
		return datatransfers.Response{Error: "Failed liking video", Code: http.StatusConflict}
	}
	return datatransfers.Response{Code: http.StatusOK}
}

// @Title Subscribe
// @Success 200 {object} models.Object
// @Param   info    body	{datatransfers.SubscribeBody}	true	"body"
// @router /subscribe [post]
func (c *VideoControllerAuth) SubscribeAuthor(info datatransfers.SubscribeBody) datatransfers.Response {
	err := c.Handler.Subscribe(c.userID, info.AuthorUsername, info.Subscribe)
	if err != nil {
		fmt.Printf("[VideoControllerAuth::SubscribeAuthor] failed subscribing author. %+v\n", err)
		return datatransfers.Response{Error: "Failed subscribing author", Code: http.StatusConflict}
	}
	return datatransfers.Response{Code: http.StatusOK}
}

// @Title Content Video
// @Success 200 {object} models.Object
// @Param   info    body	{datatransfers.LikeBody}	true	"body"
// @router /comment [post]
func (c *VideoControllerAuth) CommentVideo(info datatransfers.CommentBody) datatransfers.Response {
	comment, err := c.Handler.CommentVideo(c.userID, info.Hash, info.Content)
	if err != nil {
		fmt.Printf("[VideoControllerAuth::CommentVideo] failed commenting video. %+v\n", err)
		return datatransfers.Response{Error: "Failed commenting video", Code: http.StatusInternalServerError}
	}
	return datatransfers.Response{Data: comment, Code: http.StatusOK}
}
