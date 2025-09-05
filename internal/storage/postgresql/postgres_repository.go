package postgresql

import (
	"fmt"
	"github.com/shirr9/order-api/internal/order"
)

type Repository interface {
	AddOrder(o order.Order) error
	FindOrderById(id string) order.Order
	Close()
}

type PostgresRepository struct {
	conn Connection
}

func NewPostgresRepository(c Connection) (*PostgresRepository, error) {
	if c == nil {
		return nil, fmt.Errorf("connection is nil")
	}
	return &PostgresRepository{conn: c}, nil
}

func (r *PostgresRepository) Close() {
	if r.conn != nil {
		r.conn.Close()
		// log
	}
}

func (r *PostgresRepository) AddOrder(o order.Order) error {
	q := `WITH 
`
}

func (r *PostgresRepository) FindOrderById(id string) order.Order {
	//TODO implement me
	panic("implement me")
}
