package service

import (
	"GoServer/Entities"
	"GoServer/internal/repository"
	"context"
	"errors"
)

type ChatService struct {
	repository repository.Chat
}

func NewChatService(repository repository.Chat) *ChatService {
	return &ChatService{
		repository: repository,
	}
}

func (service *ChatService) CreateChat(ctx context.Context, userId int64, chatDTO Entities.ChatDTO) (int64, error) {
	if len(chatDTO.Members) < 2 {
		return 0, errors.New("members must be more than 1")
	}
	var isUserIdInMembers bool
	for _, member := range chatDTO.Members {
		if member == userId {
			isUserIdInMembers = true
			break
		}
	}
	if !isUserIdInMembers {
		return 0, errors.New("user not in members")
	}

	return service.repository.CreateChat(ctx, chatDTO)
}

func (service *ChatService) UpdateChat(ctx context.Context, userId int64, chatId int64, chatDTO Entities.ChatUpdateDTO) error {
	return service.repository.UpdateChat(ctx, userId, chatId, chatDTO)
}

func (service *ChatService) DeleteChat(ctx context.Context, userId int64, chatId int64) ([]int64, error) {
	return service.repository.DeleteChat(ctx, userId, chatId)
}

func (service *ChatService) GetChats(ctx context.Context, userId int64) (string, string, string, []int64, []int64, string, []Entities.Chat, error) {
	return service.repository.GetChats(ctx, userId)
}

func (service *ChatService) UpdateChatLists(ctx context.Context, id int64, newChatLists string) error {
	return service.repository.UpdateChatLists(ctx, id, newChatLists)
}
