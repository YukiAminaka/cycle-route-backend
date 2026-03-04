output "db_password" {
  description = "Database password (sensitive)"
  value       = random_password.db_password.result
  sensitive   = true
}

output "db_password_secret_id" {
  description = "Secret Manager secret ID for DB password"
  value       = google_secret_manager_secret.db_password.secret_id
}

output "kratos_cookie_secret_id" {
  description = "Secret Manager secret ID for Kratos cookie secret"
  value       = google_secret_manager_secret.kratos_cookie_secret.secret_id
}

output "kratos_cipher_secret_id" {
  description = "Secret Manager secret ID for Kratos cipher secret"
  value       = google_secret_manager_secret.kratos_cipher_secret.secret_id
}

output "kratos_smtp_secret_id" {
  description = "Secret Manager secret ID for Kratos SMTP connection URI"
  value       = google_secret_manager_secret.kratos_smtp.secret_id
}
