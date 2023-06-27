package fasthttp_utils

import (
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
)

func JSON(c *fiber.Ctx, output interface{}) error {
	body := c.Request().Body()

	if len(body) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Body is empty",
		})
	}

	if err := json.Unmarshal(body, &output); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return nil
}
