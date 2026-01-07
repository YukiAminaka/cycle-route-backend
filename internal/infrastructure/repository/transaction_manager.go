package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/YukiAminaka/cycle-route-backend/internal/infrastructure/database/dbgen"
	"github.com/YukiAminaka/cycle-route-backend/internal/usecase/transaction"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionManager struct{
	pool *pgxpool.Pool
	queries *dbgen.Queries
}

func NewTransactionManager(queries *dbgen.Queries, pool *pgxpool.Pool) transaction.TransactionManager {
	return &TransactionManager{queries: queries, pool: pool}
}

func (tm *TransactionManager) RunInTransaction(ctx context.Context, fn func(queries *dbgen.Queries) error) error {
	// トランザクションを開始
	tx, err := tm.pool.Begin(ctx)
	if err != nil {
		return err
	}

	// トランザクション用のQueriesを作成
	q := tm.queries.WithTx(tx)


	// トランザクション内の関数を実行
	err = fn(q)
	if err != nil {
		// トランザクション内でエラーが発生したらロールバック
		log.Printf("db rollback: %v\n", err)
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	// 全ての関数が成功したらコミット
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit failed: %w", err)
	}

	return nil
}