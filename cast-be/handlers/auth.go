package handlers

import (
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"math/rand"
	"time"

	"gitlab.com/daystram/cast/cast-be/config"
	"gitlab.com/daystram/cast/cast-be/constants"
	"gitlab.com/daystram/cast/cast-be/datatransfers"
	errors2 "gitlab.com/daystram/cast/cast-be/errors"
	"gitlab.com/daystram/cast/cast-be/util"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func (m *module) CheckUniqueUserField(field, value string) (err error) {
	return m.db.userOrm.CheckUnique(field, value)
}

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

func (m *module) SendVerification(user datatransfers.User) (err error) {
	token := datatransfers.Token{
		Invoker:   user.ID,
		Purpose:   constants.TokenPurposeVerification,
		Hash:      m.generateHash(constants.TokenHashDefaultLength),
		CreatedAt: time.Now(),
	}
	_ = m.db.tokenOrm.DeleteOneByInvokerPurpose(token.Invoker, token.Purpose)
	if _, err = m.db.tokenOrm.InsertToken(token); err != nil {
		fmt.Printf("[SendVerification] Failed inserting token. %+v\n", err)
		return
	}
	m.SendSingleEmail("Email Verification", user.Email, constants.EmailTemplateVerification, map[string]string{
		"name": user.Name,
		"link": fmt.Sprintf(constants.EmailLinkVerification, config.AppConfig.Hostname, token.Hash),
	})
	return
}

func (m *module) Verify(key string) (err error) {
	var token datatransfers.Token
	if token, err = m.db.tokenOrm.GetOneByHash(key); err != nil {
		fmt.Printf("[Verify] Failed retrieving token details. %+v\n", err)
		return
	}
	if err = m.db.userOrm.SetVerified(token.Invoker); err != nil {
		fmt.Printf("[Verify] Failed verifying user. %+v\n", err)
		return
	}
	if err = m.db.tokenOrm.DeleteOneByHash(key); err != nil {
		fmt.Printf("[Verify] Failed removing token. %+v\n", err)
		return
	}
	return
}

func (m *module) SendResetToken(user datatransfers.User) (err error) {
	token := datatransfers.Token{
		Invoker:   user.ID,
		Purpose:   constants.TokenPurposeVerification,
		Hash:      m.generateHash(constants.TokenHashDefaultLength),
		CreatedAt: time.Now(),
	}
	_ = m.db.tokenOrm.DeleteOneByInvokerPurpose(token.Invoker, token.Purpose)
	if _, err = m.db.tokenOrm.InsertToken(token); err != nil {
		fmt.Printf("[SendResetToken] Failed inserting token. %+v\n", err)
		return
	}
	m.SendSingleEmail("Password Reset", user.Email, constants.EmailTemplateReset, map[string]string{
		"name": user.Name,
		"link": fmt.Sprintf(constants.EmailLinkReset, config.AppConfig.Hostname, token.Hash),
	})
	return
}

func (m *module) CheckResetToken(key string) (err error) {
	if _, err = m.db.tokenOrm.GetOneByHash(key); err != nil {
		fmt.Printf("[CheckResetToken] Failed retrieving token details. %+v\n", err)
		return
	}
	return
}

func (m *module) UpdatePassword(info datatransfers.UserUpdatePassword) (err error) {
	var user datatransfers.User
	var token datatransfers.Token
	if token, err = m.db.tokenOrm.GetOneByHash(info.Key); err != nil {
		fmt.Printf("[UpdatePassword] Failed retrieving token details. %+v\n", err)
		return
	}
	if user, err = m.db.userOrm.GetOneByID(token.Invoker); err != nil {
		fmt.Printf("[UpdatePassword] Failed retrieving user from ID. %+v\n", err)
		return
	}
	hashed, _ := bcrypt.GenerateFromPassword([]byte(info.Password), bcrypt.DefaultCost)
	user.Password = string(hashed)
	if err = m.db.userOrm.EditUser(user); err != nil {
		fmt.Printf("[UpdatePassword] Failed updating user password. %+v\n", err)
		return
	}
	if err = m.db.tokenOrm.DeleteOneByHash(info.Key); err != nil {
		fmt.Printf("[Verify] Failed removing token. %+v\n", err)
		return
	}
	return
}

func (m *module) Authenticate(info datatransfers.UserLogin) (user datatransfers.User, token string, err error) {
	if user, err = m.validate(info); err != nil {
		return
	}
	if token, err = m.generateJWT(user, info.Remember); err != nil {
		return
	}
	return
}

func (m *module) validate(info datatransfers.UserLogin) (user datatransfers.User, err error) {
	if user, err = m.db.userOrm.GetOneByUsername(info.Username); err != nil {
		if err == mongo.ErrNoDocuments {
			if user, err = m.db.userOrm.GetOneByEmail(info.Username); err != nil {
				if err == mongo.ErrNoDocuments {
					return datatransfers.User{}, errors2.ErrNotRegistered
				} else {
					return
				}
			}
		} else {
			return
		}
	}
	if !user.Verified {
		return datatransfers.User{}, errors2.ErrNotVerified
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(info.Password)); err != nil {
		return datatransfers.User{}, errors2.ErrIncorrectPassword
	}
	return
}

func (m *module) generateJWT(user datatransfers.User, extended bool) (tokenString string, err error) {
	timeout := time.Now().Add(constants.AuthenticationTimeout)
	if extended {
		timeout = time.Now().Add(constants.AuthenticationTimeoutExtended)
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.ID,
		"expiry":   timeout.Unix(),
		"remember": extended,
	})
	tokenString, _ = token.SignedString([]byte(config.AppConfig.JWTSecret))
	return
}

func (m *module) generateHash(length int) string {
	hash := make([]byte, length)
	for i := range hash {
		hash[i] = constants.TokenHashCharacters[rand.Intn(len(constants.TokenHashCharacters))]
	}
	return string(hash)
}
