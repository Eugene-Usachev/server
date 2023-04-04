package handler

import (
	"GoServer/internal/service/files"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"strings"
)

func (handler *Handler) uploadFile(ctx *gin.Context) {
	userId, exists := ctx.Get("userId")
	if !exists || userId.(uint) < 1 {
		NewErrorResponse(ctx, http.StatusBadRequest, "userId is required")
		return
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	file := form.File["file"][0]
	ext := filepath.Ext(file.Filename)
	var folder string

	switch strings.ToLower(ext) {
	case ".jpg", ".jpeg", ".png", ".gif":
		folder = "Image"
	case ".mp3", ".wav":
		folder = "Music"
	case ".mp4", ".avi", ".mkv":
		folder = "Video"
	default:
		folder = "Other"
	}

	var name string //todo music
	name, err = files.UploadFile(ctx, file, fmt.Sprintf("./static/UserFiles/%d/%s/", userId.(uint), folder))
	if err != nil {
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.String(http.StatusOK, name)
}
