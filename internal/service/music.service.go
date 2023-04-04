package service

import (
	"GoServer/Entities"
	"GoServer/internal/repository"
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"strings"
)

type MusicService struct {
	repository repository.Music
}

func NewMusicService(repository repository.Music) *MusicService {
	return &MusicService{
		repository: repository,
	}
}

func (service *MusicService) GetMusics(ctx context.Context, name string, offset uint) ([]Entities.Music, error) {
	return service.repository.GetMusics(ctx, name, offset)
}

func (service *MusicService) GetMusic(ctx context.Context, id uint) (string, string, error) {
	parentUserId, title, err := service.repository.GetMusic(ctx, id)
	if err != nil {
		return "", "", err
	}

	fileExt := title[strings.LastIndex(title, ".")+1:]

	return fmt.Sprintf("./static/UserFiles/%d/Music/%d.%s", parentUserId, id, fileExt), fileExt, nil
}

func (service *MusicService) AddMusic(ctx *gin.Context, id uint, music Entities.CreateMusicDTO) error {

	file, err := ctx.FormFile("file")
	if err != nil || file == nil {
		return err
	}
	if file.Size > 1024*10*1024 {
		return errors.New("file too large")
	}

	path := fmt.Sprintf("./static/UserFiles/%d/Music/", id)

	err = os.MkdirAll(path, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return errors.New("impossible to create directory")
	}
	fileExt := file.Filename[strings.LastIndex(file.Filename, ".")+1:]

	music.Title = music.Title + "." + fileExt

	musicId, err := service.repository.AddMusic(ctx.Request.Context(), id, music)
	if err != nil {
		return err
	}

	file.Filename = fmt.Sprintf("%d.%s", musicId, fileExt)

	err = ctx.SaveUploadedFile(file, path+file.Filename)
	if err != nil {
		return err
	}

	return nil
}
