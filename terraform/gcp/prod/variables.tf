variable "project_id" {
  description = "GCP Project ID"
  type        = string
}

variable "region" {
  description = "GCP region"
  type        = string
  default     = "asia-northeast1"
}

variable "project_name" {
  description = "Project name used for resource naming"
  type        = string
  default     = "rideline"
}

variable "environment" {
  description = "Environment name"
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

variable "kratos_smtp_connection_uri" {
  description = "SMTP connection URI for Kratos courier (e.g. smtps://user:password@smtp.example.com:465)"
  type        = string
  sensitive   = true
}

variable "kratos_smtp_from_address" {
  description = "From email address for Kratos courier SMTP"
  type        = string
  sensitive   = true
}

variable "kratos_public_base_url" {
  description = "Externally accessible base URL for Kratos public API."
  type        = string
  default     = ""
}

variable "kratos_admin_base_url" {
  description = "Internally accessible base URL for Kratos admin API."
  type        = string
  default     = ""
}

variable "frontend_url" {
  description = "Frontend service URL."
  type        = string
  default     = ""
}

variable "backend_url" {
  description = "API service URL."
  type        = string
  default     = ""
}
