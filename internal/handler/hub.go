package handler

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
)

// Hub maintains the set of active clients and broadcasts messages to the clients.
type Hub struct {
	// Unauthorized clients.
	unauthClients map[*Client]bool

	// Authorized clients.
	authClients map[int64]*Client

	// Inbound messages from the clients.
	broadcast chan ParsedRequest

	// register requests from the clients.
	register chan *Client

	// unregister requests from clients.
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:     make(chan ParsedRequest),
		register:      make(chan *Client),
		unregister:    make(chan *Client),
		unauthClients: make(map[*Client]bool),
		authClients:   make(map[int64]*Client),
	}
}

func createResponse(messageType string, data string) ([]byte, error) {
	return json.Marshal(map[string]string{
		"method": messageType,
		"data":   data,
	})
}

func (hub *Hub) Run(handler *Handler) {
	logrus.Info("Hub started")
	defer handlePanic()
	for {
		select {
		case client := <-hub.register:
			clientId := client.userId
			if clientId != 0 {
				msg, _ := createResponse("userConnected", fmt.Sprintf("%d", clientId))
				for _, authClient := range hub.authClients {
					authClient.send <- msg
				}
				for unauthClient := range hub.unauthClients {
					unauthClient.send <- msg
				}
				hub.authClients[clientId] = client
			} else {
				hub.unauthClients[client] = true
			}

		case client := <-hub.unregister:
			clientId := client.userId
			if clientId != 0 {
				msg, _ := createResponse("userDisconnected", fmt.Sprintf("%d", clientId))
				for _, authClient := range hub.authClients {
					authClient.send <- msg
				}
				for unauthClient := range hub.unauthClients {
					unauthClient.send <- msg
				}
				delete(hub.authClients, clientId)
				close(client.send)
			} else {
				delete(hub.unauthClients, client)
				close(client.send)
			}

		case parsedRequest := <-hub.broadcast:
			if parsedRequest.Client == nil {
				continue
			}
			switch parsedRequest.Method {
			case "getOnlineUsers":
				parsedRequest.Client.send <- hub.getOnlineUsers(parsedRequest.Data.([]interface{}))

			case "sendMessage":
				handler.sendMessage(parsedRequest, hub)

			case "updateMessage":
				handler.updateMessage(parsedRequest, hub)

			case "deleteMessage":
				handler.deleteMessage(parsedRequest, hub)

			case "createChat":
				handler.createChat(parsedRequest, hub)

			case "updateChat":
				handler.updateChat(parsedRequest, hub)

			case "deleteChat":
				handler.deleteChat(parsedRequest, hub)

			default:
				parsedRequest.Client.send <- []byte("405 Method Not Allowed\r\n")
			}
		}
	}
}

func (hub Hub) getOnlineUsers(necessaryToGet []interface{}) []byte {
	var onlineUsers = []int64{}
	for _, userId := range necessaryToGet {
		userIdI := int64(userId.(float64))
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
