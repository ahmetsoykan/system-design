provider "aws" {
  region = "eu-west-1"
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
    region = "eu-west-1"
  }

  required_version = "~> 1.0"
}

data "aws_availability_zones" "available" {}
