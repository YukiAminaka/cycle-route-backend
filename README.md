# Cycle Route Backend

サイクリングルート管理アプリケーションのバックエンドAPI

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
│   ├── migrations/              # マイグレーションファイル
│   └── seeds/                   # シードデータ
├── .env                         # 環境変数
├── compose.yml                  # Docker Compose設定
├── go.mod
├── go.sum
└── sqlc.yaml                    # SQLC設定
```

## クリーンアーキテクチャの層

### 1. Domain層（内側）
- **責務**: ビジネスロジックの核となる部分
- **依存**: 他のどの層にも依存しない
- **内容**:
  - Entity: ビジネスルール、ドメインモデル
  - Repository Interface: データアクセスの抽象化

### 2. Usecase層
- **責務**: アプリケーション固有のビジネスロジック
- **依存**: Domain層のみに依存
- **内容**: ユースケースの実装、ビジネスフロー

### 3. Interface層
- **責務**: 外部とのインターフェース
- **依存**: UsecaseとDomain層に依存
- **内容**: HTTPハンドラー、プレゼンター、ミドルウェア

### 4. Infrastructure層（外側）
- **責務**: 技術的な実装詳細
- **依存**: すべての層に依存可能
- **内容**: DB接続、外部API、リポジトリ実装

**依存の方向**: Infrastructure → Interface → Usecase → Domain

## 開発環境のセットアップ

### 1. 環境変数の設定

`.env`ファイルを作成し、データベース接続情報を設定します。

```bash
DATABASE_URL=postgres://postgres:password@localhost:5432/postgres_db
SERVER_PORT=8080
```

### 2. データベースの起動

```bash
docker compose up -d
```

### 3. スキーマの適用

```bash
docker compose exec -T postgres psql -U postgres -d postgres_db < sqlc/schema.sql
```

### 4. SQLCでコード生成

スキーマやクエリを変更した後は、SQLCでコードを再生成します。

```bash
sqlc generate
```

### 5. アプリケーションの起動

```bash
go run cmd/api/main.go
```

## 開発ワークフロー

### スキーマ変更の手順

1. `sqlc/schema.sql` を編集
2. SQLCでコード生成: `sqlc generate`
3. Dockerのデータベースに反映: `docker compose exec -T postgres psql -U postgres -d postgres_db < sqlc/schema.sql`

### 新機能追加の手順

1. **Domain層**: エンティティとリポジトリインターフェースを定義
2. **Usecase層**: ビジネスロジックを実装
3. **Infrastructure層**: リポジトリの実装
4. **Interface層**: HTTPハンドラーを実装
5. **Router**: ルーティングを設定

## 技術スタック

- **言語**: Go 1.25.1
- **データベース**: PostgreSQL with PostGIS
- **ORマッパー**: sqlc
- **DB接続**: pgx/v5
- **地理情報処理**: paulmach/orb

## ライセンス

MIT
