package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisClient(addr string, db int, ctx context.Context) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   db,
	})

	return &RedisClient{client: rdb, ctx: ctx}
}

func (r *RedisClient) SetData(key string, data []byte) error {
	return r.client.Set(r.ctx, key, data, 1*time.Hour).Err()
}

func (r *RedisClient) GetData(key string) ([]byte, error) {
	return r.client.Get(r.ctx, key).Bytes()
}
