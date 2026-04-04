output "db_instance_name" {
  description = "Cloud SQL instance name"
  value       = google_sql_database_instance.main.name
}

output "db_connection_name" {
  description = "Cloud SQL connection name"
  value       = google_sql_database_instance.main.connection_name
}

output "db_private_ip" {
  description = "Cloud SQL instance private IP address"
  value       = google_sql_database_instance.main.private_ip_address
}

output "db_name" {
  description = "Database name"
  value       = google_sql_database.main.name
}

output "db_user" {
  description = "Database user name"
  value       = google_sql_user.main.name
}

