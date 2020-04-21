package v1

import (
	"fmt"
	"gitlab.com/daystram/cast/cast-be/config"
	"gitlab.com/daystram/cast/cast-be/constants"
	"gitlab.com/daystram/cast/cast-be/datatransfers"
	"gitlab.com/daystram/cast/cast-be/handlers"
	"net/http"

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

// @Title Update User Details
// @Success 200 {object} models.Object
// @Param   stub		query	string	false	"stub"
// @router /edit [put]
func (c *UserControllerAuth) UpdateProfile(_ string) datatransfers.Response {
	info := datatransfers.UserEditForm{}
	err := c.ParseForm(&info)
	if err != nil {
		fmt.Printf("[UserControllerAuth::UpdateProfile] failed parsing user details. %+v\n", err)
		return datatransfers.Response{Error: "Failed parsing user detail", Code: http.StatusInternalServerError}
	}
	err = c.Handler.UpdateUser(info, c.userID)
	if err != nil {
		fmt.Printf("[UserControllerAuth::UpdateProfile] user cannot be found. %+v\n", err)
		return datatransfers.Response{Error: "failed setting stream window", Code: http.StatusInternalServerError}
	}
	if _, _, err = c.GetFile("profile"); err != nil {
		if err == http.ErrMissingFile {
			return datatransfers.Response{Code: http.StatusOK}
		} else {
			fmt.Printf("[UserControllerAuth::UpdateProfile] failed retrieving profile image. %+v\n", err)
			return datatransfers.Response{Error: "Failed retrieving profile image", Code: http.StatusInternalServerError}
		}
	}
	var user datatransfers.UserDetail
	if user, err = c.Handler.UserDetails(c.userID); err != nil {
		fmt.Printf("[UserControllerAuth::UpdateProfile] failed retrieving user info. %+v\n", err)
		return datatransfers.Response{Error: "Failed retrieving user info", Code: http.StatusInternalServerError}
	}
	err = c.SaveToFile("profile", fmt.Sprintf("%s/profile/%s.ori", config.AppConfig.UploadsDirectory, user.Username))
	if err != nil {
		fmt.Printf("[UserControllerAuth::UpdateProfile] failed saving profile image. %+v\n", err)
		return datatransfers.Response{Error: "Failed saving profile image", Code: http.StatusInternalServerError}
	}
	err = c.Handler.NormalizeProfile(user.Username)
	if err != nil {
		fmt.Printf("[UserControllerAuth::UpdateProfile] failed normalizing profile image. %+v\n", err)
		return datatransfers.Response{Error: "Failed normalizing profile image", Code: http.StatusInternalServerError}
	}
	return datatransfers.Response{Code: http.StatusOK}
}
