output "network_id" {
  description = "VPC network ID"
  value       = google_compute_network.vpc_network.id
}

output "subnetwork_id" {
  description = "VPC subnetwork ID"
  value       = google_compute_subnetwork.subnetwork.id
}
