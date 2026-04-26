package dynamodb

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/alexe0110/chat-system/internal/model"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

type messageItem struct {
	ChatID         string `dynamodbav:"chat_id"`
	CreatedAt      string `dynamodbav:"created_at"`
	MessageID      string `dynamodbav:"message_id"`
	SenderID       string `dynamodbav:"sender_id"`
	ReceiverID     string `dynamodbav:"receiver_id"`
	MessageContent string `dynamodbav:"message_content"`
	UpdatedAt      string `dynamodbav:"updated_at"`
}

func buildChatID(a, b uuid.UUID) string {
	ids := []string{a.String(), b.String()}
	sort.Strings(ids)
	return strings.Join(ids, "#")
}

func (m *messageItem) toModel() (*model.Message, error) {
	id, err := uuid.Parse(m.MessageID)
	if err != nil {
		return nil, fmt.Errorf("parse message_id: %w", err)
	}
	senderID, _ := uuid.Parse(m.SenderID)
	receiverID, _ := uuid.Parse(m.ReceiverID)
	createdAt, _ := time.Parse(time.RFC3339, m.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339, m.UpdatedAt)

	return &model.Message{
		ID:             id,
		SenderID:       senderID,
		ReceiverID:     receiverID,
		MessageContent: m.MessageContent,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}, nil
}

type MessageRepo struct {
	db        *dynamodb.Client
	tableName string
}

func NewMessageRepository(db *dynamodb.Client, tableName string) *MessageRepo {
	return &MessageRepo{
		db:        db,
		tableName: tableName,
	}
}

func (r *MessageRepo) SendMessage(ctx context.Context, senderID, receiverID uuid.UUID, messageContent string) (*model.Message, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	item := messageItem{
		ChatID:         buildChatID(senderID, receiverID),
		CreatedAt:      now,
		MessageID:      uuid.New().String(),
		SenderID:       senderID.String(),
		ReceiverID:     receiverID.String(),
		MessageContent: messageContent,
		UpdatedAt:      now,
	}

	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return nil, fmt.Errorf("marshal message: %w", err)
	}

	_, err = r.db.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      av,
	})
	if err != nil {
		return nil, fmt.Errorf("put message: %w", err)
	}

	return item.toModel()
}

func (r *MessageRepo) GetConversation(ctx context.Context, senderID, receiverID uuid.UUID) ([]*model.Message, error) {
	chatID := buildChatID(senderID, receiverID)

	result, err := r.db.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		KeyConditionExpression: aws.String("chat_id = :cid"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":cid": &types.AttributeValueMemberS{Value: chatID},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("query conversation: %w", err)
	}

	messages := make([]*model.Message, 0, len(result.Items))
	for _, item := range result.Items {
		var mi messageItem
		if err := attributevalue.UnmarshalMap(item, &mi); err != nil {
			return nil, fmt.Errorf("unmarshal message: %w", err)
		}
		msg, err := mi.toModel()
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

func (r *MessageRepo) GetMessageByID(ctx context.Context, id uuid.UUID) (*model.Message, error) {
	result, err := r.db.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		IndexName:              aws.String("message-id-index"),
		KeyConditionExpression: aws.String("message_id = :mid"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":mid": &types.AttributeValueMemberS{Value: id.String()},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("query by message_id: %w", err)
	}

	if len(result.Items) == 0 {
		return nil, fmt.Errorf("message not found")
	}

	var item messageItem
	if err := attributevalue.UnmarshalMap(result.Items[0], &item); err != nil {
		return nil, fmt.Errorf("unmarshal message: %w", err)
	}

	return item.toModel()
}
