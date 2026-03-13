output "cluster_id" {
  value = aws_ecs_cluster.main.id
}

output "cluster_name" {
  value = aws_ecs_cluster.main.name
}

output "backend_services_sg_id" {
  value = aws_security_group.backend_services.id
}
