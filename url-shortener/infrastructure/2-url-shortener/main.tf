## Main Terraform ##
locals {
  vpc_cidr = "10.0.0.0/16"
  azs      = slice(data.aws_availability_zones.available.names, 0, 3)
}

## -------------------##
## Network Layer      ##

module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "~> 5.0"

  name = "demo"
  cidr = local.vpc_cidr

  azs             = local.azs
  private_subnets = [for k, v in local.azs : cidrsubnet(local.vpc_cidr, 4, k)]
  public_subnets  = [for k, v in local.azs : cidrsubnet(local.vpc_cidr, 8, k + 48)]

  enable_nat_gateway = true
  single_nat_gateway = true

  tags = var.tags
}

## -------------------##
## Data Layer         ##

module "dynamodb" {
  source         = "./modules/dynamodb"
  table_name     = var.table_name
  read_capacity  = var.read_capacity
  write_capacity = var.write_capacity
}

module "elasticache" {
  source                 = "./modules/elasticache"
  vpc_id                 = module.vpc.vpc_id
  vpc_cidr_block         = module.vpc.vpc_cidr_block
  elasticache_subnet_ids = module.vpc.private_subnets
}

## -------------------##
## Application Layer  ##

module "ecs" {
  source = "./modules/application"
  name   = var.name
  region = var.region
  az     = local.azs

  # container
  container_port  = var.container_port
  container_image = var.container_image

  # network
  vpc_id          = module.vpc.vpc_id
  vpc_cidr_block  = module.vpc.vpc_cidr_block
  private_subnets = module.vpc.private_subnets
  public_subnets  = module.vpc.public_subnets


  # application dependencies
  dynamodb_table_arn = module.dynamodb.dynamo_table_arn
  redis_host         = module.elasticache.configuration_endpoint_address

  # load balancer configurations
  target_group_name = var.target_group_name

  # listeners_configurations
  listeners_configurations = {

    // 1. HTTP Only
    ex_http = {
      port     = 80
      protocol = "HTTP"

      forward = {
        target_group_key = "ex_ecs"
      }
    }

    // 2. HTTPS and HTTPS Redirect
    # ex-http-https-redirect = {
    #   port     = 80
    #   protocol = "HTTP"
    #   redirect = {
    #     port        = "443"
    #     protocol    = "HTTPS"
    #     status_code = "HTTP_301"
    #   }
    # }
    # ex-https = {
    #   port            = 443
    #   protocol        = "HTTPS"
    #   ssl_policy      = "ELBSecurityPolicy-TLS13-1-2-Res-2021-06"
    #   certificate_arn = ${certificate_arn}

    #   forward = {
    #     target_group_key = var.target_group_name
    #   }
    # }
  }

  # r53 records creation for ALB domain name
  # route53_records = {
  #   A = {
  #     name    = var.domain_name
  #     type    = "A"
  #     zone_id = var.zone_id
  #   }
  # }

  # tags
  tags = var.tags
}
