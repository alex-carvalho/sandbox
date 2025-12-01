package consumer

import (
	"encoding/json"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sirupsen/logrus"

	"github.com/alex-carvalho/kafka-postgres-consumer/internal/database"
	"github.com/alex-carvalho/kafka-postgres-consumer/pkg/models"
)

const (
	batchSize    = 100
	batchTimeout = 5 * time.Second
)

type Consumer struct {
	consumer *kafka.Consumer
	db       *database.DB
	logger   *logrus.Logger
	done     chan bool
}

func NewConsumer(broker string, db *database.DB, logger *logrus.Logger) *Consumer {
	config := kafka.ConfigMap{
		"bootstrap.servers":       broker,
		"group.id":                "postgres-consumer-group",
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
		db:       db,
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
					c.flushBatch(batch)
				}
				c.logger.Info("Consumer stopped")
				return

			case <-timer.C:
				// Flush batch on timeout
				if len(batch) > 0 {
					c.flushBatch(batch)
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

				c.logger.Debugf("Received message: %s", string(msg.Value))

				var vote models.Vote
				err = json.Unmarshal(msg.Value, &vote)
				if err != nil {
					c.logger.Errorf("Error unmarshaling message: %v", err)
					continue
				}

				batch = append(batch, vote)

				// Flush batch when it reaches size limit
				if len(batch) >= batchSize {
					c.flushBatch(batch)
					batch = make([]models.Vote, 0, batchSize)
					timer.Reset(batchTimeout)
				}
			}
		}
	}()

	c.logger.Info("Connected to Kafka, listening on topic 'votes'")
	return nil
}

func (c *Consumer) flushBatch(batch []models.Vote) {
	err := c.db.InsertVotesBatch(batch)
	if err != nil {
		c.logger.Errorf("Error inserting batch of %d votes: %v", len(batch), err)
		return
	}

	c.logger.Infof("Batch inserted: %d votes stored", len(batch))
}

func (c *Consumer) Stop() error {
	c.done <- true
	return c.consumer.Close()
}
