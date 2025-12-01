package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/alex-carvalho/voting-api/internal/kafka"
	"github.com/alex-carvalho/voting-api/pkg/models"
)

type VoteHandler struct {
	producer *kafka.Producer
	logger   *logrus.Logger
}

func NewVoteHandler(producer *kafka.Producer, logger *logrus.Logger) *VoteHandler {
	return &VoteHandler{
		producer: producer,
		logger:   logger,
	}
}

func (h *VoteHandler) Handle(c *gin.Context) {
	var req models.VoteRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warnf("Invalid vote request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message := map[string]interface{}{
		"id":          uuid.New().String(), // Unique message ID for idempotence
		"user_id":     req.UserID,
		"voting_id":   req.VotingID,
		"vote_option": req.VoteOption,
		"timestamp":   time.Now().Unix(), // Timestamp for ordering
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		h.logger.Errorf("Failed to marshal vote message: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process vote"})
		return
	}

	err = h.producer.SendMessage("votes", messageBytes)
	if err != nil {
		h.logger.Errorf("Failed to send vote to Kafka: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store vote"})
		return
	}

	h.logger.Infof("Vote received: user_id=%d, voting_id=%d, vote_option=%d", req.UserID, req.VotingID, req.VoteOption)
	c.JSON(http.StatusOK, models.VoteResponse{
		Status:  "success",
		Message: "Vote recorded",
	})
}
