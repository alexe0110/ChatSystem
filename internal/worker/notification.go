package worker

import (
	"log"

	"github.com/alexe0110/chat-system/internal/model"
)

type NotificationWorker struct {
	ch chan model.Notification
}

func NewNotificationWorker(bufferSize int) *NotificationWorker {
	return &NotificationWorker{
		make(chan model.Notification, bufferSize),
	}
}

func (w *NotificationWorker) Send(n model.Notification) {
	w.ch <- n
}

func (w *NotificationWorker) Start() {
	go func() {
		for n := range w.ch {
			log.Printf("Notification for %s: %s", n.ReceiverID, n.Content)
		}
	}()
}
