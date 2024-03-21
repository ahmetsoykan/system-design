variable "table_name" {
  description = "define table name"
}

variable "read_capacity" {
  description = "define read capacity"
}

variable "write_capacity" {
  description = "define write capacity"
}
variable "partition_key" {
  type        = string
  description = "Write the partition key"
  default     = "id"
}
variable "partition_key_type" {
  type        = string
  description = "Write the partition key type"
  default     = "S"
}
variable "tags" {
  default = {
    "createdBy"   = "A Team"
    "environment" = "Development"
  }
  type = map(string)
}
