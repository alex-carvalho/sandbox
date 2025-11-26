package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"

	"github.com/alex-carvalho/kafka-postgres-consumer/internal/consumer"
	"github.com/alex-carvalho/kafka-postgres-consumer/internal/database"
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

	postgresURL := os.Getenv("DATABASE_URL")
	if postgresURL == "" {
		postgresURL = "postgres://postgres:postgres@localhost:5432/voting?sslmode=disable"
	}

	db, err := database.Connect(postgresURL)
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	logger.Info("Connected to PostgreSQL")

	c := consumer.NewConsumer(kafkaBroker, db, logger)
	if err := c.Start(); err != nil {
		logger.Fatalf("Failed to start consumer: %v", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	logger.Info("Shutting down kafka-postgres-consumer")
	c.Stop()
}
