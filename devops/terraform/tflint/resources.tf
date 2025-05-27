resource "aws_instance" "foo" {
  ami           = "ami-0ff8a91507f77f867"
  instance_type = "t1.2xlarge" # invalid type!
}

resource "aws_s3_bucket" "foos3" {
  bucket = "foo-bar-" # invalid bucket name!
}