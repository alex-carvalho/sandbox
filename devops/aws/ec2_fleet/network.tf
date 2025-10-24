data "aws_subnet" "selected" {
  count = 3
  filter {
    name   = "availability-zone"
    values = ["us-east-1${substr("abc", count.index, 1)}"]
  }
  filter {
    name   = "default-for-az"
    values = ["true"]
  }
}