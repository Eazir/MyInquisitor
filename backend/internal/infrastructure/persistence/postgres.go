package persistence

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/myinquisitor/backend/internal/infrastructure/persistence/sqlc"
)

type PostgresDB struct {
	Pool    *pgxpool.Pool
	Queries *sqlc.Queries
}

func NewPostgresDB(ctx context.Context, databaseURL string) (*PostgresDB, error) {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	return &PostgresDB{
		Pool:    pool,
		Queries: sqlc.New(pool),
	}, nil
}

func (db *PostgresDB) Close() {
	db.Pool.Close()
}
