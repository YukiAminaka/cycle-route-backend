output "db_password_secret_arn" {
  value = aws_secretsmanager_secret.db_password.arn
}

output "kratos_secrets_arn" {
  value = aws_secretsmanager_secret.kratos_secrets.arn
}
