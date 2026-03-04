variable "project_id" {
  description = "GCP Project ID"
  type        = string
  default = "value"
}

variable "region" {
  description = "GCP region"
  type        = string
  default     = "asia-northeast1"
}

variable "project_name" {
  description = "Project name used for resource naming"
  type        = string
  default     = "cycle-route"
}

variable "environment" {
  description = "Environment name"
  type        = string
}

variable "github_org" {
  description = "GitHub organization or username"
  type        = string
}

variable "github_repo" {
  description = "GitHub repository name"
  type        = string
  default     = "cycle-route-backend"
}

variable "kratos_smtp_connection_uri" {
  description = "SMTP connection URI for Kratos courier (e.g. smtps://user:password@smtp.example.com:465)"
  type        = string
  sensitive   = true
}

variable "kratos_public_base_url" {
  description = "Externally accessible base URL for Kratos public API. Set after first deployment using terraform output kratos_public_url."
  type        = string
  default     = ""
}

variable "kratos_admin_base_url" {
  description = "Internally accessible base URL for Kratos admin API. Set after first deployment using terraform output kratos_admin_url."
  type        = string
  default     = ""
}

variable "frontend_url" {
  description = "Frontend service URL. Set after first deployment using terraform output frontend_url."
  type        = string
  default     = ""
}

variable "backend_url" {
  description = "API service URL. Set after first deployment using terraform output api_url."
  type        = string
  default     = ""
}
