output "db_password" {
  description = "Database password (sensitive)"
  value       = random_password.db_password.result
  sensitive   = true
}

output "db_password_secret_id" {
  description = "Secret Manager secret ID for DB password"
  value       = google_secret_manager_secret.db_password.secret_id
}

output "kratos_secrets_secret_id" {
  description = "Secret Manager secret ID for Kratos secrets"
  value       = google_secret_manager_secret.kratos_secrets.secret_id
}
