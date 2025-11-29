package handler

import (
	userUsecase "github.com/YukiAminaka/cycle-route-backend/internal/usecase/user"
)

// UserHandler はユーザー関連のHTTPハンドラー
type UserHandler struct {
	createUserUsecase userUsecase.ICreateUserUsecase
	getUserUsecase    userUsecase.IGetUserByIDUsecase
}

// NewUserHandler はUserHandlerを作成する
func NewUserHandler(
	createUserUsecase userUsecase.ICreateUserUsecase,
	getUserUsecase userUsecase.IGetUserByIDUsecase,
) *UserHandler {
	return &UserHandler{
		createUserUsecase: createUserUsecase,
		getUserUsecase:    getUserUsecase,
	}
}

// TODO: HTTPハンドラーメソッドを実装
// 例:
// func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) { ... }
// func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) { ... }
