package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/alex-carvalho/kafka-redis-consumer/pkg/models"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(addr string) (*RedisClient, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:         addr,
		PoolSize:     50, // Connection pool size
		MaxIdleConns: 10, // Max idle connections
		MinIdleConns: 5,  // Min idle connections
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return &RedisClient{client: rdb}, nil
}

func (rc *RedisClient) Incr(ctx context.Context, key string) error {
	return rc.client.Incr(ctx, key).Err()
}

// IncrBatch uses Redis pipeline for batch operations
// This is much faster than individual INCR operations
func (rc *RedisClient) IncrBatch(ctx context.Context, votes []models.Vote) error {
	if len(votes) == 0 {
		return nil
	}

	// Use pipeline for batch operations
	pipe := rc.client.Pipeline()

	// Group votes by voting_id for better cache locality
	voteMap := make(map[int]map[int]int) // voting_id -> vote_option -> count

	for _, vote := range votes {
		if voteMap[vote.VotingID] == nil {
			voteMap[vote.VotingID] = make(map[int]int)
		}
		voteMap[vote.VotingID][vote.VoteOption]++
	}

	// Add commands to pipeline using INCRBY for better performance
	for votingID, options := range voteMap {
		for option, count := range options {
			key := fmt.Sprintf("votes:%d:%d", votingID, option)
			pipe.IncrBy(ctx, key, int64(count))
		}
	}

	// Execute all commands at once
	_, err := pipe.Exec(ctx)
	return err
}

func (rc *RedisClient) Close() error {
	return rc.client.Close()
}
