package websocket

import (
	"GoServer/internal/service"
)

type handlerFunc func(request ParsedRequest)

type Handler struct {
	service *service.Service
	hub     *Hub
	router  []handlerFunc
}

func newHandler(service *service.Service) *Handler {
	var router = make([]handlerFunc, 256)
	handler := &Handler{
		service: service,
		router:  router,
	}

	router[getOnlineUsersMethod] = handler.getOnlineUsers
	return handler
}

// handle return false if method not found and true otherwise.
func (handler *Handler) handle(request ParsedRequest) bool {
	function := handler.router[request.Method]
	if function != nil {
		function(request)
		return true
	}
	return false
}
