package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(addr string) (*RedisClient, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
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

func (rc *RedisClient) Close() error {
	return rc.client.Close()
}
