package service

import (
	"GoServer/Entities"
	"GoServer/internal/repository"
	"context"
)

type MessageService struct {
	repository repository.Message
}

func NewMessageService(repository repository.Message) *MessageService {
	return &MessageService{
		repository: repository,
	}
}

func (service *MessageService) SaveMessage(ctx context.Context, userId int64, messageDTO Entities.MessageDTO) (int64, []int64, string, error) {
	messageDTO.Date = NewDate()
	id, members, err := service.repository.SaveMessage(ctx, userId, messageDTO)
	return id, members, messageDTO.Date, err
}
func (service *MessageService) UpdateMessage(ctx context.Context, messageId int64, userId int64, newData string) ([]int64, error) {
	return service.repository.UpdateMessage(ctx, messageId, userId, newData)
}
func (service *MessageService) DeleteMessage(ctx context.Context, messageId int64, userId int64) ([]int64, error) {
	return service.repository.DeleteMessage(ctx, messageId, userId)
}

func (service *MessageService) GetLastMessages(ctx context.Context, userId uint, chatsId string) ([]Entities.Message, error) {
	return service.repository.GetLastMessages(ctx, userId, chatsId)
}

func (service *MessageService) GetMessages(ctx context.Context, chatId, offset uint) ([20]Entities.Message, error) {
	return service.repository.GetMessages(ctx, chatId, offset)
}
