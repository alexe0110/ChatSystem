package repository

import (
	"context"
	"github.com/alexe0110/chat-system/internal/model"
	"github.com/google/uuid"
)

type MessageRepository interface {
	SendMessage(ctx context.Context, senderID, receiverID uuid.UUID, messageContent string) (*model.Message, error)
	GetConversation(ctx context.Context, senderID, receiverID uuid.UUID) ([]*model.Message, error)
	GetMessageByID(ctx context.Context, id uuid.UUID) (*model.Message, error)
}
