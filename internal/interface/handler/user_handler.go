package handler

import (
	"github.com/YukiAminaka/cycle-route-backend/internal/usecase"
)

// UserHandler はユーザー関連のHTTPハンドラー
type UserHandler struct {
	userUsecase usecase.IUserUsecase
}

// NewUserHandler はUserHandlerを作成する
func NewUserHandler(userUsecase usecase.IUserUsecase) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
	}
}

// TODO: HTTPハンドラーメソッドを実装
// 例:
// func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) { ... }
// func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) { ... }
