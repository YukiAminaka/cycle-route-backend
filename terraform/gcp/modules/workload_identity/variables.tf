variable "project_id" {
  description = "GCP Project ID"
  type        = string
}

variable "project_name" {
  description = "Project name used for resource naming"
  type        = string
}

variable "github_repository" {
  description = "GitHub repository in 'org/repo' format (e.g. 'myorg/myrepo')"
  type        = string
}

variable "frontend_github_repository" {
  description = "Frontend GitHub repository in 'org/repo' format (e.g. 'myorg/myrepo-frontend')"
  type        = string
}

variable "cloud_run_service_account_emails" {
  description = "Emails of Cloud Run service accounts that GitHub Actions (backend) needs to act as"
  type        = list(string)
  default     = []
}

variable "frontend_cloud_run_service_account_email" {
  description = "Email of the frontend Cloud Run service account that GitHub Actions (frontend) needs to act as"
  type        = string
}

variable "db_service_account_email" {
  description = "Email of database service account for migration"
  type        = string
}

variable "terraform_state_bucket" {
  description = "Terraform state を保存する GCS バケット名"
  type        = string
}
