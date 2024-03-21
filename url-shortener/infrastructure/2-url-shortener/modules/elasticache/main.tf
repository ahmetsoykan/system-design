resource "aws_security_group" "default" {
  name_prefix = "url-shortener-cache-sg"
  vpc_id      = var.vpc_id

  ingress {
    from_port   = 6379
    to_port     = 6379
    protocol    = "tcp"
    cidr_blocks = [var.vpc_cidr_block]
  }

  egress {
    from_port   = 6379
    to_port     = 6379
    protocol    = "tcp"
    cidr_blocks = [var.vpc_cidr_block]
  }
}

resource "aws_elasticache_subnet_group" "default" {
  name       = "url-shortener-cache-subnet"
  subnet_ids = var.elasticache_subnet_ids
}

resource "aws_elasticache_replication_group" "default" {
  replication_group_id = "url-shortener-cache"
  description          = "url-shortener cache"
  node_type            = "cache.r7g.large"
  port                 = 6379
  parameter_group_name = "default.redis7.cluster.on"
  engine_version       = "7.1"

  snapshot_retention_limit = 5
  snapshot_window          = "00:00-05:00"

  subnet_group_name          = aws_elasticache_subnet_group.default.name
  automatic_failover_enabled = true

  security_group_ids = [
    aws_security_group.default.id
  ]

  tags = var.tags
}
