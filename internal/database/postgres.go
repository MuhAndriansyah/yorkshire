package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresqlPool(ctx context.Context, databaseUrl string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(databaseUrl)

	if err != nil {
		return nil, fmt.Errorf("failed to parse configuration: %w", err)
	}

	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnIdleTime = 30 * time.Minute
	config.MaxConnLifetime = time.Hour
	config.HealthCheckPeriod = time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, config)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	ctxPing, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := pool.Ping(ctxPing); err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	return pool, nil
}
