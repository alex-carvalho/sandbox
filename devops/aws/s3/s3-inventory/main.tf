provider "aws" {
  region = "us-east-1"
}

resource "aws_s3_bucket" "source" {
  bucket = "my-source-bucket-example"
}

resource "aws_s3_bucket" "destination" {
  bucket = "my-destination-bucket-example"
}

resource "aws_s3_bucket_policy" "destination_policy" {
  bucket = aws_s3_bucket.destination.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Service = "s3.amazonaws.com"
        }
        Action = "s3:PutObject"
        Resource = "${aws_s3_bucket.destination.arn}/*"
        Condition = {
          StringEquals = {
            "aws:SourceAccount" = data.aws_caller_identity.current.account_id
          }
        }
      }
    ]
  })
}

data "aws_caller_identity" "current" {}

resource "aws_s3_bucket_inventory" "inventory" {
  bucket = aws_s3_bucket.source.id
  name   = "daily-inventory"

  included_object_versions = "Current"
  schedule {
    frequency = "Daily"
  }

  destination {
    bucket {
      format     = "CSV"
      bucket_arn = aws_s3_bucket.destination.arn
      prefix     = "inventory-reports/"
    }
  }

  optional_fields = [
    "Size",
    "LastModifiedDate",
    "StorageClass",
    "ETag",
    "EncryptionStatus",
    "ReplicationStatus"
  ]
}
