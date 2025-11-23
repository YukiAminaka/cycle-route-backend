package repository

import (
	"context"

	"github.com/YukiAminaka/cycle-route-backend/internal/domain/entity"
	"github.com/YukiAminaka/cycle-route-backend/internal/domain/repository"
	"github.com/YukiAminaka/cycle-route-backend/internal/infrastructure/database/sqlc"
)

type userRepositoryImpl struct {
	queries *sqlc.Queries
}

// NewUserRepository はユーザーリポジトリの実装を作成する
func NewUserRepository(queries *sqlc.Queries) repository.IUserRepository {
	return &userRepositoryImpl{queries: queries}
}

func (r *userRepositoryImpl) GetUserByID(ctx context.Context, id int64) (*entity.User, error) {
	// TODO: 実装
	return nil, nil
}

func (r *userRepositoryImpl) CreateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	// TODO: 実装
	return nil, nil
}
