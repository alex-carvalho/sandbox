data "aws_vpc" "default" {
  default = true
}

data "aws_subnets" "default" {
  filter {
    name   = "vpc-id"
    values = [data.aws_vpc.default.id]
  }
}

resource "aws_security_group" "batch_sg" {
  name        = "batch-example-sg"
  description = "Security group for example AWS Batch"
  vpc_id      = data.aws_vpc.default.id

  tags = {
    Name = "batch-example-sg"
  }
}

resource "aws_iam_role" "batch_service" {
  name = "batch-service-role-example"

  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect    = "Allow",
        Principal = { Service = "batch.amazonaws.com" },
        Action    = "sts:AssumeRole"
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "batch_service_attach" {
  role       = aws_iam_role.batch_service.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSBatchServiceRole"
}

resource "aws_iam_role" "ecs_task_execution" {
  name = "ecs-task-execution-role-batch-example"

  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect    = "Allow",
        Principal = { Service = "ecs-tasks.amazonaws.com" },
        Action    = "sts:AssumeRole"
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "ecs_task_execution_attach" {
  role       = aws_iam_role.ecs_task_execution.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

resource "aws_iam_role_policy_attachment" "ecs_task_ecr_read" {
  role       = aws_iam_role.ecs_task_execution.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"
}


## aws_batch resources

resource "aws_batch_compute_environment" "fargate" {
  name         = "example-fargate-ce"
  service_role = aws_iam_role.batch_service.arn
  type         = "MANAGED"

  compute_resources {
    type               = "FARGATE"
    max_vcpus          = 4
  subnets            = data.aws_subnets.default.ids
    security_group_ids = [aws_security_group.batch_sg.id]
  }
}

resource "aws_batch_job_queue" "queue" {
  name     = "example-batch-queue"
  state    = "ENABLED"
  priority = 1

  compute_environment_order {
    order            = 1
    compute_environment = aws_batch_compute_environment.fargate.arn
  }
}

resource "aws_batch_job_definition" "hello" {
  name                  = "example-hello"
  type                  = "container"
  platform_capabilities = ["FARGATE"]

  container_properties = jsonencode({
    image            = "public.ecr.aws/docker/library/busybox:latest"
    vcpus            = 1
    memory           = 1024
    command          = ["echo", "Hello from AWS Batch (Fargate)!"]
    executionRoleArn = aws_iam_role.ecs_task_execution.arn
  })
}
