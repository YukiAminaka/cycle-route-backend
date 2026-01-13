# Cycle Route Backend

サイクリングルート管理アプリケーションのバックエンド API

## プロジェクト構成

このプロジェクトはクリーンアーキテクチャに基づいて設計されています。

```
cycle-route-backend/
├── cmd/
│   └── api/
│       └── main.go              # アプリケーションのエントリーポイント
├── internal/
│   ├── domain/                  # Domain層（最も内側・ビジネスロジックの核）
│   │   ├── entity/              # エンティティ（ドメインモデル）
│   │   │   └── user.go
│   │   └── repository/          # リポジトリインターフェース定義
│   │       └── user_repository.go
│   ├── usecase/                 # Usecase層（アプリケーションのビジネスロジック）
│   │   └── user_usecase.go
│   ├── interface/               # Interface層（外部とのやり取り）
│   │   ├── handler/             # HTTPハンドラー
│   │   │   └── user_handler.go
│   │   ├── middleware/          # HTTPミドルウェア
│   │   │   └── logger.go
│   │   └── presenter/           # レスポンス整形
│   │       └── response.go
│   └── infrastructure/          # Infrastructure層（最も外側・技術的詳細）
│       ├── database/
│       │   ├── postgres.go      # DB接続管理
│       │   └── sqlc/            # SQLCで生成されたコード
│       │       ├── db.go
│       │       ├── models.go
│       │       ├── query.sql.go
│       │       └── custom_types.go
│       ├── repository/          # リポジトリの実装
│       │   └── user_repository_impl.go
│       └── router/
│           └── router.go        # ルーティング設定
├── config/                      # 設定管理
│   └── config.go
├── sqlc/                        # SQL定義ファイル
│   ├── schema.sql
│   └── query.sql
├── db/                          # データベース関連
│   ├── migrations/              # Atlasマイグレーションファイル
│   └── seeds/                   # シードデータ
├── .env                         # 環境変数
├── atlas.hcl                    # Atlas設定ファイル
├── compose.yml                  # Docker Compose設定
├── go.mod
├── go.sum
└── sqlc.yaml                    # SQLC設定
```

## クリーンアーキテクチャの層

### 1. Domain 層（内側）

- **責務**: ビジネスロジックの核となる部分
- **依存**: 他のどの層にも依存しない
- **内容**:
  - Entity: ビジネスルール、ドメインモデル
  - Repository Interface: データアクセスの抽象化

### 2. Usecase 層

- **責務**: アプリケーション固有のビジネスロジック
- **依存**: Domain 層のみに依存
- **内容**: ユースケースの実装、ビジネスフロー

### 3. Interface 層

- **責務**: 外部とのインターフェース
- **依存**: Usecase と Domain 層に依存
- **内容**: HTTP ハンドラー、プレゼンター、ミドルウェア

### 4. Infrastructure 層（外側）

- **責務**: 技術的な実装詳細
- **依存**: すべての層に依存可能
- **内容**: DB 接続、外部 API、リポジトリ実装

**依存の方向**: Infrastructure → Interface → Usecase → Domain

## 開発環境のセットアップ

### 1. 環境変数の設定

`.env`ファイルを作成し、データベース接続情報を設定します。

### 2. アプリケーションの起動

```bash
docker compose up -d
# or
GO_ENV=dev go run cmd/api/main.go
```

### 3. スキーマの適用

```bash
atlas migrate apply --env dev
```

### 4. シードデータ投入

```
docker compose exec -T postgres psql -U postgres -d postgres_db < db/seeds/dev_seed.sql
```

### 5. SQLC でコード生成

スキーマやクエリを変更した後は、SQLC でコードを再生成します。

```bash
sqlc generate
```

## テストの実行

```bash
go test ./...
```

## 認証が必要な API のテスト

このプロジェクトでは Ory Kratos を使用した認証を実装しています。認証が必要なエンドポイントをテストするには、セッションクッキーが必要です。

### 簡単な方法: Admin API でテストユーザーを作成する

Kratos Admin API を使用してテストユーザーを作成し、セッションクッキーを取得するスクリプトを用意しています。

```bash
# デフォルト（ランダムなメールアドレス）
./scripts/create-test-user.sh

# カスタムメールアドレスとパスワード
./scripts/create-test-user.sh your@email.com YourPassword123!

# 環境変数でKratosのURLを指定
KRATOS_ADMIN_URL=http://127.0.0.1:4434 ./scripts/create-test-user.sh
```

このスクリプトは以下の処理を自動で行います：

1. Kratos Admin API を使ってユーザーを作成（既存の場合はスキップ）
2. 作成したユーザーでログイン
3. セッショントークンとテストユーザーの情報を表示

### 取得したクッキーの使い方

#### curl で API を呼び出す

```bash
# セッショントークンを使用
curl -H 'Cookie: ory_kratos_session=YOUR_SESSION_TOKEN' \
  http://localhost:8080/api/v1/users/USER_ID

# ルートを作成
curl -H 'Cookie: ory_kratos_session=YOUR_SESSION_TOKEN' \
  http://localhost:8080/api/v1/routes \
  -X POST \
  -H 'Content-Type: application/json' \
  -d '{"name":"Test Route",...}'
```

#### Swagger UI で使用

1. `./scripts/create-test-user.sh` を実行
2. 表示された **Session Token** をコピー
3. Swagger UI（http://localhost:8080/api/v1/swagger/index.html）を開く
4. 右上の「Authorize」ボタン（🔒 アイコン）をクリック
5. `CookieAuth` の欄にセッショントークンを貼り付け
6. 「Authorize」をクリックして「Close」

これで認証が必要なエンドポイント（🔒 マーク付き）を Swagger UI から試せます。

### 手動でクッキーを取得する場合

<details>
<summary>クリックして展開</summary>

```bash
# ログインフローを開始
FLOW=$(curl -s 'http://127.0.0.1:4433/self-service/login/api' -c cookies.txt)
FLOW_ID=$(echo $FLOW | jq -r '.id')
CSRF_TOKEN=$(echo $FLOW | jq -r '.ui.nodes[] | select(.attributes.name=="csrf_token") | .attributes.value')

# ログイン
curl -X POST "http://127.0.0.1:4433/self-service/login?flow=$FLOW_ID" \
  -H 'Content-Type: application/json' \
  -b cookies.txt \
  -c cookies.txt \
  -d '{
    "method": "password",
    "csrf_token": "'$CSRF_TOKEN'",
    "identifier": "test@example.com",
    "password": "testpassword123"
  }'

# クッキーを確認
cat cookies.txt | grep ory_kratos_session
```

</details>

## API ドキュメント

ブラウザで以下の URL にアクセスすると、Swagger UI が表示され、API のドキュメントを確認できます。

```
http://localhost:8080/api/v1/swagger/index.html
```

### API ドキュメントの生成

- 整形

```
swag fmt
```

- 生成（エントリポイント指定）

```
swag init -g ./cmd/api/main.go
```

## 開発ワークフロー

### Atlas を使ったマイグレーション管理

このプロジェクトでは[Atlas](https://atlasgo.io/)を使用してデータベースマイグレーションを管理します。

#### Atlas のインストール

```bash
# Linux/macOS
curl -sSf https://atlasgo.sh | sh

# または Go経由でインストール
go install ariga.io/atlas/cmd/atlas@latest
```

#### スキーマ変更の基本フロー（atlas.hcl を使用）

```bash
# 1. sqlc/schema.sql を編集
vim sqlc/schema.sql

# 2. マイグレーションファイルを自動生成
atlas migrate diff migration_name --env dev

# 3. 生成されたマイグレーションを確認
cat db/migrations/[最新のファイル].sql

# 4. マイグレーションを適用
atlas migrate apply --env dev

# 5. SQLCでGoコードを生成
sqlc generate
```

**`atlas.hcl`を使うことで**、長いコマンドが `--env dev` だけで済む

#### 便利な Atlas コマンド

```bash
# マイグレーション状態の確認
atlas migrate status --env dev

# スキーマの差分を確認（マイグレーション生成前にチェック）
atlas schema diff --env dev

# 現在のデータベーススキーマを表示
atlas schema inspect --env dev

# Dry run（実際には適用せずに確認）
atlas migrate apply --env dev --dry-run

# 特定のバージョンまでマイグレーション
atlas migrate apply --env dev --to 20240101000001
```

#### マイグレーションファイルの管理

- マイグレーションファイルは `db/migrations/` に自動生成されます
- ファイル名形式: `20240101000001_migration_name.sql`
- Atlas が自動的にバージョン管理とチェックサムを管理します

### 新機能追加の手順

1. **Domain 層**: エンティティとリポジトリインターフェースを定義
2. **Usecase 層**: ビジネスロジックを実装
3. **Infrastructure 層**: リポジトリの実装
4. **Interface 層**: HTTP ハンドラーを実装
5. **Router**: ルーティングを設定

## 技術スタック

- **言語**: Go 1.25.1
- **データベース**: PostgreSQL with PostGIS
- **マイグレーション**: Atlas
- **OR マッパー**: sqlc
- **DB 接続**: pgx/v5
- **地理情報処理**: paulmach/orb
