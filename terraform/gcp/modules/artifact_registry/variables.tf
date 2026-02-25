variable "project_name" {
  description = "Project name"
  type        = string
}

variable "region" {
  description = "GCP region"
  type        = string
}

variable "repositories" {
  description = "List of repository names to create"
  type        = list(string)
  default     = ["frontend", "api", "kratos"]
}
