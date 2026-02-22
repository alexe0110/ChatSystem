package service

import (
	"context"
	"github.com/alexe0110/chat-system/internal/model"
	"github.com/alexe0110/chat-system/internal/repository"
	"github.com/google/uuid"
)

type MessageService struct {
	repo repository.MessageRepository
}

func NewMessageService(repo repository.MessageRepository) *MessageService {
	return &MessageService{
		repo: repo,
	}
}

func (r MessageService) SendMessage(ctx context.Context, senderID, receiverID uuid.UUID, messageContent string) (*model.Message, error) {
	message, err := r.repo.SendMessage(ctx, senderID, receiverID, messageContent)
	if err != nil {
		return nil, err
	}

	return message, nil
}

func (r *MessageService) GetConversation(ctx context.Context, senderID, receiverID uuid.UUID) ([]*model.Message, error) {
	message, err := r.repo.GetConversation(ctx, senderID, receiverID)
	if err != nil {
		return nil, err
	}

	return message, nil
}

func (r *MessageService) GetMessageByID(ctx context.Context, ID uuid.UUID) (*model.Message, error) {
	message, err := r.repo.GetMessageByID(ctx, ID)
	if err != nil {
		return nil, err
	}

	return message, nil
}
