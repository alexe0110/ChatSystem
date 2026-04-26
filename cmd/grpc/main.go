package main

import (
	"database/sql"
	"log"
	"net"
	"os"

	"github.com/alexe0110/chat-system/internal/hub"
	"github.com/alexe0110/chat-system/internal/middleware"
	"github.com/alexe0110/chat-system/internal/repository/postgres"
	"github.com/alexe0110/chat-system/internal/service"
	"github.com/alexe0110/chat-system/pb"
	"github.com/alexe0110/chat-system/pkg/storage"
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

	if pgDSN == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	db, err := sql.Open("postgres", pgDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.AuthInterceptor(secret)),
		grpc.StreamInterceptor(middleware.AuthStreamInterceptor(secret)),
	)

	userRepo := postgres.NewUserRepository(db)
	messageRepo := postgres.NewMessageRepository(db)

	userService := service.NewUserService(userRepo)
	messageService := service.NewMessageService(messageRepo)

	minioStorage := storage.NewMinioStorage(minioEndpoint, minioUser, minioPassword, "files")

	chatHub := hub.NewHub()
	userServiceServer := myGPRC.NewUserServiceServer(userService)
	chatServiceServer := myGPRC.NewChatServiceServer(messageService, chatHub, minioStorage)

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
