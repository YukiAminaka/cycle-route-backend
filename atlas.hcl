// Atlas configuration for cycle-route-backend


env "dev" {
  // データベース接続URL
  url = "postgres://postgres:password@localhost:5432/postgres_db?sslmode=disable&search_path=public"

  // スキーマの定義元（理想の状態）
  src = "file://internal/infrastructure/database/sqlc/schema.sql"

  // マイグレーションファイルの保存先
  migration {
    dir = "file://db/migrations"
  }

  // 開発用データベース（diff計算用の一時DB）
  dev = "docker://postgis/postgis/18-3.6/dev"

  // スキーマ名
  schemas = ["public"]

  // マイグレーションのフォーマット
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}
