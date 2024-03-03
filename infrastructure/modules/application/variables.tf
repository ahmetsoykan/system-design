variable "name" {
  description = "Application name"
  type        = string
}

variable "region" {
  type = string
}

variable "az" {
  type = any
}

variable "vpc_id" {
  type = string
}

variable "vpc_cidr_block" {
  type = string
}

variable "container_port" {
  type = number
}

variable "private_subnets" {
  description = "List of subnets to associate with the task or service"
  type        = list(string)
}

variable "tags" {
  type = map(string)
}

variable "container_tag" {
  type = string
}

variable "public_subnets" {
  description = "Subnets ids"
  type        = list(any)
}

variable "redis_host" {
  description = "Redis host for connection string"
  type        = string
}

variable "dynamodb_table_arn" {
  description = "DynamoDB Table ARN"
  type        = string
}

variable "target_group_name" {
  type = string
}

variable "listeners_configurations" {
  type = any
}

variable "route53_records" {
  type    = any
  default = {}
}
