package queue

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/IBM/sarama"
	"time"
)

var (
	ErrUnknownTaskType = errors.New("unknown task type")
)

type KafkaConfig struct {
	Brokers []string
	Topic   string
	Group   string
	// Additional Kafka-specific configurations
	RetryMax     int
	RequiredAcks sarama.RequiredAcks
	Timeout      time.Duration
}

type KafkaQueue struct {
	producer sarama.SyncProducer
	consumer sarama.ConsumerGroup
	topic    string
	group    string
}

func NewKafkaQueue(cfg KafkaConfig) (*KafkaQueue, error) {
	config := sarama.NewConfig()
	
	// Producer configurations
	config.Producer.RequiredAcks = cfg.RequiredAcks
	config.Producer.Retry.Max = cfg.RetryMax
	config.Producer.Return.Successes = true
	config.Producer.Timeout = cfg.Timeout
	
	// Create producer
	producer, err := sarama.NewSyncProducer(cfg.Brokers, config)
	if err != nil {
		return nil, err
	}

	// Create consumer
	consumer, err := sarama.NewConsumerGroup(cfg.Brokers, cfg.Group, config)
	if err != nil {
		producer.Close()
		return nil, err
	}

	return &KafkaQueue{
		producer: producer,
		consumer: consumer,
		topic:    cfg.Topic,
		group:    cfg.Group,
	}, nil
}

func (kq *KafkaQueue) Close() error {
	if err := kq.producer.Close(); err != nil {
		return err
	}
	return kq.consumer.Close()
}

func (kq *KafkaQueue) PublishTask(ctx context.Context, task interface{}) error {
	body, err := json.Marshal(task)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic:     kq.topic,
		Value:     sarama.ByteEncoder(body),
		Timestamp: time.Now(),
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		_, _, err = kq.producer.SendMessage(msg)
		return err
	}
}

func (kq *KafkaQueue) ConsumeTask(ctx context.Context, handler func([]byte) error) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			err := kq.consumer.Consume(ctx, []string{kq.topic}, &ConsumerGroupHandler{
				handler: handler,
			})
			if err != nil {
				return err
			}
		}
	}
}

// ConsumerGroupHandler implements sarama.ConsumerGroupHandler
type ConsumerGroupHandler struct {
	handler func([]byte) error
}

func (h *ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		select {
		case <-session.Context().Done():
			return nil
		default:
			if err := h.handler(msg.Value); err != nil {
				// Handle error (retry logic, dead letter queue, etc.)
				continue
			}
			session.MarkMessage(msg, "")
		}
	}
	return nil
}

// Example task processor
type TaskProcessor struct {
	queue *KafkaQueue
}

func NewTaskProcessor(cfg KafkaConfig) (*TaskProcessor, error) {
	queue, err := NewKafkaQueue(cfg)
	if err != nil {
		return nil, err
	}

	return &TaskProcessor{
		queue: queue,
	}, nil
}

func (tp *TaskProcessor) Start(ctx context.Context) error {
	return tp.queue.ConsumeTask(ctx, tp.processTask)
}

func (tp *TaskProcessor) processTask(data []byte) error {
	var task Task
	if err := json.Unmarshal(data, &task); err != nil {
		return err
	}

	switch task.Type {
	case "lookup":
		return tp.handleLookup(task)
	case "update":
		return tp.handleUpdate(task)
	default:
		return ErrUnknownTaskType
	}
}

type Task struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

func (tp *TaskProcessor) handleLookup(task Task) error {
	// Implement lookup logic
	return nil
}

func (tp *TaskProcessor) handleUpdate(task Task) error {
	// Implement update logic
	return nil
}