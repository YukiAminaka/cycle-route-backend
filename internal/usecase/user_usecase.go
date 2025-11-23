package usecase

import (
	"context"

	"github.com/YukiAminaka/cycle-route-backend/internal/domain/entity"
	"github.com/YukiAminaka/cycle-route-backend/internal/domain/repository"
)

// IUserUsecase はユーザーユースケースのインターフェース
type IUserUsecase interface {
	GetUserByID(ctx context.Context, id int64) (*entity.User, error)
	CreateUser(ctx context.Context, user *entity.User) (*entity.User, error)
}

type userUsecase struct {
	userRepo repository.IUserRepository
}

// NewUserUsecase はユーザーユースケースを作成する
func NewUserUsecase(userRepo repository.IUserRepository) IUserUsecase {
	return &userUsecase{
		userRepo: userRepo,
	}
}

func (u *userUsecase) GetUserByID(ctx context.Context, id int64) (*entity.User, error) {
	return u.userRepo.GetUserByID(ctx, id)
}

func (u *userUsecase) CreateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	// TODO: ビジネスロジックの実装（バリデーション、追加処理など）
	return u.userRepo.CreateUser(ctx, user)
}
