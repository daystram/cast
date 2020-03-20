package handlers

import (
	"fmt"
	errors2 "gitlab.com/daystram/cast/cast-be/errors"
	"time"

	"gitlab.com/daystram/cast/cast-be/config"
	"gitlab.com/daystram/cast/cast-be/constants"
	"gitlab.com/daystram/cast/cast-be/datatransfers"

	"github.com/astaxie/beego/orm"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func (m *module) CheckUniqueUserField(field, value string) (err error) {
	return m.db().userOrm.CheckUnique(field, value)
}

func (m *module) Register(info datatransfers.UserRegister) (err error) {
	if err = m.db().userOrm.CheckUnique("Name", info.Name); err != nil {
		fmt.Printf("[Register] Name %s already exists. %+v\n", info.Name, err)
		return
	}
	if err = m.db().userOrm.CheckUnique("Username", info.Username); err != nil {
		fmt.Printf("[Register] Username %s already exists. %+v\n", info.Username, err)
		return
	}
	if err = m.db().userOrm.CheckUnique("Email", info.Email); err != nil {
		fmt.Printf("[Register] Email %s already exists. %+v\n", info.Email, err)
		return
	}
	hashed, _ := bcrypt.GenerateFromPassword([]byte(info.Password), bcrypt.DefaultCost)
	var userID primitive.ObjectID
	if userID, err = m.db().userOrm.InsertUser(datatransfers.User{
		Name:      info.Name,
		Username:  info.Username,
		Email:     info.Email,
		Password:  string(hashed),
		CreatedAt: time.Now(),
	}); err != nil {
		fmt.Printf("[Register] Failed adding %s user entry. %+v\n", info.Username, err)
		return
	}
	if _, err = m.db().videoOrm.InsertVideo(datatransfers.VideoInsert{
		ID:          primitive.NewObjectID(),
		Hash:        info.Username,
		Type:        constants.VideoTypeLive,
		Author:      userID,
		Title:       fmt.Sprintf("%s's Livestream", info.Name),
		Description: "",
		IsLive:      false,
		CreatedAt:   time.Now(),
	}); err != nil {
		_ = m.db().userOrm.DeleteOneByID(userID)
		fmt.Printf("[Register] Failed adding %s live video entry. %+v\n", info.Username, err)
		return
	}

	return
}

func (m *module) Authenticate(info datatransfers.UserLogin) (token string, err error) {
	var user datatransfers.User
	if user, err = m.validate(info); err != nil {
		return
	}
	if token, err = m.generateToken(user); err != nil {
		return
	}
	return
}

func (m *module) validate(info datatransfers.UserLogin) (user datatransfers.User, err error) {
	db := m.db()
	if user, err = db.userOrm.GetOneByUsername(info.Username); err != nil {
		if err == orm.ErrNoRows {
			return datatransfers.User{}, errors2.ErrNotRegistered
		}
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(info.Password)); err != nil {
		return datatransfers.User{}, errors2.ErrIncorrectPassword
	}
	return
}

func (m *module) generateToken(user datatransfers.User) (tokenString string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":     user.ID,
		"expiry": time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenString, _ = token.SignedString([]byte(config.AppConfig.JWTSecret))
	return
}
