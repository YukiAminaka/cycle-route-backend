output "db_instance_name" {
  description = "Cloud SQL instance name"
  value       = google_sql_database_instance.main.name
}

output "db_connection_name" {
  description = "Cloud SQL connection name (used by Cloud SQL Auth Proxy)"
  value       = google_sql_database_instance.main.connection_name
}

output "db_name" {
  description = "Database name"
  value       = google_sql_database.main.name
}

output "db_user" {
  description = "Database user name"
  value       = google_sql_user.main.name
}

output "db_service_account_email" {
  description = "Database service account email"
  value       = google_service_account.database.email
}
