package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func New(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("database: erro ao configurar o pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("database: erro ao conectar no banco: %w", err)
	}

	return pool, nil
}

func WithTx(ctx context.Context, pool *pgxpool.Pool, work func(pgx.Tx) error) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("database: erro ao iniciar transação: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := work(tx); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("database: erro ao confirmar transação: %w", err)
	}
	return nil
}
