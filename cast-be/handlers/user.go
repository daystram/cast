package handlers

import (
	data "gitlab.com/daystram/cast/cast-be/datatransfers"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (m *module) UserDetails(userID primitive.ObjectID) (user data.User, err error) {
	return m.db().userOrm.GetOneByID(userID)
}

func (m *module) UpdateUser(user data.UserEdit, ID primitive.ObjectID) (err error) {
	return m.db().userOrm.EditUser(data.User{
		ID:    ID,
		Name:  user.Name,
		Email: user.Email,
	})
}
