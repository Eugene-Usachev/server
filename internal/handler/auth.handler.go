package handler

import (
	. "GoServer/Entities"
	"GoServer/internal/repository"
	utils "GoServer/pkg/fasthttp_utils"
	"github.com/gofiber/fiber/v2"
	"log"
	"strconv"
)

func (handler *Handler) signUp(c *fiber.Ctx) error {
	var input UserDTO
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	err := checkPassword(input.Password)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	id, err, Tokens := handler.services.Authorization.CreateUser(c.Context(), input)
	if err != nil {
		// TODO r
		log.Println(err)
		switch err {
		case repository.EmailBusy:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		case repository.LoginBusy:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id":            id,
		"access_token":  Tokens.AccessToken,
		"refresh_token": Tokens.RefreshToken,
	})
}

func (handler *Handler) signIn(c *fiber.Ctx) error {
	var input SignInDTO

	if err := c.BodyParser(&input); err != nil || (input.Email == "" && input.Login == "") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	user, tokens, err := handler.services.Authorization.SignIn(c.Context(), input)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"id":            user.ID,
		"email":         user.Email,
		"login":         user.Login,
		"avatar":        user.Avatar,
		"name":          user.Name,
		"surname":       user.Surname,
		"access_token":  tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
	})
}

func (handler *Handler) refresh(c *fiber.Ctx) error {
	var input RefreshDTO

	err := utils.JSON(c, &input)
	if err != nil {
		return err
	}

	if input.Id < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Id is empty",
		})
	}

	dto, err := handler.services.Authorization.Refresh(c.Context(), input.Id, input.Token)
	if err != nil {
		// TODO r
		log.Println("refresh error", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"access_token":  dto.AccessToken,
		"refresh_token": dto.RefreshToken,
		"avatar":        dto.Avatar,
		"name":          dto.Name,
		"surname":       dto.Surname,
	})
}

func (handler *Handler) refreshTokens(c *fiber.Ctx) error {
	refreshToken := c.Query("token")
	id := c.Params("id")
	uid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	var tokens AllTokenResponse
	tokens, err = handler.services.Authorization.RefreshTokens(c.Context(), uint(uid), refreshToken)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"access_token":  tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
	})
}
