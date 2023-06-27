package websocket

import (
	"encoding/json"
	"fmt"
	fb "github.com/Eugene-Usachev/fastbytes"
	"github.com/Eugene-Usachev/fst"
	"github.com/redis/rueidis"
	"log"
	"strconv"
)

var (
	methodNotAllowed = fb.S2B("405 Method Not Allowed\r\n")
)

// Hub maintains the set of active clients and broadcasts messages to the clients.
type Hub struct {
	// Unauthorized clients (!) in this node.
	unauthClients map[*Client]bool

	// Authorized clients (!) in this node.
	authClients map[int]*Client

	// Inbound messages from the clients.
	broadcast chan ParsedRequest

	// register requests from the clients.
	register chan *Client

	// unregister requests from clients.
	unregister chan *Client

	accessConverter *fst.Converter
}

func NewHub(accessConverter *fst.Converter) *Hub {
	return &Hub{
		broadcast:       make(chan ParsedRequest),
		register:        make(chan *Client),
		unregister:      make(chan *Client),
		unauthClients:   make(map[*Client]bool),
		authClients:     make(map[int]*Client),
		accessConverter: accessConverter,
	}
}

func createResponse(messageType string, data string) ([]byte, error) {
	return json.Marshal(map[string]string{
		"method": messageType,
		"data":   data,
	})
}

func (websocketClient *WebsocketClient) Run() {
	log.Println("hub started")
	defer func() {
		if reason := recover(); reason != nil {
			// TODO comment in production
			log.Println("Handled panic, reason: ", reason)
		}
	}()

	for {
		select {
		case client := <-websocketClient.hub.register:
			clientId := client.userId
			if clientId != -1 {
				websocketClient.hub.authClients[clientId] = client
				_, _ = websocketClient.redis.Do(client.ctx, websocketClient.redis.B().Smembers().Key(strconv.Itoa(clientId)).Build()).ToAny()
			} else {
				websocketClient.hub.unauthClients[client] = true
			}

		case client := <-websocketClient.hub.unregister:
			client.cancel()
			client.conn.Close()
			close(client.send)
			if client.userId != -1 {
				strId := strconv.Itoa(client.userId)
				var needToDo []rueidis.Completed
				needToDo = append(needToDo, websocketClient.redis.B().Smembers().Key(strId).Build())
				needToDo = append(needToDo, websocketClient.redis.B().Srem().Key(onlineUsers).Member(strId).Build())
				for _, userId := range client.subscriptions {
					needToDo = append(needToDo, websocketClient.redis.B().Srem().Key(userId).Member(strId).Build())
				}
				resp := websocketClient.redis.DoMulti(websocketClient.ctx, needToDo...)
				fmt.Sprintf("Type: %T", resp[0])
				delete(websocketClient.hub.authClients, client.userId)
			} else {
				delete(websocketClient.hub.unauthClients, client)
			}

		case parsedRequest := <-websocketClient.hub.broadcast:
			if parsedRequest.Client == nil {
				continue
			}
			websocketClient.handler.handle(parsedRequest)
		}
	}
}

// TODO r
func (hub Hub) getOnlineUsers(necessaryToGet []interface{}) []byte {
	var onlineUsers = []int{}
	for _, userId := range necessaryToGet {
		userIdI := int(userId.(float64))
		if hub.authClients[userIdI] != nil {
			onlineUsers = append(onlineUsers, userIdI)
		}
	}

	onlineUsersJSON, _ := json.Marshal(map[string]interface{}{
		"data":   onlineUsers,
		"method": "getOnlineUsers",
	})

	return onlineUsersJSON
}
