package postgresql

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shirr9/order-api/internal/config"
)

type Connection interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Close()
}

// mb will add
//type Storage struct {
//	pool *pgxpool.Pool
//}

func NewConnection(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	// dbdriver://username:password@host:port/dbname?param1=true&param2=false
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DbName, cfg.SSlMode)
	//safeDSN := fmt.Sprintf("postgresql://%s:***@%s:%s/%s?sslmode=%s",
	//	cfg.Username, cfg.Host, cfg.Port, cfg.DbName, cfg.SSlMode)
	// log with safeDSN
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}
	if e := pool.Ping(ctx); e != nil {
		return nil, fmt.Errorf("failed to ping database: %w", e)
	}
	// log: successfully connected to database
	return pool, nil
}
