// Atlas configuration for cycle-route-backend

variable "db_url" {
  type    = string
  default = getenv("DATABASE_URL")
}

env "dev" {
  // データベース接続URL
  url = var.db_url

  // スキーマの定義元（理想の状態）
  src = "file://sqlc/schema.sql"

  // マイグレーションファイルの保存先
  migration {
    dir = "file://db/migrations"
  }

  // 開発用データベース（diff計算用の一時DB）
  dev = "docker://postgres/15/dev?search_path=public"

  // スキーマ名
  schemas = ["public"]

  // マイグレーションのフォーマット
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}
