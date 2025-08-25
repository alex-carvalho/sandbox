# AWS Batch (Fargate) Terraform example

This example creates a minimal AWS Batch setup using Fargate:
- IAM roles (Batch service role, ECS task execution role)
- A managed Fargate compute environment
- A job queue
- A simple job definition that runs a container which echoes a message

Assumptions
- Uses the default VPC and its subnets in the target AWS region.

Quick start

```bash
terraform init
terraform validate
terraform plan
terraform apply
```

- After apply, note the job queue and job definition ARN from the outputs. Submit a job using the AWS CLI (example):

```bash
aws batch submit-job \
  --region us-east-1 \
  --job-name example-hello-job \
  --job-queue <JOB_QUEUE_ARN_OR_NAME> \
  --job-definition <JOB_DEFINITION_ARN_OR_NAME> \
  --platform-capabilities FARGATE
```

Replace `<JOB_QUEUE_ARN_OR_NAME>` and `<JOB_DEFINITION_ARN_OR_NAME>` with the outputs from Terraform.


Security & static scanning

```bash
terrascan scan -t aws
trivy config --exit-code 0 --no-progress .
```
