package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	DSN          string
	MaxIdleTime  string
	MaxIdleConns int
	MaxOpenConns int
}

func NewConn(c *Config) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(c.DSN)
	if err != nil {
		return nil, err
	}

	duration, err := time.ParseDuration(c.MaxIdleTime)
	if err != nil {
		return nil, err
	}

	poolConfig.MaxConnIdleTime = duration
	poolConfig.MaxConns = int32(c.MaxOpenConns)
	poolConfig.MinConns = int32(c.MaxIdleConns)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to start connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("database unreachable %w", err)
	}
	return pool, nil
}
