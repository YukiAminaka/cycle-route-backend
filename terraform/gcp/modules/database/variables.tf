variable "project_id" {
  description = "GCP project ID"
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

variable "db_password" {
  description = "Database master password"
  type        = string
  sensitive   = true
}

variable "instance_tier" {
  description = "Cloud SQL instance tier"
  type        = string
  default     = "db-f1-micro"
}

variable "vpc_network_id" {
  description = "VPC network ID (Cloud SQL Private IP用)"
  type        = string
}

variable "private_vpc_connection_id" {
  description = "Private Services Access connection ID (depends_on用)"
  type        = string
}
