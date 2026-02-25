output "workload_identity_provider" {
  description = "Workload Identity Provider のリソース名 (GitHub Secrets の WIF_PROVIDER に設定)"
  value       = module.workload_identity.workload_identity_provider
}

output "github_actions_service_account" {
  description = "GitHub Actions SA のメールアドレス (GitHub Secrets の WIF_SERVICE_ACCOUNT に設定)"
  value       = module.workload_identity.service_account_email
}

output "artifact_registry_urls" {
  description = "Artifact Registry repository URLs (push images here)"
  value       = module.artifact_registry.repository_urls
}

output "db_connection_name" {
  description = "Cloud SQL connection name (for Cloud SQL Auth Proxy)"
  value       = module.database.db_connection_name
}

output "frontend_url" {
  description = "Frontend Cloud Run URL"
  value       = module.cloud_run.frontend_url
}

output "api_url" {
  description = "API Cloud Run URL (internal only)"
  value       = module.cloud_run.api_url
}

output "kratos_public_url" {
  description = "Kratos public Cloud Run URL"
  value       = module.cloud_run.kratos_public_url
}
