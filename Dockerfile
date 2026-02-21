FROM golang:1.25-alpine

WORKDIR /app

# 非rootユーザーを作成
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser && \
    chown -R appuser:appuser /app

# 非rootユーザーに切り替え
USER appuser

# GOPATHをappuserのホームディレクトリに設定
ENV GOPATH=/home/appuser/go
ENV PATH=$PATH:$GOPATH/bin

# ホットリロードツールAirをインストール（appuserとして）
RUN go install github.com/air-verse/air@latest

# 依存関係ファイルをコピーしてキャッシュ
COPY --chown=appuser:appuser go.mod go.sum ./
RUN go mod download

# 環境変数を設定
ENV GO_ENV=dev

# コンテナ起動時にAirを実行する
CMD ["air", "-c", ".air.toml"]