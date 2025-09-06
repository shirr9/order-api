package service

import (
	"context"
	"github.com/shirr9/order-api/internal/order"
	"github.com/shirr9/order-api/internal/storage/postgresql"
	"log/slog"
)

type Service struct {
	repo postgresql.Repository
	log  *slog.Logger
}

func NewService(repo postgresql.Repository, logger *slog.Logger) *Service {
	return &Service{repo: repo, log: logger}
}

func (s *Service) AddOrder(ctx context.Context, o order.Order) error {
	if err := s.repo.AddOrder(ctx, o); err != nil {
		s.log.LogAttrs(
			ctx,
			slog.LevelError,
			"can't add order",
			slog.String("order_uid", o.OrderUID),
			slog.Any("err", err),
		)
		return err
	}
	s.log.LogAttrs(
		ctx,
		slog.LevelInfo,
		"order created",
		slog.String("order_uid", o.OrderUID),
	)
	return nil
}

func (s *Service) FindOrderById(ctx context.Context, id string) (*order.Order, error) {
	o, err := s.repo.FindOrderById(ctx, id)
	if err != nil {
		s.log.LogAttrs(
			ctx,
			slog.LevelError,
			"can't find order",
			slog.String("order_uid", id),
			slog.Any("err", err),
		)
		return nil, err
	}
	s.log.LogAttrs(
		ctx,
		slog.LevelInfo,
		"order found",
		slog.String("order_uid", id),
	)
	return o, nil
}
