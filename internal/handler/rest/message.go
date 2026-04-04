package rest

import (
	"net/http"

	"github.com/alexe0110/chat-system/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MessageHandler struct {
	service *service.MessageService
}

func NewMessageHandler(messageService *service.MessageService) *MessageHandler {
	return &MessageHandler{
		messageService,
	}
}

type SendMessageRequest struct {
	SenderID       uuid.UUID `json:"sender_id" binding:"required"`
	ReceiverID     uuid.UUID `json:"receiver_id" binding:"required"`
	MessageContent string    `json:"message_content" binding:"required"`
}

func (handler *MessageHandler) SendMessage(c *gin.Context) {
	var req SendMessageRequest

	err := c.ShouldBindJSON(&req)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	msg, err := handler.service.SendMessage(c.Request.Context(), req.SenderID, req.ReceiverID, req.MessageContent)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, msg)

}

func (handler *MessageHandler) GetByID(c *gin.Context) {
	msgID := c.Param("id")
	msgIDUUID, err := uuid.Parse(msgID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	msg, err := handler.service.GetMessageByID(c.Request.Context(), msgIDUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, msg)
}

func (handler *MessageHandler) GetConversation(c *gin.Context) {
	senderID := c.Query("sender_id")
	receiverID := c.Query("receiver_id")

	senderIDUUID, err := uuid.Parse(senderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	receiverIDUUID, err := uuid.Parse(receiverID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	conversation, err := handler.service.GetConversation(c.Request.Context(), senderIDUUID, receiverIDUUID)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, conversation)
}
