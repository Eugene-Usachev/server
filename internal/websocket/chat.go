package websocket

import (
	"GoServer/Entities"
	"context"
	"github.com/Eugene-Usachev/fastbytes"
	"github.com/goccy/go-json"
	"strconv"
)

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
	jsonResponse, err = createResponse("newChat", dto)
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
	ChatId uint `json:"id" binding:"required"`
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
	chatId, err := strconv.ParseUint(fastbytes.B2S(request.Data), 10, 64)
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
