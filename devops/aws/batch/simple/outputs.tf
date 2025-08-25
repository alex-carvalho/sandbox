output "job_queue_arn" {
  description = "ARN of the Batch job queue"
  value       = aws_batch_job_queue.queue.arn
}

output "job_definition_arn" {
  description = "ARN of the Batch job definition"
  value       = aws_batch_job_definition.hello.arn
}

output "compute_environment_arn" {
  description = "ARN of the Batch compute environment"
  value       = aws_batch_compute_environment.fargate.arn
}
