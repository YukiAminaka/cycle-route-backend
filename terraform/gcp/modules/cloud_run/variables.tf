variable "project_id" {
  description = "GCP Project ID"
  type        = string
}

variable "project_name" {
  description = "Project name"
  type        = string
}

variable "environment" {
  description = "Environment name"
  type        = string
}

variable "region" {
  description = "GCP region"
  type        = string
}

variable "db_connection_name" {
  description = "Cloud SQL connection name (project:region:instance)"
  type        = string
}

variable "db_name" {
  description = "Database name"
  type        = string
}

variable "db_user" {
  description = "Database user name"
  type        = string
}

variable "db_password" {
  description = "Database password (used to construct Kratos DSN secret)"
  type        = string
  sensitive   = true
}

variable "db_password_secret_id" {
  description = "Secret Manager secret ID for the DB password"
  type        = string
}

variable "kratos_cookie_secret_id" {
  description = "Secret Manager secret ID for Kratos cookie secret"
  type        = string
}

variable "kratos_cipher_secret_id" {
  description = "Secret Manager secret ID for Kratos cipher secret"
  type        = string
}

variable "kratos_smtp_secret_id" {
  description = "Secret Manager secret ID for Kratos SMTP connection URI"
  type        = string
}

variable "kratos_public_base_url" {
  description = "Externally accessible base URL for Kratos public API (set after first deployment)"
  type        = string
  default     = ""
}

variable "kratos_admin_base_url" {
  description = "Internally accessible base URL for Kratos admin API (set after first deployment)"
  type        = string
  default     = ""
}

variable "frontend_url" {
  description = "Frontend service URL (set after first deployment)"
  type        = string
  default     = ""
}

variable "backend_url" {
  description = "API service URL (set after first deployment)"
  type        = string
  default     = ""
}

variable "frontend_image" {
  description = "Frontend container image (including tag)"
  type        = string
}

variable "api_image" {
  description = "API container image (including tag)"
  type        = string
}

variable "kratos_image" {
  description = "Kratos container image (including tag)"
  type        = string
}
