package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/astaxie/beego/context"

	"github.com/daystram/cast/cast-be/config"
	"github.com/daystram/cast/cast-be/constants"
	"github.com/daystram/cast/cast-be/datatransfers"
)

func AuthenticateAccessToken(ctx *context.Context) {
	var accessToken string
	if accessToken = strings.TrimPrefix(ctx.Input.Header("Authorization"), "Bearer "); accessToken == "" {
		if accessToken = ctx.Input.Query("access_token"); accessToken == "" {
			ctx.ResponseWriter.WriteHeader(http.StatusUnauthorized)
			return
		}
	}
	var err error
	var info datatransfers.AccessTokenInfo
	if info, err = verifyAccessToken(accessToken); err != nil || !info.Active {
		log.Printf("[AuthFilter] Invalid access_token. %+v\n", err)
		errMessage, _ := json.Marshal(map[string]interface{}{"message": "invalid access_token"})
		ctx.ResponseWriter.WriteHeader(http.StatusUnauthorized)
		_, _ = ctx.ResponseWriter.Write(errMessage)
		return
	}
 	ctx.Input.SetParam(constants.ContextParamUserID, info.Subject)
}

func verifyAccessToken(accessToken string) (info datatransfers.AccessTokenInfo, err error) {
	var response *http.Response
	if response, err = http.Post(fmt.Sprintf("%s/oauth/introspect", config.AppConfig.RatifyIssuer),
		"application/x-www-form-urlencoded",
		bytes.NewBuffer([]byte(fmt.Sprintf(
			"token=%s&client_id=%s&client_secret=%s&token_type_hint=access_token",
			accessToken, config.AppConfig.RatifyClientID, config.AppConfig.RatifyClientSecret,
		))),
	); err != nil {
		return
	}
	var body []byte
	if body, err = ioutil.ReadAll(response.Body); err != nil {
		return
	}
	err = json.Unmarshal(body, &info)
	return
}
