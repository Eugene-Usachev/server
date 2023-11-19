package handler

import (
	"github.com/gofiber/fiber/v2"
)

//func (handler *Handler) createChat(request websocket2.ParsedRequest, hub *websocket2.Hub) {
//	defer handlePanic()
//	if !prepareAuthRequest(request) {
//		return
//	}
//
//	var dto Entities.ChatDTO
//	err := json.Unmarshal([]byte(request.Data.(string)), &dto)
//	if err != nil {
//		return
//	}
//
//	id, err := handler.services.CreateChat(context.Background(), request.Client.userId, dto)
//	if err != nil {
//		request.Client.send <- []byte(err.Error())
//		return
//	}
//
//	dto.Id = id
//	var (
//		jsonResponse []byte
//	)
//	jsonResponse, err = json.Marshal(map[string]interface{}{
//		"method": "newChat",
//		"data":   dto,
//	})
//	if err != nil {
//		request.Client.send <- jsonResponse
//		return
//	}
//
//	sendDataToMembers(hub, jsonResponse, dto.Members)
//}
//
//type UpdateDTO struct {
//	Entities.ChatUpdateDTO
//	ChatId uint `json:"chat_id" binding:"required"`
//}
//
//func (handler *Handler) updateChat(request websocket2.ParsedRequest, hub *websocket2.Hub) {
//	defer handlePanic()
//	if !prepareAuthRequest(request) {
//		return
//	}
//	var dto UpdateDTO
//	err := json.Unmarshal([]byte(request.Data.(string)), &dto)
//	if err != nil {
//		return
//	}
//	chatDTO := Entities.ChatUpdateDTO{
//		Name:    dto.Name,
//		Avatar:  dto.Avatar,
//		Members: dto.Members,
//	}
//	err = handler.services.UpdateChat(context.Background(), request.Client.userId, dto.ChatId, chatDTO)
//	if err != nil {
//		request.Client.send <- []byte(err.Error())
//		return
//	}
//	var (
//		jsonResponse []byte
//	)
//	jsonResponse, err = json.Marshal(map[string]interface{}{
//		"method": "updateChat",
//		"data":   dto,
//	})
//	if err != nil {
//		log.Printf("error encoding json: %v\n", err)
//		return
//	}
//	sendDataToMembers(hub, jsonResponse, dto.Members)
//}
//
//func (handler *Handler) deleteChat(request websocket2.ParsedRequest, hub *websocket2.Hub) {
//	defer handlePanic()
//	if !prepareAuthRequest(request) {
//		return
//	}
//	var dto uint
//	err := json.Unmarshal([]byte(request.Data.(string)), &dto)
//	if err != nil {
//		return
//	}
//	members, err := handler.services.DeleteChat(context.Background(), request.Client.userId, dto)
//	if err != nil {
//		request.Client.send <- []byte(err.Error())
//		return
//	}
//	var jsonResponse []byte
//	jsonResponse, err = json.Marshal(map[string]interface{}{
//		"method": "deleteChat",
//		"data":   dto,
//	})
//	if err != nil {
//		log.Printf("error encoding json: %v\n", err)
//		return
//	}
//	sendDataToMembers(hub, jsonResponse, members)
//}

func (handler *Handler) getChatsList(ctx *fiber.Ctx) error {
	userId, ok := ctx.Locals("userId").(uint)
	if !ok || userId == 0 {
		return NewErrorResponse(ctx, fiber.StatusUnauthorized, "invalid auth token")
	}

	friends, subscribers, chatLists, err := handler.services.Chat.GetChatsListAndInfoForUser(ctx.Context(), userId)
	if err != nil {
		return NewErrorResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"chatsList":   chatLists,
		"friends":     friends,
		"subscribers": subscribers,
	})
}

type chatListUpdateDTO struct {
	NewChatLists string `json:"newChatLists"`
}

func (handler *Handler) UpdateChatLists(ctx *fiber.Ctx) error {
	userId := ctx.Locals("userId")
	userIdI := userId.(uint)
	if userIdI < 1 {
		return NewErrorResponse(ctx, fiber.StatusUnauthorized, "invalid auth token")
	}

	var dto chatListUpdateDTO

	err := ctx.BodyParser(&dto)
	if err != nil {
		return NewErrorResponse(ctx, fiber.StatusBadRequest, err.Error())
	}

	err = handler.services.Chat.UpdateChatLists(ctx.Context(), userIdI, string(dto.NewChatLists))
	if err != nil {
		return NewErrorResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}
