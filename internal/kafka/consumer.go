package kafka

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/Shopify/sarama"
)

type Consumer struct {
	consumer sarama.ConsumerGroup
	topics   []string
	handler  MessageHandler
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
	}, nil
}

func (c *Consumer) Start(ctx context.Context) error {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			if err := c.consumer.Consume(ctx, c.topics, &consumerGroupHandler{handler: c.handler}); err != nil {
				log.Printf("Error from consumer: %v", err)
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
}

func (h *consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var kafkaMsg Message
		if err := json.Unmarshal(msg.Value, &kafkaMsg); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			continue
		}

		if err := h.handler.HandleMessage(&kafkaMsg); err != nil {
			log.Printf("Error handling message: %v", err)
			continue
		}

		session.MarkMessage(msg, "")
	}
	return nil
}
