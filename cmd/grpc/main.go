package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"os"

	"github.com/alexe0110/chat-system/internal/hub"
	"github.com/alexe0110/chat-system/internal/middleware"
	"github.com/alexe0110/chat-system/internal/repository"
	dynamorepo "github.com/alexe0110/chat-system/internal/repository/dynamodb"
	"github.com/alexe0110/chat-system/internal/repository/postgres"
	"github.com/alexe0110/chat-system/internal/service"
	"github.com/alexe0110/chat-system/pb"
	"github.com/alexe0110/chat-system/pkg/storage"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"google.golang.org/grpc"

	myGPRC "github.com/alexe0110/chat-system/internal/handler/grpc"
	_ "github.com/lib/pq"
)

func main() {
	const secret = "qwerty"

	pgDSN := os.Getenv("DATABASE_URL")
	minioEndpoint := os.Getenv("MINIO_ENDPOINT")
	minioUser := os.Getenv("MINIO_USER")
	minioPassword := os.Getenv("MINIO_PASSWORD")
	storageType := os.Getenv("STORAGE_TYPE")
	dbType := os.Getenv("DB_TYPE")
	AWSRegion := os.Getenv("AWS_REGION")
	bucketName := "chatsystem-files-dev"

	var userRepo repository.UserRepository
	var messageRepo repository.MessageRepository

	if dbType == "dynamodb" {
		cfg, err := config.LoadDefaultConfig(context.Background(),
			config.WithRegion(AWSRegion),
		)
		if err != nil {
			log.Fatalf("unable to load AWS config: %v", err)
		}

		dynamoClient := dynamodb.NewFromConfig(cfg)

		userRepo = dynamorepo.NewUserRepository(dynamoClient, "users")
		messageRepo = dynamorepo.NewMessageRepository(dynamoClient, "messages")
	} else {
		if pgDSN == "" {
			log.Fatal("DATABASE_URL is not set")
		}

		db, err := sql.Open("postgres", pgDSN)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		userRepo = postgres.NewUserRepository(db)
		messageRepo = postgres.NewMessageRepository(db)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.AuthInterceptor(secret)),
		grpc.StreamInterceptor(middleware.AuthStreamInterceptor(secret)),
	)

	userService := service.NewUserService(userRepo)
	messageService := service.NewMessageService(messageRepo)

	var fileStorage storage.FileStorage
	if storageType == "s3" {
		fileStorage = storage.NewS3Storage(bucketName, AWSRegion)
	} else {
		fileStorage = storage.NewMinioStorage(minioEndpoint, minioUser, minioPassword, bucketName)
	}

	chatHub := hub.NewHub()
	userServiceServer := myGPRC.NewUserServiceServer(userService)
	chatServiceServer := myGPRC.NewChatServiceServer(messageService, chatHub, fileStorage)

	pb.RegisterUserServiceServer(grpcServer, userServiceServer)
	pb.RegisterChatServiceServer(grpcServer, chatServiceServer)

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("gRPC server starting on :50051")

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal(err)
	}
}
