package v1

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/astaxie/beego"

	"github.com/daystram/cast/cast-be/config"
	"github.com/daystram/cast/cast-be/constants"
	"github.com/daystram/cast/cast-be/datatransfers"
	"github.com/daystram/cast/cast-be/errors"
	"github.com/daystram/cast/cast-be/handlers"
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

// @Title Verify Email
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

// @Title Resend Verification
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
	user, token, err := c.Handler.Authenticate(info)
	switch err {
	case nil:
		timeout := int(constants.AuthenticationTimeout.Seconds())
		if info.Remember {
			timeout = int(constants.AuthenticationTimeoutExtended.Seconds())
		}
		c.Ctx.SetCookie(constants.AuthenticationCookieKey, fmt.Sprintf("%s|Bearer %s", user.Username, token), timeout, "/", config.AppConfig.Domain, !config.AppConfig.Debug)
		return datatransfers.Response{Code: http.StatusOK}
	case errors.ErrNotRegistered:
		log.Printf("[AuthController::PostAuthenticate] failed authenticating %s. %+v\n", info.Username, err)
		return datatransfers.Response{Error: "Username or email not registered", Code: http.StatusNotFound}
	case errors.ErrIncorrectPassword:
		log.Printf("[AuthController::PostAuthenticate] failed authenticating %s. %+v\n", info.Username, err)
		return datatransfers.Response{Error: "Incorrect password", Code: http.StatusForbidden}
	case errors.ErrNotVerified:
		log.Printf("[AuthController::PostAuthenticate] failed authenticating %s. %+v\n", info.Username, err)
		return datatransfers.Response{Error: "User not verified", Code: http.StatusNotAcceptable}
	default:
		log.Printf("[AuthController::PostAuthenticate] failed authenticating %s. %+v\n", info.Username, err)
		return datatransfers.Response{Error: "An error has occurred", Code: http.StatusNotFound}
	}
}

// @Title Forget Password
// @Param   info    body	{datatransfers.UserResend}	true	"email"
// @Success 200 success
// @router /forget [post]
func (c *AuthController) PostResetPassword(info datatransfers.UserResend) datatransfers.Response {
	user, err := c.Handler.GetUserByEmail(info.Email)
	if err != nil {
		log.Printf("[AuthController::PostResetPassword] cannot find user with email %s. %+v\n", info.Email, err)
		return datatransfers.Response{Error: "Email not registered", Code: http.StatusNotFound}
	}
	if !user.Verified {
		log.Printf("[AuthController::PostResetPassword] user not verified\n")
		return datatransfers.Response{Error: "Email not verified", Code: http.StatusForbidden}
	}
	err = c.Handler.SendResetToken(user)
	if err != nil {
		log.Printf("[AuthController::PostResetPassword] cannot send password reset email. %+v\n", err)
		return datatransfers.Response{Error: "Failed sending email", Code: http.StatusInternalServerError}
	}
	return datatransfers.Response{Code: http.StatusOK}
}

// @Title Validate Reset Token
// @Param   info    body	{datatransfers.UserVerify}	true	"verification token"
// @Success 200 success
// @router /validate_reset [post]
func (c *AuthController) PostValidateReset(info datatransfers.UserVerify) datatransfers.Response {
	err := c.Handler.CheckResetToken(info.Key)
	if err != nil {
		log.Printf("[AuthController::PostValidateReset] cannot validate reset token. %+v\n", err)
		return datatransfers.Response{Error: "Reset token invalid", Code: http.StatusUnauthorized}
	}
	return datatransfers.Response{Code: http.StatusOK}
}

// @Title Update Password
// @Param   info    body	{datatransfers.UserUpdatePassword}	true	"updated password"
// @Success 200 success
// @router /update [put]
func (c *AuthController) PostUpdatePassword(info datatransfers.UserUpdatePassword) datatransfers.Response {
	err := c.Handler.UpdatePassword(info)
	if err != nil {
		log.Printf("[AuthController::PostUpdatePassword] cannot find update password. %+v\n", err)
		return datatransfers.Response{Error: "Password update failed", Code: http.StatusInternalServerError}
	}
	return datatransfers.Response{Code: http.StatusOK}
}

// @Title Logout
// @Success 200 success
// @router /logout [post]
func (c *AuthController) PostDeAuthenticate() {
	c.Ctx.SetCookie(constants.AuthenticationCookieKey, "", -1)
}
