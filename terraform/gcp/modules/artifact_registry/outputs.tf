output "repository_urls" {
  description = "Map of repository names to their full image URLs (without tag)"
  value = {
    for k, v in google_artifact_registry_repository.repos :
    k => "${v.location}-docker.pkg.dev/${data.google_project.current.project_id}/${v.repository_id}"
  }
}
