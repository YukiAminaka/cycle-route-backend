package repository

import (
	"context"

	"github.com/YukiAminaka/cycle-route-backend/internal/domain/entity"
)

// IUserRepository はユーザーリポジトリのインターフェース
type IUserRepository interface {
	GetUserByID(ctx context.Context, id int64) (*entity.User, error)
	CreateUser(ctx context.Context, user *entity.User) (*entity.User, error)
}
