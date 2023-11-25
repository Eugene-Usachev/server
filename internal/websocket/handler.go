package websocket

import (
	"GoServer/Entities"
	"GoServer/internal/service"
	"context"
	"github.com/Eugene-Usachev/fastbytes"
	"github.com/goccy/go-json"
	"strconv"
)

type handlerFunc func(request ParsedRequest)

type Handler struct {
	services *service.Service
	hub      *Hub
	router   []handlerFunc
}

func newHandler(service *service.Service) *Handler {
	var router = make([]handlerFunc, 256)
	handler := &Handler{
		services: service,
		router:   router,
	}

	router[getOnlineUsers] = handler.getOnlineUsers
	router[createChat] = handler.createChat
	return handler
}

func (handler *Handler) createChat(request ParsedRequest) {
	if request.Client.userId == -1 {
		return
	}

	var dto Entities.ChatDTO
	err := json.Unmarshal(request.Data, &dto)
	if err != nil {
		return
	}

	id, err := handler.services.CreateChat(context.Background(), uint(request.Client.userId), dto)
	if err != nil {
		request.Client.send <- fastbytes.S2B(err.Error())
		if err != nil {
			request.Client.hub.logger.Error("WS: updateChat service error:", err)
			return
		}
		return
	}

	dto.ID = id
	var (
		jsonResponse []byte
	)
	jsonResponse, err = json.Marshal(map[string]interface{}{
		"method": "newChat",
		"data":   dto,
	})
	if err != nil {
		request.Client.send <- jsonResponse
		if err != nil {
			request.Client.hub.logger.Error("WS: createChat json error:", err)
			return
		}
		return
	}

	sendDataToMembers(request.Client.hub, jsonResponse, dto.Members)
}

type UpdateDTO struct {
	Entities.ChatUpdateDTO
	ChatId uint `json:"chat_id" binding:"required"`
}

func (handler *Handler) updateChat(request ParsedRequest) {
	if request.Client.userId == -1 {
		return
	}
	var dto UpdateDTO
	err := json.Unmarshal(request.Data, &dto)
	if err != nil {
		return
	}
	chatDTO := Entities.ChatUpdateDTO{
		Name:    dto.Name,
		Avatar:  dto.Avatar,
		Members: dto.Members,
	}
	err = handler.services.UpdateChat(context.Background(), uint(request.Client.userId), dto.ChatId, chatDTO)
	if err != nil {
		request.Client.send <- fastbytes.S2B(err.Error())
		request.Client.hub.logger.Error("WS: updateChat service error:", err)
		return
	}
	var (
		jsonResponse []byte
	)
	jsonResponse, err = json.Marshal(map[string]interface{}{
		"method": "updateChat",
		"data":   dto,
	})
	if err != nil {
		request.Client.hub.logger.Error("WS: updateChat json error:", err)
		return
	}
	sendDataToMembers(request.Client.hub, jsonResponse, dto.Members)
}

func (handler *Handler) deleteChat(request ParsedRequest) {
	if request.Client.userId == -1 {
		return
	}

	chatId, err := strconv.ParseUint(fastbytes.B2S(request.Data), 10, 32)
	if err != nil {
		return
	}
	members, err := handler.services.DeleteChat(context.Background(), uint(request.Client.userId), uint(chatId))
	if err != nil {
		request.Client.send <- fastbytes.S2B(err.Error())
		request.Client.hub.logger.Error("WS: deleteChat service error:", err)
		return
	}
	var jsonResponse []byte
	jsonResponse, err = json.Marshal(map[string]interface{}{
		"method": "deleteChat",
		"data":   chatId,
	})
	if err != nil {
		request.Client.hub.logger.Error("WS: deleteChat json error:", err)
		return
	}
	sendDataToMembers(request.Client.hub, jsonResponse, members)
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
	function := handler.router[request.Method]
	if function != nil {
		function(request)
		return true
	}
	return false
}
