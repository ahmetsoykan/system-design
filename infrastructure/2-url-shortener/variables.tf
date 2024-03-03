variable "name" {
  description = "Application name"
  type        = string
  default     = "url-shortener"
}

variable "region" {
  type    = string
  default = "eu-west-1"
}

variable "domain_name" {
  type    = string
  default = ""
}

variable "zone_id" {
  type    = string
  default = ""
}

variable "tags" {
  default = {
    "createdBy"   = "A Team"
    "environment" = "Development"
  }
  type = map(string)
}

## Application Layer
variable "container_port" {
  type    = number
  default = 8080
}

variable "container_image" {
  type = string
}
variable "target_group_name" {
  type    = string
  default = "ex_ecs"
}

## Data Layer
variable "table_name" {
  type        = string
  description = "Name of the table to create in DynamoDB"
  default     = "urls"
}

variable "read_capacity" {
  type        = string
  description = "Read capacity of the table"
  default     = 20
}

variable "write_capacity" {
  type        = string
  description = "Write capacity of the table"
  default     = 20
}
