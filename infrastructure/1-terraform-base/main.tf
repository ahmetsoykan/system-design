######################################################################
### ECR

resource "aws_ecr_repository" "ecr" {
  name                 = "url-shortener"
  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }

  force_delete = true

  tags = var.tags
}

######################################################################
### S3

resource "aws_s3_bucket" "state" {
  bucket = "ahmetsoykan-terraform-state"

  lifecycle {
    prevent_destroy = true
  }
  tags = var.tags
}

resource "aws_s3_bucket_versioning" "terraform_state" {
  bucket = aws_s3_bucket.state.id

  versioning_configuration {
    status = "Enabled"
  }
}
