FROM golang:1.25-alpine

WORKDIR /app

# ホットリロードツール Air をインストール
RUN go install github.com/air-verse/air@latest

# 依存関係ファイルをコピーしてキャッシュ
COPY go.mod go.sum ./
RUN go mod download

# 環境変数を設定
ENV GO_ENV=dev

# コンテナ起動時に Air を実行する
CMD ["air", "-c", ".air.toml"]