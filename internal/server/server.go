package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/YukiAminaka/cycle-route-backend/config"
	"github.com/YukiAminaka/cycle-route-backend/internal/infrastructure/database/dbgen"
	"github.com/YukiAminaka/cycle-route-backend/internal/server/route"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Run(ctx context.Context, conf *config.Config, q *dbgen.Queries, pool *pgxpool.Pool) error {
	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}

	router.Use(cors.New(config))
	// Recovery ミドルウェアは panic が発生しても 500 エラーを返してくれる
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	route.InitRoute(router, q, pool)

	address := conf.Server.Address + ":" + conf.Server.Port
	log.Printf("Starting server on %s...\n", address)
	srv := &http.Server{
		Addr:              address,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       10 * time.Minute,
		WriteTimeout:      10 * time.Minute,
	}
	// サーバーを起動
	go func() {
		// srv.Shutdownが呼ばれるとhttp.ErrServerClosedを返すのでこれは無視する
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	//割り込み信号を待って、5秒のタイムアウトでサーバーを正常にシャットダウンします。
	quit := make(chan os.Signal, 1)
	// kill (パラメータなし) はデフォルトで syscall.SIGTERM を送信します
	// kill -2 は syscall.SIGINT です
	// kill -9 は syscall.SIGKILL ですが、捕捉できないため、追加する必要はありません
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}
	return nil
}