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

func (service *ChatService) CreateChat(ctx context.Context, userId uint, chatDTO Entities.ChatDTO) (uint, error) {
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

func (service *ChatService) UpdateChat(ctx context.Context, userId uint, chatId uint, chatDTO Entities.ChatUpdateDTO) error {
	return service.repository.UpdateChat(ctx, userId, chatId, chatDTO)
}

func (service *ChatService) DeleteChat(ctx context.Context, userId uint, chatId uint) ([]uint, error) {
	return service.repository.DeleteChat(ctx, userId, chatId)
}

func (service *ChatService) GetChatsListAndInfoForUser(ctx context.Context, userId uint) (friends []uint, chatLists string, err error) {
	return service.repository.GetChatsListAndInfoForUser(ctx, userId)
}

func (service *ChatService) GetChats(ctx context.Context, userId uint, chatsId string) ([]Entities.Chat, error) {
	return service.repository.GetChats(ctx, userId, chatsId)
}
func (service *ChatService) UpdateChatLists(ctx context.Context, id uint, newChatLists string) error {
	return service.repository.UpdateChatLists(ctx, id, newChatLists)
}
