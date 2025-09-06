package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/shirr9/order-api/internal/config"
)

type Redis struct {
	cache *redis.Client
}

func NewRedis(ctx *context.Context, cfg *config.Config) *Redis {

}
