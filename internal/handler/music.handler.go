package handler

import (
	"GoServer/Entities"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func (handler *Handler) getMusics(c *fiber.Ctx) error {
	name := c.Query("name")
	offset, err := strconv.ParseUint(c.Query("offset"), 10, 64)
	if err != nil {
		return NewErrorResponse(c, fiber.StatusBadRequest, "invalid offset")
	}

	musics, err := handler.services.Music.GetMusics(c.Context(), name, uint(offset))
	if err != nil {
		return NewErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"data": musics,
	})
}

func (handler *Handler) getMusic(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return NewErrorResponse(c, fiber.StatusBadRequest, "invalid id")
	}
	pathToMusic, contentType, err := handler.services.Music.GetMusic(c.Context(), uint(id))
	if err != nil {
		return NewErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}
	c.Set("Content-Type", contentType)
	return c.SendFile(pathToMusic)
}

func (handler *Handler) addMusic(c *fiber.Ctx) error {
	id, exist := c.Locals("userId").(uint)
	if !exist || id < 1 {
		return NewErrorResponse(c, fiber.StatusBadRequest, "no authorized user")
	}
	var input Entities.CreateMusicDTO
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("invalid body request")
	}
	err := handler.services.Music.AddMusic(c, id, input)
	if err != nil {
		return NewErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}
	return c.SendStatus(fiber.StatusCreated)
}
