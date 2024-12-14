package kafka

import (
	"log"
)

type defaultMessageHandler struct{}

func NewMessageHandler() MessageHandler {
	return &defaultMessageHandler{}
}

func (h *defaultMessageHandler) HandleMessage(msg *Message) error {
	log.Printf("Received message - Type: %s, Query: %s, CacheHit: %v", msg.Type, msg.Query, msg.CacheHit)
	return nil
}
