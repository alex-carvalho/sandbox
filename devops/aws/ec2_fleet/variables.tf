variable "region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}

variable "fleet_capacity" {
  description = "Target capacity for EC2 fleet"
  type        = number
  default     = 1
}

variable "NODE_NAME" {
  description = "Base name for Jenkins agent nodes. Instance ID will be appended for uniqueness."
  type        = string
  default     = "jenkins-agent"
}