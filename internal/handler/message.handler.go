package handler

import (
	"GoServer/Entities"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

const socketBasRequest = "400 Invalid message type\n"

func (handler *Handler) sendMessage(request ParsedRequest, hub *Hub) {
	defer handlePanic()
	if !prepareAuthRequest(request) {
		return
	}
	var message Entities.MessageDTO
	err := json.Unmarshal([]byte(request.Data.(string)), &message)
	if err != nil {
		log.Println(err.Error())
		return
	}
	id, members, date, err := handler.services.SaveMessage(context.Background(), request.Client.userId, message)
	if err != nil {
		request.Client.send <- []byte(err.Error())
		return
	}

	message.ID = id
	message.Date = date
	var (
		jsonResponse []byte
	)
	jsonResponse, err = json.Marshal(map[string]interface{}{
		"method": "newMessage",
		"data":   message,
	})
	if err != nil {
		request.Client.send <- []byte(err.Error())
		return
	}

	sendDataToMembers(hub, jsonResponse, members)
}

type UpdateMessageDTO struct {
	MessageId int64  `json:"message_id" binding:"required"`
	Data      string `json:"data" binding:"required"`
}

func (handler *Handler) getLastMessages(ctx *gin.Context) {
	userId, exists := ctx.Get("userId")
	if !exists {
		NewErrorResponse(ctx, http.StatusBadRequest, "userId not found")
	}
	chatsId := ctx.Query("chats_id")

	messages, err := handler.services.GetLastMessages(ctx.Request.Context(), userId.(uint), chatsId)
	if err != nil {
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	ctx.JSON(http.StatusOK, gin.H{"messages": messages})
}

func (handler *Handler) updateMessage(request ParsedRequest, hub *Hub) {
	defer handlePanic()
	if !prepareAuthRequest(request) {
		return
	}
	var dto UpdateMessageDTO
	err := json.Unmarshal([]byte(request.Data.(string)), &dto)
	if err != nil {
		return
	}
	members, err := handler.services.UpdateMessage(context.Background(), dto.MessageId, request.Client.userId, dto.Data)
	if err != nil {
		request.Client.send <- []byte(err.Error())
		return
	}
	var jsonResponse []byte
	jsonResponse, err = json.Marshal(map[string]any{
		"method": "updateMessage",
		"data":   []any{dto.Data, dto.MessageId},
	})
	if err != nil {
		log.Printf("error encoding json: %v\n", err)
		return
	}
	sendDataToMembers(hub, jsonResponse, members)
}

func (handler *Handler) deleteMessage(request ParsedRequest, hub *Hub) {
	defer handlePanic()
	if !prepareAuthRequest(request) {
		return
	}
	var dto int64
	err := json.Unmarshal([]byte(request.Data.(string)), &dto)
	if err != nil {
		return
	}
	members, err := handler.services.DeleteMessage(context.Background(), dto, request.Client.userId)
	if err != nil {
		request.Client.send <- []byte(err.Error())
		return
	}
	var jsonResponse []byte
	jsonResponse, err = json.Marshal(map[string]interface{}{
		"method": "deleteMessage",
		"data":   dto,
	})
	if err != nil {
		log.Printf("error encoding json: %v\n", err)
		return
	}
	sendDataToMembers(hub, jsonResponse, members)
}

func handlePanic() {
	if err := recover(); err != nil {
		log.Printf("panic recovered: %v", err)
	}
}

func prepareAuthRequest(parsedRequest ParsedRequest) bool {
	if parsedRequest.Client.userId == 0 {
		parsedRequest.Client.send <- []byte("401 Not Allowed\r\n")
		return false
	}
	return true
}

func sendDataToMembers(hub *Hub, data []byte, members []int64) {
	defer handlePanic()
	for _, member := range members {
		client := hub.authClients[member]
		if client != nil {
			client.send <- data
		}
	}
}

func (handler *Handler) getMessages(ctx *gin.Context) {
	chatId, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil || chatId < 1 {
		NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}
	var offset uint64
	offset, err = strconv.ParseUint(ctx.Query("offset"), 10, 64)
	if err != nil || offset < 0 {
		NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	var messages = [20]Entities.Message{}
	messages, err = handler.services.GetMessages(ctx.Request.Context(), uint(chatId), uint(offset))
	if err != nil {
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
	}
	ctx.JSON(http.StatusOK, gin.H{"messages": messages})
}
