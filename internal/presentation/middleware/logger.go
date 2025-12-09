package middleware

import (
	"log"
	"net/http"
	"time"
)

// Logger はHTTPリクエストをログに記録するミドルウェア
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// リクエストを処理
		next.ServeHTTP(w, r)

		// ログ出力
		log.Printf(
			"%s %s %s",
			r.Method,
			r.RequestURI,
			time.Since(start),
		)
	})
}
