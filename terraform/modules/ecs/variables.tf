variable "project_name" {
  type = string
}

variable "environment" {
  type = string
}

variable "vpc_id" {
  type = string
}

variable "private_subnet_ids" {
  type = list(string)
}

variable "alb_target_group_arns" {
  type = map(string)
  description = "ALB target group ARNs (frontend, kratos only)"
}

variable "alb_security_group_id" {
  type = string
}

variable "db_endpoint" {
  type = string
}

variable "db_name" {
  type = string
}

variable "db_password_secret_arn" {
  type = string
}

variable "kratos_secrets_arn" {
  type = string
}

variable "ecr_repositories" {
  type = map(string)
}
