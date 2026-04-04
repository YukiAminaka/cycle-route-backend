output "network_id" {
  description = "VPC network ID"
  value       = google_compute_network.vpc_network.id
}

output "subnetwork_id" {
  description = "VPC subnetwork ID"
  value       = google_compute_subnetwork.subnetwork.id
}

output "private_vpc_connection_id" {
  description = "Private Services Access connection ID (Cloud SQLへのPrivate IP接続用)"
  value       = google_service_networking_connection.private_vpc_connection.id
}
