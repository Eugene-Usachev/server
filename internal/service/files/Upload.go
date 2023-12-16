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
	"sync"
)

const (
	MB               = 2 << 20
	KB               = 2 << 10
	UserStorageSize  = MB * 256
	FileMaxSize      = MB * 10
	PostFilesMaxSize = FileMaxSize * 10
)

var (
	UserStorageLoadNow = func() uint64 {
		var size uint64 = 0
		filepath.Walk("../static/UserFiles", func(_ string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				size += uint64(info.Size())
			}
			return nil
		})
		return size
	}()
	UserStorageLoadNowMutex sync.Mutex
)

var (
	ErrFileTooLarge           = errors.New("file too large")
	ErrImpossibleToCreateDir  = errors.New("impossible to create directory")
	ErrFileFormatNotSupported = errors.New("file format is not supported")
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
	if file.Size > FileMaxSize {
		return "", ErrFileTooLarge
	}
	UserStorageLoadNowMutex.Lock()
	if UserStorageLoadNow+uint64(file.Size) > UserStorageSize {
		return "", ErrFileTooLarge
	}
	UserStorageLoadNow += uint64(file.Size)
	UserStorageLoadNowMutex.Unlock()

	err := os.MkdirAll(path, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return "", ErrImpossibleToCreateDir
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

func UploadPostFiles(ctx *fiber.Ctx, files []*multipart.FileHeader, path string) ([]string, error) {
	filesSize := uint64(0)
	for _, file := range files {
		filesSize += uint64(file.Size)
	}

	if filesSize > FileMaxSize {
		return []string{}, ErrFileTooLarge
	}
	UserStorageLoadNowMutex.Lock()
	if UserStorageLoadNow+filesSize > UserStorageSize {
		return []string{}, ErrFileTooLarge
	}
	UserStorageLoadNow += filesSize
	UserStorageLoadNowMutex.Unlock()

	err := os.MkdirAll(path, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return []string{}, ErrImpossibleToCreateDir
	}

	names := make([]string, 0, len(files))
	for _, file := range files {
		fileExt := file.Filename[strings.LastIndex(file.Filename, ".")+1:]
		var folder string
		switch strings.ToLower(fileExt) {
		case ".jpg", ".jpeg", ".png", ".gif":
			folder = "Image"
		case ".mp3", ".wav":
			folder = "Music"
		case ".mp4", ".avi", ".mkv":
			folder = "Video"
		default:
			folder = "Other"
		}

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
		fullPath := path + folder + "/" + file.Filename
		err = ctx.SaveFile(file, fullPath)
		if err != nil {
			return []string{}, err
		}
		names = append(names, fullPath)
	}

	return names, nil
}

func UploadImage(c *fiber.Ctx, path string) (string, error) {
	userAvatar, err := c.FormFile("avatar")
	if err != nil || userAvatar == nil {
		return "", err
	}
	if userAvatar.Size > FileMaxSize {
		return "", ErrFileTooLarge
	}

	UserStorageLoadNowMutex.Lock()
	if UserStorageLoadNow+uint64(userAvatar.Size) > UserStorageSize {
		return "", ErrFileTooLarge
	}
	UserStorageLoadNow += uint64(userAvatar.Size)
	UserStorageLoadNowMutex.Unlock()

	fileExt := userAvatar.Filename[strings.LastIndex(userAvatar.Filename, ".")+1:]
	if fileExt != "png" && fileExt != "jpg" && fileExt != "jpeg" && fileExt != "gif" && fileExt != "webp" {
		return "", ErrFileFormatNotSupported
	}

	err = os.MkdirAll(path, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return "", ErrImpossibleToCreateDir
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
