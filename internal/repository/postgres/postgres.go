package postgres

import (
	"context"
	"fmt"
	"time"

	"effective-mobile-task/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	
	connectTimeout = 5 * time.Second
	pingTimeout    = 5 * time.Second
)


func NewPostgresDB(cfg config.PostgresConfig) (*pgxpool.Pool, error) {

	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout)
	defer cancel()


	pool, err := pgxpool.New(ctx, cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к postgres: %w", err)
	}


	if err := pingDatabase(pool); err != nil {
		pool.Close() 
		return nil, err
	}

	return pool, nil
}


func pingDatabase(pool *pgxpool.Pool) error {
	
	ctx, cancel := context.WithTimeout(context.Background(), pingTimeout)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		return fmt.Errorf("не удалось пинговать postgres: %w", err)
	}

	return nil
}
