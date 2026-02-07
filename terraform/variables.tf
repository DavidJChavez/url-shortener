variable "project_id" {
  description = "GCP Project ID"
  type        = string
  default     = "djrc-url-shortener-prod"
}

variable "region" {
  description = "GCP Region"
  type        = string
  default     = "us-central1"
}

variable "database_url" {
  description = "PostgreSQL connection string"
  type        = string
  sensitive   = true
}
