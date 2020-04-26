package handlers

import (
	"errors"
	"fmt"
	"gitlab.com/daystram/cast/cast-be/constants"

	data "gitlab.com/daystram/cast/cast-be/datatransfers"
	"gitlab.com/daystram/cast/cast-be/util"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (m *module) UserDetails(userID primitive.ObjectID) (detail data.UserDetail, err error) {
	var user data.User
	var videos []data.Video
	views := 0
	if user, err = m.db.userOrm.GetOneByID(userID); err != nil {
		return data.UserDetail{}, errors.New(fmt.Sprintf("[UserDetails] user not found. %+v\n", err))
	}
	if videos, err = m.db.videoOrm.GetAllVODByAuthor(userID); err != nil {
		return data.UserDetail{}, errors.New(fmt.Sprintf("[UserDetails] cannot retrieve all user videos. %+v\n", err))
	}
	for _, video := range videos {
		views += video.Views
	}
	detail = data.UserDetail{
		Name:        user.Name,
		Username:    user.Username,
		Email:       user.Email,
		Subscribers: user.Subscribers,
		Views:       views,
		Uploads:     len(videos),
	}
	return
}

func (m *module) GetUserByEmail(email string) (user data.User, err error) {
	return m.db.userOrm.GetOneByEmail(email)
}

func (m *module) UpdateUser(user data.UserEditForm, ID primitive.ObjectID) (err error) {
	return m.db.userOrm.EditUser(data.User{
		ID:    ID,
		Name:  user.Name,
		Email: user.Email,
	})
}

func (m *module) NormalizeProfile(username string) (err error) {
	return util.NormalizeImage(constants.ProfileRootDir, username, constants.ProfileWidth, constants.ProfileHeight)
}
