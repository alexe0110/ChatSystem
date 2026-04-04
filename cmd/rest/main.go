package main

import (
	"database/sql"
	"log"

	restHandler "github.com/alexe0110/chat-system/internal/handler/rest"
	"github.com/alexe0110/chat-system/internal/middleware"
	"github.com/alexe0110/chat-system/internal/repository/postgres"
	"github.com/alexe0110/chat-system/internal/service"
	"github.com/gin-gonic/gin"

	_ "github.com/lib/pq"
)

func main() {
	const secret = "qwerty"

	db, err := sql.Open("postgres", "postgresql://postgres:postgres@localhost:5432/chat_db?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	userRepo := postgres.NewUserRepository(db)
	messageRepo := postgres.NewMessageRepository(db)

	userService := service.NewUserService(userRepo)
	messageService := service.NewMessageService(messageRepo)

	userHandler := restHandler.NewUserHandler(userService, secret)
	messageHandler := restHandler.NewMessageHandler(messageService)

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

	_ = router.Run(":8080")

}
