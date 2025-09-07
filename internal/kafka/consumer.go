package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"

	"github.com/segmentio/kafka-go"
	"github.com/shirr9/order-api/internal/config"
	"github.com/shirr9/order-api/internal/order"
)

type OrderSaver interface {
	AddOrder(ctx context.Context, o order.Order) error
}

type Consumer struct {
	reader *kafka.Reader
	logger *slog.Logger
	saver  OrderSaver
}

func NewConsumer(cfg config.KafkaConfig, logger *slog.Logger, saver OrderSaver) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  cfg.Brokers,
		Topic:    cfg.Topic,
		GroupID:  "order-api-group",
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})

	return &Consumer{
		reader: r,
		logger: logger,
		saver:  saver,
	}
}

func (c *Consumer) Run(ctx context.Context) {
	c.logger.Info("kafka consumer started", slog.String("topic", c.reader.Config().Topic))
	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				c.logger.Info("context canceled, stopping kafka consumer.")
				break
			}
			c.logger.Error("failed to read message from kafka", slog.String("err", err.Error()))
			continue
		}

		var orderData order.Order
		if err := json.Unmarshal(m.Value, &orderData); err != nil {
			c.logger.Error("failed to unmarshal order, skipping message", slog.String("err", err.Error()))
			continue
		}

		c.logger.Info("received an order", slog.String("order_uid", orderData.OrderUID))
		if err := c.saver.AddOrder(ctx, orderData); err != nil {
			c.logger.Error("failed to save order", slog.String("err", err.Error()),
				slog.String("order_uid", orderData.OrderUID))
		}
	}
}

func (c *Consumer) Close() error {
	c.logger.Info("closing kafka consumer")
	return c.reader.Close()
}
