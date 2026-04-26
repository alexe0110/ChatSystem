package dynamodb

import (
	"context"
	"fmt"
	"time"

	"github.com/alexe0110/chat-system/internal/model"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

type userItem struct {
	UserID         string `dynamodbav:"user_id"`
	Login          string `dynamodbav:"login"`
	Name           string `dynamodbav:"name"`
	HashedPassword string `dynamodbav:"hashed_password"`
	CreatedAt      string `dynamodbav:"created_at"`
	UpdatedAt      string `dynamodbav:"updated_at"`
}

func (i *userItem) toModel() (*model.User, error) {
	id, err := uuid.Parse(i.UserID)
	if err != nil {
		return nil, fmt.Errorf("parse user_id: %w", err)
	}

	createdAt, _ := time.Parse(time.RFC3339, i.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339, i.UpdatedAt)

	return &model.User{
		ID:             id,
		Login:          i.Login,
		Name:           i.Name,
		HashedPassword: i.HashedPassword,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}, nil
}

type UserRepo struct {
	db        *dynamodb.Client
	tableName string
}

func NewUserRepository(db *dynamodb.Client, tableName string) *UserRepo {
	return &UserRepo{
		db:        db,
		tableName: tableName,
	}
}

func (r *UserRepo) CreateUser(ctx context.Context, login, name, hashedPassword string) (*model.User, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	item := userItem{
		UserID:         uuid.New().String(),
		Login:          login,
		Name:           name,
		HashedPassword: hashedPassword,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return nil, fmt.Errorf("marshal user: %w", err)
	}

	_, err = r.db.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           aws.String(r.tableName),
		Item:                av,
		ConditionExpression: aws.String("attribute_not_exists(user_id)"),
	})

	if err != nil {
		return nil, fmt.Errorf("put user: %w", err)
	}

	return item.toModel()
}

func (r *UserRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	result, err := r.db.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"user_id": &types.AttributeValueMemberS{Value: id.String()},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	if result.Item == nil {
		return nil, fmt.Errorf("user not found")
	}

	var item userItem
	if err := attributevalue.UnmarshalMap(result.Item, &item); err != nil {
		return nil, fmt.Errorf("unmarshal user: %w", err)
	}

	return item.toModel()
}

func (r *UserRepo) GetByLogin(ctx context.Context, login string) (*model.User, error) {
	result, err := r.db.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		IndexName:              aws.String("login-index"),
		KeyConditionExpression: aws.String("login = :login"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":login": &types.AttributeValueMemberS{Value: login},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("query by login: %w", err)
	}

	if len(result.Items) == 0 {
		return nil, fmt.Errorf("user not found")
	}

	var item userItem
	if err := attributevalue.UnmarshalMap(result.Items[0], &item); err != nil {
		return nil, fmt.Errorf("unmarshal user: %w", err)
	}

	return item.toModel()
}
