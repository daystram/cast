package v1

import (
	"log"
	"net/http"

	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"

	"github.com/daystram/cast/cast-be/constants"
	"github.com/daystram/cast/cast-be/handlers"
)

// WebSocket Chat Controller
type WebSocketController struct {
	beego.Controller
	Handler handlers.Handler
}

// @Title Connect
// @Success 200 success
// @router /chat/:hash [get]
func (c *WebSocketController) Connect(hash string, _ string) {
	var err error
	err = c.Handler.ConnectChatWS(c.Ctx, hash)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(c.Ctx.ResponseWriter, "Not a websocket handshake!", http.StatusBadRequest)
		log.Printf("[WebSocketControllerAuth::Connect] handshake error. %+v\n", err)
		return
	} else if err != nil {
		log.Printf("[WebSocketControllerAuth::Connect] failed upgrading to WSS. %+v\n", err)
		return
	}
}

// WebSocket Chat Controller
type WebSocketControllerAuth struct {
	beego.Controller
	Handler handlers.Handler
	userID  string
}

func (c *WebSocketControllerAuth) Prepare() {
	c.userID = c.Ctx.Input.Param(constants.ContextParamUserID)
}

// @Title Connect Notification
// @Success 200 success
// @Param	access_token	query	string	false	"Bearer token"	""
// @router /notification [get]
func (c *WebSocketControllerAuth) ConnectNotification(_ string) {
	var err error
	err = c.Handler.ConnectNotificationWS(c.Ctx, c.userID)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(c.Ctx.ResponseWriter, "Not a websocket handshake!", http.StatusBadRequest)
		log.Printf("[WebSocketControllerAuth::ConnectNotification] handshake error. %+v\n", err)
		return
	} else if err != nil {
		log.Printf("[WebSocketControllerAuth::ConnectNotification] failed upgrading to WSS. %+v\n", err)
		return
	}
}

// @Title Connect Chat
// @Success 200 success
// @Param	access_token	query	string	false	"Bearer token"	""
// @router /chat/:hash [get]
func (c *WebSocketControllerAuth) ConnectChat(hash string, _ string) {
	var err error
	err = c.Handler.ConnectChatWS(c.Ctx, hash, c.userID)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(c.Ctx.ResponseWriter, "Not a websocket handshake!", http.StatusBadRequest)
		log.Printf("[WebSocketControllerAuth::ConnectChat] handshake error. %+v\n", err)
		return
	} else if err != nil {
		log.Printf("[WebSocketControllerAuth::ConnectChat] failed upgrading to WSS. %+v\n", err)
		return
	}
}
