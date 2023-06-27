package files

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"mime"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

const (
	MB               = 2 << 20
	KB               = 2 << 10
	FileMaxSize      = MB * 10
	PostFilesMaxSize = FileMaxSize * 10
)

func IsFileMusic(file *multipart.FileHeader) bool {

	fileType := mime.TypeByExtension(filepath.Ext(file.Filename))

	if fileType == "audio/mpeg" || fileType == "audio/ogg" || fileType == "audio/wav" {
		return true
	} else {
		return false
	}
}

func UploadFile(ctx *fiber.Ctx, file *multipart.FileHeader, path string) (string, error) {
	return "", nil
	if file.Size > FileMaxSize {
		return "", errors.New("file too large")
	}

	err := os.MkdirAll(path, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return "", errors.New("impossible to create directory")
	}
	fileExt := file.Filename[strings.LastIndex(file.Filename, ".")+1:]

	for i := 0; ; i++ {
		_, err = os.Stat(path + file.Filename)
		if err == nil && !os.IsNotExist(err) {
			var index int
			if i == 0 {
				index = strings.LastIndex(file.Filename, ".")
			} else {
				index = strings.LastIndex(file.Filename, ".") - 1
			}
			file.Filename = fmt.Sprintf(file.Filename[:index]+"%d.%s", i, fileExt)
		} else {
			break
		}
	}

	err = ctx.SaveFile(file, path+file.Filename)
	if err != nil {
		return "", err
	}
	return file.Filename, nil
}

func UploadPostFiles(c *fiber.Ctx) ([]string, error) {
	return []string{""}, nil
	//form, err := c.MultipartForm()
	return nil, nil
}

func UploadImage(c *fiber.Ctx, path string) (string, error) {
	return "", nil
	userAvatar, err := c.FormFile("avatar")
	if err != nil || userAvatar == nil {
		return "", err
	}
	if userAvatar.Size > FileMaxSize {
		return "", errors.New("file too large")
	}

	fileExt := userAvatar.Filename[strings.LastIndex(userAvatar.Filename, ".")+1:]
	if fileExt != "png" && fileExt != "jpg" && fileExt != "jpeg" && fileExt != "gif" && fileExt != "webp" {
		return "", errors.New("file format not supported")
	}

	err = os.MkdirAll(path, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return "", errors.New("impossible to create image directory")
	}

	for i := 0; ; i++ {
		_, err = os.Stat(path + userAvatar.Filename)
		if err == nil && !os.IsNotExist(err) {
			var index int
			if i == 0 {
				index = strings.LastIndex(userAvatar.Filename, ".")
			} else {
				index = strings.LastIndex(userAvatar.Filename, ".") - 1
			}
			userAvatar.Filename = fmt.Sprintf(userAvatar.Filename[:index]+"%d.%s", i, fileExt)
		} else {
			break
		}
	}

	err = c.SaveFile(userAvatar, path+userAvatar.Filename)
	err = c.SaveFile(userAvatar, path+userAvatar.Filename)
	if err != nil {
		return "", err
	}
	return userAvatar.Filename, nil
}
