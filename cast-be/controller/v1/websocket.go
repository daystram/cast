package v1

import (
	"log"
	"net/http"

	"gitlab.com/daystram/cast/cast-be/constants"
	"gitlab.com/daystram/cast/cast-be/handlers"

	"github.com/astaxie/beego"
	"github.com/daystram/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// WebSocket Chat Controller
type WebSocketController struct {
	beego.Controller
	Handler handlers.Handler
}

// @Title Connect
// @Success 200 success
// @router /:hash [get]
func (c *WebSocketController) Connect(hash string, _ string) {
	var err error
	err = c.Handler.ConnectWebSocket(c.Ctx, hash)
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
	userID  primitive.ObjectID
}

func (c *WebSocketControllerAuth) Prepare() {
	c.userID, _ = primitive.ObjectIDFromHex(c.Ctx.Input.Param(constants.ContextParamUserID))
}

// @Title Connect
// @Success 200 success
// @Param	access_token	query	string	false	"Bearer token"	""
// @router /:hash [get]
func (c *WebSocketControllerAuth) Connect(hash string, _ string) {
	var err error
	err = c.Handler.ConnectWebSocket(c.Ctx, hash, c.userID)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(c.Ctx.ResponseWriter, "Not a websocket handshake!", http.StatusBadRequest)
		log.Printf("[WebSocketControllerAuth::Connect] handshake error. %+v\n", err)
		return
	} else if err != nil {
		log.Printf("[WebSocketControllerAuth::Connect] failed upgrading to WSS. %+v\n", err)
		return
	}
}
