data "aws_vpc" "default" {
  default = true
}

data "aws_subnets" "default" {
  filter {
    name   = "vpc-id"
    values = [data.aws_vpc.default.id]
  }
}

data "aws_security_group" "default" {
  filter {
    name   = "vpc-id"
    values = [data.aws_vpc.default.id]
  }

  filter {
    name   = "group-name"
    values = ["default"]
  }
  
}

data "aws_iam_role" "batch_service" {
  name = "AWSServiceRoleForBatch"

}

resource "aws_iam_role" "ecs_task_execution" {
  name = "ecs-task-execution-role-batch-poc"

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


## aws_batch resources

resource "aws_batch_compute_environment" "fargate" {
  name         = "poc-fargate-ce"
  service_role = data.aws_iam_role.batch_service.arn
  type         = "MANAGED"

  compute_resources {
    type               = "FARGATE"
    max_vcpus          = 4
  subnets            = data.aws_subnets.default.ids
    security_group_ids = [data.aws_security_group.default.id]
  }
}

resource "aws_batch_job_queue" "queue" {
  name     = "poc-batch-queue"
  state    = "ENABLED"
  priority = 0

  compute_environment_order {
    order               = 1
    compute_environment = aws_batch_compute_environment.fargate.arn
  }
}

resource "aws_batch_job_definition" "hello" {
  name                  = "poc-hello"
  type                  = "container"
  platform_capabilities = ["FARGATE"]

  container_properties = jsonencode({
    image = "public.ecr.aws/amazonlinux/amazonlinux:latest"
    networkConfiguration = {
      assignPublicIp = "ENABLED"
    },
    resourceRequirements = [
      {
        type  = "VCPU"
        value = "0.25"
      },
      {
        type  = "MEMORY"
        value = "512"
      }
    ]
    command          = ["echo", "Hello from AWS Batch (Fargate)!"]
    executionRoleArn = aws_iam_role.ecs_task_execution.arn
  })
}
