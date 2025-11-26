package consumer

import (
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"

	"github.com/alex-carvalho/kafka-postgres-consumer/internal/database"
	"github.com/alex-carvalho/kafka-postgres-consumer/pkg/models"
)

type Consumer struct {
	reader *kafka.Reader
	db     *database.DB
	logger *logrus.Logger
	done   chan bool
}

func NewConsumer(broker string, db *database.DB, logger *logrus.Logger) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker},
		GroupID: "postgres-consumer-group",
		Topic:   "votes",
	})

	return &Consumer{
		reader: reader,
		db:     db,
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

				err = c.db.InsertVote(vote.UserID, vote.VotingID, vote.VoteOption)
				if err != nil {
					c.logger.Errorf("Error inserting vote: %v", err)
					continue
				}

				c.logger.Infof("Vote stored: user_id=%d, voting_id=%d, vote_option=%d", vote.UserID, vote.VotingID, vote.VoteOption)
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
