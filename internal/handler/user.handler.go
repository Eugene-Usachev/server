package handler

import (
	"GoServer/Entities"
	"GoServer/pkg/fasthttp_utils"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func (handler *Handler) getUser(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil || id < 1 {
		return NewErrorResponse(c, fiber.StatusBadRequest, "nothing to get")
	}

	user, err := handler.services.User.GetUserById(c.Context(), uint(id))
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return NewErrorResponse(c, fiber.StatusNotFound, "user is not exist")
		}
		handler.Logger.Error("get user error: " + err.Error())
		return NewErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(user)
}

func (handler *Handler) getUserSubsIds(c *fiber.Ctx) error {
	id := c.Locals("userId")
	if id == nil {
		return NewErrorResponse(c, fiber.StatusBadRequest, "no id")
	}
	if id.(uint) < 1 {
		return NewErrorResponse(c, fiber.StatusBadRequest, "no id")
	}

	ids, err := handler.services.User.GetUserSubsIds(c.Context(), id.(uint))
	if err != nil {
		handler.Logger.Error("get user subs ids error: " + err.Error())
		return NewErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(ids)
}

func (handler *Handler) getFriendsAndSubs(c *fiber.Ctx) error {
	client := c.Query("userId")
	clientUint, err := strconv.ParseUint(client, 10, 64)
	if err != nil || clientUint < 1 {
		clientUint = 0
	}
	userId, err := strconv.ParseUint(c.Params("userId"), 10, 64)
	if err != nil || userId < 1 {
		return NewErrorResponse(c, fiber.StatusBadRequest, "nothing to get")
	}

	friendsAndSubs, err := handler.services.User.GetFriendsAndSubs(c.Context(), uint(clientUint), uint(userId))
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return NewErrorResponse(c, fiber.StatusNotFound, "user is not exist")
		}
		handler.Logger.Error("get friends and subs error: " + err.Error())
		return NewErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"user":   friendsAndSubs.User,
		"client": friendsAndSubs.Client,
	})
}

func (handler *Handler) getUsersForFriendPage(c *fiber.Ctx) error {
	idOfUsers := c.Query("idOfUsers")

	users, err := handler.services.User.GetUsersForFriendsPage(c.Context(), idOfUsers)
	if err != nil {
		handler.Logger.Error("get users for friend page error: " + err.Error())
		return NewErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{"users": users})
}

func (handler *Handler) getUsers(c *fiber.Ctx) error {
	idOfUsers := c.Query("idOfUsers")

	users, err := handler.services.User.GetUsers(c.Context(), idOfUsers)
	if err != nil {
		handler.Logger.Error("get users error: " + err.Error())
		return NewErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{"users": users})
}

func (handler *Handler) updateUser(c *fiber.Ctx) error {
	userId := c.Locals("userId")
	if userId.(uint) < 1 {
		return NewErrorResponse(c, fiber.StatusBadRequest, "impossible to get user id")
	}

	var input Entities.UpdateUserDTO
	if err := fasthttp_utils.JSON(c, &input); err != nil {
		return NewErrorResponse(c, fiber.StatusBadRequest, "invalid request body")
	}

	err := handler.services.User.UpdateUser(c.Context(), userId.(uint), input)
	if err != nil {
		handler.Logger.Error("update user error: " + err.Error())
		return NewErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (handler *Handler) changeAvatar(ctx *fiber.Ctx) error {
	userId := ctx.Locals("userId")
	if userId.(uint) < 1 {
		return NewErrorResponse(ctx, fiber.StatusBadRequest, "impossible to get user id")
	}

	fileName, err := handler.services.User.ChangeAvatar(ctx, userId.(uint))
	if err != nil {
		handler.Logger.Error("change avatar error: " + err.Error())
		return NewErrorResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"fileName": fileName,
	})
}

func (handler *Handler) addToFriends(c *fiber.Ctx) error {
	userId := c.Locals("userId")
	if userId.(uint) < 1 {
		return NewErrorResponse(c, fiber.StatusUnauthorized, "impossible to get user id")
	}

	bodyId, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil || bodyId < 1 {
		return NewErrorResponse(c, fiber.StatusBadRequest, "invalid body id")
	}

	err = handler.services.User.AddToFriends(c.Context(), userId.(uint), uint(bodyId))
	if err != nil {
		handler.Logger.Error("add to friends error: " + err.Error())
		return NewErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (handler *Handler) deleteFromFriends(c *fiber.Ctx) error {
	userId := c.Locals("userId")
	if userId.(uint) < 1 {
		return NewErrorResponse(c, fiber.StatusBadRequest, "impossible to get user id")
	}
	bodyId, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil || bodyId < 1 {
		return NewErrorResponse(c, fiber.StatusBadRequest, "invalid body id")
	}

	err = handler.services.User.DeleteFromFriends(c.Context(), userId.(uint), uint(bodyId))
	if err != nil {
		handler.Logger.Error("server delete from friends error: " + err.Error())
		return NewErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (handler *Handler) addToSubs(c *fiber.Ctx) error {
	userId := c.Locals("userId")
	if userId.(uint) < 1 {
		return NewErrorResponse(c, fiber.StatusBadRequest, "impossible to get user id")
	}
	bodyId, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil || bodyId < 1 {
		return NewErrorResponse(c, fiber.StatusBadRequest, "invalid body id")
	}

	err = handler.services.User.AddToSubs(c.Context(), userId.(uint), uint(bodyId))
	if err != nil {
		handler.Logger.Error("add to subs error: " + err.Error())
		return NewErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (handler *Handler) deleteFromSubs(c *fiber.Ctx) error {
	userId := c.Locals("userId")
	if userId.(uint) < 1 {
		return NewErrorResponse(c, fiber.StatusBadRequest, "impossible to get user id")
	}
	bodyId, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil || bodyId < 1 {
		return NewErrorResponse(c, fiber.StatusBadRequest, "invalid body id")
	}

	err = handler.services.User.DeleteFromSubs(c.Context(), userId.(uint), uint(bodyId))
	if err != nil {
		handler.Logger.Error("server delete from subs error: " + err.Error())
		return NewErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (handler *Handler) deleteUser(c *fiber.Ctx) error {
	userId := c.Locals("userId")
	if userId.(uint) < 1 {
		return NewErrorResponse(c, fiber.StatusBadRequest, "impossible to get user id")
	}
	err := handler.services.User.DeleteUser(c.Context(), userId.(uint))
	if err != nil {
		handler.Logger.Error("server delete user error: " + err.Error())
		return NewErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}
	return c.SendStatus(fiber.StatusNoContent)
}
