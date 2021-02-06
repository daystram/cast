package handlers

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/daystram/cast/cast-be/config"
	"github.com/daystram/cast/cast-be/constants"
	"github.com/daystram/cast/cast-be/datatransfers"
)

func (m *module) Register(idToken datatransfers.UserRegister) (err error) {
	var user datatransfers.User
	if user, err = parseIDToken(idToken.IDToken); err != nil {
		fmt.Printf("[Register] Failed parsing id_token. %+v\n", err)
		return
	}
	exists := false
	if err = m.db.userOrm.CheckUnique("_id", user.ID); err != nil {
		exists = true
	}
	if err = m.db.userOrm.CheckUnique("username", user.Username); err != nil { // double check
		exists = true
	}
	if exists {
		if err = m.db.userOrm.UpdateUser(user); err != nil {
			fmt.Printf("[Register] Failed updating user entry. %+v\n", err)
			return
		}
	} else {
		if err = m.db.userOrm.InsertUser(user); err != nil {
			fmt.Printf("[Register] Failed adding %s user entry. %+v\n", user.Username, err)
			return
		}
		if _, err = m.db.videoOrm.InsertVideo(datatransfers.VideoInsert{
			ID:          primitive.NewObjectID(),
			Hash:        user.Username,
			Type:        constants.VideoTypeLive,
			Author:      user.ID,
			Title:       fmt.Sprintf("%s's Livestream", user.Username),
			Tags:        []string{"live", "first"},
			Description: "Welcome to my stream!",
			Resolutions: -1,
			IsLive:      false,
		}); err != nil {
			_ = m.db.userOrm.DeleteOneByID(user.ID)
			fmt.Printf("[Register] Failed adding %s live video entry. %+v\n", user.Username, err)
			return
		}
		_, _ = m.s3.CopyObject(&s3.CopyObjectInput{
			Bucket:     aws.String(config.AppConfig.S3Bucket),
			CopySource: aws.String(fmt.Sprintf("%s/%s.jpg", config.AppConfig.S3Bucket, constants.ThumbnailDefault)),
			Key:        aws.String(fmt.Sprintf("%s/%s.jpg", constants.ThumbnailRootDir, user.Username)),
		})
	}
	return
}

func parseIDToken(idToken string) (user datatransfers.User, err error) {
	claims := jwt.MapClaims{}
	if _, _, err = new(jwt.Parser).ParseUnverified(idToken, claims); err != nil {
		return
	}
	user = datatransfers.User{
		ID:        claims["sub"].(string),
		Username:  claims["preferred_username"].(string),
		Name:      fmt.Sprintf("%s %s", claims["given_name"], claims["family_name"]),
		CreatedAt: time.Now(),
	}
	return
}
