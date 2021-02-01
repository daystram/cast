package handlers

import (
	"fmt"
	"github.com/daystram/cast/cast-be/config"
	"github.com/daystram/cast/cast-be/util"
	"golang.org/x/crypto/bcrypt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/daystram/cast/cast-be/constants"
	"github.com/daystram/cast/cast-be/datatransfers"
)


func (m *module) Register(info datatransfers.UserRegister) (err error) {
	if err = m.db.userOrm.CheckUnique("Name", info.Name); err != nil {
		fmt.Printf("[Register] Name %s already exists. %+v\n", info.Name, err)
		return
	}
	if err = m.db.userOrm.CheckUnique("Username", info.Username); err != nil {
		fmt.Printf("[Register] Username %s already exists. %+v\n", info.Username, err)
		return
	}
	if err = m.db.userOrm.CheckUnique("Email", info.Email); err != nil {
		fmt.Printf("[Register] Email %s already exists. %+v\n", info.Email, err)
		return
	}
	hashed, _ := bcrypt.GenerateFromPassword([]byte(info.Password), bcrypt.DefaultCost)
	var userID primitive.ObjectID
	user := datatransfers.User{
		Name:      info.Name,
		Username:  info.Username,
		Email:     info.Email,
		Password:  string(hashed),
		CreatedAt: time.Now(),
	}
	if userID, err = m.db.userOrm.InsertUser(user); err != nil {
		fmt.Printf("[Register] Failed adding %s user entry. %+v\n", info.Username, err)
		return
	}
	user.ID = userID
	if _, err = m.db.videoOrm.InsertVideo(datatransfers.VideoInsert{
		ID:          primitive.NewObjectID(),
		Hash:        info.Username,
		Type:        constants.VideoTypeLive,
		Author:      userID,
		Title:       fmt.Sprintf("%s's Livestream", info.Name),
		Tags:        []string{"live", "first"},
		Description: "Welcome to my stream!",
		Resolutions: -1,
		IsLive:      false,
	}); err != nil {
		_ = m.db.userOrm.DeleteOneByID(userID)
		fmt.Printf("[Register] Failed adding %s live video entry. %+v\n", info.Username, err)
		return
	}
	if err = m.SendVerification(user); err != nil {
		fmt.Printf("[Register] Failed sending %s verification mail. %+v\n", info.Username, err)
	}
	_ = util.Copy(
		fmt.Sprintf("%s/%s/%s.jpg", config.AppConfig.UploadsDirectory, constants.ProfileRootDir, constants.ProfileDefault),
		fmt.Sprintf("%s/%s/%s.jpg", config.AppConfig.UploadsDirectory, constants.ProfileRootDir, user.Username),
	)
	_ = util.Copy(
		fmt.Sprintf("%s/%s/%s.jpg", config.AppConfig.UploadsDirectory, constants.ThumbnailRootDir, constants.ThumbnailDefault),
		fmt.Sprintf("%s/%s/%s.jpg", config.AppConfig.UploadsDirectory, constants.ThumbnailRootDir, user.Username),
	)
	return
}

