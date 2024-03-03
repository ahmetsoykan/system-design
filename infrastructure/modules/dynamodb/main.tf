resource "aws_dynamodb_table" "dynamodb_table" {
  name           = "${var.table_name}"
  billing_mode   = "PROVISIONED"
  read_capacity  = var.read_capacity
  write_capacity = var.write_capacity
  hash_key       = "id"

  attribute {
    name = "id"
    type = "S"
  }

  attribute {
    name = "longurl"
    type = "S"
  }

  ttl {
    attribute_name = "ttl"
    enabled        = true
  }

  # reverse index
  global_secondary_index {
    name               = "gsi1"
    hash_key           = "longurl"
    write_capacity     = var.read_capacity
    read_capacity      = var.write_capacity
    projection_type    = "ALL"
  }

  tags = var.tags
}