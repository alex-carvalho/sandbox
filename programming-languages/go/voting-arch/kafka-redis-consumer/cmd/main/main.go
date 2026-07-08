package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"

	"github.com/alex-carvalho/kafka-redis-consumer/internal/cache"
	"github.com/alex-carvalho/kafka-redis-consumer/internal/consumer"
)

var logger = logrus.New()

func init() {
	logger.SetFormatter(&logrus.JSONFormatter{})
}

func main() {
	kafkaBroker := os.Getenv("KAFKA_BROKER")
	if kafkaBroker == "" {
		kafkaBroker = "localhost:9092"
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	redisClient, err := cache.NewRedisClient(redisAddr)
	if err != nil {
		logger.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()

	logger.Info("Connected to Redis")

	c := consumer.NewConsumer(kafkaBroker, redisClient, logger)
	if err := c.Start(); err != nil {
		logger.Fatalf("Failed to start consumer: %v", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	logger.Info("Shutting down kafka-redis-consumer")
	c.Stop()
}
