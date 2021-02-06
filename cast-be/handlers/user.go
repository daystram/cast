package handlers

import (
	"errors"
	"fmt"

	data "github.com/daystram/cast/cast-be/datatransfers"
)

func (m *module) UserDetails(userID string) (detail data.UserDetail, err error) {
	var user data.User
	var videos []data.Video
	var subscriberCount int
	views := 0
	if user, err = m.db.userOrm.GetOneByID(userID); err != nil {
		return data.UserDetail{}, errors.New(fmt.Sprintf("[UserDetails] user not found. %+v\n", err))
	}
	if videos, err = m.db.videoOrm.GetAllVODByAuthor(userID, true); err != nil {
		return data.UserDetail{}, errors.New(fmt.Sprintf("[UserDetails] cannot retrieve all user videos. %+v\n", err))
	}
	if subscriberCount, err = m.db.subscriptionOrm.GetCountByAuthorID(user.ID); err != nil {
		return data.UserDetail{}, errors.New(fmt.Sprintf("[UserDetails] cannot retrieve user subscriber count. %+v\n", err))
	}
	for _, video := range videos {
		views += video.Views
	}
	detail = data.UserDetail{
		ID:          user.ID,
		Username:    user.Username,
		Name:        user.Name,
		Subscribers: subscriberCount,
		Views:       views,
		Uploads:     len(videos),
	}
	return
}

func (m *module) UserGetOneByID(userID string) (user data.User, err error) {
	if user, err = m.db.userOrm.GetOneByID(userID); err != nil {
		return data.User{}, fmt.Errorf("[UserGetOneByID] user not found. %+v\n", err)
	}
	return
}
