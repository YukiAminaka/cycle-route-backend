package route

import (
	"context"

	"github.com/YukiAminaka/cycle-route-backend/internal/infrastructure/database/dbgen"
	"github.com/YukiAminaka/cycle-route-backend/internal/infrastructure/repository"
	"github.com/YukiAminaka/cycle-route-backend/internal/usecase/transaction"
)

type IDeleteRouteUsecase interface {
	DeleteRoute(ctx context.Context, routeID string) error
}

type deleteRouteUsecase struct {
	txManager transaction.TransactionManager
}

func NewDeleteRouteUsecase(txManager transaction.TransactionManager) IDeleteRouteUsecase {
	return &deleteRouteUsecase{
		txManager: txManager,
	}
}

func (u *deleteRouteUsecase) DeleteRoute(ctx context.Context, routeID string) error {
	
	// トランザクション内で削除操作を実行
	return u.txManager.RunInTransaction(ctx, func(q *dbgen.Queries) error {
		routeRepo := repository.NewRouteRepository(q)
		return routeRepo.DeleteRoute(ctx, routeID)
	})
}


