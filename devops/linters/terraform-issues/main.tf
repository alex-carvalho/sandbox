
# This variable has no type constraint - TFLint will warn about this
variable "instance_type" {
  description = "EC2 instance type"
  default = "t2.micro"
}

resource "aws_instance" "problematic_ec2" {
  # Missing required ami - terraform validate will catch this
  instance_type = var.instance_type
  
  associate_public_ip_address = true   # Checkov will flag this as a security issue

  root_block_device {
    encrypted = false # Checkov will flag unencrypted volumes
    volume_size = 100 
  }  
  
  vpc_security_group_ids = [aws_security_group.wide_open.id] # Security group with all ports open - Checkov will flag this
}

# Security group with dangerous rules
resource "aws_security_group" "wide_open" {
  name        = "allow_all"
  description = "Allow all inbound traffic"

  # Dangerous ingress rule - Checkov will flag this
  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}