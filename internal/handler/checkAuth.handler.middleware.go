package handler

import (
	"GoServer/pkg/fasthttp_utils"
	"github.com/Eugene-Usachev/fastbytes"
	"github.com/gofiber/fiber/v2"
)

func (handler *Handler) CheckAuth(c *fiber.Ctx) error {
	accessToken := fastbytes.B2S(fasthttp_utils.GetAuthorizationHeader(c.Context()))
	if accessToken == "" {
		return c.SendStatus(fiber.StatusUnauthorized)
	} else {
		userId, err := handler.accessConverter.ParseToken(accessToken)
		if err != nil {
			return c.SendStatus(fiber.StatusUnauthorized)
		} else {
			id := fastbytes.B2U(userId)
			c.Locals("userId", id)
		}
	}
	return c.Next()
}
