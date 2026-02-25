output "service_account_email" {
  description = "GitHub Actions Service Account のメールアドレス (GitHub Secrets の WIF_SERVICE_ACCOUNT に設定)"
  value       = google_service_account.github_actions.email
}

output "workload_identity_provider" {
  description = "Workload Identity Pool Provider のリソース名 (GitHub Secrets の WIF_PROVIDER に設定)"
  value       = google_iam_workload_identity_pool_provider.github.name
}
