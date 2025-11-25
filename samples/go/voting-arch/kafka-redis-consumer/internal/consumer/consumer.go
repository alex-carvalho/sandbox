package consumer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"

	"github.com/alex-carvalho/kafka-redis-consumer/internal/cache"
	"github.com/alex-carvalho/kafka-redis-consumer/pkg/models"
)

type Consumer struct {
	reader *kafka.Reader
	cache  *cache.RedisClient
	logger *logrus.Logger
	done   chan bool
}

func NewConsumer(broker string, redisClient *cache.RedisClient, logger *logrus.Logger) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker},
		GroupID: "redis-consumer-group",
		Topic:   "votes",
	})

	return &Consumer{
		reader: reader,
		cache:  redisClient,
		logger: logger,
		done:   make(chan bool),
	}
}

func (c *Consumer) Start() error {
	go func() {
		for {
			select {
			case <-c.done:
				return
			default:
				msg, err := c.reader.ReadMessage(context.Background())
				if err != nil {
					c.logger.Errorf("Error reading message: %v", err)
					continue
				}

				var vote models.Vote
				err = json.Unmarshal(msg.Value, &vote)
				if err != nil {
					c.logger.Errorf("Error unmarshaling message: %v", err)
					continue
				}

				key := fmt.Sprintf("votes:%d:%d", vote.VotingID, vote.VoteOption)
				err = c.cache.Incr(context.Background(), key)
				if err != nil {
					c.logger.Errorf("Error incrementing Redis counter: %v", err)
					continue
				}

				c.logger.Infof("Vote counted in Redis: voting_id=%d, vote_option=%d", vote.VotingID, vote.VoteOption)
			}
		}
	}()

	c.logger.Info("Connected to Kafka, listening on topic 'votes'")
	return nil
}

func (c *Consumer) Stop() error {
	c.done <- true
	return c.reader.Close()
}
