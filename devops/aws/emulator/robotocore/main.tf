terraform {
  required_version = ">= 1.5.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  access_key = "123456789012"
  secret_key = "test"
  region     = "us-east-1"

  endpoints {
    s3  = "http://localhost:4566"
    sqs = "http://localhost:4566"
    rds = "http://localhost:4566"
    iam = "http://localhost:4566"
    sts = "http://localhost:4566"
  }

  # Skip real AWS validations and override host AWS_PROFILE
  profile                     = "default"
  skip_credentials_validation = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true
  s3_use_path_style           = true
}

# 1. S3 Bucket
resource "aws_s3_bucket" "app_bucket" {
  bucket        = "robotocore-poc-bucket"
  force_destroy = true
}

# 2. SQS Dead-Letter Queue & Main Queue
resource "aws_sqs_queue" "app_dlq" {
  name = "robotocore-poc-dlq"
}

resource "aws_sqs_queue" "app_queue" {
  name                      = "robotocore-poc-queue"
  message_retention_seconds = 86400
  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.app_dlq.arn
    maxReceiveCount     = 3
  })
}

# 3. RDS PostgreSQL Instance
resource "aws_db_instance" "app_db" {
  identifier          = "robotocore-poc-db"
  allocated_storage   = 20
  engine              = "postgres"
  engine_version      = "15.4"
  instance_class      = "db.t3.micro"
  db_name             = "appdb"
  username            = "dbadmin"
  password            = "localdevpassword123!"
  skip_final_snapshot = true
}

# 4. IAM Role & Policy (S3 + SQS + RDS permissions)
data "aws_iam_policy_document" "assume_role" {
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRole"]
    principals {
      type        = "Service"
      identifiers = ["ec2.amazonaws.com", "lambda.amazonaws.com"]
    }
    principals {
      type        = "AWS"
      identifiers = ["arn:aws:iam::123456789012:root"]
    }
  }
}

resource "aws_iam_role" "app_role" {
  name               = "robotocore-poc-role"
  assume_role_policy = data.aws_iam_policy_document.assume_role.json
}

data "aws_iam_policy_document" "app_permissions" {
  statement {
    sid       = "S3BucketAccess"
    effect    = "Allow"
    actions   = ["s3:GetObject", "s3:PutObject", "s3:DeleteObject", "s3:ListBucket"]
    resources = [aws_s3_bucket.app_bucket.arn, "${aws_s3_bucket.app_bucket.arn}/*"]
  }

  statement {
    sid       = "SQSQueueAccess"
    effect    = "Allow"
    actions   = ["sqs:SendMessage", "sqs:ReceiveMessage", "sqs:DeleteMessage", "sqs:GetQueueAttributes", "sqs:GetQueueUrl"]
    resources = [aws_sqs_queue.app_queue.arn]
  }

  statement {
    sid       = "RDSDescribeConnect"
    effect    = "Allow"
    actions   = ["rds:DescribeDBInstances", "rds-db:connect"]
    resources = [aws_db_instance.app_db.arn]
  }

  statement {
    sid       = "ExplicitDenyDeleteQueueAndBucket"
    effect    = "Deny"
    actions   = ["s3:DeleteBucket", "sqs:DeleteQueue"]
    resources = [aws_s3_bucket.app_bucket.arn, aws_sqs_queue.app_queue.arn]
  }
}

resource "aws_iam_policy" "app_policy" {
  name        = "robotocore-poc-policy"
  description = "Least-privilege policy for S3, SQS, and RDS in Robotocore PoC"
  policy      = data.aws_iam_policy_document.app_permissions.json
}

resource "aws_iam_role_policy_attachment" "app_attach" {
  role       = aws_iam_role.app_role.name
  policy_arn = aws_iam_policy.app_policy.arn
}

output "s3_bucket_name" {
  value = aws_s3_bucket.app_bucket.id
}

output "sqs_queue_url" {
  value = aws_sqs_queue.app_queue.url
}

output "rds_endpoint" {
  value = aws_db_instance.app_db.endpoint
}

output "iam_role_name" {
  value = aws_iam_role.app_role.name
}

output "iam_role_arn" {
  value = aws_iam_role.app_role.arn
}

output "iam_policy_arn" {
  value = aws_iam_policy.app_policy.arn
}
