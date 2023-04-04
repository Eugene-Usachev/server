package handler

import (
	"GoServer/Entities"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func (handler *Handler) createChat(request ParsedRequest, hub *Hub) {
	defer handlePanic()
	if !prepareAuthRequest(request) {
		return
	}

	var dto Entities.ChatDTO
	err := json.Unmarshal([]byte(request.Data.(string)), &dto)
	if err != nil {
		return
	}

	id, err := handler.services.CreateChat(context.Background(), request.Client.userId, dto)
	if err != nil {
		request.Client.send <- []byte(err.Error())
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
		return
	}

	sendDataToMembers(hub, jsonResponse, dto.Members)
}

type UpdateDTO struct {
	Entities.ChatUpdateDTO
	ChatId int64 `json:"chat_id" binding:"required"`
}

func (handler *Handler) updateChat(request ParsedRequest, hub *Hub) {
	defer handlePanic()
	if !prepareAuthRequest(request) {
		return
	}
	var dto UpdateDTO
	err := json.Unmarshal([]byte(request.Data.(string)), &dto)
	if err != nil {
		return
	}
	chatDTO := Entities.ChatUpdateDTO{
		Name:    dto.Name,
		Avatar:  dto.Avatar,
		Members: dto.Members,
	}
	err = handler.services.UpdateChat(context.Background(), request.Client.userId, dto.ChatId, chatDTO)
	if err != nil {
		request.Client.send <- []byte(err.Error())
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
		log.Printf("error encoding json: %v\n", err)
		return
	}
	sendDataToMembers(hub, jsonResponse, dto.Members)
}

func (handler *Handler) deleteChat(request ParsedRequest, hub *Hub) {
	defer handlePanic()
	if !prepareAuthRequest(request) {
		return
	}
	var dto int64
	err := json.Unmarshal([]byte(request.Data.(string)), &dto)
	if err != nil {
		return
	}
	members, err := handler.services.DeleteChat(context.Background(), request.Client.userId, dto)
	if err != nil {
		request.Client.send <- []byte(err.Error())
		return
	}
	var jsonResponse []byte
	jsonResponse, err = json.Marshal(map[string]interface{}{
		"method": "deleteChat",
		"data":   dto,
	})
	if err != nil {
		log.Printf("error encoding json: %v\n", err)
		return
	}
	sendDataToMembers(hub, jsonResponse, members)
}

// getChats returns a list of chats and chats.
func (handler *Handler) getChats(ctx *gin.Context) {
	userId, exists := ctx.Get("userId")
	userIdI := int64(userId.(uint))
	if !exists {
		NewErrorResponse(ctx, http.StatusUnauthorized, "invalid auth token")
		return
	}

	avatar, name, surname, friends, subscribers, chatLists, chats, err := handler.services.GetChats(ctx.Request.Context(), userIdI)
	if err != nil {
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"chats":       chats,
		"chatLists":   chatLists,
		"avatar":      avatar,
		"friends":     friends,
		"subscribers": subscribers,
		"name":        name,
		"surname":     surname,
	})
}

type chatListUpdateDTO struct {
	NewChatLists string `json:"newChatLists"`
}

func (handler *Handler) UpdateChatLists(ctx *gin.Context) {
	userId, exists := ctx.Get("userId")
	userIdI := int64(userId.(uint))
	if !exists || userIdI < 1 {
		NewErrorResponse(ctx, http.StatusUnauthorized, "invalid auth token")
		return
	}

	var dto chatListUpdateDTO

	err := ctx.BindJSON(&dto)
	if err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	err = handler.services.UpdateChatLists(ctx.Request.Context(), userIdI, string(dto.NewChatLists))
	if err != nil {
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}
