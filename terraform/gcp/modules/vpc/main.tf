resource "google_compute_network" "vpc_network" {
  name                    = "${var.project_name}-${var.environment}-vpc"
  auto_create_subnetworks = false
  mtu                     = 1460
}

resource "google_compute_subnetwork" "subnetwork" {
  name                     = "${var.project_name}-${var.environment}-subnet"
  network                  = google_compute_network.vpc_network.id
  region                   = var.region
  ip_cidr_range            = "10.0.0.0/24"
  private_ip_google_access = true
}



# CloudSQLでPrivate IP接続を行うためのプライベートサービスアクセスを構成
# サービスプロデューサーネットワーク向けにIPアドレス範囲を予約
resource "google_compute_global_address" "private_ip_range" {
  name          = "${var.project_name}-${var.environment}-private-ip-range"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 16
  network       = google_compute_network.vpc_network.id
}

# ユーザーVPCとService Producer VPCをVPCピアリングする
resource "google_service_networking_connection" "private_vpc_connection" {
  network                 = google_compute_network.vpc_network.id
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.private_ip_range.name]
}

resource "google_dns_managed_zone" "dns_zone" {
  name     = "${var.project_name}-${var.environment}-dns-zone"
  dns_name = "run.app."

  visibility = "private"

  private_visibility_config {
    networks {
      network_url = google_compute_network.vpc_network.id
    }
  }
}

resource "google_dns_record_set" "dns_record_set_a" {
  name         = "run.app."
  type         = "A"
  ttl          = 60
  managed_zone = google_dns_managed_zone.dns_zone.name
  rrdatas      = ["199.36.153.4", "199.36.153.5", "199.36.153.6", "199.36.153.7"] # restricted.googleapis.com
}

resource "google_dns_record_set" "dns_record_set_cname" {
  name         = "*.run.app."
  type         = "CNAME"
  ttl          = 60
  managed_zone = google_dns_managed_zone.dns_zone.name
  rrdatas      = ["run.app."]
}
