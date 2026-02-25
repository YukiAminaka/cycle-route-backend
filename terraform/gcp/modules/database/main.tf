resource "google_sql_database_instance" "main" {
  name             = "${var.project_name}-${var.environment}-db"
  database_version = "POSTGRES_16"
  region           = var.region

  settings {
    tier = var.instance_tier # マシンタイプ

    ip_configuration {
      ipv4_enabled = true # このCloud SQLインスタンスにパブリックIPv4アドレスを割り当てるかどうか
    }

    database_flags {
      name  = "log_min_duration_statement"
      value = "1000"
    }

    backup_configuration {
      enabled                        = true
      start_time                     = "03:00"
      point_in_time_recovery_enabled = true

      backup_retention_settings {
        retained_backups = 7
      }
    }

    insights_config {
      query_insights_enabled  = true
      record_application_tags = true
      record_client_address   = false
    }
    # インスタンスが自動的に再起動してアップデートを適用できる1時間のメンテナンスウィンドウを宣言
    maintenance_window {
      day  = 7 # Sunday
      hour = 3
    }
  }

  deletion_protection = true
}

resource "google_sql_database" "main" {
  name     = var.project_name
  instance = google_sql_database_instance.main.name
}

resource "google_sql_user" "main" {
  name     = var.project_name
  instance = google_sql_database_instance.main.name
  password = var.db_password
}
