variable "project_name" {
  description = "Project name"
  type        = string
}

variable "environment" {
  description = "Environment name"
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
