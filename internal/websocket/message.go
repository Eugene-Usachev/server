package websocket

import (
	"GoServer/Entities"
	"context"
	"github.com/Eugene-Usachev/fastbytes"
	"github.com/goccy/go-json"
	"strconv"
)

func (handler *Handler) sendMessage(request ParsedRequest) {
	if request.Client.userId == -1 {
		return
	}

	var message Entities.MessageDTO
	err := json.Unmarshal(request.Data, &message)
	if err != nil {
		return
	}
	id, members, date, err := handler.services.SaveMessage(context.Background(), uint(request.Client.userId), message)
	if err != nil {
		request.Client.hub.logger.Error("WS: sendMessage service error:", err.Error())
		request.Client.send <- []byte(err.Error())
		handler.hub.logger.Error("WS: sendMessage service error:", err.Error())
		return
	}
	message.ID = id
	message.Date = date

	var jsonResponse []byte
	jsonResponse, err = createResponse("newMessage", message)
	if err != nil {
		request.Client.hub.logger.Error("WS: sendMessage json error:", err)
		request.Client.send <- jsonResponse
		if err != nil {
			request.Client.hub.logger.Error("WS: createChat json error:", err)
			return
		}
		return
	}

	sendDataToMembers(request.Client.hub, jsonResponse, members)
}

type UpdateMessageDTO struct {
	MessageId uint   `json:"message_id" binding:"required"`
	Data      string `json:"data" binding:"required"`
}

func (handler *Handler) updateMessage(request ParsedRequest) {
	if request.Client.userId == -1 {
		return
	}

	var dto UpdateMessageDTO
	err := json.Unmarshal(request.Data, &dto)
	if err != nil {
		return
	}

	members, err := handler.services.UpdateMessage(context.Background(), dto.MessageId, uint(request.Client.userId), dto.Data)
	if err != nil {
		request.Client.send <- []byte(err.Error())
		handler.hub.logger.Error("WS: updateMessage service error:", err.Error())
		return
	}

	var jsonResponse []byte
	jsonResponse, err = createResponse("updateMessage", dto)
	if err != nil {
		request.Client.hub.logger.Error("WS: createChat json error:", err)
		request.Client.send <- jsonResponse
		return
	}

	sendDataToMembers(request.Client.hub, jsonResponse, members)
}

func (handler *Handler) deleteMessage(request ParsedRequest) {
	if request.Client.userId == -1 {
		return
	}

	dto, err := strconv.Atoi(fastbytes.B2S(request.Data))
	if err != nil {
		return
	}

	members, err := handler.services.DeleteMessage(context.Background(), uint(dto), uint(request.Client.userId))
	if err != nil {
		request.Client.send <- []byte(err.Error())
		handler.hub.logger.Error("WS: deleteMessage service error:", err.Error())
		return
	}

	var jsonResponse []byte
	jsonResponse, err = createResponse("deleteMessage", dto)
	if err != nil {
		request.Client.hub.logger.Error("WS: createChat json error:", err)
		request.Client.send <- jsonResponse
		return
	}

	sendDataToMembers(request.Client.hub, jsonResponse, members)
}
