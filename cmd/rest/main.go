package main

import (
	"context"
	"database/sql"
	"log"
	"os"

	restHandler "github.com/alexe0110/chat-system/internal/handler/rest"
	"github.com/alexe0110/chat-system/internal/middleware"
	"github.com/alexe0110/chat-system/internal/repository"
	dynamorepo "github.com/alexe0110/chat-system/internal/repository/dynamodb"
	"github.com/alexe0110/chat-system/internal/repository/postgres"
	"github.com/alexe0110/chat-system/internal/service"
	"github.com/alexe0110/chat-system/internal/worker"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gin-gonic/gin"

	_ "github.com/lib/pq"
)

func main() {
	const secret = "qwerty"

	dbType := os.Getenv("DB_TYPE")
	pgDSN := os.Getenv("DATABASE_URL")
	AWSRegion := os.Getenv("AWS_REGION")

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

	notificationWorker := worker.NewNotificationWorker(5)
	notificationWorker.Start()

	userService := service.NewUserService(userRepo)
	messageService := service.NewMessageService(messageRepo)

	userHandler := restHandler.NewUserHandler(userService, secret)
	messageHandler := restHandler.NewMessageHandler(messageService, notificationWorker)

	router := gin.Default()

	userRouter := router.Group("user")
	userRouter.POST("/register", userHandler.Register)
	userRouter.POST("/login", userHandler.Login)
	userRouter.GET("/:id", middleware.AuthMiddleware(secret), userHandler.GetByID)

	messageRouter := router.Group("message")
	messageRouter.Use(middleware.AuthMiddleware(secret))
	messageRouter.POST("/send", messageHandler.SendMessage)
	messageRouter.GET("/conversation", messageHandler.GetConversation)
	messageRouter.GET("/:id", messageHandler.GetByID)

	serviceRouter := router.Group("")
	serviceRouter.GET("/health")

	_ = router.Run(":8080")
}
