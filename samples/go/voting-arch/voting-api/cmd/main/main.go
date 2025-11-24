package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"github.com/alex-carvalho/voting-api/internal/handlers"
	"github.com/alex-carvalho/voting-api/internal/kafka"
	"github.com/alex-carvalho/voting-api/internal/middleware"
	"github.com/alex-carvalho/voting-api/internal/redis"
)

var logger = logrus.New()

func init() {
	godotenv.Load()
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

	redisClient, err := redis.NewRedisClient(redisAddr)
	if err != nil {
		logger.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()

	producer, err := kafka.NewProducer(kafkaBroker)
	if err != nil {
		logger.Fatalf("Failed to create Kafka producer: %v", err)
	}
	defer producer.Close()

	router := gin.Default()
	router.Use(middleware.PanicRecovery())

	voteHandler := handlers.NewVoteHandler(producer, logger)
	resultHandler := handlers.NewResultHandler(redisClient, logger)

	router.POST("/vote", voteHandler.Handle)
	router.GET("/results", resultHandler.Handle)

	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8081"
	}

	go func() {
		if err := router.Run(":" + port); err != nil {
			logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	logger.Info("Shutting down voting-api")
}
