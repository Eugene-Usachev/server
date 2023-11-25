package websocket

import (
	"context"
	fb "github.com/Eugene-Usachev/fastbytes"
	loggerLib "github.com/Eugene-Usachev/logger"
	"github.com/gofiber/contrib/websocket"
	"sync"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// CreateResponse pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 1024
)

const (
	onlineUsers = "ou"
)

var Config *websocket.Config = nil

func InitConfig(logger *loggerLib.FastLogger) {
	Config = &websocket.Config{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		WriteBufferPool: &sync.Pool{},
		RecoverHandler: func(conn *websocket.Conn) {
			if reason := recover(); reason != nil {
				logger.Error("recover websocket panic, err: ", reason)
			}
		},
	}
}

// TODO use pool
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	//if userId is -1 then user is unauthorized
	userId int

	//ctx is context.Context. It is used to cancel the context
	ctx context.Context

	//cancel is function to cancel the ctx.
	cancel context.CancelFunc

	// subscriptions keeps list of subscriptions on Redis
	subscriptions []string
}

// startReadPump pumps messages from the websocket connection to the hub.
func (client *Client) startReadPump() {
	defer func() {
		client.Close()
	}()
	var err error
	client.conn.SetReadLimit(maxMessageSize)
	err = client.conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		client.hub.logger.Error("set read deadline error: " + err.Error())
		return
	}
	client.conn.SetPongHandler(func(string) error {
		err = client.conn.SetReadDeadline(time.Now().Add(pongWait))
		if err != nil {
			client.hub.logger.Error("set read deadline error: " + err.Error())
			return err
		}
		return nil
	})
	for {
		var request []byte
		_, request, err = client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				client.hub.logger.Error("websocket.IsUnexpectedCloseError: " + err.Error())
			}
			break
		}
		parsedRequest := parseRequest(request, client)
		client.hub.broadcast <- parsedRequest
	}
}

// startWritePump pumps messages from the hub to the websocket connection.
func (client *Client) startWritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		client.Close()
		ticker.Stop()
	}()
	for {
		select {
		case message, ok := <-client.send:
			err := client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				client.hub.logger.Error("set write deadline error: " + err.Error())
				return
			}
			if !ok {
				// The hub closed the channel.
				err = client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					client.hub.logger.Error("writeMessage error: " + err.Error())
					return
				}
				return
			}

			w, err := client.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			_, err = w.Write(message)
			if err != nil {
				client.hub.logger.Error("write error: " + err.Error())
				return
			}

			// Add queued chat messages to the current websocket message.
			n := len(client.send)
			for i := 0; i < n; i++ {
				_, err = w.Write(<-client.send)
				if err != nil {
					client.hub.logger.Error("write error: " + err.Error())
					return
				}
			}

			if err = w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			err := client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				client.hub.logger.Error("set write deadline error: " + err.Error())
				return
			}
			if err = client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (client *Client) Close() {
	client.cancel()
	client.hub.unregister <- client
	client.conn.Close()
	client.subscriptions = nil
	client.hub.clientPool.Put(client)
}

// We don't need to refreshTokens tokens, because user in first send request and get token and in second send request to websocket.
func (client *Client) verify(conn *websocket.Conn) int {
	accessToken := conn.Query("auth")
	if accessToken == "" {
		return -1
	} else {
		userId, err := client.hub.accessConverter.ParseToken(accessToken)
		if err != nil {
			return -1
		} else {
			return fb.B2I(userId)
		}
	}
}
