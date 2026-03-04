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

resource "random_password" "kratos_cipher_secret" {
  length  = 32
  special = false
}

resource "google_secret_manager_secret" "kratos_cookie_secret" {
  secret_id = "${var.project_name}-${var.environment}-kratos-cookie-secret"

  replication {
    auto {}
  }

  labels = {
    project     = var.project_name
    environment = var.environment
  }
}

resource "google_secret_manager_secret_version" "kratos_cookie_secret" {
  secret      = google_secret_manager_secret.kratos_cookie_secret.id
  secret_data = random_password.kratos_cookie_secret.result
}

resource "google_secret_manager_secret" "kratos_cipher_secret" {
  secret_id = "${var.project_name}-${var.environment}-kratos-cipher-secret"

  replication {
    auto {}
  }

  labels = {
    project     = var.project_name
    environment = var.environment
  }
}

resource "google_secret_manager_secret_version" "kratos_cipher_secret" {
  secret      = google_secret_manager_secret.kratos_cipher_secret.id
  secret_data = random_password.kratos_cipher_secret.result
}

resource "google_secret_manager_secret" "kratos_smtp" {
  secret_id = "${var.project_name}-${var.environment}-kratos-smtp"

  replication {
    auto {}
  }

  labels = {
    project     = var.project_name
    environment = var.environment
  }
}

resource "google_secret_manager_secret_version" "kratos_smtp" {
  secret      = google_secret_manager_secret.kratos_smtp.id
  secret_data = var.kratos_smtp_connection_uri
}
