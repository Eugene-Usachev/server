package handler

import (
	"github.com/Eugene-Usachev/fastbytes"
	"github.com/gofiber/fiber/v2"
)

func (handler *Handler) getChatsList(ctx *fiber.Ctx) error {
	userId, ok := ctx.Locals("userId").(uint)
	if !ok || userId == 0 {
		return NewErrorResponse(ctx, fiber.StatusUnauthorized, "invalid auth token")
	}

	friends, chatLists, err := handler.services.Chat.GetChatsListAndInfoForUser(ctx.Context(), userId)
	if err != nil {
		return NewErrorResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"chatsList": chatLists,
		"friends":   friends,
	})
}

func (handler *Handler) getChats(ctx *fiber.Ctx) error {
	userId, ok := ctx.Locals("userId").(uint)
	if !ok || userId == 0 {
		return NewErrorResponse(ctx, fiber.StatusUnauthorized, "invalid auth token")
	}

	chatsIds := ctx.Query("chatsIds")
	if chatsIds == "" {
		return NewErrorResponse(ctx, fiber.StatusBadRequest, "bad request")
	}

	chats, err := handler.services.Chat.GetChats(ctx.Context(), userId, chatsIds)
	if err != nil {
		return NewErrorResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(chats)
}

func (handler *Handler) UpdateChatLists(ctx *fiber.Ctx) error {
	userId, ok := ctx.Locals("userId").(uint)
	if !ok || userId == 0 {
		return NewErrorResponse(ctx, fiber.StatusUnauthorized, "invalid auth token")
	}

	dto := fastbytes.B2S(ctx.Body())
	if dto == "" {
		return NewErrorResponse(ctx, fiber.StatusBadRequest, "bad request")
	}
	err := handler.services.Chat.UpdateChatLists(ctx.Context(), userId, dto)
	if err != nil {
		return NewErrorResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}
