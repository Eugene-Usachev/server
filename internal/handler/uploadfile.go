package handler

import (
	"GoServer/internal/service/files"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"path/filepath"
	"strings"
)

func (handler *Handler) uploadFile(ctx *fiber.Ctx) error {
	userId := ctx.Locals("userId")
	if userId.(uint) < 1 {
		return NewErrorResponse(ctx, fiber.StatusBadRequest, "userId is required")
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		return NewErrorResponse(ctx, fiber.StatusInternalServerError, err.Error())
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
		return NewErrorResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}

	return ctx.Status(fiber.StatusOK).SendString(name)
}
