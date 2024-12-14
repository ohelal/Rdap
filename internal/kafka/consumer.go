package kafka

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/Shopify/sarama"
)

type Consumer struct {
	consumer sarama.ConsumerGroup
	topics   []string
	handler  MessageHandler
	cb       *CircuitBreaker
}

type MessageHandler interface {
	HandleMessage(msg *Message) error
}

func NewConsumer(brokers []string, groupID string, topics []string, handler MessageHandler) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumer, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		consumer: consumer,
		topics:   topics,
		handler:  handler,
		cb:       NewCircuitBreaker(5, 1*time.Minute),
	}, nil
}

func (c *Consumer) Start(ctx context.Context) error {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			if !c.cb.AllowRequest() {
				log.Println("Circuit breaker is open, waiting before retry...")
				time.Sleep(5 * time.Second)
				continue
			}

			if err := c.consumer.Consume(ctx, c.topics, &consumerGroupHandler{
				handler: c.handler,
				cb:      c.cb,
			}); err != nil {
				log.Printf("Error from consumer: %v", err)
				c.cb.OnFailure()
				continue
			}

			if ctx.Err() != nil {
				return
			}
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down consumer...")
	wg.Wait()
	return c.consumer.Close()
}

type consumerGroupHandler struct {
	handler MessageHandler
	cb      *CircuitBreaker
}

func (h *consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error { return nil }
func (h *consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var kafkaMsg Message
		if err := json.Unmarshal(msg.Value, &kafkaMsg); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			h.cb.OnFailure()
			continue
		}

		if err := h.handler.HandleMessage(&kafkaMsg); err != nil {
			log.Printf("Error handling message: %v", err)
			h.cb.OnFailure()
			continue
		}

		session.MarkMessage(msg, "")
		h.cb.OnSuccess()
	}
	return nil
}
