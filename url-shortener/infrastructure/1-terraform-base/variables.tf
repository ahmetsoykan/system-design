variable "name" {
  description = "Application name"
  type        = string
  default     = "url-shortener"
}

variable "tags" {
  default = {
    "createdBy"   = "A Team"
    "environment" = "Development"
  }
  type = map(string)
}
