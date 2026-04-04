terraform {
  required_version = ">= 1.0"
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 6.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.0"
    }
  }

  backend "gcs" {
    bucket = "rideline-489422-terraform-state"
    prefix = "terraform/state"
  }
}

provider "google" {
  project = var.project_id
  region  = var.region
}

# Google Cloudの指定した公開APIを有効化する
resource "google_project_service" "apis" {
  for_each = toset([
    "run.googleapis.com",              # Cloud Run Admin API を有効にする
    "sqladmin.googleapis.com",         # SQL Admin API を有効にする
    "secretmanager.googleapis.com",    # Secret Manager API を有効にする
    "artifactregistry.googleapis.com", # Artifact Registry API を有効にする
    "compute.googleapis.com",
    "dns.googleapis.com",
    "cloudresourcemanager.googleapis.com",
    "iam.googleapis.com",
    "iamcredentials.googleapis.com",
    "servicenetworking.googleapis.com", # Private Services Access用
  ])

  service            = each.value # 有効にするサービス
  disable_on_destroy = false      # Terraformリソースが破棄されたときにサービスを無効にするか
}

module "vpc" {
  source = "../modules/vpc"

  project_name = var.project_name
  environment  = var.environment
  region       = var.region

  depends_on = [google_project_service.apis]
}

module "secrets" {
  source = "../modules/secrets"

  project_name               = var.project_name
  environment                = var.environment
  kratos_smtp_connection_uri = var.kratos_smtp_connection_uri
  kratos_smtp_from_address   = var.kratos_smtp_from_address

  depends_on = [google_project_service.apis] # google_project_service.apisが作成されてから作成する
}

module "database" {
  source = "../modules/database"

  project_id                = var.project_id
  project_name              = var.project_name
  environment               = var.environment
  region                    = var.region
  db_password               = module.secrets.db_password
  vpc_network_id            = module.vpc.network_id
  private_vpc_connection_id = module.vpc.private_vpc_connection_id

  depends_on = [google_project_service.apis]
}

module "artifact_registry" {
  source = "../modules/artifact_registry"

  project_name = var.project_name
  region       = var.region
  repositories = ["frontend", "api", "kratos", "atlas"]

  depends_on = [google_project_service.apis]
}

module "workload_identity" {
  source = "../modules/workload_identity"

  project_id                 = var.project_id
  project_name               = var.project_name
  github_repository          = var.github_repository
  frontend_github_repository = var.frontend_github_repository

  cloud_run_service_account_emails = [
    module.cloud_run.frontend_service_account_email,
    module.cloud_run.api_service_account_email,
    module.cloud_run.kratos_service_account_email,
    module.cloud_run.migration_service_account_email,
  ]

  frontend_cloud_run_service_account_email = module.cloud_run.frontend_service_account_email

  terraform_state_bucket = "rideline-489422-terraform-state"

  depends_on = [google_project_service.apis]
}

module "cloud_run" {
  source = "../modules/cloud_run"

  project_id   = var.project_id
  project_name = var.project_name
  environment  = var.environment
  region       = var.region

  db_private_ip         = module.database.db_private_ip
  db_name               = module.database.db_name
  db_user               = module.database.db_user
  db_password           = module.secrets.db_password
  db_password_secret_id = module.secrets.db_password_secret_id

  kratos_cookie_secret_id            = module.secrets.kratos_cookie_secret_id
  kratos_cipher_secret_id            = module.secrets.kratos_cipher_secret_id
  kratos_smtp_secret_id              = module.secrets.kratos_smtp_secret_id
  kratos_smtp_from_address_secret_id = module.secrets.kratos_smtp_from_address_secret_id

  kratos_public_base_url = var.kratos_public_base_url
  kratos_admin_base_url  = var.kratos_admin_base_url
  frontend_url           = var.frontend_url
  backend_url            = var.backend_url

  frontend_image = "${module.artifact_registry.repository_urls["frontend"]}:latest"
  api_image      = "${module.artifact_registry.repository_urls["api"]}:latest"
  kratos_image   = "${module.artifact_registry.repository_urls["kratos"]}:latest"

  vpc_network_id    = module.vpc.network_id
  vpc_subnetwork_id = module.vpc.subnetwork_id
}
