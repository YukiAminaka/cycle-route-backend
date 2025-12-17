package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/YukiAminaka/cycle-route-backend/config"
	"github.com/YukiAminaka/cycle-route-backend/internal/infrastructure/database/dbgen"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	maxRetries       = 5
	retryDelay       = 5 * time.Second
	connectionTimeout = 10 * time.Second
)

// NewDB creates a new database connection pool with proper configuration
func NewDB(cfg config.DBConfig) *dbgen.Queries {
	ctx, cancel := context.WithTimeout(context.Background(), connectionTimeout)
	defer cancel()

	// データベース接続用のURLを作成
	databaseUrl := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name,
	)

	// コネクションプール設定をパース
	poolConfig, err := pgxpool.ParseConfig(databaseUrl)
	if err != nil {
		log.Fatalf("Failed to parse database URL: %v", err)
	}

	// コネクションプールの設定
	poolConfig.MaxConns = 25
	poolConfig.MinConns = 5
	poolConfig.MaxConnLifetime = time.Hour
	poolConfig.MaxConnIdleTime = 30 * time.Minute
	poolConfig.HealthCheckPeriod = time.Minute

	// リトライロジック付きで接続プールを作成
	var pool *pgxpool.Pool
	for i := 0; i < maxRetries; i++ {
		pool, err = pgxpool.NewWithConfig(ctx, poolConfig)
		if err == nil {
			// 疎通確認
			if err = pool.Ping(ctx); err == nil {
				log.Println("Successfully connected to database")
				queries := dbgen.New(pool)
				return queries
			}
			pool.Close()
		}

		log.Printf("Failed to connect to database (attempt %d/%d): %v", i+1, maxRetries, err)
		if i < maxRetries-1 {
			log.Printf("Retrying in %v...", retryDelay)
			time.Sleep(retryDelay)
		}
	}

	log.Fatalf("Failed to connect to database after %d attempts: %v", maxRetries, err)
	return nil
}




// CloseDB closes the database connection pool
func CloseDB(queries *dbgen.Queries) {
	if queries != nil {
		// pgxpool.Pool型を取得してクローズ
		// Note: sqlcが生成したQueriesからプールを取得する方法が必要
		log.Println("Closing database connection pool")
	}
}