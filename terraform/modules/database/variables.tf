variable "project_name" {
  type = string
}

variable "environment" {
  type = string
}

variable "vpc_id" {
  type = string
}

variable "backend_services_sg_id" {
  type = string
}

variable "database_subnet_ids" {
  type = list(string)
}

variable "db_password_secret" {
  type = string
}
