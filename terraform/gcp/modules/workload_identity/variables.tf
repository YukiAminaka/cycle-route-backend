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

variable "cloud_run_service_account_emails" {
  description = "Emails of Cloud Run service accounts that GitHub Actions needs to act as"
  type        = list(string)
  default     = []
}

variable "db_service_account_email" {
  description = "Email of database service account for migration"
  type        = string
}

variable "terraform_state_bucket" {
  description = "Terraform state を保存する GCS バケット名"
  type        = string
}
