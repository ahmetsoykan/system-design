locals {
  container_name = "${var.name}-container"
}

module "ecs_cluster" {
  source = "./cluster"

  cluster_name = var.name

  # Capacity provider
  fargate_capacity_providers = {
    FARGATE = {
      default_capacity_provider_strategy = {
        weight = 50
        base   = 20
      }
    }
    FARGATE_SPOT = {
      default_capacity_provider_strategy = {
        weight = 50
      }
    }
  }

  tags = var.tags
}

module "ecs_service" {
  source = "./service"

  name        = var.name
  cluster_arn = module.ecs_cluster.arn

  desired_count = 3

  cpu    = 2048
  memory = 4096

  tasks_iam_role_statements = [{
    effect    = "Allow"
    actions   = ["dynamodb:*"]
    resources = ["${var.dynamodb_table_arn}", "${var.dynamodb_table_arn}/*"]
  }]

  # Enables ECS Exec
  enable_execute_command = true

  # Container definition(s)
  container_definitions = {

    awscollector = {
      cpu                = 512
      memory             = 1024
      essential          = true
      image              = "public.ecr.aws/aws-observability/aws-otel-collector:latest"
      memory_reservation = 100
      port_mappings = [
        {
          name          = "awscollector"
          containerPort = 4317
          hostPort      = 4317
          protocol      = "tcp"
        }
      ]
      environment = [
        {
          name  = "AWS_REGION"
          value = "${var.region}"
        }
      ]
    }

    (local.container_name) = {
      cpu       = 1024
      memory    = 3072
      essential = true
      image     = var.container_image
      port_mappings = [
        {
          name          = local.container_name
          containerPort = var.container_port
          hostPort      = var.container_port
          protocol      = "tcp"
        }
      ]
      environment = [
        {
          name  = "APP_REGION"
          value = "${var.region}"
        },
        {
          name  = "APP_PORT"
          value = "8080"
        },
        {
          name  = "APP_REDISHOST"
          value = "${var.redis_host}:6379"
        }
      ]
      # Example image used requires access to write to root filesystem
      readonly_root_filesystem = false

      #   dependencies = [{
      #     containerName = "fluent-bit"
      #     condition     = "START"
      #   }]

      enable_cloudwatch_logging = true
      log_configuration = {
        logDriver = "awslogs"
        options = {
          awslogs-region        = var.region
          awslogs-group         = "/aws/service/${var.name}"
          awslogs-stream-prefix = "ecs"
        }
      }

      # linux_parameters = {
      #   capabilities = {
      #     drop = [
      #       "NET_RAW"
      #     ]
      #   }
      # }

      memory_reservation = 100
    }
  }

  service_connect_configuration = {
    namespace = aws_service_discovery_http_namespace.this.arn
    service = {
      client_alias = {
        port     = var.container_port
        dns_name = local.container_name
      }
      port_name      = local.container_name
      discovery_name = local.container_name
    }
  }

  load_balancer = {
    service = {
      target_group_arn = module.alb.target_groups["ex_ecs"].arn
      container_name   = local.container_name
      container_port   = var.container_port
    }
  }

  subnet_ids = var.private_subnets
  security_group_rules = {
    alb_ingress_container_port = {
      type                     = "ingress"
      from_port                = var.container_port
      to_port                  = var.container_port
      protocol                 = "tcp"
      description              = "Service port"
      source_security_group_id = module.alb.security_group_id
    }

    egress_all = {
      type        = "egress"
      from_port   = 0
      to_port     = 0
      protocol    = "-1"
      cidr_blocks = ["0.0.0.0/0"]
    }
  }

  service_tags = var.tags

  tags = var.tags
}

################################################################################
# Supporting Resources
################################################################################

# data "aws_ssm_parameter" "fluentbit" {
#   name = "/aws/service/aws-for-fluent-bit/stable"
# }

resource "aws_service_discovery_http_namespace" "this" {
  name        = var.name
  description = "CloudMap namespace for ${var.name}"
  tags        = var.tags
}

module "alb" {
  source  = "terraform-aws-modules/alb/aws"
  version = "~> 9.0"

  name = var.name

  load_balancer_type = "application"

  # network
  vpc_id  = var.vpc_id
  subnets = var.public_subnets

  # For example only
  enable_deletion_protection = false

  # Security Group
  security_group_ingress_rules = {
    all_http = {
      from_port   = 80
      to_port     = 80
      ip_protocol = "tcp"
      cidr_ipv4   = "0.0.0.0/0"
    },
    all_https = {
      from_port   = 443
      to_port     = 443
      ip_protocol = "tcp"
      cidr_ipv4   = "0.0.0.0/0"
    }
  }

  security_group_egress_rules = {
    all = {
      ip_protocol = "-1"
      cidr_ipv4   = var.vpc_cidr_block
    }
  }

  listeners = var.listeners_configurations

  target_groups = {
    "${var.target_group_name}" = {
      backend_protocol                  = "HTTP"
      backend_port                      = var.container_port
      target_type                       = "ip"
      deregistration_delay              = 5
      load_balancing_cross_zone_enabled = true

      health_check = {
        enabled             = true
        healthy_threshold   = 5
        interval            = 30
        matcher             = "200"
        path                = "/"
        port                = "traffic-port"
        protocol            = "HTTP"
        timeout             = 5
        unhealthy_threshold = 2
      }

      # There's nothing to attach here in this definition. Instead,
      # ECS will attach the IPs of the tasks to this target group
      create_attachment = false
    }
  }

  route53_records = var.route53_records

  tags = var.tags
}

# ECS Task Log Group
resource "aws_cloudwatch_log_group" "logs" {
  name              = "/aws/service/${var.name}"
  retention_in_days = 7
  tags              = var.tags
}
