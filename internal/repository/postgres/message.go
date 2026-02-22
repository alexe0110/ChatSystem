package postgres

import (
	"context"
	"database/sql"
	"github.com/alexe0110/chat-system/internal/model"
	"github.com/google/uuid"
)

const (
	queryCreateMessage = `INSERT INTO messages (id, sender_id, receiver_id, message_content)
						VALUES ($1, $2, $3, $4)
						RETURNING id, sender_id, receiver_id, message_content, created_at, updated_at`
	queryGetMessageByID  = `SELECT id, sender_id, receiver_id, message_content, created_at, updated_at FROM messages WHERE id=$1`
	queryGetConversation = `SELECT id, sender_id, receiver_id, message_content, created_at, updated_at FROM messages 
						   WHERE (sender_id=$1 and receiver_id=$2) or (sender_id=$2 and receiver_id=$1) ORDER BY created_at`
)

type MessageRepo struct {
	db *sql.DB
}

func NewMessageRepository(db *sql.DB) *MessageRepo {
	return &MessageRepo{
		db: db,
	}
}

func (repo *MessageRepo) SendMessage(ctx context.Context, senderID, receiverID uuid.UUID, messageContent string) (*model.Message, error) {
	id := uuid.New()
	message := &model.Message{}
	err := repo.db.QueryRowContext(ctx, queryCreateMessage, id, senderID, receiverID, messageContent).
		Scan(&message.ID, &message.SenderID, &message.ReceiverID, &message.MessageContent, &message.CreatedAt, &message.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return message, nil
}

func (repo *MessageRepo) GetMessageByID(ctx context.Context, id uuid.UUID) (*model.Message, error) {
	message := &model.Message{}
	err := repo.db.QueryRowContext(ctx, queryGetMessageByID, id).
		Scan(&message.ID, &message.SenderID, &message.ReceiverID, &message.MessageContent, &message.CreatedAt, &message.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return message, nil
}

func (repo *MessageRepo) GetConversation(ctx context.Context, senderID, receiverID uuid.UUID) ([]*model.Message, error) {
	var messages []*model.Message

	rows, err := repo.db.QueryContext(ctx, queryGetConversation, senderID, receiverID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		message := &model.Message{}
		err := rows.Scan(&message.ID, &message.SenderID, &message.ReceiverID, &message.MessageContent, &message.CreatedAt, &message.UpdatedAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}
