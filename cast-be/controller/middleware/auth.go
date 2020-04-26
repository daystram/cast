package middleware

import (
	"encoding/json"
	"fmt"
	"gitlab.com/daystram/cast/cast-be/datatransfers"
	"log"
	"net/http"
	"strings"
	"time"

	"gitlab.com/daystram/cast/cast-be/config"
	"gitlab.com/daystram/cast/cast-be/constants"

	"github.com/astaxie/beego/context"
	"github.com/dgrijalva/jwt-go"
)

type JwtAuthorization struct {
	jwtTokenStr string
	secret      []byte
	JWTClaims   datatransfers.JWTClaims
}

func NewJWTAuthorization(secret string, bearerTokenStr string) JwtAuthorization {
	jwtTokenStr := parseBearerToken(bearerTokenStr)
	return JwtAuthorization{jwtTokenStr, []byte(secret), datatransfers.JWTClaims{}}
}

func (j *JwtAuthorization) ExtractClaimsFromToken() (id string, expiry int64, remember bool, err error) {
	claims := jwt.MapClaims{}
	var token *jwt.Token
	if token, err = jwt.ParseWithClaims(j.jwtTokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return j.secret, nil
	}); err != nil || !token.Valid {
		return
	}
	err = j.parseJWTClaims(claims)
	if err != nil {
		return
	}
	return j.JWTClaims.ID, j.JWTClaims.Expiry, j.JWTClaims.Remember, nil
}

func parseBearerToken(bearerTokenStr string) string {
	splitToken := strings.Split(bearerTokenStr, "Bearer ")
	if len(splitToken) < 2 {
		return ""
	}
	return splitToken[1]
}

func (j *JwtAuthorization) parseJWTClaims(claims jwt.MapClaims) (err error) {
	id, ok := claims["id"]
	if !ok {
		return fmt.Errorf("key 'id' is not contained within JWT's claims")
	}
	j.JWTClaims.ID, ok = id.(string)
	if !ok {
		return fmt.Errorf("key 'id' is in wrong format")
	}

	expiryClaim, ok := claims["expiry"]
	if !ok {
		return fmt.Errorf("key 'expiry' is not contained within JWT's claims")
	}
	expiry, ok := expiryClaim.(float64)
	j.JWTClaims.Expiry = int64(expiry)
	if !ok {
		return fmt.Errorf("key 'expiry' is in wrong format")
	}

	rememberClaim, ok := claims["remember"]
	if !ok {
		return fmt.Errorf("key 'remember' is not contained within JWT's claims")
	}
	j.JWTClaims.Remember, ok = rememberClaim.(bool)
	if !ok {
		return fmt.Errorf("key 'remember' is in wrong format")
	}
	return
}

func AuthenticateJWT(ctx *context.Context) {
	bearerTokenStr := ctx.Input.Query("access_token")
	cookie := ctx.GetCookie(constants.AuthenticationCookieKey)
	if bearerTokenStr == "" {
		if cookie == "" {
			log.Println("[AuthFilter] invalid auth cookie", )
			ctx.SetCookie(constants.AuthenticationCookieKey, "", -1)
			ctx.ResponseWriter.WriteHeader(http.StatusForbidden)
			return
		}
		if split := strings.Split(cookie, "|"); len(split) == 2 {
			bearerTokenStr = split[1]
		}
	}
	jwtAuthorization := NewJWTAuthorization(config.AppConfig.JWTSecret, bearerTokenStr)
	id, expiry, remember, err := jwtAuthorization.ExtractClaimsFromToken()
	if expiry < time.Now().Unix() {
		err = fmt.Errorf("JWT invalid")
	}
	if err != nil {
		log.Printf("[AuthFilter] failed to get JWT token. %+v\n", err)
		errMessage, _ := json.Marshal(map[string]interface{}{"message": "JWT is invalid"})
		ctx.ResponseWriter.WriteHeader(http.StatusForbidden)
		_, _ = ctx.ResponseWriter.Write(errMessage)
		return
	}
	timeout := constants.AuthenticationTimeout
	if remember {
		timeout = constants.AuthenticationTimeoutExtended
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       id,
		"expiry":   time.Now().Add(timeout).Unix(),
		"remember": remember,
	})
	tokenString, err := token.SignedString([]byte(config.AppConfig.JWTSecret))
	if err != nil {
		return
	}
	ctx.SetCookie(constants.AuthenticationCookieKey, fmt.Sprintf("%s|Bearer %s", strings.Split(cookie, "|")[0], tokenString), int(timeout.Seconds()), "/", config.AppConfig.Domain, !config.AppConfig.Debug)
	ctx.Input.SetParam(constants.ContextParamUserID, id)
}
