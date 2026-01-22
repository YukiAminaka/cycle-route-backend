terraform {
  required_version = ">= 1.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = var.aws_region
}

module "network" {
  source = "../modules/network"

  project_name = var.project_name
  environment  = var.environment
  vpc_cidr     = var.vpc_cidr
}

module "secrets" {
  source = "../modules/secrets"

  project_name = var.project_name
  environment  = var.environment
}

module "database" {
  source = "../modules/database"

  project_name        = var.project_name
  environment         = var.environment
  vpc_id              = module.network.vpc_id
  database_subnet_ids = module.network.database_subnet_ids
  backend_services_sg_id = module.ecs.backend_services_sg_id
  db_password_secret  = module.secrets.db_password_secret_arn
}

module "ecr" {
  source = "../modules/ecr"

  project_name = var.project_name
  repositories = ["frontend", "api", "kratos"]
}

module "alb" {
  source = "../modules/alb"

  project_name       = var.project_name
  environment        = var.environment
  vpc_id             = module.network.vpc_id
  public_subnet_ids  = module.network.public_subnet_ids
}

module "ecs" {
  source = "../modules/ecs"

  project_name           = var.project_name
  environment            = var.environment
  vpc_id                 = module.network.vpc_id
  private_subnet_ids     = module.network.private_subnet_ids
  alb_target_group_arns  = module.alb.target_group_arns
  alb_security_group_id  = module.alb.alb_security_group_id
  db_endpoint            = module.database.db_endpoint
  db_name                = module.database.db_name
  db_password_secret_arn = module.secrets.db_password_secret_arn
  kratos_secrets_arn     = module.secrets.kratos_secrets_arn
  ecr_repositories       = module.ecr.repository_urls
}
