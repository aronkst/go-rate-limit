package storage

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(options *redis.Options) (*RedisClient, error) {
	ctx := context.Background()

	client := redis.NewClient(options)
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return &RedisClient{client: client}, nil
}

func (r *RedisClient) IncrementRequestCount(identifier string, expiry time.Duration) (int, error) {
	ctx := context.Background()

	val, err := r.client.Incr(ctx, identifier).Result()
	if err != nil {
		return 0, err
	}

	if val == 1 {
		r.client.Expire(ctx, identifier, expiry)
	}

	return int(val), nil
}

func (r *RedisClient) GetRequestCount(identifier string) (int, error) {
	ctx := context.Background()

	val, err := r.client.Get(ctx, identifier).Result()
	if err == redis.Nil {
		return 0, nil
	} else if err != nil {
		return 0, err
	}

	count, err := strconv.Atoi(val)
	if err != nil {
		return 0, err
	}

	return count, nil
}
