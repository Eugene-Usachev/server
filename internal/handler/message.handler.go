package handler

import (
	"github.com/gofiber/fiber/v2"
	"strconv"
)

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

	return c.Status(fiber.StatusOK).JSON(messages)
}
