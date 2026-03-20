# ============================================================
# Kratos DSN secret 
# ============================================================

resource "google_secret_manager_secret" "kratos_dsn" {
  secret_id = "${var.project_name}-${var.environment}-kratos-dsn"

  replication {
    auto {}
  }

  labels = {
    project     = var.project_name
    environment = var.environment
  }
}

resource "google_secret_manager_secret_version" "kratos_dsn" {
  secret      = google_secret_manager_secret.kratos_dsn.id
  secret_data = "postgres://${var.db_user}:${var.db_password}@/${var.db_name}?host=/cloudsql/${var.db_connection_name}&sslmode=disable&max_conns=20&max_idle_conns=4"
}

# ============================================================
# API DSN secret
# ============================================================

resource "google_secret_manager_secret" "api_dsn" {
  secret_id = "${var.project_name}-${var.environment}-api-dsn"

  replication {
    auto {}
  }

  labels = {
    project     = var.project_name
    environment = var.environment
  }
}

resource "google_secret_manager_secret_version" "api_dsn" {
  secret      = google_secret_manager_secret.api_dsn.id
  secret_data = "postgres://${var.db_user}:${var.db_password}@/${var.db_name}?host=/cloudsql/${var.db_connection_name}&sslmode=disable"
}

# ============================================================
# Service Accounts
# ============================================================

# Cloud Runのサービスアカウントを作成
resource "google_service_account" "frontend" {
  account_id   = "${var.project_name}-${var.environment}-frontend"
  display_name = "${var.project_name} Frontend Service Account"
}

resource "google_service_account" "api" {
  account_id   = "${var.project_name}-${var.environment}-api"
  display_name = "${var.project_name} API Service Account"
}

resource "google_service_account" "kratos" {
  account_id   = "${var.project_name}-${var.environment}-kratos"
  display_name = "${var.project_name} Kratos Service Account"
}

# ============================================================
# Cloud SQL IAM (Cloud SQL Auth Proxy via volume mount)
# ============================================================

# Cloud SQL インスタンスに接続する場合はroles/cloudsql.clientが必要
resource "google_project_iam_member" "api_cloudsql_client" {
  project = var.project_id
  role    = "roles/cloudsql.client"
  member  = "serviceAccount:${google_service_account.api.email}"
}

resource "google_project_iam_member" "kratos_cloudsql_client" {
  project = var.project_id
  role    = "roles/cloudsql.client"
  member  = "serviceAccount:${google_service_account.kratos.email}"
}

# ============================================================
# Secret Manager IAM
# ============================================================

# Secret ManagerにアクセスできるIAMの設定
resource "google_secret_manager_secret_iam_member" "api_dsn" {
  secret_id = google_secret_manager_secret.api_dsn.secret_id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${google_service_account.api.email}"
}

resource "google_secret_manager_secret_iam_member" "kratos_db_password" {
  secret_id = var.db_password_secret_id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${google_service_account.kratos.email}"
}

resource "google_secret_manager_secret_iam_member" "kratos_cookie_secret" {
  secret_id = var.kratos_cookie_secret_id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${google_service_account.kratos.email}"
}

resource "google_secret_manager_secret_iam_member" "kratos_cipher_secret" {
  secret_id = var.kratos_cipher_secret_id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${google_service_account.kratos.email}"
}

resource "google_secret_manager_secret_iam_member" "kratos_smtp" {
  secret_id = var.kratos_smtp_secret_id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${google_service_account.kratos.email}"
}

resource "google_secret_manager_secret_iam_member" "kratos_smtp_from_address" {
  secret_id = var.kratos_smtp_from_address_secret_id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${google_service_account.kratos.email}"
}

resource "google_secret_manager_secret_iam_member" "kratos_dsn" {
  secret_id = google_secret_manager_secret.kratos_dsn.secret_id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${google_service_account.kratos.email}"
}

# ============================================================
# Kratos Public (port 4433) —　publicly accessible
# ============================================================

resource "google_cloud_run_v2_service" "kratos_public" {
  name                = "${var.project_name}-${var.environment}-kratos-public"
  location            = var.region
  ingress             = "INGRESS_TRAFFIC_ALL"
  deletion_protection = false

  lifecycle {
    ignore_changes = [template[0].containers[0].image]
  }

  template {
    service_account = google_service_account.kratos.email

    scaling {
      min_instance_count = 0
      max_instance_count = 2
    }

    volumes {
      name = "cloudsql"
      cloud_sql_instance {
        instances = [var.db_connection_name]
      }
    }

    containers {
      image = var.kratos_image

      ports {
        container_port = 4433
      }

      resources {
        limits = {
          cpu    = "1"
          memory = "512Mi"
        }
      }

      volume_mounts {
        name       = "cloudsql"
        mount_path = "/cloudsql"
      }

      env {
        name = "DSN"
        value_source {
          secret_key_ref {
            secret  = google_secret_manager_secret.kratos_dsn.secret_id
            version = "latest"
          }
        }
      }

      env {
        name = "SECRETS_COOKIE_0"
        value_source {
          secret_key_ref {
            secret  = var.kratos_cookie_secret_id
            version = "latest"
          }
        }
      }

      env {
        name = "SECRETS_CIPHER_0"
        value_source {
          secret_key_ref {
            secret  = var.kratos_cipher_secret_id
            version = "latest"
          }
        }
      }

      env {
        name = "COURIER_SMTP_CONNECTION_URI"
        value_source {
          secret_key_ref {
            secret  = var.kratos_smtp_secret_id
            version = "latest"
          }
        }
      }

      env {
        name = "COURIER_SMTP_FROM_ADDRESS"
        value_source {
          secret_key_ref {
            secret  = var.kratos_smtp_from_address_secret_id
            version = "latest"
          }
        }
      }

      env {
        name  = "SERVE_PUBLIC_BASE_URL"
        value = var.kratos_public_base_url
      }

      env {
        name  = "SERVE_ADMIN_BASE_URL"
        value = var.kratos_admin_base_url
      }

      env {
        name  = "SERVE_PUBLIC_CORS_ALLOWED_ORIGINS_0"
        value = var.frontend_url
      }

      env {
        name  = "SELFSERVICE_DEFAULT_BROWSER_RETURN_URL"
        value = "${var.frontend_url}/"
      }

      env {
        name  = "SELFSERVICE_ALLOWED_RETURN_URLS_0"
        value = var.frontend_url
      }

      env {
        name  = "SELFSERVICE_FLOWS_ERROR_UI_URL"
        value = "${var.frontend_url}/error"
      }

      env {
        name  = "SELFSERVICE_FLOWS_SETTINGS_UI_URL"
        value = "${var.frontend_url}/settings"
      }

      env {
        name  = "SELFSERVICE_FLOWS_RECOVERY_UI_URL"
        value = "${var.frontend_url}/recovery"
      }

      env {
        name  = "SELFSERVICE_FLOWS_VERIFICATION_UI_URL"
        value = "${var.frontend_url}/verification"
      }

      env {
        name  = "SELFSERVICE_FLOWS_VERIFICATION_AFTER_DEFAULT_BROWSER_RETURN_URL"
        value = "${var.frontend_url}/"
      }

      env {
        name  = "SELFSERVICE_FLOWS_LOGOUT_AFTER_DEFAULT_BROWSER_RETURN_URL"
        value = "${var.frontend_url}/login"
      }

      env {
        name  = "SELFSERVICE_FLOWS_LOGIN_UI_URL"
        value = "${var.frontend_url}/login"
      }

      env {
        name  = "SELFSERVICE_FLOWS_REGISTRATION_UI_URL"
        value = "${var.frontend_url}/register"
      }

      env {
        name  = "SELFSERVICE_FLOWS_REGISTRATION_AFTER_PASSWORD_HOOKS_0_CONFIG_URL"
        value = "${var.backend_url}/api/v1/users"
      }
    }
  }

  depends_on = [
    google_secret_manager_secret_version.kratos_dsn,
    google_secret_manager_secret_iam_member.kratos_dsn,
    google_secret_manager_secret_iam_member.kratos_cookie_secret,
    google_secret_manager_secret_iam_member.kratos_cipher_secret,
    google_secret_manager_secret_iam_member.kratos_smtp,
    google_secret_manager_secret_iam_member.kratos_smtp_from_address,
    google_project_iam_member.kratos_cloudsql_client,
  ]
}

# IAM: Allow public access to Kratos public
resource "google_cloud_run_v2_service_iam_member" "kratos_public_invoker" {
  project  = google_cloud_run_v2_service.kratos_public.project
  location = google_cloud_run_v2_service.kratos_public.location
  name     = google_cloud_run_v2_service.kratos_public.name
  role     = "roles/run.invoker"
  member   = "allUsers"
}

# ============================================================
# Kratos Admin (port 4434) — internal only
# ============================================================

resource "google_cloud_run_v2_service" "kratos_admin" {
  name                = "${var.project_name}-${var.environment}-kratos-admin"
  location            = var.region
  ingress             = "INGRESS_TRAFFIC_INTERNAL_ONLY"
  deletion_protection = false

  lifecycle {
    ignore_changes = [template[0].containers[0].image]
  }

  template {
    service_account = google_service_account.kratos.email

    scaling {
      min_instance_count = 0
      max_instance_count = 2
    }

    volumes {
      name = "cloudsql"
      cloud_sql_instance {
        instances = [var.db_connection_name]
      }
    }

    containers {
      image = var.kratos_image

      ports {
        container_port = 4434
      }

      resources {
        limits = {
          cpu    = "1"
          memory = "512Mi"
        }
      }

      volume_mounts {
        name       = "cloudsql"
        mount_path = "/cloudsql"
      }

      env {
        name = "DSN"
        value_source {
          secret_key_ref {
            secret  = google_secret_manager_secret.kratos_dsn.secret_id
            version = "latest"
          }
        }
      }

      env {
        name = "SECRETS_COOKIE_0"
        value_source {
          secret_key_ref {
            secret  = var.kratos_cookie_secret_id
            version = "latest"
          }
        }
      }

      env {
        name = "SECRETS_CIPHER_0"
        value_source {
          secret_key_ref {
            secret  = var.kratos_cipher_secret_id
            version = "latest"
          }
        }
      }

      env {
        name = "COURIER_SMTP_CONNECTION_URI"
        value_source {
          secret_key_ref {
            secret  = var.kratos_smtp_secret_id
            version = "latest"
          }
        }
      }

      env {
        name = "COURIER_SMTP_FROM_ADDRESS"
        value_source {
          secret_key_ref {
            secret  = var.kratos_smtp_from_address_secret_id
            version = "latest"
          }
        }
      }

      env {
        name  = "SERVE_PUBLIC_BASE_URL"
        value = var.kratos_public_base_url
      }

      env {
        name  = "SERVE_ADMIN_BASE_URL"
        value = var.kratos_admin_base_url
      }

      env {
        name  = "SERVE_PUBLIC_CORS_ALLOWED_ORIGINS_0"
        value = var.frontend_url
      }

      env {
        name  = "SELFSERVICE_DEFAULT_BROWSER_RETURN_URL"
        value = "${var.frontend_url}/"
      }

      env {
        name  = "SELFSERVICE_ALLOWED_RETURN_URLS_0"
        value = var.frontend_url
      }

      env {
        name  = "SELFSERVICE_FLOWS_ERROR_UI_URL"
        value = "${var.frontend_url}/error"
      }

      env {
        name  = "SELFSERVICE_FLOWS_SETTINGS_UI_URL"
        value = "${var.frontend_url}/settings"
      }

      env {
        name  = "SELFSERVICE_FLOWS_RECOVERY_UI_URL"
        value = "${var.frontend_url}/recovery"
      }

      env {
        name  = "SELFSERVICE_FLOWS_VERIFICATION_UI_URL"
        value = "${var.frontend_url}/verification"
      }

      env {
        name  = "SELFSERVICE_FLOWS_VERIFICATION_AFTER_DEFAULT_BROWSER_RETURN_URL"
        value = "${var.frontend_url}/"
      }

      env {
        name  = "SELFSERVICE_FLOWS_LOGOUT_AFTER_DEFAULT_BROWSER_RETURN_URL"
        value = "${var.frontend_url}/login"
      }

      env {
        name  = "SELFSERVICE_FLOWS_LOGIN_UI_URL"
        value = "${var.frontend_url}/login"
      }

      env {
        name  = "SELFSERVICE_FLOWS_REGISTRATION_UI_URL"
        value = "${var.frontend_url}/register"
      }

      env {
        name  = "SELFSERVICE_FLOWS_REGISTRATION_AFTER_PASSWORD_HOOKS_0_CONFIG_URL"
        value = "${var.backend_url}/api/v1/users"
      }
    }
  }

  depends_on = [
    google_secret_manager_secret_version.kratos_dsn,
    google_secret_manager_secret_iam_member.kratos_dsn,
    google_secret_manager_secret_iam_member.kratos_cookie_secret,
    google_secret_manager_secret_iam_member.kratos_cipher_secret,
    google_secret_manager_secret_iam_member.kratos_smtp,
    google_secret_manager_secret_iam_member.kratos_smtp_from_address,
    google_project_iam_member.kratos_cloudsql_client,
  ]
}

# ============================================================
# Kratos Migration Job
# ============================================================

resource "google_cloud_run_v2_job" "kratos_migrate" {
  name     = "${var.project_name}-${var.environment}-kratos-migrate"
  location = var.region

  lifecycle {
    ignore_changes = [template[0].template[0].containers[0].image]
  }

  template {
    template {
      service_account = google_service_account.kratos.email

      max_retries = 1

      volumes {
        name = "cloudsql"
        cloud_sql_instance {
          instances = [var.db_connection_name]
        }
      }

      containers {
        image = var.kratos_image
        args  = ["-c", "/etc/config/kratos/kratos.prod.yml", "migrate", "sql", "up", "-e", "--yes"]

        resources {
          limits = {
            cpu    = "1"
            memory = "512Mi"
          }
        }

        volume_mounts {
          name       = "cloudsql"
          mount_path = "/cloudsql"
        }

        env {
          name = "DSN"
          value_source {
            secret_key_ref {
              secret  = google_secret_manager_secret.kratos_dsn.secret_id
              version = "latest"
            }
          }
        }

        env {
          name = "SECRETS_COOKIE_0"
          value_source {
            secret_key_ref {
              secret  = var.kratos_cookie_secret_id
              version = "latest"
            }
          }
        }

        env {
          name = "SECRETS_CIPHER_0"
          value_source {
            secret_key_ref {
              secret  = var.kratos_cipher_secret_id
              version = "latest"
            }
          }
        }

        env {
          name = "COURIER_SMTP_CONNECTION_URI"
          value_source {
            secret_key_ref {
              secret  = var.kratos_smtp_secret_id
              version = "latest"
            }
          }
        }

        env {
          name = "COURIER_SMTP_FROM_ADDRESS"
          value_source {
            secret_key_ref {
              secret  = var.kratos_smtp_from_address_secret_id
              version = "latest"
            }
          }
        }

        env {
          name  = "SERVE_PUBLIC_BASE_URL"
          value = var.kratos_public_base_url
        }

        env {
          name  = "SERVE_ADMIN_BASE_URL"
          value = var.kratos_admin_base_url
        }

        env {
          name  = "SELFSERVICE_FLOWS_REGISTRATION_AFTER_PASSWORD_HOOKS_2_CONFIG_URL"
          value = "${var.backend_url}/api/v1/users"
        }
      }
    }
  }

  depends_on = [
    google_secret_manager_secret_version.kratos_dsn,
    google_secret_manager_secret_iam_member.kratos_dsn,
    google_secret_manager_secret_iam_member.kratos_cookie_secret,
    google_secret_manager_secret_iam_member.kratos_cipher_secret,
    google_secret_manager_secret_iam_member.kratos_smtp,
    google_secret_manager_secret_iam_member.kratos_smtp_from_address,
    google_project_iam_member.kratos_cloudsql_client,
  ]
}

# ============================================================
# API Service (port 8080) — internal only
# ============================================================

resource "google_cloud_run_v2_service" "api" {
  name                = "${var.project_name}-${var.environment}-api"
  location            = var.region
  ingress             = "INGRESS_TRAFFIC_INTERNAL_ONLY"
  deletion_protection = false

  lifecycle {
    ignore_changes = [template[0].containers[0].image]
  }

  template {
    service_account = google_service_account.api.email

    scaling {
      min_instance_count = 0
      max_instance_count = 2
    }

    volumes {
      name = "cloudsql"
      cloud_sql_instance {
        instances = [var.db_connection_name]
      }
    }

    containers {
      image = var.api_image

      ports {
        container_port = 8080
      }

      resources {
        limits = {
          cpu    = "1"
          memory = "512Mi"
        }
      }

      volume_mounts {
        name       = "cloudsql"
        mount_path = "/cloudsql"
      }

      env {
        name  = "GO_ENV"
        value = var.environment
      }

      env {
        name  = "KRATOS_ADMIN_URL"
        value = google_cloud_run_v2_service.kratos_admin.uri
      }

      env {
        name  = "FRONTEND_ORIGIN"
        value = var.frontend_url
      }

      env {
        name = "DATABASE_URL"
        value_source {
          secret_key_ref {
            secret  = google_secret_manager_secret.api_dsn.secret_id
            version = "latest"
          }
        }
      }
    }
  }

  depends_on = [
    google_secret_manager_secret_iam_member.api_dsn,
    google_secret_manager_secret_version.api_dsn,
    google_project_iam_member.api_cloudsql_client,
  ]
}

# IAM: Allow API service to invoke Kratos admin
resource "google_cloud_run_v2_service_iam_member" "api_invokes_kratos_admin" {
  project  = google_cloud_run_v2_service.kratos_admin.project
  location = google_cloud_run_v2_service.kratos_admin.location
  name     = google_cloud_run_v2_service.kratos_admin.name
  role     = "roles/run.invoker"
  member   = "serviceAccount:${google_service_account.api.email}"
}

# ============================================================
# Frontend Service (port 3000) — publicly accessible
# ============================================================

resource "google_cloud_run_v2_service" "frontend" {
  name                = "${var.project_name}-${var.environment}-frontend"
  location            = var.region
  ingress             = "INGRESS_TRAFFIC_ALL"
  deletion_protection = false

  lifecycle {
    ignore_changes = [template[0].containers[0].image]
  }

  template {
    service_account = google_service_account.frontend.email

    scaling {
      min_instance_count = 0
      max_instance_count = 2
    }

    containers {
      image = var.frontend_image

      ports {
        container_port = 3000
      }

      resources {
        limits = {
          cpu    = "1"
          memory = "512Mi"
        }
      }

      env {
        name  = "API_URL"
        value = google_cloud_run_v2_service.api.uri
      }

      env {
        name  = "ORY_SDK_URL"
        value = google_cloud_run_v2_service.kratos_public.uri
      }
    }
  }
}

# IAM: Allow Frontend to invoke API
resource "google_cloud_run_v2_service_iam_member" "frontend_invokes_api" {
  project  = google_cloud_run_v2_service.api.project
  location = google_cloud_run_v2_service.api.location
  name     = google_cloud_run_v2_service.api.name
  role     = "roles/run.invoker"
  member   = "serviceAccount:${google_service_account.frontend.email}"
}

# IAM: Allow public access to Frontend
resource "google_cloud_run_v2_service_iam_member" "frontend_invoker" {
  project  = google_cloud_run_v2_service.frontend.project
  location = google_cloud_run_v2_service.frontend.location
  name     = google_cloud_run_v2_service.frontend.name
  role     = "roles/run.invoker"
  member   = "allUsers"
}
