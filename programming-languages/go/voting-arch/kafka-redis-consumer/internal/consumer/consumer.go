package consumer

import (
	"context"
	"encoding/json"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sirupsen/logrus"

	"github.com/alex-carvalho/kafka-redis-consumer/internal/cache"
	"github.com/alex-carvalho/kafka-redis-consumer/pkg/models"
)

const (
	batchSize    = 1000 // Batch size for Redis pipeline operations
	batchTimeout = 2 * time.Second
)

type Consumer struct {
	consumer *kafka.Consumer
	cache    *cache.RedisClient
	logger   *logrus.Logger
	done     chan bool
}

func NewConsumer(broker string, redisClient *cache.RedisClient, logger *logrus.Logger) *Consumer {
	config := kafka.ConfigMap{
		"bootstrap.servers":       broker,
		"group.id":                "redis-consumer-group",
		"auto.offset.reset":       "earliest",
		"enable.auto.commit":      true,
		"auto.commit.interval.ms": 1000,
	}

	consumer, err := kafka.NewConsumer(&config)
	if err != nil {
		logger.Fatalf("Failed to create Kafka consumer: %v", err)
	}

	return &Consumer{
		consumer: consumer,
		cache:    redisClient,
		logger:   logger,
		done:     make(chan bool),
	}
}

func (c *Consumer) Start() error {
	err := c.consumer.SubscribeTopics([]string{"votes"}, nil)
	if err != nil {
		c.logger.Errorf("Failed to subscribe to topic: %v", err)
		return err
	}

	go func() {
		batch := make([]models.Vote, 0, batchSize)
		timer := time.NewTimer(batchTimeout)
		defer timer.Stop()

		for {
			select {
			case <-c.done:
				// Flush remaining votes
				if len(batch) > 0 {
					c.flushBatchToRedis(batch)
				}
				c.logger.Info("Consumer stopped")
				return

			case <-timer.C:
				// Flush batch on timeout
				if len(batch) > 0 {
					c.flushBatchToRedis(batch)
					batch = make([]models.Vote, 0, batchSize)
				}
				timer.Reset(batchTimeout)

			default:
				msg, err := c.consumer.ReadMessage(100 * time.Millisecond)
				if err != nil {
					// Ignore timeout errors
					kafkaErr, ok := err.(kafka.Error)
					if ok && kafkaErr.Code() != kafka.ErrTimedOut {
						c.logger.Errorf("Error reading message: %v", err)
					}
					continue
				}

				var vote models.Vote
				err = json.Unmarshal(msg.Value, &vote)
				if err != nil {
					c.logger.Errorf("Error unmarshaling message: %v", err)
					continue
				}

				batch = append(batch, vote)

				// Flush batch when it reaches size limit
				if len(batch) >= batchSize {
					c.flushBatchToRedis(batch)
					batch = make([]models.Vote, 0, batchSize)
					timer.Reset(batchTimeout)
				}
			}
		}
	}()

	c.logger.Info("Connected to Kafka, listening on topic 'votes'")
	return nil
}

func (c *Consumer) flushBatchToRedis(batch []models.Vote) {
	err := c.cache.IncrBatch(context.Background(), batch)
	if err != nil {
		c.logger.Errorf("Error processing batch of %d votes: %v", len(batch), err)
		return
	}

	c.logger.Infof("Batch processed: %d votes counted in Redis", len(batch))
}

func (c *Consumer) Stop() error {
	c.done <- true
	return c.consumer.Close()
}
