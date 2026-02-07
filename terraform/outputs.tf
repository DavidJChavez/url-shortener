output "api_url" {
  description = "URL of the deployed API"
  value       = google_cloud_run_v2_service.api.uri
}

output "artifact_registry" {
  description = "Artifact Registry URL"
  value       = "${var.region}-docker.pkg.dev/${var.project_id}/url-shortener"
}
