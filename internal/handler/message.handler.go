package handler

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"strconv"
)

//func (handler *Handler) sendMessage(request websocket2.ParsedRequest, hub *websocket2.Hub) {
//	defer handlePanic()
//	if !prepareAuthRequest(request) {
//		return
//	}
//	var message Entities.MessageDTO
//	err := json.Unmarshal([]byte(request.Data.(string)), &message)
//	if err != nil {
//		log.Println(err.Error())
//		return
//	}
//	id, members, date, err := handler.services.SaveMessage(context.Background(), request.Client.userId, message)
//	if err != nil {
//		request.Client.send <- []byte(err.Error())
//		return
//	}
//
//	message.ID = id
//	message.Date = date
//	var (
//		jsonResponse []byte
//	)
//	jsonResponse, err = json.Marshal(map[string]interface{}{
//		"method": "newMessage",
//		"data":   message,
//	})
//	if err != nil {
//		request.Client.send <- []byte(err.Error())
//		return
//	}
//
//	sendDataToMembers(hub, jsonResponse, members)
//}
//
//type UpdateMessageDTO struct {
//	MessageId uint   `json:"message_id" binding:"required"`
//	Data      string `json:"data" binding:"required"`
//}
//
//func (handler *Handler) updateMessage(request websocket2.ParsedRequest, hub *websocket2.Hub) {
//	defer handlePanic()
//	if !prepareAuthRequest(request) {
//		return
//	}
//	var dto UpdateMessageDTO
//	err := json.Unmarshal([]byte(request.Data.(string)), &dto)
//	if err != nil {
//		return
//	}
//	members, err := handler.services.UpdateMessage(context.Background(), dto.MessageId, request.Client.userId, dto.Data)
//	if err != nil {
//		request.Client.send <- []byte(err.Error())
//		return
//	}
//	var jsonResponse []byte
//	jsonResponse, err = json.Marshal(map[string]any{
//		"method": "updateMessage",
//		"data":   []any{dto.Data, dto.MessageId},
//	})
//	if err != nil {
//		log.Printf("error encoding json: %v\n", err)
//		return
//	}
//	sendDataToMembers(hub, jsonResponse, members)
//}
//
//func (handler *Handler) deleteMessage(request websocket2.ParsedRequest, hub *websocket2.Hub) {
//	defer handlePanic()
//	if !prepareAuthRequest(request) {
//		return
//	}
//	var dto uint
//	err := json.Unmarshal([]byte(request.Data.(string)), &dto)
//	if err != nil {
//		return
//	}
//	members, err := handler.services.DeleteMessage(context.Background(), dto, request.Client.userId)
//	if err != nil {
//		request.Client.send <- []byte(err.Error())
//		return
//	}
//	var jsonResponse []byte
//	jsonResponse, err = json.Marshal(map[string]interface{}{
//		"method": "deleteMessage",
//		"data":   dto,
//	})
//	if err != nil {
//		log.Printf("error encoding json: %v\n", err)
//		return
//	}
//	sendDataToMembers(hub, jsonResponse, members)
//}
//
//func prepareAuthRequest(parsedRequest websocket2.ParsedRequest) bool {
//	if parsedRequest.Client.userId == 0 {
//		parsedRequest.Client.send <- []byte("401 Not Allowed\r\n")
//		return false
//	}
//	return true
//}
//
//func sendDataToMembers(hub *websocket2.Hub, data []byte, members []uint) {
//	defer handlePanic()
//	for _, member := range members {
//		client := hub.authClients[member]
//		if client != nil {
//			client.send <- data
//		}
//	}
//}

func handlePanic() {
	if err := recover(); err != nil {
		log.Printf("panic recovered: %v", err)
	}
}

func (handler *Handler) getLastMessages(c *fiber.Ctx) error {
	userId, exists := c.Locals("userId").(uint)
	if !exists {
		return NewErrorResponse(c, fiber.StatusBadRequest, "userId not found")
	}

	chatsId := c.Query("chats_id")
	messages, err := handler.services.Message.GetLastMessages(c.Context(), userId, chatsId)
	if err != nil {
		return NewErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"messages": messages,
	})
}

func (handler *Handler) getMessages(c *fiber.Ctx) error {
	chatId, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil || chatId < 1 {
		return NewErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	offset, err := strconv.ParseUint(c.Query("offset"), 10, 64)
	if err != nil || offset < 0 {
		return NewErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	messages, err := handler.services.Message.GetMessages(c.Context(), uint(chatId), uint(offset))
	if err != nil {
		return NewErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"messages": messages,
	})
}
