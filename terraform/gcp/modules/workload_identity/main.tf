# ============================================================
# Workload Identity Pool
# ============================================================

# 外部ワークロードID の集合を表します
resource "google_iam_workload_identity_pool" "github" {
  workload_identity_pool_id = "${var.project_name}-github-pool"
  display_name              = "GitHub Actions Pool"
  description               = "Workload Identity Pool for GitHub Actions CI/CD"
  disabled                  = false
}

# ============================================================
# Workload Identity Pool Provider (GitHub OIDC)
# ============================================================

resource "google_iam_workload_identity_pool_provider" "github" {
  workload_identity_pool_id          = google_iam_workload_identity_pool.github.workload_identity_pool_id
  workload_identity_pool_provider_id = "${var.project_name}-github-provider"
  display_name                       = "GitHub Actions Provider"
  description                        = "GitHub Actions OIDC provider"

  oidc {
    issuer_uri = "https://token.actions.githubusercontent.com"
  }

  # GitHub Actionsから送られてくる身分証明書（OIDCトークン）のどの項目を、Google Cloud側でどのように扱うかを定義
  attribute_mapping = {
    "google.subject"       = "assertion.sub"
    "attribute.actor"      = "assertion.actor"
    "attribute.repository" = "assertion.repository"
    "attribute.ref"        = "assertion.ref"
  }

  # これらのリポジトリからのリクエストのみ許可
  attribute_condition = "attribute.repository in ['${var.github_repository}', '${var.frontend_github_repository}']"
}

# ============================================================
# Service Account for GitHub Actions
# ============================================================

# Actionsで借用したい権限を持つサービスアカウント
resource "google_service_account" "github_actions" {
  account_id   = "${var.project_name}-github-actions"
  display_name = "${var.project_name} GitHub Actions Service Account"
  description  = "Service Account used by GitHub Actions for CI/CD"
}

# ============================================================
# Allow GitHub Actions (WIF) to impersonate the Service Account
# ============================================================

# 外部ID(メンバー)に対してGitHub Actionsが指定したサービスアカウントになりすますロールが付与される
resource "google_service_account_iam_member" "workload_identity_binding" {
  service_account_id = google_service_account.github_actions.name
  role               = "roles/iam.workloadIdentityUser"
  member             = "principalSet://iam.googleapis.com/${google_iam_workload_identity_pool.github.name}/attribute.repository/${var.github_repository}"
}

# ============================================================
# Allow GitHub Actions (WIF) to impersonate the Database SA for migration
# ============================================================

resource "google_service_account_iam_member" "db_workload_identity_binding" {
  service_account_id = "projects/${var.project_id}/serviceAccounts/${var.db_service_account_email}"
  role               = "roles/iam.workloadIdentityUser"
  member             = "principalSet://iam.googleapis.com/${google_iam_workload_identity_pool.github.name}/attribute.repository/${var.github_repository}"
}

# ============================================================
# Service Account for Frontend GitHub Actions
# ============================================================

resource "google_service_account" "github_actions_frontend" {
  account_id   = "${var.project_name}-github-actions-fe"
  display_name = "${var.project_name} GitHub Actions Frontend Service Account"
  description  = "Service Account used by Frontend GitHub Actions for CI/CD"
}

# ============================================================
# Allow Frontend GitHub Actions (WIF) to impersonate the Frontend SA
# ============================================================

resource "google_service_account_iam_member" "workload_identity_binding_frontend" {
  service_account_id = google_service_account.github_actions_frontend.name
  role               = "roles/iam.workloadIdentityUser"
  member             = "principalSet://iam.googleapis.com/${google_iam_workload_identity_pool.github.name}/attribute.repository/${var.frontend_github_repository}"
}

# ============================================================
# IAM: Frontend SA → Frontend Cloud Run SA として動作する権限
# ============================================================

resource "google_service_account_iam_member" "github_actions_frontend_sa_user" {
  service_account_id = "projects/${var.project_id}/serviceAccounts/${var.frontend_cloud_run_service_account_email}"
  role               = "roles/iam.serviceAccountUser"
  member             = "serviceAccount:${google_service_account.github_actions_frontend.email}"
}

# ============================================================
# IAM: Artifact Registry への push 権限
# ============================================================

# 指定したプロジェクトの特定のプリンシパルの特定のロールを一括で管理するリソース
# 今回は、GitHub Actionsのサービスアカウントに対して、Artifact Registryへのpush権限を付与するために使用
resource "google_project_iam_member" "github_actions_ar_writer" {
  project = var.project_id
  role    = "roles/artifactregistry.repoAdmin"
  member  = "serviceAccount:${google_service_account.github_actions.email}"
}

# ============================================================
# IAM: Cloud Run のデプロイ権限
# ============================================================

resource "google_project_iam_member" "github_actions_run_developer" {
  project = var.project_id
  role    = "roles/run.developer"
  member  = "serviceAccount:${google_service_account.github_actions.email}"
}

# ============================================================
# IAM: Frontend SA の Artifact Registry push 権限 / Cloud Run デプロイ権限
# ============================================================

resource "google_project_iam_member" "github_actions_frontend_ar_writer" {
  project = var.project_id
  role    = "roles/artifactregistry.repoAdmin"
  member  = "serviceAccount:${google_service_account.github_actions_frontend.email}"
}

resource "google_project_iam_member" "github_actions_frontend_run_developer" {
  project = var.project_id
  role    = "roles/run.developer"
  member  = "serviceAccount:${google_service_account.github_actions_frontend.email}"
}

# ============================================================
# IAM: Cloud Run SA として動作する権限
# (gcloud run deploy --service-account 指定時に必要)
# ============================================================

resource "google_service_account_iam_member" "github_actions_sa_user" {
  for_each = toset(var.cloud_run_service_account_emails)

  service_account_id = "projects/${var.project_id}/serviceAccounts/${each.value}"
  role               = "roles/iam.serviceAccountUser"
  member             = "serviceAccount:${google_service_account.github_actions.email}"
}

# cloud runサービスアカウントを、Terraformのサービスアカウントが借用する権限が必要 
resource "google_service_account_iam_member" "terraform_sa_user" {
  for_each = toset(var.cloud_run_service_account_emails)

  service_account_id = "projects/${var.project_id}/serviceAccounts/${each.value}"
  role               = "roles/iam.serviceAccountUser" # サービスアカウントとして操作を実行するロール
  member             = "serviceAccount:${google_service_account.terraform.email}"
}

# ============================================================
# Service Account for Terraform
# ============================================================

resource "google_service_account" "terraform" {
  account_id   = "${var.project_name}-terraform"
  display_name = "${var.project_name} Terraform Service Account"
  description  = "Service Account used by GitHub Actions for Terraform operations"
}

# ============================================================
# Allow GitHub Actions (WIF) to impersonate the Terraform SA
# ============================================================

resource "google_service_account_iam_member" "terraform_wif_binding" {
  service_account_id = google_service_account.terraform.name
  role               = "roles/iam.workloadIdentityUser"
  member             = "principalSet://iam.googleapis.com/${google_iam_workload_identity_pool.github.name}/attribute.repository/${var.github_repository}"
}

# ============================================================
# IAM: Terraform SA に必要なプロジェクトレベル権限
# ============================================================

locals {
  terraform_project_roles = toset([
    "roles/serviceusage.serviceUsageAdmin",  # google_project_service でAPIを有効化
    "roles/iam.serviceAccountAdmin",         # サービスアカウントの作成・管理
    "roles/iam.workloadIdentityPoolAdmin",   # google_iam_workload_identity_pool の管理
    "roles/resourcemanager.projectIamAdmin", # google_project_iam_member の設定
    "roles/run.admin",                       # Cloud Run サービス作成・IAM設定
    "roles/storage.admin",                   # google_storage_bucket の管理
    "roles/secretmanager.admin",             # google_secret_manager_secret の管理
    "roles/cloudsql.admin",                  # google_sql_database_instance の管理
    "roles/artifactregistry.admin",          # google_artifact_registry_repository の作成
  ])
}

resource "google_project_iam_member" "terraform_roles" {
  for_each = local.terraform_project_roles

  project = var.project_id
  role    = each.value
  member  = "serviceAccount:${google_service_account.terraform.email}"
}
