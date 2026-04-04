package rest

import (
	"net/http"

	"github.com/alexe0110/chat-system/internal/service"
	"github.com/alexe0110/chat-system/pkg/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	service *service.UserService
	secret  string
}

func NewUserHandler(userService *service.UserService, secret string) *UserHandler {
	return &UserHandler{
		service: userService,
		secret:  secret,
	}
}

type RegisterRequest struct {
	Login    string `json:"login" binding:"required"`
	Name     string `json:"name"`
	Password string `json:"password" binding:"required"`
}

type LoginRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (handler *UserHandler) Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := handler.service.Login(c.Request.Context(), req.Login, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	token, err := auth.GenerateToken(result.ID, handler.secret)
	if err != nil {
		c.JSON(500, gin.H{"error": "token generation failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (handler *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := handler.service.Register(c.Request.Context(), req.Login, req.Name, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (handler *UserHandler) GetByID(c *gin.Context) {
	IDStr := c.Param("id")
	ID, err := uuid.Parse(IDStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := handler.service.GetByID(c.Request.Context(), ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}
