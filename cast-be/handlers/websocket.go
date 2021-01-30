package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/daystram/cast/cast-be/constants"
	"github.com/daystram/cast/cast-be/datatransfers"

	"github.com/astaxie/beego/context"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (m *module) ConnectNotificationWS(ctx *context.Context, userID primitive.ObjectID) (err error) {
	var user datatransfers.User
	user, err = m.db.userOrm.GetOneByID(userID)
	if err != nil {
		fmt.Printf("[ConnectNotificationWS] failed retrieving user info for %s. %+v\n", userID.Hex(), err)
		return
	}
	var ws *websocket.Conn
	ctx.Request.Header.Set("Sec-Websocket-Version", "13")
	ctx.Request.Header.Del("Sec-Websocket-Extensions")
	if ws, err = m.notification.upgrader.Upgrade(ctx.ResponseWriter, ctx.Request, ctx.Request.Header); err != nil {
		return err
	}
	m.notification.sockets[user.ID.Hex()] = ws
	go m.NotificationPingWorker(ws)
	return
}

func (m *module) ConnectChatWS(ctx *context.Context, hash string, userID ...primitive.ObjectID) (err error) {
	var video datatransfers.Video
	if video, err = m.db.videoOrm.GetOneByHash(hash); err != nil {
		fmt.Printf("[ConnectChatWS] unkown video with hash %s. %+v\n", hash, err)
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
			fmt.Printf("[ConnectChatWS] failed retrieving user info for %s. %+v\n", userID[0].Hex(), err)
			return
		}
	}
	go m.ChatReaderWorker(ws, hash, user, video.Type == constants.VideoTypeLive && video.Author.ID != user.ID)
	return
}

func (m *module) NotificationPingWorker(conn *websocket.Conn) {
	for {
		message := datatransfers.WebSocketMessage{}
		if err := conn.ReadJSON(&message); err != nil {
			if !websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				fmt.Printf("[NotificationPingWorker] failed reading message. %+v\n", err)
			}
			_ = conn.WriteJSON(datatransfers.WebSocketMessage{
				Type: constants.MessageTypeError,
				Data: "Failed reading message!",
				Code: http.StatusInternalServerError,
			})
			break
		}
		switch message.Type {
		case constants.MessageTypePing:
			_ = conn.WriteJSON(datatransfers.WebSocketMessage{
				Type: constants.MessageTypePing,
				Data: "pong",
				Code: http.StatusOK,
			})
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

func (m *module) PushNotification(userID primitive.ObjectID, message datatransfers.NotificationOutgoing) {
	if conn, exists := m.notification.sockets[userID.Hex()]; exists {
		_ = conn.WriteJSON(datatransfers.WebSocketMessage{
			Type: constants.MessageTypeNotification,
			Data: message,
			Code: http.StatusOK,
		})
	}
}

func (m *module) BroadcastNotificationSubscriber(authorID primitive.ObjectID, message datatransfers.NotificationOutgoing) {
	var err error
	var subscriptions []datatransfers.Subscription
	if subscriptions, err = m.db.subscriptionOrm.GetSubscriptionsByAuthorID(authorID); err != nil {
		fmt.Printf("[BroadcastNotificationSubscriber] failed to retrieve subscribers\n")
	}
	fmt.Println(len(subscriptions))
	for _, subscription := range subscriptions {
		if subscriptions, err = m.db.subscriptionOrm.GetSubscriptionsByAuthorID(authorID); err != nil {
			fmt.Printf("[BroadcastNotificationSubscriber] failed to retrieve subscribers\n")
		}
		m.PushNotification(subscription.UserID, message)
	}
}
