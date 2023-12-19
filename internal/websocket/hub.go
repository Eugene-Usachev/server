package websocket

import (
	"GoServer/internal/service"
	"GoServer/pkg/dashMap"
	"context"
	fb "github.com/Eugene-Usachev/fastbytes"
	"github.com/Eugene-Usachev/fst"
	loggerLib "github.com/Eugene-Usachev/logger"
	"github.com/goccy/go-json"
	"github.com/gofiber/contrib/websocket"
	"github.com/redis/rueidis"
	"sync"
)

type subscription struct {
	m     sync.RWMutex
	slice []*Client
}

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

	// All subscriptions to online clients. Only for authorized clients.
	subscriptions *dashMap.DashMap[*subscription]

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
		subscriptions:   dashMap.NewDashMap[*subscription](512),
		accessConverter: accessConverter,
		logger:          logger,
		redis:           redisClient,
	}
	self.handler = newHandler(self, service)
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

const redisOnlineUsers = "onlineUsers"

func (hub *Hub) Run() {
	hub.logger.Info("hub started")
	defer func() {
		if p := recover(); p != nil {
			hub.logger.Error("hub panic: ", p)
		}
	}()

	for {
		select {
		case client := <-hub.register:
			clientId := client.userId
			if clientId != -1 {
				hub.AuthClients[clientId] = client
				_, _ = hub.redis.Do(client.ctx, hub.redis.B().Smembers().Key(fb.B2S(fb.I2B(clientId))).Build()).ToAny()
				sub := hub.subscriptions.Get(fb.B2S(fb.I2B(clientId)))
				if sub != nil {
					processSubscriptionOnOnline(sub, clientId)
				}
			} else {
				hub.unauthClients[client] = true
			}

		case client := <-hub.unregister:
			if client.userId != -1 {
				strId := fb.B2S(fb.I2B(client.userId))
				if len(client.subscriptions) > 0 {
					for _, sub_ := range client.subscriptions {
						sub := hub.subscriptions.Get(sub_)
						if sub != nil {
							sub.m.Lock()
							l := len(sub.slice)
							if l == 1 {
								hub.subscriptions.Delete(sub_)
								continue
							}
							slice := make([]*Client, l-1, cap(sub.slice))
							index := 0
							for _, client_ := range sub.slice {
								if client_.userId == client.userId {
									break
								}
								index++
							}
							copy(slice, sub.slice[:index])
							copy(slice[index:], sub.slice[index+1:])
						}
					}

					// this code for scale
					//var needToDo []rueidis.Completed
					//needToDo = append(needToDo, hub.redis.B().Smembers().Key(strId).Build())
					//needToDo = append(needToDo, hub.redis.B().Srem().Key(redisOnlineUsers).Member(strId).Build())
					//for _, userId := range client.subscriptions {
					//	needToDo = append(needToDo, hub.redis.B().Srem().Key(userId).Member(strId).Build())
					//}
					//// We call Error only for wait the operation
					//_ = hub.redis.DoMulti(context.Background(), needToDo...)[0].Error()
				}
				subs := hub.subscriptions.Get(strId)
				if subs != nil {
					processSubscriptionOnOffline(subs, client.userId)
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

var welcomeMessage = "Welcome"

// ServeWs handles websocket requests from the peer.
func (hub *Hub) ServeWs(conn *websocket.Conn) {
	client := hub.clientPool.Get().(*Client)
	client.conn = conn
	client.userId = client.verify(conn)
	if client.userId != -1 {
		c := hub.AuthClients[client.userId]
		if c != nil {
			c.Close()
		}
	}
	client.hub.register <- client
	clientCtx := context.Background()
	client.ctx, client.cancel = context.WithCancel(clientCtx)
	hub.redis.Do(client.ctx, hub.redis.B().Sadd().Key(redisOnlineUsers).Member(fb.B2S(fb.I2B(client.userId))).Build())

	go client.startWritePump()
	//go hub.subscribeClient(client)
	res, _ := createResponse(welcomeMessage, []byte{})
	client.send <- res
	client.startReadPump()

}

func (hub *Hub) subscribeToClient(client *Client, id string) {
	if client.userId == -1 {
		return
	}
	client.subscriptions = append(client.subscriptions, id)
	sub := hub.subscriptions.Get(id)
	if sub != nil {
		sub.m.Lock()
		sub.slice = append(sub.slice, client)
		sub.m.Unlock()
	} else {
		hub.subscriptions.Set(id, &subscription{
			m:     sync.RWMutex{},
			slice: []*Client{client},
		})
	}
	// for scale
	//dedicatedClient, cancel := hub.redis.Dedicate()
	//defer func() {
	//	cancel()
	//	dedicatedClient.Close()
	//}()
	//dedicatedClient.Receive(client.ctx, hub.redis.B().Ssubscribe().Channel(fb.B2S(fb.I2B(client.userId))).Build(), func(msg rueidis.PubSubMessage) {
	//
	//	client.send <- fb.S2B(msg.Message)
	//})
}
