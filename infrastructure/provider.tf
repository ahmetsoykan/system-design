provider "aws" {
  region = var.region
}

terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.33.0"
    }
  }

  backend "s3" {
    bucket = "ahmetsoykan-terraform-state"
    key    = "dev/url-shortener/terraform.tfstate"
    region = var.region
  }

  required_version = "~> 1.0"
}

data "aws_availability_zones" "available" {}
