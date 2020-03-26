package v1

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"gitlab.com/daystram/cast/cast-be/datatransfers"
	"gitlab.com/daystram/cast/cast-be/errors"
	"gitlab.com/daystram/cast/cast-be/handlers"

	"github.com/astaxie/beego"
)

type AuthController struct {
	beego.Controller
	Handler handlers.Handler
}

// @Title Register
// @Param   info    body	{datatransfers.UserRegister}	true	"registration info"
// @Success 200 success
// @router /signup [post]
func (c *AuthController) PostRegister(info datatransfers.UserRegister) datatransfers.Response {
	err := c.Handler.Register(info)
	if err != nil {
		log.Printf("[AuthController::PostRegister] failed registering %s. %+v\n", info.Username, err)
		return datatransfers.Response{Error: "failed registering", Code: http.StatusInternalServerError}
	}
	return datatransfers.Response{Code: http.StatusOK}
}

// @Title Check Field
// @Param   info    body	{datatransfers.UserFieldCheck}	true	"user field info"
// @Success 200 success
// @router /check [post]
func (c *AuthController) PostCheckUnique(info datatransfers.UserFieldCheck) datatransfers.Response {
	err := c.Handler.CheckUniqueUserField(info.Field, info.Value)
	if err != nil {
		log.Printf("[AuthController::PostCheckUnique] user %s field already used. %+v\n", info.Field, err)
		return datatransfers.Response{Error: fmt.Sprintf("%s already used", strings.Title(info.Field)), Code: http.StatusConflict}
	}
	return datatransfers.Response{Code: http.StatusOK}
}

// @Title Verify
// @Param   info    body	{datatransfers.UserVerify}	true	"verification token"
// @Success 200 success
// @router /verify [post]
func (c *AuthController) PostVerify(info datatransfers.UserVerify) datatransfers.Response {
	err := c.Handler.Verify(info.Key)
	if err != nil {
		log.Printf("[AuthController::PostVerify] cannot verify user. %+v\n", err)
		return datatransfers.Response{Error: "Verification key invalid", Code: http.StatusUnauthorized}
	}
	return datatransfers.Response{Code: http.StatusOK}
}

// @Title Resend
// @Param   info    body	{datatransfers.UserResend}	true	"email"
// @Success 200 success
// @router /resend [post]
func (c *AuthController) PostResend(info datatransfers.UserResend) datatransfers.Response {
	user, err := c.Handler.GetUserByEmail(info.Email)
	if err != nil {
		log.Printf("[AuthController::PostResend] cannot find user with email %s. %+v\n", info.Email, err)
		return datatransfers.Response{Error: "Email not registered", Code: http.StatusNotFound}
	}
	if user.Verified {
		log.Printf("[AuthController::PostResend] user already verified\n")
		return datatransfers.Response{Error: "Already verified", Code: http.StatusConflict}
	}
	err = c.Handler.SendVerification(user)
	if err != nil {
		log.Printf("[AuthController::PostResend] cannot re-send verification email. %+v\n", err)
		return datatransfers.Response{Error: "Failed sending email", Code: http.StatusInternalServerError}
	}
	return datatransfers.Response{Code: http.StatusOK}
}

// @Title Login
// @Param   info    body	{datatransfers.UserLogin}	true	"login info"
// @Success 200 success
// @router /login [post]
func (c *AuthController) PostAuthenticate(info datatransfers.UserLogin) datatransfers.Response {
	token, err := c.Handler.Authenticate(info)
	switch err {
	case nil:
		c.Ctx.Output.Header("Authorization", fmt.Sprintf("Bearer %s", token))
		return datatransfers.Response{Data: fmt.Sprintf("Bearer %s", token), Code: http.StatusOK}
	case errors.ErrNotRegistered:
		log.Printf("[AuthController::PostAuthenticate] failed authenticating %s. %+v\n", info.Username, err)
		return datatransfers.Response{Error: "Username not registered", Code: http.StatusNotFound}
	case errors.ErrIncorrectPassword:
		log.Printf("[AuthController::PostAuthenticate] failed authenticating %s. %+v\n", info.Username, err)
		return datatransfers.Response{Error: "Incorrect password", Code: http.StatusForbidden}
	case errors.ErrNotVerified:
		log.Printf("[AuthController::PostAuthenticate] failed authenticating %s. %+v\n", info.Username, err)
		return datatransfers.Response{Error: "User not verified", Code: http.StatusNotAcceptable}
	default:
		log.Printf("[AuthController::PostAuthenticate] failed authenticating %s. %+v\n", info.Username, err)
		return datatransfers.Response{Error: "Username not registered", Code: http.StatusNotFound}
	}
}

// @Title Logout
// @Success 200 success
// @router /logout [post]
func (c *AuthController) PostDeAuthenticate() {
	// TODO: stop stream
	c.Ctx.Output.Header("Authorization", "")
}
