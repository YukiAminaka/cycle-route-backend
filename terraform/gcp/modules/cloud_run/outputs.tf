output "frontend_service_name" {
  description = "Frontend Cloud Run service name"
  value       = google_cloud_run_v2_service.frontend.name
}

output "frontend_url" {
  description = "Frontend Cloud Run service URL"
  value       = google_cloud_run_v2_service.frontend.uri
}

output "api_service_name" {
  description = "API Cloud Run service name"
  value       = google_cloud_run_v2_service.api.name
}

output "api_url" {
  description = "API Cloud Run service URL (internal only)"
  value       = google_cloud_run_v2_service.api.uri
}

output "kratos_public_service_name" {
  description = "Kratos public Cloud Run service name"
  value       = google_cloud_run_v2_service.kratos_public.name
}

output "kratos_public_url" {
  description = "Kratos public Cloud Run service URL"
  value       = google_cloud_run_v2_service.kratos_public.uri
}

output "kratos_admin_url" {
  description = "Kratos admin Cloud Run service URL (internal only)"
  value       = google_cloud_run_v2_service.kratos_admin.uri
}

output "frontend_service_account_email" {
  description = "Frontend Cloud Run service account email"
  value       = google_service_account.frontend.email
}

output "api_service_account_email" {
  description = "API Cloud Run service account email"
  value       = google_service_account.api.email
}

output "kratos_service_account_email" {
  description = "Kratos Cloud Run service account email"
  value       = google_service_account.kratos.email
}
