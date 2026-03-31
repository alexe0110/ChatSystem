package model

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	SenderID   uuid.UUID
	ReceiverID uuid.UUID
	Content    string
	CreatedAt  time.Time
}
