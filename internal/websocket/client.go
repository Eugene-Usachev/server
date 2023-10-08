package websocket

import (
	"GoServer/pkg/fasthttp_utils"
	"context"
	fb "github.com/Eugene-Usachev/fastbytes"
	loggerLib "github.com/Eugene-Usachev/logger"
	"github.com/fasthttp/websocket"
	"github.com/redis/rueidis"
	"github.com/valyala/fasthttp"
	"strconv"
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

var upgrader = websocket.FastHTTPUpgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// TODO use pool
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	logger *loggerLib.FastLogger

	// Buffered channel of outbound messages.
	send chan []byte

	//if UserId is -1 then user is unauthorized
	userId int

	//ctx is context.Context. It is used to cancel the context
	ctx context.Context

	//cancel is function to cancel the ctx.
	cancel context.CancelFunc

	// subscriptions keeps list of subscriptions on Redis
	subscriptions []string
}

// readPump pumps messages from the websocket connection to the hub.
func (client *Client) readPump() {
	defer func() {
		client.hub.unregister <- client
	}()
	client.conn.SetReadLimit(maxMessageSize)
	err := client.conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		client.logger.Error("set read deadline error: " + err.Error())
		return
	}
	client.conn.SetPongHandler(func(string) error {
		err = client.conn.SetReadDeadline(time.Now().Add(pongWait))
		if err != nil {
			client.logger.Error("set read deadline error: " + err.Error())
			return err
		}
		return nil
	})
	for {
		var request []byte
		_, request, err = client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				client.logger.Error("websocket.IsUnexpectedCloseError: " + err.Error())
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
		client.hub.unregister <- client
	}()
	for {
		select {
		case message, ok := <-client.send:
			err := client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				client.logger.Error("set write deadline error: " + err.Error())
				return
			}
			if !ok {
				// The hub closed the channel.
				err = client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					client.logger.Error("writeMessage error: " + err.Error())
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
				client.logger.Error("write error: " + err.Error())
				return
			}

			// Add queued chat messages to the current websocket message.
			n := len(client.send)
			for i := 0; i < n; i++ {
				_, err = w.Write(<-client.send)
				if err != nil {
					client.logger.Error("write error: " + err.Error())
					return
				}
			}

			if err = w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			err := client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				client.logger.Error("set write deadline error: " + err.Error())
				return
			}
			if err = client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

var welcomeMessage = fb.S2B("Welcome")

// ServeWs handles websocket requests from the peer.
func (websocketClient *WebsocketClient) ServeWs(ctx *fasthttp.RequestCtx) error {
	err := upgrader.Upgrade(ctx, func(conn *websocket.Conn) {
		client := &Client{hub: websocketClient.hub, conn: conn, send: make(chan []byte, 256)}
		client.userId = client.verify(ctx)
		client.hub.register <- client
		clientCtx := context.Background()
		client.ctx, client.cancel = context.WithCancel(clientCtx)
		websocketClient.redis.Do(client.ctx, websocketClient.redis.B().Sadd().Key(onlineUsers).Member(strconv.Itoa(client.userId)).Build())

		// Allow collection of memory referenced by the caller by doing all work in new goroutines.
		go client.writePump()
		go client.readPump()
		go websocketClient.subscribeClient(client)

		//client.send <- welcomeMessage
	})
	if err != nil {
		websocketClient.logger.Error("upgrade error: " + err.Error())
		return err
	}

	return nil
}

func (websocketClient *WebsocketClient) subscribeClient(client *Client) {
	dedicatedClient, err := websocketClient.redis.Dedicate()
	if err != nil {
		websocketClient.hub.unregister <- client
		return
	}
	defer dedicatedClient.Close()
	dedicatedClient.Receive(client.ctx, websocketClient.redis.B().Ssubscribe().Channel(strconv.Itoa(client.userId)).Build(), func(msg rueidis.PubSubMessage) {
		client.send <- fb.S2B(msg.Message)
	})
}

// We don't need to refreshTokens tokens, because user in first send request and get token and in second send request to websocket.
func (client *Client) verify(ctx *fasthttp.RequestCtx) int {
	accessToken := fb.B2S(fasthttp_utils.GetAuthorizationHeader(ctx))
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
