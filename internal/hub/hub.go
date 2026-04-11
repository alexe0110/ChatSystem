package hub

import (
	"fmt"
	"sync"

	"github.com/alexe0110/chat-system/pb"
	"github.com/google/uuid"
)

type Hub struct {
	clients map[uuid.UUID]chan *pb.ChatMessage
	mu      sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[uuid.UUID]chan *pb.ChatMessage),
	}
}

func (h *Hub) Register(userID uuid.UUID) chan *pb.ChatMessage {
	ch := make(chan *pb.ChatMessage, 10)
	h.mu.Lock()
	h.clients[userID] = ch
	h.mu.Unlock()

	return ch
}

func (h *Hub) Unregister(userID uuid.UUID) {
	h.mu.Lock()
	ch, ok := h.clients[userID]
	if ok {
		delete(h.clients, userID)
	}
	h.mu.Unlock()

	if ok {
		close(ch)
	}

}

func (h *Hub) Send(receiverID uuid.UUID, msg *pb.ChatMessage) error {
	h.mu.RLock()
	ch, ok := h.clients[receiverID]
	h.mu.RUnlock()

	if !ok {
		return fmt.Errorf("user %s not connected", receiverID)
	}

	ch <- msg
	return nil
}
