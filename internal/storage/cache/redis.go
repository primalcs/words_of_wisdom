package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/primalcs/words_of_wisdom/internal/interfaces"
	"github.com/redis/go-redis/v9"
)

var _ interfaces.Cache = &RedisClient{}

type RedisClient struct {
	client *redis.Client
}

func NewRedisCache(ctx context.Context, address string) (*RedisClient, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: address,
	})

	status := rdb.Ping(ctx)
	if status.Err() != nil {
		return nil, fmt.Errorf("couldn't ping redis cache: %v", status.Err())
	}

	return &RedisClient{
		client: rdb,
	}, nil
}

func (rc *RedisClient) InsertClientToken(ctx context.Context, key, value string, duration time.Duration) error {
	return rc.client.Set(ctx, key, value, duration).Err()
}

func (rc *RedisClient) GetClientToken(ctx context.Context, key string) (string, error) {
	res, err := rc.client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return res, nil
}

func (rc *RedisClient) Delete(ctx context.Context, key string) error {
	return rc.client.Del(ctx, key).Err()
}

func (rc *RedisClient) Close(ctx context.Context) error {
	return rc.client.Close()
}
