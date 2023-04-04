package handler

import (
	"GoServer/pkg/jwt"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 4096
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type ParsedRequest struct {
	Method string `json:"method"`
	Data   any    `json:"data"`
	Client *Client
}

type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	//if UserId is 0 then user is unauthorized
	userId int64
}

// readPump pumps messages from the websocket connection to the hub.
func (client *Client) readPump() {
	defer func() {
		client.hub.unregister <- client
		err := client.conn.Close()
		if err != nil {
			logrus.Error(err)
			return
		}
	}()
	client.conn.SetReadLimit(maxMessageSize)
	err := client.conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		logrus.Error(err)
		return
	}
	client.conn.SetPongHandler(func(string) error {
		err = client.conn.SetReadDeadline(time.Now().Add(pongWait))
		if err != nil {
			logrus.Error(err)
			return err
		}
		return nil
	})
	for {
		var request []byte
		_, request, err = client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		parsedRequest := parseRequest(request, client)
		client.hub.broadcast <- parsedRequest
	}
}

// writePump pumps messages from the hub to the websocket connection.
func (client *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		err := client.conn.Close()
		if err != nil {
			logrus.Error(err)
			return
		}
	}()
	for {
		select {
		case message, ok := <-client.send:
			err := client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				logrus.Error(err)
				return
			}
			if !ok {
				// The hub closed the channel.
				err = client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					logrus.Error(err)
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
				logrus.Error(err)
				return
			}

			// Add queued chat messages to the current websocket message.
			n := len(client.send)
			for i := 0; i < n; i++ {
				_, err = w.Write(<-client.send)
				if err != nil {
					logrus.Error(err)
					return
				}
			}

			if err = w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			err := client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				logrus.Error(err)
				return
			}
			if err = client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// ServeWs handles websocket requests from the peer.
func (hub *Hub) ServeWs(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.userId = client.verify(ctx)
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in new goroutines.
	go client.writePump()
	go client.readPump()

	client.send <- []byte("Welcome")
}

// We don't need to refresh tokens, because user in first send request and get token and in second send request to websocket.
func (client *Client) verify(ctx *gin.Context) int64 {
	accessToken, err := ctx.Cookie("accessToken")
	if err != nil {
		return 0
	}
	if accessToken == "" {
		return 0
	} else {
		userId, err := jwt.ParseAccessToken(accessToken)
		if err != nil {
			return 0
		}
		return int64(userId)
	}
}

func parseRequest(request []byte, client *Client) ParsedRequest {
	var parsedRequest ParsedRequest
	err := json.Unmarshal(request, &parsedRequest)
	if err != nil {
		log.Printf("error: %v", err)
		return ParsedRequest{}
	}
	parsedRequest.Client = client

	return parsedRequest
}
