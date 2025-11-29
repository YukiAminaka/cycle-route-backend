package user

import (
	"context"
)

// IUserRepository はユーザーリポジトリのインターフェース
type IUserRepository interface {
	GetUserByID(ctx context.Context, id int64) (*User, error)
	CreateUser(ctx context.Context, user *User) (*User, error)
}
