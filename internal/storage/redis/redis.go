package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/shirr9/order-api/internal/config"
	"time"
)

type CacheRepository interface {
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Get(ctx context.Context, key string) ([]byte, error)
	Delete(ctx context.Context, key string) error
	Close() error
}

type Redis struct {
	client *redis.Client
}

func NewRedis(ctx context.Context, cfg *config.Config) (*Redis, error) {
	// url = "redis://<user>:<pass>@localhost:6379/<db>"
	adr := fmt.Sprintf("%s:%s", cfg.RedisDB.Host, cfg.RedisDB.Port)
	rdb := redis.NewClient(&redis.Options{
		Addr:     adr,
		Password: cfg.RedisDB.Password,
		DB:       0,
	})

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		return nil, err
	}
	return &Redis{client: rdb}, nil
}

func (r Redis) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

func (r Redis) Get(ctx context.Context, key string) ([]byte, error) {
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r Redis) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r Redis) Close() error {
	return r.client.Close()
}
