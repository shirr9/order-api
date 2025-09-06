package postgresql

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/shirr9/order-api/internal/config"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

type Connection interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Close()
}

type Storage struct {
	pool *pgxpool.Pool
	db   *bun.DB
}

func New(ctx context.Context, cfg *config.Config) (*Storage, error) {
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
	return &Storage{pool: pool, db: bun.NewDB(stdlib.OpenDBFromPool(pool), pgdialect.New())}, nil
}

func (s *Storage) NewPostgresRepository() *PostgresRepository {
	return &PostgresRepository{pool: s.pool, DB: s.db}
}

func (s *Storage) Close() {
	if s.pool != nil {
		s.pool.Close()
	}
	if s.db != nil {
		_ = s.db.Close()
	}
}
