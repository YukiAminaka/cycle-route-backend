package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// NewPostgresConnection はPostgreSQLへの接続を確立する
func NewPostgresConnection(ctx context.Context, databaseURL string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	// 接続確認
	if err := conn.Ping(ctx); err != nil {
		conn.Close(ctx)
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	return conn, nil
}
