variable "vpc_id" {
  description = "The ID of the VPC"
  type        = string
}
variable "vpc_cidr_block" {
  description = "CIDR of the VPC"
  type        = string
}
variable "elasticache_subnet_ids" {
  description = "Subnets ids for elasticache"
  type        = list(any)
}
variable "tags" {
  default = {
    "createdBy"   = "A Team"
    "environment" = "Development"
  }
  type = map(string)
}
