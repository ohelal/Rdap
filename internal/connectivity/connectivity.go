package connectivity

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Shopify/sarama"
	"github.com/go-redis/redis/v8"
)

func RunTests() error {
	// Test Redis Connection
	fmt.Println("Testing Redis connection...")
	if err := testRedis(); err != nil {
		log.Printf("Redis test failed: %v", err)
	} else {
		fmt.Println("Redis connection successful!")
	}

	// Test Kafka Connection
	fmt.Println("\nTesting Kafka connection...")
	if err := testKafka(); err != nil {
		log.Printf("Kafka test failed: %v", err)
		return err
	}
	fmt.Println("Kafka connection successful!")
	return nil
}

func testRedis() error {
	// Get Redis URL from environment variable or use default
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "localhost:6379"
	}

	// Create Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr: redisURL,
	})

	// Test connection with ping
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis ping failed: %w", err)
	}

	// Test set and get
	testKey := "test_key"
	testValue := "test_value"

	// Set value
	if err := rdb.Set(ctx, testKey, testValue, 1*time.Minute).Err(); err != nil {
		return fmt.Errorf("redis set failed: %w", err)
	}

	// Get value
	val, err := rdb.Get(ctx, testKey).Result()
	if err != nil {
		return fmt.Errorf("redis get failed: %w", err)
	}

	if val != testValue {
		return fmt.Errorf("redis value mismatch: got %s, want %s", val, testValue)
	}

	// Clean up
	if err := rdb.Del(ctx, testKey).Err(); err != nil {
		return fmt.Errorf("redis cleanup failed: %w", err)
	}

	return nil
}

func testKafka() error {
	// Get Kafka brokers from environment variable or use default
	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokers == "" {
		kafkaBrokers = "localhost:9092"
	}

	// Create Kafka config
	config := sarama.NewConfig()
	config.Version = sarama.V2_8_0_0
	config.Producer.Return.Successes = true

	// Create producer
	producer, err := sarama.NewSyncProducer([]string{kafkaBrokers}, config)
	if err != nil {
		return fmt.Errorf("failed to create producer: %w", err)
	}
	defer producer.Close()

	// Create test topic
	testTopic := "rdap_test_topic"
	testMessage := "test_message"

	// Send test message
	msg := &sarama.ProducerMessage{
		Topic: testTopic,
		Value: sarama.StringEncoder(testMessage),
	}

	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	fmt.Printf("Message sent to partition %d at offset %d\n", partition, offset)

	// Create consumer
	consumer, err := sarama.NewConsumer([]string{kafkaBrokers}, config)
	if err != nil {
		return fmt.Errorf("failed to create consumer: %w", err)
	}
	defer consumer.Close()

	// Create partition consumer
	partitionConsumer, err := consumer.ConsumePartition(testTopic, partition, offset)
	if err != nil {
		return fmt.Errorf("failed to create partition consumer: %w", err)
	}
	defer partitionConsumer.Close()

	// Wait for message
	select {
	case msg := <-partitionConsumer.Messages():
		if string(msg.Value) != testMessage {
			return fmt.Errorf("message mismatch: got %s, want %s", string(msg.Value), testMessage)
		}
		fmt.Printf("Message received: %s\n", string(msg.Value))
	case err := <-partitionConsumer.Errors():
		return fmt.Errorf("consumer error: %w", err)
	case <-time.After(5 * time.Second):
		return fmt.Errorf("timeout waiting for message")
	}

	return nil
}
