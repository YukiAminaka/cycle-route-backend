package router

import (
	"net/http"

	"github.com/YukiAminaka/cycle-route-backend/internal/interface/handler"
	"github.com/YukiAminaka/cycle-route-backend/internal/interface/middleware"
)

// SetupRouter はルーティングを設定する
func SetupRouter(userHandler *handler.UserHandler) http.Handler {
	mux := http.NewServeMux()

	// TODO: ルートを追加
	// 例:
	// mux.HandleFunc("/api/users", userHandler.GetUsers)
	// mux.HandleFunc("/api/users/{id}", userHandler.GetUser)

	// ミドルウェアを適用
	return middleware.Logger(mux)
}
