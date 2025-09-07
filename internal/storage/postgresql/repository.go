package postgresql

import (
	"context"
	"github.com/shirr9/order-api/internal/order"
	"github.com/uptrace/bun"
)

type Repository interface {
	AddOrder(ctx context.Context, o order.Order) error
	FindOrderById(ctx context.Context, id string) (*order.Order, error)
	Close()
}

type PostgresRepository struct {
	Pool Connection
	DB   *bun.DB
}

func (r *PostgresRepository) AddOrder(ctx context.Context, o order.Order) error {
	return r.DB.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewInsert().Model(&o).Exec(ctx); err != nil {
			return err
		}
		if o.Delivery != nil {
			o.Delivery.OrderUID = o.OrderUID
			if _, err := tx.NewInsert().Model(o.Delivery).Exec(ctx); err != nil {
				return err
			}
		}
		if o.Payment != nil {
			o.Payment.OrderUID = o.OrderUID
			if _, err := tx.NewInsert().Model(o.Payment).Exec(ctx); err != nil {
				return err
			}
		}
		if len(o.Items) > 0 {
			for i := range o.Items {
				o.Items[i].OrderUID = o.OrderUID
			}
			if _, err := tx.NewInsert().Model(&o.Items).Exec(ctx); err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *PostgresRepository) FindOrderById(ctx context.Context, id string) (*order.Order, error) {
	var o order.Order
	if err := r.DB.
		NewSelect().
		Model(&o).
		Relation("Delivery").
		Relation("Payment").
		Relation("Items").
		Where("o.order_uid = ?", id).
		Scan(ctx); err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *PostgresRepository) Close() {
	if r.Pool != nil {
		r.Pool.Close()
	}
	if r.DB != nil {
		_ = r.DB.Close()
	}
}
