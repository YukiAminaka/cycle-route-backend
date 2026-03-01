FROM golang:1.26 AS builder

WORKDIR /app

# 依存関係キャッシュを効かせるために先に go.mod だけコピー
COPY go.mod go.sum ./

RUN go mod download

# アプリコードをコピー
COPY . .

# 本番用バイナリをビルド
RUN CGO_ENABLED=0 GOOS=linux go build -o server



FROM gcr.io/distroless/static-debian12

WORKDIR /app

# distrolessでは65532が非rootユーザーとして推奨されている
USER 65532

# builderからビルド成果物だけコピー
COPY --from=builder /app/server .

ENV GIN_MODE="release"

# メインコマンド（固定）
ENTRYPOINT ["./server"]