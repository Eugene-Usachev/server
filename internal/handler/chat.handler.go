package handler

import (
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
)

func (handler *Handler) getChatsList(ctx *fiber.Ctx) error {
	userId, ok := ctx.Locals("userId").(uint)
	if !ok || userId == 0 {
		return NewErrorResponse(ctx, fiber.StatusUnauthorized, "invalid auth token")
	}

	friends, chatLists, rawChats, err := handler.services.Chat.GetChatsListAndInfoForUser(ctx.Context(), userId)
	if err != nil {
		return NewErrorResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"chatsList": chatLists,
		"friends":   friends,
		"raw_chats": rawChats,
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

type UpdateChatsListsDTO struct {
	IsSetRawChatsToEmpty bool   `json:"is_set_raw_chats_to_empty"`
	NewChatsLists        string `json:"new_chats_lists"`
}

func (handler *Handler) UpdateChatsLists(ctx *fiber.Ctx) error {
	userId, ok := ctx.Locals("userId").(uint)
	if !ok || userId == 0 {
		return NewErrorResponse(ctx, fiber.StatusUnauthorized, "invalid auth token")
	}

	body := ctx.Body()
	var dto UpdateChatsListsDTO

	err := json.Unmarshal(body, &dto)
	if err != nil {
		return NewErrorResponse(ctx, fiber.StatusBadRequest, "bad request")
	}
	err = handler.services.Chat.UpdateChatLists(ctx.Context(), userId, dto.NewChatsLists, dto.IsSetRawChatsToEmpty)
	if err != nil {
		return NewErrorResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}
