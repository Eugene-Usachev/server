package handler

import (
	"GoServer/Entities"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func (handler *Handler) getMusics(ctx *gin.Context) {
	name := ctx.Query("name")
	offset, err := strconv.ParseUint(ctx.Query("offset"), 10, 64)
	if err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, "invalid offset")
		return
	}

	musics, err := handler.services.GetMusics(ctx.Request.Context(), name, uint(offset))
	if err != nil {
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": musics,
	})
}

func (handler *Handler) getMusic(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)

	if err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, "invalid id")
	}

	pathToMusic, contentType, err := handler.services.GetMusic(ctx.Request.Context(), uint(id))
	if err != nil {
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	log.Println(pathToMusic)
	ctx.Header("Content-Type", contentType)
	ctx.Status(http.StatusOK)
	ctx.File(pathToMusic)
}

func (handler *Handler) addMusic(ctx *gin.Context) {
	id, exist := ctx.Get("userId")
	if !exist || id.(uint) < 1 {
		NewErrorResponse(ctx, http.StatusBadRequest, "no authorized user")
		return
	}
	var input Entities.CreateMusicDTO
	if err := ctx.ShouldBind(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, "invalid body request")
		return
	}
	err := handler.services.AddMusic(ctx, id.(uint), input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.String(http.StatusCreated, "created")
}
