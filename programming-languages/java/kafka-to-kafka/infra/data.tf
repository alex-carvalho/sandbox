data "aws_vpc" "default" {
  default = true
}

// get a subnet id from the default VPC
data "aws_subnets" "default" {
  filter {
    name   = "vpc-id"
    values = [data.aws_vpc.default.id]
  }
}
