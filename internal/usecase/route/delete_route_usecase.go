package route

import (
	"context"
	"errors"

	routeDomain "github.com/YukiAminaka/cycle-route-backend/internal/domain/route"
	"github.com/YukiAminaka/cycle-route-backend/internal/domain/user"
	"github.com/YukiAminaka/cycle-route-backend/internal/infrastructure/database/dbgen"
	"github.com/YukiAminaka/cycle-route-backend/internal/infrastructure/repository"
	"github.com/YukiAminaka/cycle-route-backend/internal/usecase/transaction"
)

type IDeleteRouteUsecase interface {
	DeleteRoute(ctx context.Context, routeID string, kratosID string) error
}

type deleteRouteUsecase struct {
	userRepository user.IUserRepository
	txManager      transaction.TransactionManager
	routeRepo      routeDomain.IRouteRepository
}

func NewDeleteRouteUsecase(userRepository user.IUserRepository, txManager transaction.TransactionManager, routeRepo routeDomain.IRouteRepository) IDeleteRouteUsecase {
	return &deleteRouteUsecase{
		userRepository: userRepository,
		txManager:      txManager,
		routeRepo:      routeRepo,
	}
}

func (u *deleteRouteUsecase) DeleteRoute(ctx context.Context, routeID string, kratosID string) error {
	// KratosIDからユーザー情報を取得
	userEntity, err := u.userRepository.GetUserByKratosID(ctx, kratosID)
	if err != nil {
		return err
	}

	// 既存のルートを取得
	route, err := u.routeRepo.GetRouteByID(ctx, routeID)
	if err != nil {
		return err
	}

	// 権限確認: ルートの所有者とリクエストのユーザーIDが一致するか
	if route.UserID() != userEntity.ID().String() {
		return errors.New("unauthorized: user does not own the route")
	}

	// トランザクション内で削除操作を実行
	return u.txManager.RunInTransaction(ctx, func(q *dbgen.Queries) error {
		routeRepo := repository.NewRouteRepository(q)
		return routeRepo.DeleteRoute(ctx, routeID)
	})
}
