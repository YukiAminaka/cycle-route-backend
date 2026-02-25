# ============================================================
# Workload Identity Pool
# ============================================================

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

  # このリポジトリからのリクエストのみ許可
  attribute_condition = "attribute.repository == '${var.github_org}/${var.github_repo}'"
}

# ============================================================
# Service Account for GitHub Actions
# ============================================================

resource "google_service_account" "github_actions" {
  account_id   = "${var.project_name}-github-actions"
  display_name = "${var.project_name} GitHub Actions Service Account"
  description  = "Service Account used by GitHub Actions for CI/CD"
}

# ============================================================
# Allow GitHub Actions (WIF) to impersonate the Service Account
# ============================================================

resource "google_service_account_iam_member" "workload_identity_binding" {
  service_account_id = google_service_account.github_actions.name
  role               = "roles/iam.workloadIdentityUser"
  member             = "principalSet://iam.googleapis.com/${google_iam_workload_identity_pool.github.name}/attribute.repository/${var.github_org}/${var.github_repo}"
}

# ============================================================
# IAM: Artifact Registry への push 権限
# ============================================================

resource "google_project_iam_member" "github_actions_ar_writer" {
  project = var.project_id
  role    = "roles/artifactregistry.writer"
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
# IAM: Cloud Run SA として動作する権限
# (gcloud run deploy --service-account 指定時に必要)
# ============================================================

resource "google_service_account_iam_member" "github_actions_sa_user" {
  for_each = toset(var.cloud_run_service_account_emails)

  service_account_id = "projects/${var.project_id}/serviceAccounts/${each.value}"
  role               = "roles/iam.serviceAccountUser"
  member             = "serviceAccount:${google_service_account.github_actions.email}"
}
