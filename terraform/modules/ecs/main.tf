resource "aws_ecs_cluster" "main" {
  name = "${var.project_name}-${var.environment}-cluster"

  setting {
    name  = "containerInsights"
    value = "enabled"
  }
}

# Security Groups
resource "aws_security_group" "frontend" {
  name        = "${var.project_name}-${var.environment}-frontend-sg"
  description = "Security group for Frontend ECS tasks"
  vpc_id      = var.vpc_id

  tags = {
    Name = "${var.project_name}-${var.environment}-frontend-sg"
  }
}

resource "aws_vpc_security_group_ingress_rule" "frontend_from_alb" {
  security_group_id            = aws_security_group.frontend.id
  referenced_security_group_id = var.alb_security_group_id
  from_port                    = 3000
  to_port                      = 3000
  ip_protocol                  = "tcp"
  description                  = "Allow traffic from ALB"
}

resource "aws_vpc_security_group_egress_rule" "frontend_to_api" {
  security_group_id            = aws_security_group.frontend.id
  referenced_security_group_id = aws_security_group.api.id
  from_port                    = 8080
  to_port                      = 8080
  ip_protocol                  = "tcp"
  description                  = "Allow traffic to API"
}

resource "aws_vpc_security_group_egress_rule" "frontend_to_kratos" {
  security_group_id            = aws_security_group.frontend.id
  referenced_security_group_id = aws_security_group.kratos.id
  from_port                    = 4433
  to_port                      = 4433
  ip_protocol                  = "tcp"
  description                  = "Allow traffic to Kratos"
}

# trivy:ignore:AVD-AWS-0104
resource "aws_vpc_security_group_egress_rule" "frontend_https" {
  security_group_id = aws_security_group.frontend.id
  from_port         = 443
  to_port           = 443
  ip_protocol       = "tcp"
  cidr_ipv4         = "0.0.0.0/0"
  description       = "Allow HTTPS for external APIs"
}

resource "aws_security_group" "api" {
  name        = "${var.project_name}-${var.environment}-api-sg"
  description = "Security group for API ECS tasks"
  vpc_id      = var.vpc_id

  tags = {
    Name = "${var.project_name}-${var.environment}-api-sg"
  }
}

resource "aws_vpc_security_group_ingress_rule" "api_from_frontend" {
  security_group_id            = aws_security_group.api.id
  referenced_security_group_id = aws_security_group.frontend.id
  from_port                    = 8080
  to_port                      = 8080
  ip_protocol                  = "tcp"
  description                  = "Allow traffic from Frontend"
}

resource "aws_vpc_security_group_egress_rule" "api_to_db" {
  security_group_id            = aws_security_group.api.id
  referenced_security_group_id = var.db_security_group_id
  from_port                    = 5432
  to_port                      = 5432
  ip_protocol                  = "tcp"
  description                  = "Allow traffic to RDS"
}

resource "aws_vpc_security_group_egress_rule" "api_to_kratos" {
  security_group_id            = aws_security_group.api.id
  referenced_security_group_id = aws_security_group.kratos.id
  from_port                    = 4434
  to_port                      = 4434
  ip_protocol                  = "tcp"
  description                  = "Allow traffic to Kratos admin"
}

# trivy:ignore:AVD-AWS-0104
resource "aws_vpc_security_group_egress_rule" "api_https" {
  security_group_id = aws_security_group.api.id
  from_port         = 443
  to_port           = 443
  ip_protocol       = "tcp"
  cidr_ipv4         = "0.0.0.0/0"
  description       = "Allow HTTPS for external APIs"
}

resource "aws_security_group" "kratos" {
  name        = "${var.project_name}-${var.environment}-kratos-sg"
  description = "Security group for Kratos ECS tasks"
  vpc_id      = var.vpc_id

  tags = {
    Name = "${var.project_name}-${var.environment}-kratos-sg"
  }
}

resource "aws_vpc_security_group_ingress_rule" "kratos_public_from_alb" {
  security_group_id            = aws_security_group.kratos.id
  referenced_security_group_id = var.alb_security_group_id
  from_port                    = 4433
  to_port                      = 4433
  ip_protocol                  = "tcp"
  description                  = "Allow public API traffic from ALB"
}

resource "aws_vpc_security_group_ingress_rule" "kratos_public_from_frontend" {
  security_group_id            = aws_security_group.kratos.id
  referenced_security_group_id = aws_security_group.frontend.id
  from_port                    = 4433
  to_port                      = 4433
  ip_protocol                  = "tcp"
  description                  = "Allow public API traffic from Frontend"
}

resource "aws_vpc_security_group_ingress_rule" "kratos_admin_from_api" {
  security_group_id            = aws_security_group.kratos.id
  referenced_security_group_id = aws_security_group.api.id
  from_port                    = 4434
  to_port                      = 4434
  ip_protocol                  = "tcp"
  description                  = "Allow admin API traffic from API"
}

resource "aws_vpc_security_group_egress_rule" "kratos_to_db" {
  security_group_id            = aws_security_group.kratos.id
  referenced_security_group_id = var.db_security_group_id
  from_port                    = 5432
  to_port                      = 5432
  ip_protocol                  = "tcp"
  description                  = "Allow traffic to RDS"
}

# trivy:ignore:AVD-AWS-0104
resource "aws_vpc_security_group_egress_rule" "kratos_https" {
  security_group_id = aws_security_group.kratos.id
  from_port         = 443
  to_port           = 443
  ip_protocol       = "tcp"
  cidr_ipv4         = "0.0.0.0/0"
  description       = "Allow HTTPS for email/SMS providers"
}

# IAM Roles
resource "aws_iam_role" "ecs_task_execution" {
  name = "${var.project_name}-${var.environment}-ecs-task-execution"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
        Service = "ecs-tasks.amazonaws.com"
      }
    }]
  })
}

resource "aws_iam_role_policy_attachment" "ecs_task_execution" {
  role       = aws_iam_role.ecs_task_execution.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

resource "aws_iam_role_policy" "secrets_access" {
  name = "secrets-access"
  role = aws_iam_role.ecs_task_execution.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect = "Allow"
      Action = [
        "secretsmanager:GetSecretValue"
      ]
      Resource = [
        var.db_password_secret_arn,
        var.kratos_secrets_arn
      ]
    }]
  })
}

# CloudWatch Logs
resource "aws_cloudwatch_log_group" "ecs" {
  for_each = toset(["frontend", "api", "kratos"])

  name              = "/ecs/${var.project_name}-${var.environment}/${each.key}"
  retention_in_days = 7
}

# Service Discovery
resource "aws_service_discovery_private_dns_namespace" "main" {
  name = "${var.project_name}-${var.environment}.local"
  vpc  = var.vpc_id
}

resource "aws_service_discovery_service" "api" {
  name = "api"

  dns_config {
    namespace_id = aws_service_discovery_private_dns_namespace.main.id

    dns_records {
      ttl  = 10
      type = "A"
    }
  }

  health_check_custom_config {
    failure_threshold = 1
  }
}

resource "aws_service_discovery_service" "kratos" {
  name = "kratos"

  dns_config {
    namespace_id = aws_service_discovery_private_dns_namespace.main.id

    dns_records {
      ttl  = 10
      type = "A"
    }
  }

  health_check_custom_config {
    failure_threshold = 1
  }
}

# Frontend Service
resource "aws_ecs_task_definition" "frontend" {
  family                   = "${var.project_name}-${var.environment}-frontend"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = "256"
  memory                   = "512"
  execution_role_arn       = aws_iam_role.ecs_task_execution.arn

  container_definitions = jsonencode([{
    name  = "frontend"
    image = "${var.ecr_repositories["frontend"]}:latest"
    portMappings = [{
      containerPort = 3000
      protocol      = "tcp"
    }]
    logConfiguration = {
      logDriver = "awslogs"
      options = {
        "awslogs-group"         = aws_cloudwatch_log_group.ecs["frontend"].name
        "awslogs-region"        = data.aws_region.current.name
        "awslogs-stream-prefix" = "ecs"
      }
    }
    environment = [
      { name = "API_URL", value = "http://api.${var.project_name}-${var.environment}.local:8080" },
      { name = "KRATOS_PUBLIC_URL", value = "http://kratos.${var.project_name}-${var.environment}.local:4433" }
    ]
  }])
}

resource "aws_ecs_service" "frontend" {
  name            = "${var.project_name}-${var.environment}-frontend"
  cluster         = aws_ecs_cluster.main.id
  task_definition = aws_ecs_task_definition.frontend.arn
  desired_count   = 1
  launch_type     = "FARGATE"

  network_configuration {
    subnets          = var.private_subnet_ids
    security_groups  = [aws_security_group.frontend.id]
    assign_public_ip = false
  }

  load_balancer {
    target_group_arn = var.alb_target_group_arns["frontend"]
    container_name   = "frontend"
    container_port   = 3000
  }
}

# API Service
resource "aws_ecs_task_definition" "api" {
  family                   = "${var.project_name}-${var.environment}-api"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = "256"
  memory                   = "512"
  execution_role_arn       = aws_iam_role.ecs_task_execution.arn

  container_definitions = jsonencode([{
    name  = "api"
    image = "${var.ecr_repositories["api"]}:latest"
    portMappings = [{
      containerPort = 8080
      protocol      = "tcp"
    }]
    logConfiguration = {
      logDriver = "awslogs"
      options = {
        "awslogs-group"         = aws_cloudwatch_log_group.ecs["api"].name
        "awslogs-region"        = data.aws_region.current.name
        "awslogs-stream-prefix" = "ecs"
      }
    }
    environment = [
      { name = "DB_HOST", value = split(":", var.db_endpoint)[0] },
      { name = "DB_NAME", value = var.db_name },
      { name = "DB_USER", value = "postgres" },
      { name = "KRATOS_ADMIN_URL", value = "http://kratos.${var.project_name}-${var.environment}.local:4434" }
    ]
    secrets = [{
      name      = "DB_PASSWORD"
      valueFrom = "${var.db_password_secret_arn}:password::"
    }]
  }])
}

resource "aws_ecs_service" "api" {
  name            = "${var.project_name}-${var.environment}-api"
  cluster         = aws_ecs_cluster.main.id
  task_definition = aws_ecs_task_definition.api.arn
  desired_count   = 1
  launch_type     = "FARGATE"

  network_configuration {
    subnets          = var.private_subnet_ids
    security_groups  = [aws_security_group.api.id]
    assign_public_ip = false
  }

  service_registries {
    registry_arn = aws_service_discovery_service.api.arn
  }
}

# Kratos Service
resource "aws_ecs_task_definition" "kratos" {
  family                   = "${var.project_name}-${var.environment}-kratos"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = "256"
  memory                   = "512"
  execution_role_arn       = aws_iam_role.ecs_task_execution.arn

  container_definitions = jsonencode([{
    name  = "kratos"
    image = "${var.ecr_repositories["kratos"]}:latest"
    portMappings = [
      { containerPort = 4433, protocol = "tcp" },
      { containerPort = 4434, protocol = "tcp" }
    ]
    logConfiguration = {
      logDriver = "awslogs"
      options = {
        "awslogs-group"         = aws_cloudwatch_log_group.ecs["kratos"].name
        "awslogs-region"        = data.aws_region.current.name
        "awslogs-stream-prefix" = "ecs"
      }
    }
    environment = [
      { name = "DSN", value = "postgres://postgres@${var.db_endpoint}/${var.db_name}" }
    ]
    secrets = [{
      name      = "DB_PASSWORD"
      valueFrom = "${var.db_password_secret_arn}:password::"
    }]
  }])
}

resource "aws_ecs_service" "kratos" {
  name            = "${var.project_name}-${var.environment}-kratos"
  cluster         = aws_ecs_cluster.main.id
  task_definition = aws_ecs_task_definition.kratos.arn
  desired_count   = 1
  launch_type     = "FARGATE"

  network_configuration {
    subnets          = var.private_subnet_ids
    security_groups  = [aws_security_group.kratos.id]
    assign_public_ip = false
  }

  load_balancer {
    target_group_arn = var.alb_target_group_arns["kratos"]
    container_name   = "kratos"
    container_port   = 4433
  }

  service_registries {
    registry_arn = aws_service_discovery_service.kratos.arn
  }
}

data "aws_region" "current" {}
