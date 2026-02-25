variable "project_id" {
  description = "GCP Project ID"
  type        = string
}

variable "project_name" {
  description = "Project name used for resource naming"
  type        = string
}

variable "github_org" {
  description = "GitHub organization or username"
  type        = string
}

variable "github_repo" {
  description = "GitHub repository name"
  type        = string
}

variable "cloud_run_service_account_emails" {
  description = "Emails of Cloud Run service accounts that GitHub Actions needs to act as"
  type        = list(string)
  default     = []
}
