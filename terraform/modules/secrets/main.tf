resource "aws_secretsmanager_secret" "db_password" {
  name = "${var.project_name}-${var.environment}-db-password"
}

resource "aws_secretsmanager_secret_version" "db_password" {
  secret_id     = aws_secretsmanager_secret.db_password.id
  secret_string = jsonencode({
    password = random_password.db_password.result
  })
}

resource "random_password" "db_password" {
  length  = 32
  special = true
}

resource "aws_secretsmanager_secret" "kratos_secrets" {
  name = "${var.project_name}-${var.environment}-kratos-secrets"
}

resource "aws_secretsmanager_secret_version" "kratos_secrets" {
  secret_id = aws_secretsmanager_secret.kratos_secrets.id
  secret_string = jsonencode({
    cookie_secret      = random_password.cookie_secret.result
    csrf_cookie_secret = random_password.csrf_cookie_secret.result
  })
}

resource "random_password" "cookie_secret" {
  length  = 32
  special = false
}

resource "random_password" "csrf_cookie_secret" {
  length  = 32
  special = false
}
