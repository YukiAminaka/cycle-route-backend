package user

import (
	"context"
)

// IUserRepository はユーザーリポジトリのインターフェース
type IUserRepository interface {
	GetUserByID(ctx context.Context, id string) (*User, error)
	CreateUser(ctx context.Context, user *User) (*User, error)
}
