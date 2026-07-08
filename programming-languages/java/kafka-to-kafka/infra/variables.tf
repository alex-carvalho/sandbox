
variable "aws_region" {
  type    = string
  default = "us-east-1"
}

variable "cluster_name" {
  type    = string
  default = "kafka-poc"
}

variable "allowed_cidr" {
  type    = string
  default = "0.0.0.0/0"
}

variable "ebs_volume_gb" {
  type    = number
  default = 10
}