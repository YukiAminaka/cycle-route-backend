resource "google_artifact_registry_repository" "repos" {
  for_each = toset(var.repositories)

  repository_id = "${var.project_name}-${each.value}"
  location      = var.region
  format        = "DOCKER"
  description   = "Docker repository for ${var.project_name} ${each.value}"
  cleanup_policy_dry_run = false

  cleanup_policies {
    id     = "keep-minimum-versions"
    action = "KEEP"

    most_recent_versions {
      keep_count = 2
    }
  }

  cleanup_policies {
    id     = "delete-old-versions"
    action = "DELETE"

    condition {
      older_than = "86400s" # 1日以上経過したものを削除対象とする（KEEPポリシーで保護されているものは削除されません）
    }
  }
}
