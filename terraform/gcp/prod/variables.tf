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
  default     = "prod"
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

