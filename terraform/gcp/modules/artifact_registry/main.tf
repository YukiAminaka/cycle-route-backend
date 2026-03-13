
# Artifact Registryのアクセス権限はリポジトリ単位で設定できる。
# 例えばCloud Run(A)はリポジトリAだけ読み取れればよく、リポジトリBへのアクセスは不要
# 1つにまとめると過剰な権限を与えることになる
# クリーンアップポリシーもリポジトリ単位で設定する
# サービスごとに保持数や期間を変えたい場合、分かれていた方が柔軟に対応できる
# Artifact Registry / Docker Hubともに「1リポジトリ = 1イメージ種別」が一般的

resource "google_artifact_registry_repository" "repos" {
  for_each = toset(var.repositories)

  repository_id          = "${var.project_name}-${each.value}"
  location               = var.region
  format                 = "DOCKER"
  description            = "Docker repository for ${var.project_name} ${each.value}"
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
      older_than = "86400s" # 1日以上経過したものを削除対象とする（KEEPポリシーで保護されているものは削除されない）
    }
  }
}
