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
    bucket = "cycle-route-488410-terraform-state"
    prefix  = "terraform/state"
  }
}

provider "google" {
  project = var.project_id
  region  = var.region
}

# ============================================================
# Enable required GCP APIs
# ============================================================

# Google Cloudの指定した公開APIを有効化する
resource "google_project_service" "apis" {
  for_each = toset([
    "run.googleapis.com", # Cloud Run Admin API を有効にする
    "sqladmin.googleapis.com", # SQL Admin API を有効にする
    "secretmanager.googleapis.com", # Secret Manager API を有効にする
    "artifactregistry.googleapis.com",
    "compute.googleapis.com",
    "cloudresourcemanager.googleapis.com",
    "iam.googleapis.com",
    "iamcredentials.googleapis.com",
  ])

  service            = each.value # 有効にするサービス
  disable_on_destroy = false # Terraformリソースが破棄されたときにサービスを無効にするか
}

# ============================================================
# Modules
# ============================================================

module "secrets" {
  source = "../modules/secrets"

  project_name = var.project_name
  environment  = var.environment

  depends_on = [google_project_service.apis] # google_project_service.apisが作成されてから作成する
}

module "database" {
  source = "../modules/database"

  project_name = var.project_name
  environment  = var.environment
  region       = var.region
  db_password  = module.secrets.db_password

  depends_on = [google_project_service.apis]
}

module "artifact_registry" {
  source = "../modules/artifact_registry"

  project_name = var.project_name
  region       = var.region
  repositories = ["frontend", "api", "kratos"]

  depends_on = [google_project_service.apis]
}

module "workload_identity" {
  source = "../modules/workload_identity"

  project_id   = var.project_id
  project_name = var.project_name
  github_org   = var.github_org
  github_repo  = var.github_repo

  cloud_run_service_account_emails = [
    module.cloud_run.frontend_service_account_email,
    module.cloud_run.api_service_account_email,
    module.cloud_run.kratos_service_account_email,
  ]

  depends_on = [google_project_service.apis]
}

module "cloud_run" {
  source = "../modules/cloud_run"

  project_id   = var.project_id
  project_name = var.project_name
  environment  = var.environment
  region       = var.region

  db_connection_name    = module.database.db_connection_name
  db_name               = module.database.db_name
  db_user               = module.database.db_user
  db_password           = module.secrets.db_password
  db_password_secret_id = module.secrets.db_password_secret_id

  kratos_secrets_secret_id = module.secrets.kratos_secrets_secret_id

  frontend_image = "${module.artifact_registry.repository_urls["frontend"]}:latest"
  api_image      = "${module.artifact_registry.repository_urls["api"]}:latest"
  kratos_image   = "${module.artifact_registry.repository_urls["kratos"]}:latest"
}
