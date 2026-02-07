terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
  }
}

provider "google" {
  project = var.project_id
  region  = var.region
}

resource "google_artifact_registry_repository" "api" {
  location      = var.region
  repository_id = "url-shortener"
  format        = "DOCKER"
}

resource "google_cloud_run_v2_service" "api" {
  name     = "url-shortener-api"
  location = var.region

  template {
    containers {
      image = "${var.region}-docker.pkg.dev/${var.project_id}/url-shortener/api:latest"

      env {
        name  = "DATABASE_URL"
        value = var.database_url
      }

      env {
        name  = "BASE_URL"
        value = "https://api.davidjchavez.com/shortener"
      }

      ports {
        container_port = 8080
      }
    }

    scaling {
      min_instance_count = 0
      max_instance_count = 1
    }
  }
}

resource "google_cloud_run_v2_service_iam_member" "public" {
  project  = var.project_id
  location = var.region
  name     = google_cloud_run_v2_service.api.name
  role     = "roles/run.invoker"
  member   = "allUsers"
}
