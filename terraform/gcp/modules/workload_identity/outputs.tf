output "service_account_email" {
  description = "GitHub Actions Service Account のメールアドレス (GitHub Secrets の BUILD_ACCOUNT に設定)"
  value       = google_service_account.github_actions.email
}

output "terraform_service_account_email" {
  description = "Terraform Service Account のメールアドレス (GitHub Secrets の OPERATION_ACCOUNT に設定)"
  value       = google_service_account.terraform.email
}

output "db_migration_service_account_email" {
  description = "Database Migration Service Account のメールアドレス (GitHub Secrets の MIGRATION_ACCOUNT に設定)"
  value       = var.db_service_account_email
}

output "workload_identity_provider" {
  description = "Workload Identity Pool Provider のリソース名 (GitHub Secrets の WORKLOAD_IDENTITY_PROVIDER に設定)"
  value       = google_iam_workload_identity_pool_provider.github.name
}

output "frontend_service_account_email" {
  description = "Frontend GitHub Actions Service Account のメールアドレス (フロントエンドリポジトリの GitHub Secrets に設定)"
  value       = google_service_account.github_actions_frontend.email
}
