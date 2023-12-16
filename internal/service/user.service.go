package service

import (
	"GoServer/Entities"
	"GoServer/internal/repository"
	"GoServer/internal/service/files"
	"context"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

type UserService struct {
	repository repository.User
}

func NewUserService(repository repository.User) *UserService {
	return &UserService{
		repository: repository,
	}
}

func (service *UserService) GetUserById(ctx context.Context, id uint) (Entities.GetUserDTO, error) {
	return service.repository.GetUserById(ctx, id)
}

func (service *UserService) GetUserSubsIds(ctx context.Context, id uint) ([]uint, error) {
	return service.repository.GetUserSubsIds(ctx, id)
}

func (service *UserService) GetFriendsAndSubs(ctx context.Context, clientId, userId uint) (Entities.GetFriendsAndSubsDTO, error) {
	return service.repository.GetFriendsAndSubs(ctx, clientId, userId)
}

func (service *UserService) UpdateUser(ctx context.Context, id uint, UpdateUserDTO Entities.UpdateUserDTO) error {
	return service.repository.UpdateUser(ctx, id, UpdateUserDTO)
}

func (service *UserService) ChangeAvatar(ctx *fiber.Ctx, userId uint) (string, error) {
	path := fmt.Sprintf("../static/UserFiles/%d/Image/", userId)

	fileName, err := files.UploadImage(ctx, path)
	if err != nil {
		return "", err
	}

	err = service.repository.ChangeAvatar(ctx.Context(), userId, fileName)
	if err != nil {
		return "", errors.New("impossible to change avatar")
	}

	return fileName, nil
}

func (service *UserService) AddToFriends(ctx context.Context, id, body uint) error {
	return service.repository.AddToFriends(ctx, id, body)
}

func (service *UserService) DeleteFromFriends(ctx context.Context, id, body uint) error {
	return service.repository.DeleteFromFriends(ctx, id, body)
}

func (service *UserService) AddToSubs(ctx context.Context, id, body uint) error {
	return service.repository.AddToSubs(ctx, id, body)
}

func (service *UserService) DeleteFromSubs(ctx context.Context, id, body uint) error {
	return service.repository.DeleteFromSubs(ctx, id, body)
}

func (service *UserService) DeleteUser(ctx context.Context, id uint) error {
	return service.repository.DeleteUser(ctx, id)
}

func (service *UserService) GetUsers(ctx context.Context, idOfUsers string) ([]Entities.MiniUser, error) {
	return service.repository.GetUsers(ctx, idOfUsers)
}

func (service *UserService) GetUsersForFriendsPage(ctx context.Context, idOfUsers string, clientId uint) ([]Entities.FriendUser, error) {
	return service.repository.GetUsersForFriendsPage(ctx, idOfUsers, clientId)
}

func (service *UserService) GetOnlineUsers(ctx context.Context, slice []string) ([]int, error) {
	// TODO
	panic("implement me")
	return nil, nil
}

func (service *UserService) SubscribeOnUsers(ctx context.Context, slice []string, clientId string) error {
	// TODO
	panic("implement me")
	return nil
}
