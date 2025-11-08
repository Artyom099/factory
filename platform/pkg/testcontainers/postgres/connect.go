package postgres

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/jackc/pgx/v5/pgxpool"
)

func connectPostgresClient(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, errors.Errorf("failed to connect to postgres: %v", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, errors.Errorf("failed to ping postgres: %v", err)
	}

	return pool, nil
}
