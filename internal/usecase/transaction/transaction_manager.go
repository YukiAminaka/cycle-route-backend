package transaction

import (
	"context"

	"github.com/YukiAminaka/cycle-route-backend/internal/infrastructure/database/dbgen"
)

type TransactionManager interface {
	RunInTransaction(ctx context.Context, fn func(q *dbgen.Queries) error) error
}