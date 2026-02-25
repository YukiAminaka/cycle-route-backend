resource "random_password" "db_password" {
  length  = 32
  special = true
}

resource "google_secret_manager_secret" "db_password" {
  secret_id = "${var.project_name}-${var.environment}-db-password"

  replication {
    auto {}
  }

  labels = {
    project     = var.project_name
    environment = var.environment
  }
}

resource "google_secret_manager_secret_version" "db_password" {
  secret      = google_secret_manager_secret.db_password.id
  secret_data = random_password.db_password.result
}

resource "random_password" "kratos_cookie_secret" {
  length  = 32
  special = false
}

resource "random_password" "kratos_csrf_cookie_secret" {
  length  = 32
  special = false
}

resource "google_secret_manager_secret" "kratos_secrets" {
  secret_id = "${var.project_name}-${var.environment}-kratos-secrets"

  replication {
    auto {}
  }

  labels = {
    project     = var.project_name
    environment = var.environment
  }
}

resource "google_secret_manager_secret_version" "kratos_secrets" {
  secret = google_secret_manager_secret.kratos_secrets.id
  secret_data = jsonencode({
    cookie_secret      = random_password.kratos_cookie_secret.result
    csrf_cookie_secret = random_password.kratos_csrf_cookie_secret.result
  })
}
