package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"github.com/shirr9/order-api/internal/order"
	"github.com/shirr9/order-api/internal/storage/cache"
	"github.com/shirr9/order-api/internal/storage/postgresql"
	"log/slog"
)

type Service struct {
	repo  postgresql.Repository
	log   *slog.Logger
	cache cache.CacheRepository
}

func NewService(repo postgresql.Repository, logger *slog.Logger, cache cache.CacheRepository) *Service {
	return &Service{repo: repo, log: logger, cache: cache}
}

func (s *Service) AddOrder(ctx context.Context, o order.Order) error {
	// add to database
	if err := s.repo.AddOrder(ctx, o); err != nil {
		s.log.LogAttrs(
			ctx, slog.LevelError, "can't add order", slog.String("order_uid", o.OrderUID),
			slog.Any("err", err))
		return err
	}
	s.log.LogAttrs(ctx, slog.LevelInfo, "add order to database", slog.String("order_uid", o.OrderUID))
	return nil
}

func (s *Service) FindOrderById(ctx context.Context, id string) (*order.Order, error) {
	// check data in cache
	data, err := s.cache.Get(ctx, id)
	if errors.Is(err, redis.Nil) { // key does not exist
		s.log.LogAttrs(ctx, slog.LevelInfo, "cache doesn't contain order", slog.String("order_uid", id))
		// find in repo
		o, err := s.repo.FindOrderById(ctx, id)
		if err != nil {
			s.log.LogAttrs(
				ctx, slog.LevelError, "finding order in database failed", slog.String("order_uid", id),
				slog.Any("err", err))
			return nil, err
		}
		s.log.LogAttrs(ctx, slog.LevelInfo, "order found in database", slog.String("order_uid", id))

		// add to cache
		jData, err := json.Marshal(o)
		if err != nil {
			s.log.LogAttrs(ctx, slog.LevelError, "failed to marshal data to JSON", slog.String("order_uid", id),
				slog.Any("err", err))
			return nil, err
		}
		err = s.cache.Set(ctx, id, jData)
		if err != nil {
			s.log.LogAttrs(ctx, slog.LevelError, "saving in cache failed", slog.String("order_uid", id),
				slog.Any("err", err))
			return nil, err
		}
		s.log.LogAttrs(ctx, slog.LevelInfo, "order was saved in cache", slog.String("order_uid", id))

		return o, nil

	} else if err != nil {
		s.log.LogAttrs(
			ctx, slog.LevelError, "Get from cache failed", slog.String("order_uid", id),
			slog.Any("err", err),
		)
		return nil, err
	}

	if len(data) == 0 {
		s.log.LogAttrs(ctx, slog.LevelInfo, "value is empty", slog.String("order_uid", id))
	}

	var o order.Order
	if err := json.Unmarshal(data, &o); err != nil {
		s.log.LogAttrs(ctx, slog.LevelError, "unmarshal failed", slog.String("order_uid", id), slog.Any("err", err))
	}
	s.log.LogAttrs(ctx, slog.LevelInfo, "order was found in cache", slog.String("order_uid", id))
	return &o, nil
}
