package main

import (
    "context"
    "log"
    "os"
    "os/signal"
    "syscall"
    "rdap_service/internal/queue"
)

func main() {
    cfg := queue.KafkaConfig{
        Brokers: []string{"kafka-1:9092", "kafka-2:9092", "kafka-3:9092"},
        Topic:   "rdap_tasks",
        Group:   "rdap_workers",
        RetryMax: 3,
        RequiredAcks: sarama.WaitForAll,
        Timeout: 10 * time.Second,
    }

    processor, err := queue.NewTaskProcessor(cfg)
    if err != nil {
        log.Fatalf("Failed to create task processor: %v", err)
    }

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // Handle graceful shutdown
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

    go func() {
        <-sigChan
        cancel()
    }()

    if err := processor.Start(ctx); err != nil {
        log.Fatalf("Task processor failed: %v", err)
    }
} 