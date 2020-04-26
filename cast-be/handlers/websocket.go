package handlers

import (
	"fmt"
	"github.com/astaxie/beego/context"
	"github.com/daystram/websocket"
	"gitlab.com/daystram/cast/cast-be/constants"
	"gitlab.com/daystram/cast/cast-be/datatransfers"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (m *module) ConnectWebSocket(ctx *context.Context, hash string, userID ...primitive.ObjectID) (err error) {
	var video datatransfers.Video
	if video, err = m.db.videoOrm.GetOneByHash(hash); err != nil {
		fmt.Printf("[ConnectWebSocket] unkown video with hash %s. %+v\n", hash, err)
		return
	}
	var ws *websocket.Conn
	ctx.Request.Header.Set("Sec-Websocket-Version", "13")
	ctx.Request.Header.Del("Sec-Websocket-Extensions")
	if ws, err = m.chat.upgrader.Upgrade(ctx.ResponseWriter, ctx.Request, ctx.Request.Header); err != nil {
		return err
	}
	m.chat.sockets[hash] = append(m.chat.sockets[hash], ws)
	var user datatransfers.User
	if len(userID) != 0 {
		user, err = m.db.userOrm.GetOneByID(userID[0])
		if err != nil {
			fmt.Printf("[ConnectWebSocket] failed retrieving user info for %s. %+v\n", userID[0].Hex(), err)
			return
		}
	}
	go m.ChatReaderWorker(ws, hash, user, video.Type == constants.VideoTypeLive && video.Author.ID != user.ID)
	return
}

func (m *module) ChatReaderWorker(conn *websocket.Conn, hash string, user datatransfers.User, live bool) {
	if live {
		if err := m.db.videoOrm.IncrementViews(hash); err != nil {
			fmt.Printf("[ChatReaderWorker] failed incrementing views for %s. %+v\n", hash, err)
		}
	}
	for {
		message := datatransfers.WebSocketMessage{}
		if err := conn.ReadJSON(&message); err != nil {
			if !websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				fmt.Printf("[ChatReaderWorker] failed reading message for %s. %+v\n", hash, err)
			}
			_ = conn.WriteJSON(datatransfers.WebSocketMessage{
				Type: constants.MessageTypeError,
				Data: "Failed reading chat!",
				Code: http.StatusInternalServerError,
			})
			if live {
				if err = m.db.videoOrm.IncrementViews(hash, true); err != nil {
					fmt.Printf("[ChatReaderWorker] failed decrementing views for %s. %+v\n", hash, err)
				}
			}
			break
		}
		if user.Username == "" {
			continue
		}
		switch message.Type {
		case constants.MessageTypeChat:
			// TODO: insert into DB?
			chat, ok := message.Data.(string)
			if !ok {
				fmt.Printf("[ChatReaderWorker] failed parsing chat for %s\n", hash)
				_ = conn.WriteJSON(datatransfers.WebSocketMessage{
					Type: constants.MessageTypeError,
					Data: "Failed sending chat!",
					Code: http.StatusInternalServerError,
				})
				break
			}
			currentTime := time.Now()
			for _, conn := range m.chat.sockets[hash] {
				_ = conn.WriteJSON(datatransfers.WebSocketMessage{
					Type: constants.MessageTypeChat,
					Data: datatransfers.ChatOutgoing{
						Author:    user.Username,
						Chat:      chat,
						CreatedAt: currentTime,
					},
					Code: http.StatusOK,
				})
			}
			break
		default:
			fmt.Printf("[ChatReaderWorker] unknown message type: %s\n", message.Type)
			_ = conn.WriteJSON(datatransfers.WebSocketMessage{
				Type: constants.MessageTypeError,
				Data: "Unknown request!",
				Code: http.StatusBadRequest,
			})
			break
		}
	}
}
