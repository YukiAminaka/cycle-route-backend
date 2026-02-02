terraform {
  required_version = ">= 1.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
  # 本番環境用のBackend設定（例: S3 + DynamoDB）
  # backend "s3" {
  #   bucket         = "cycle-route-tfstate-prod"
  #   key            = "prod/terraform.tfstate"
  #   region         = "ap-northeast-1"
  #   encrypt        = true
  #   dynamodb_table = "terraform_locks"
  # }
}

provider "aws" {
  region = var.aws_region

  default_tags {
    tags = {
      Project     = var.project_name
      Environment = var.environment
      ManagedBy   = "Terraform"
    }
  }
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
  db_password_secret  = module.secrets.db_password_secret_arn
}

module "ecr" {
  source = "../modules/ecr"

  project_name = var.project_name
  repositories = ["frontend", "api", "kratos"]
}

module "dns" {
  source = "../modules/dns"
  project_name = var.project_name
  environment  = var.environment
  domain_name  = var.domain_name
  alb_dns_name = module.alb.alb_dns_name
  alb_zone_id  = module.alb.alb_zone_id
}

module "alb" {
  source = "../modules/alb"
  domain_name        = var.domain_name
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

# 循環参照回避のためのセキュリティグループルール
# databaseモジュールとecsモジュールが相互に依存しないように、ルートで紐付けを行う
resource "aws_security_group_rule" "db_ingress_from_ecs" {
  type                     = "ingress"
  from_port                = 5432
  to_port                  = 5432
  protocol                 = "tcp"
  source_security_group_id = module.ecs.backend_services_sg_id
  security_group_id        = module.database.security_group_id
}
