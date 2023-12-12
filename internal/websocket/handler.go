package websocket

import (
	"GoServer/internal/service"
	"strconv"
)

type handlerFunc func(request ParsedRequest)

type Handler struct {
	services *service.Service
	hub      *Hub
	router   []handlerFunc
}

func newHandler(hub *Hub, service *service.Service) *Handler {
	var router = make([]handlerFunc, size)
	handler := &Handler{
		services: service,
		router:   router,
		hub:      hub,
	}

	router[getOnlineUsers] = handler.getOnlineUsers
	router[createChat] = handler.createChat
	router[updateChat] = handler.updateChat
	router[deleteChat] = handler.deleteChat
	router[sendMessage] = handler.sendMessage
	router[updateMessage] = handler.updateMessage
	router[deleteMessage] = handler.deleteMessage
	return handler
}

func sendDataToMembers(hub *Hub, jsonResponse []byte, members []uint) {
	for _, memberId := range members {
		// TODO sendTo
		if member, ok := hub.AuthClients[int(memberId)]; ok {
			member.send <- jsonResponse
		}
	}
}

// handle return false if method not found and true otherwise.
func (handler *Handler) handle(request ParsedRequest) bool {
	method, err := strconv.Atoi(string(request.Method))
	if err != nil {
		handler.hub.logger.Error("Unexpected method: ", request.Method)
		return false
	}
	if uint8(method) >= size {
		handler.hub.logger.Error("Unexpected method: ", request.Method)
		return false
	}
	function := handler.router[method]
	if function != nil {
		function(request)
		return true
	}
	return false
}
