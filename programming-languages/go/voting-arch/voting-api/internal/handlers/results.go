package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/alex-carvalho/voting-api/internal/redis"
	"github.com/alex-carvalho/voting-api/pkg/models"
)

type ResultHandler struct {
	redisClient *redis.Client
	logger      *logrus.Logger
}

func NewResultHandler(redisClient *redis.Client, logger *logrus.Logger) *ResultHandler {
	return &ResultHandler{
		redisClient: redisClient,
		logger:      logger,
	}
}

func (h *ResultHandler) Handle(c *gin.Context) {
	votingIDStr := c.Query("voting_id")
	if votingIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "voting_id parameter required"})
		return
	}

	votingID, err := strconv.Atoi(votingIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid voting_id"})
		return
	}

	ctx := c.Request.Context()
	keys, err := h.redisClient.Keys(ctx, fmt.Sprintf("votes:%d:*", votingID))
	if err != nil {
		h.logger.Errorf("Failed to fetch results from Redis: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch results"})
		return
	}

	results := make(map[int]int)
	totalVotes := 0

	for _, key := range keys {
		val, err := h.redisClient.GetInt(ctx, key)
		if err != nil {
			h.logger.Warnf("Failed to get value for key %s: %v", key, err)
			continue
		}

		var option int
		fmt.Sscanf(key, fmt.Sprintf("votes:%d:%%d", votingID), &option)
		results[option] = val
		totalVotes += val
	}

	h.logger.Infof("Results fetched for voting_id=%d", votingID)
	c.JSON(http.StatusOK, models.ResultsResponse{
		VotingID:   votingID,
		Results:    results,
		TotalVotes: totalVotes,
	})
}
