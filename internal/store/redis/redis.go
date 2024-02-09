package redis

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisClient(host, port string, db int, ctx context.Context) *RedisClient {
	if os.Getenv("REDIS_HOST") != "" {
		host = os.Getenv("REDIS_HOST")
	}
	if string(os.Getenv("REDIS_PORT")) != "" {
		port = string(os.Getenv("REDIS_PORT"))
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: host + ":" + port,
		DB:   db,
	})
	fmt.Println("Redis: ", host+":"+port)
	return &RedisClient{client: rdb, ctx: ctx}
}

func (r *RedisClient) SetData(key string, data []byte) error {
	return r.client.Set(r.ctx, key, data, 1*time.Hour).Err()
}

func (r *RedisClient) GetData(key string) ([]byte, error) {
	return r.client.Get(r.ctx, key).Bytes()
}
