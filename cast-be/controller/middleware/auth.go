package middleware

import (
	"encoding/json"
	"fmt"
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
	JwtClaims   JwtClaims
}

type JwtClaims struct {
	ID     string
	Expiry int64
}

func NewJwtAuthorization(secret string, bearerTokenStr string) JwtAuthorization {
	jwtTokenStr := parseBearerToken(bearerTokenStr)
	return JwtAuthorization{jwtTokenStr, []byte(secret), JwtClaims{}}
}

func (j *JwtAuthorization) ExtractClaimsFromToken() (id string, expiry int64, err error) {
	claims := jwt.MapClaims{}
	var token *jwt.Token
	if token, err = jwt.ParseWithClaims(j.jwtTokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return j.secret, nil
	}); err != nil || !token.Valid {
		return
	}
	err = j.parseJwtClaims(claims)
	if err != nil {
		return
	}

	return j.JwtClaims.ID, j.JwtClaims.Expiry, nil
}

func parseBearerToken(bearerTokenStr string) string {
	splitToken := strings.Split(bearerTokenStr, "Bearer ")
	if len(splitToken) < 2 {
		return ""
	}
	return splitToken[1]
}

func (j *JwtAuthorization) parseJwtClaims(claims jwt.MapClaims) (err error) {
	id, ok := claims["id"]
	if !ok {
		return fmt.Errorf("key 'id' is not contained within JWT's claims")
	}
	j.JwtClaims.ID, ok = id.(string)
	if !ok {
		return fmt.Errorf("key 'id' is in wrong format")
	}

	expiryClaim, ok := claims["expiry"]
	if !ok {
		return fmt.Errorf("key 'expiry' is not contained within JWT's claims")
	}
	expiry, ok := expiryClaim.(float64)
	j.JwtClaims.Expiry = int64(expiry)
	if !ok {
		return fmt.Errorf("key 'expiry' is in wrong format")
	}
	return
}

func AuthenticateJWT(ctx *context.Context) {
	bearerTokenStr := ctx.Input.Query("access_token")
	if bearerTokenStr == "" {
		bearerTokenStr = ctx.Input.Header("Authorization")
	}
	jwtAuthorization := NewJwtAuthorization(config.AppConfig.JWTSecret, bearerTokenStr)
	id, expiry, err := jwtAuthorization.ExtractClaimsFromToken()
	if expiry < time.Now().Unix() {
		err = fmt.Errorf("JWT invalid")
	}
	if err != nil {
		log.Printf("[AuthFilter] failed to get jwt token. %+v\n", err)
		errMessage, _ := json.Marshal(map[string]interface{}{"message": "JWT is invalid"})
		ctx.ResponseWriter.WriteHeader(http.StatusForbidden)
		ctx.ResponseWriter.Header().Set("Authorization", "NONE")
		_, _ = ctx.ResponseWriter.Write(errMessage)
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":     id,
		"expiry": time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenString, err := token.SignedString([]byte(config.AppConfig.JWTSecret))
	if err != nil {
		return
	}
	ctx.ResponseWriter.Header().Add("Authorization", fmt.Sprintf("Bearer %v", tokenString))
	ctx.Input.SetParam(constants.ContextParamUserID, id)
}
