package kafka

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Shopify/sarama"
)

type Producer struct {
	producer sarama.SyncProducer
	topic    string
	cb       *CircuitBreaker
}

func NewProducer(brokers []string, topic string) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &Producer{
		producer: producer,
		topic:    topic,
		cb:       NewCircuitBreaker(5, 1*time.Minute),
	}, nil
}

func (p *Producer) SendMessage(msg *Message) error {
	if !p.cb.AllowRequest() {
		return fmt.Errorf("circuit breaker is open")
	}

	jsonData, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	message := &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.StringEncoder(jsonData),
		Key:   sarama.StringEncoder(msg.Type), // Use the query type as the key for partitioning
	}

	partition, offset, err := p.producer.SendMessage(message)
	if err != nil {
		p.cb.OnFailure()
		return err
	}

	p.cb.OnSuccess()
	log.Printf("Message sent to partition %d at offset %d\n", partition, offset)
	return nil
}

func (p *Producer) Close() error {
	return p.producer.Close()
}
