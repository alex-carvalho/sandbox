variable "region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}

variable "fleet_capacity" {
  description = "Target capacity for EC2 fleet"
  type        = number
  default     = 0
}

variable "NODE_NAME" {
  description = "Base name for Jenkins agent nodes. Instance ID will be appended for uniqueness."
  type        = string
  default     = "jenkins-agent"
}

variable "agent_ssh_pubkey" {
  description = "Public SSH key string for Jenkins controller to SSH to agents"
  type        = string
  default     = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIJ1KVixfXK/TTe4x74R40kK9NRE2VYi1PQ8STHuU/zro"
}