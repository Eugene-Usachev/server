package handler

import (
	"GoServer/Entities"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (handler *Handler) getUser(ctx *gin.Context) {

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil || id < 1 {
		NewErrorResponse(ctx, http.StatusBadRequest, "nothing to get")
		return
	}

	requestOwner, err := strconv.ParseUint(ctx.Query("requestOwner"), 10, 64)
	if err != nil || requestOwner < 1 {
		requestOwner = 0
	}

	user, subs, err := handler.services.GetUserById(ctx.Request.Context(), uint(id), uint(requestOwner))
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			NewErrorResponse(ctx, http.StatusNotFound, "user is not exist")
			return
		}
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"user": user, "subs": subs})
}

func (handler *Handler) getFriendsAndSubs(ctx *gin.Context) {

	client := ctx.Query("userId")
	clientUint, err := strconv.ParseUint(client, 10, 64)
	if err != nil || clientUint < 1 {
		clientUint = 0
	}
	userId, err := strconv.ParseUint(ctx.Param("userId"), 10, 64)
	if err != nil || userId < 1 {
		NewErrorResponse(ctx, http.StatusBadRequest, "nothing to get")
		return
	}

	friendsAndSubs, err := handler.services.GetFriendsAndSubs(ctx.Request.Context(), uint(clientUint), uint(userId))
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			NewErrorResponse(ctx, http.StatusNotFound, "user is not exist")
			return
		}
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user":   friendsAndSubs.User,
		"client": friendsAndSubs.Client,
	})
}

func (handler *Handler) getUsersForFriendPage(ctx *gin.Context) {
	idOfUsers := ctx.Query("idOfUsers")

	users, err := handler.services.GetUsersForFriendsPage(ctx.Request.Context(), idOfUsers)
	if err != nil {
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	ctx.JSON(http.StatusOK, gin.H{"users": users})
}

func (handler *Handler) getUsers(ctx *gin.Context) {
	idOfUsers := ctx.Query("idOfUsers")

	users, err := handler.services.GetUsers(ctx.Request.Context(), idOfUsers)
	if err != nil {
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	ctx.JSON(http.StatusOK, gin.H{"users": users})
}

func (handler *Handler) updateUser(ctx *gin.Context) {
	userId, isExist := ctx.Get("userId")
	if !isExist || userId.(uint) < 1 {
		NewErrorResponse(ctx, http.StatusBadRequest, "impossible to get user id")
		return
	}

	var input Entities.UpdateUserDTO
	if err := ctx.BindJSON(&input); err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, "invalid request body")
		return
	}

	err := handler.services.UpdateUser(ctx.Request.Context(), userId.(uint), input)
	if err != nil {
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (handler *Handler) changeAvatar(ctx *gin.Context) {
	userId, isExist := ctx.Get("userId")
	if !isExist || userId.(uint) < 1 {
		NewErrorResponse(ctx, http.StatusBadRequest, "impossible to get user id")
		return
	}

	fileName, err := handler.services.ChangeAvatar(ctx, userId.(uint))
	if err != nil {
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"fileName": fileName,
	})
}

func (handler *Handler) addToFriends(ctx *gin.Context) {
	userId, isExist := ctx.Get("userId")
	if !isExist || userId.(uint) < 1 {
		NewErrorResponse(ctx, http.StatusUnauthorized, "impossible to get user id")
		return
	}

	bodyId, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil || bodyId < 1 {
		NewErrorResponse(ctx, http.StatusBadRequest, "invalid body id")
	}

	err = handler.services.AddToFriends(ctx.Request.Context(), userId.(uint), uint(bodyId))
	if err != nil {
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (handler *Handler) deleteFromFriends(ctx *gin.Context) {
	userId, isExist := ctx.Get("userId")
	if !isExist || userId.(uint) < 1 {
		NewErrorResponse(ctx, http.StatusBadRequest, "impossible to get user id")
		return
	}
	bodyId, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil || bodyId < 1 {
		NewErrorResponse(ctx, http.StatusBadRequest, "invalid body id")
	}

	err = handler.services.DeleteFromFriends(ctx.Request.Context(), userId.(uint), uint(bodyId))
	if err != nil {
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (handler *Handler) addToSubs(ctx *gin.Context) {
	userId, isExist := ctx.Get("userId")
	if !isExist || userId.(uint) < 1 {
		NewErrorResponse(ctx, http.StatusBadRequest, "impossible to get user id")
		return
	}
	bodyId, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil || bodyId < 1 {
		NewErrorResponse(ctx, http.StatusBadRequest, "invalid body id")
	}

	err = handler.services.AddToSubs(ctx.Request.Context(), userId.(uint), uint(bodyId))
	if err != nil {
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (handler *Handler) deleteFromSubs(ctx *gin.Context) {
	userId, isExist := ctx.Get("userId")
	if !isExist || userId.(uint) < 1 {
		NewErrorResponse(ctx, http.StatusBadRequest, "impossible to get user id")
		return
	}
	bodyId, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil || bodyId < 1 {
		NewErrorResponse(ctx, http.StatusBadRequest, "invalid body id")
	}

	err = handler.services.DeleteFromSubs(ctx.Request.Context(), userId.(uint), uint(bodyId))
	if err != nil {
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (handler *Handler) deleteUser(ctx *gin.Context) {
	userId, isExist := ctx.Get("userId")
	if !isExist || userId.(uint) < 1 {
		NewErrorResponse(ctx, http.StatusBadRequest, "impossible to get user id")
		return
	}
	err := handler.services.DeleteUser(ctx.Request.Context(), userId.(uint))
	if err != nil {
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusNoContent, gin.H{})
}
