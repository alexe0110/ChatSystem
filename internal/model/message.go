package model

import (
	"github.com/google/uuid"
	"time"
)

type Message struct {
	ID             uuid.UUID `json:"id"`
	SenderID       uuid.UUID `json:"sender_id"`
	ReceiverID     uuid.UUID `json:"receiver_id"`
	MessageContent string    `json:"message_content"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
