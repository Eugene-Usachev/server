package websocket

import (
	"GoServer/internal/service"
	"context"
	fb "github.com/Eugene-Usachev/fastbytes"
	"github.com/Eugene-Usachev/fst"
	loggerLib "github.com/Eugene-Usachev/logger"
	"github.com/goccy/go-json"
	"github.com/gofiber/contrib/websocket"
	"github.com/redis/rueidis"
	"strconv"
	"sync"
)

var (
	methodNotAllowed = fb.S2B("405 Method Not Allowed\r\n")
)

// Hub maintains the set of active clients and broadcasts messages to the clients.
type Hub struct {
	// Unauthorized clients (!) in this node.
	unauthClients map[*Client]bool
	// Authorized clients (!) in this node.
	AuthClients map[int]*Client
	// Inbound messages from the clients.
	broadcast chan ParsedRequest
	// register requests from the clients.
	register chan *Client
	// unregister requests from clients.
	unregister chan *Client

	// TODO redis is not SOLID
	redis rueidis.Client

	handler         *Handler
	logger          *loggerLib.FastLogger
	clientPool      sync.Pool
	accessConverter *fst.Converter
}

func NewHub(service *service.Service, redisClient rueidis.Client, accessConverter *fst.Converter, logger *loggerLib.FastLogger) *Hub {
	InitConfig(logger)
	self := &Hub{
		broadcast:       make(chan ParsedRequest),
		register:        make(chan *Client),
		unregister:      make(chan *Client),
		unauthClients:   make(map[*Client]bool),
		AuthClients:     make(map[int]*Client),
		accessConverter: accessConverter,
		logger:          logger,
		handler:         newHandler(service),
		redis:           redisClient,
	}

	self.clientPool = sync.Pool{
		New: func() any {
			return self.newEmptyClient()
		},
	}

	return self
}

func (hub *Hub) newEmptyClient() *Client {
	return &Client{
		hub:           hub,
		conn:          nil,
		send:          make(chan []byte, 5),
		userId:        -1,
		ctx:           nil,
		cancel:        nil,
		subscriptions: nil,
	}
}

func createResponse(method any, data any) ([]byte, error) {
	return json.Marshal(map[string]any{
		"method": method,
		"data":   data,
	})
}

func (hub *Hub) Run() {
	hub.logger.Info("hub started")
	defer func() {
		if reason := recover(); reason != nil {
			hub.logger.Error("Handled panic, reason: ", reason)
		}
	}()

	for {
		select {
		case client := <-hub.register:
			clientId := client.userId
			if clientId != -1 {
				hub.AuthClients[clientId] = client
				_, _ = hub.redis.Do(client.ctx, hub.redis.B().Smembers().Key(strconv.Itoa(clientId)).Build()).ToAny()
			} else {
				hub.unauthClients[client] = true
			}

		case client := <-hub.unregister:
			if client.userId != -1 {
				if len(client.subscriptions) > 0 {
					strId := strconv.Itoa(client.userId)
					var needToDo []rueidis.Completed
					needToDo = append(needToDo, hub.redis.B().Smembers().Key(strId).Build())
					needToDo = append(needToDo, hub.redis.B().Srem().Key(onlineUsers).Member(strId).Build())
					for _, userId := range client.subscriptions {
						needToDo = append(needToDo, hub.redis.B().Srem().Key(userId).Member(strId).Build())
					}
					// We call Error only for wait the operation
					_ = hub.redis.DoMulti(context.Background(), needToDo...)[0].Error()
				}
				delete(hub.AuthClients, client.userId)
			} else {
				delete(hub.unauthClients, client)
			}

		case parsedRequest := <-hub.broadcast:
			if parsedRequest.Client == nil {
				continue
			}
			hub.handler.handle(parsedRequest)
		}
	}
}

func (hub *Hub) Close() {
	hub.redis.Close()
}

var welcomeMessage = fb.S2B("Welcome")

// ServeWs handles websocket requests from the peer.
func (hub *Hub) ServeWs(conn *websocket.Conn) {
	client := hub.clientPool.Get().(*Client)
	client.conn = conn
	client.userId = client.verify(conn)
	client.hub.register <- client
	clientCtx := context.Background()
	client.ctx, client.cancel = context.WithCancel(clientCtx)
	hub.redis.Do(client.ctx, hub.redis.B().Sadd().Key(onlineUsers).Member(strconv.Itoa(client.userId)).Build())

	// Allow collection of memory referenced by the caller by doing all work in new goroutines.
	go client.startWritePump()
	go hub.subscribeClient(client)
	client.send <- welcomeMessage
	client.startReadPump()

}

func (hub *Hub) subscribeClient(client *Client) {
	dedicatedClient, cancel := hub.redis.Dedicate()
	defer func() {
		cancel()
		dedicatedClient.Close()
	}()
	dedicatedClient.Receive(client.ctx, hub.redis.B().Ssubscribe().Channel(strconv.Itoa(client.userId)).Build(), func(msg rueidis.PubSubMessage) {
		client.send <- fb.S2B(msg.Message)
	})
}

// TODO r THIS IS A SERVICE
func (hub *Hub) getOnlineUsers(necessaryToGet []interface{}) []byte {
	var onlineUsers = []int{}
	for _, userId := range necessaryToGet {
		userIdI := int(userId.(float64))
		if hub.AuthClients[userIdI] != nil {
			onlineUsers = append(onlineUsers, userIdI)
		}
	}

	onlineUsersJSON, _ := createResponse("getOnlineUsers", onlineUsers)

	return onlineUsersJSON
}
