package handlers

import (
	"errors"
	"fmt"
	data "gitlab.com/daystram/cast/cast-be/datatransfers"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (m *module) UserDetails(userID primitive.ObjectID) (detail data.UserDetail, err error) {
	var user data.User
	var videos []data.Video
	views := 0
	if user, err = m.db().userOrm.GetOneByID(userID); err != nil {
		return data.UserDetail{}, errors.New(fmt.Sprintf("[UserDetails] user not found. %+v\n", err))
	}
	if videos, err = m.db().videoOrm.GetAllVODByAuthor(userID); err != nil {
		return data.UserDetail{}, errors.New(fmt.Sprintf("[UserDetails] cannot retrieve all user videos. %+v\n", err))
	}
	for _, video := range videos {
		views += video.Views
	}
	detail = data.UserDetail{
		Name:        user.Name,
		Email:       user.Email,
		Password:    user.Password,
		Subscribers: user.Subscribers,
		Views:       views,
		Uploads:     len(videos),
	}
	return
}

func (m *module) UpdateUser(user data.UserEdit, ID primitive.ObjectID) (err error) {
	return m.db().userOrm.EditUser(data.User{
		ID:    ID,
		Name:  user.Name,
		Email: user.Email,
	})
}
