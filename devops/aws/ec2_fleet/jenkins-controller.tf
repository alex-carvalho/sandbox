resource "aws_network_interface" "jenkins_controller_private" {
  subnet_id       = data.aws_subnet.selected[0].id
  security_groups = [aws_security_group.jenkins_controller.id]

  tags = {
    Name = "jenkins-controller-private-eni"
  }
}

resource "aws_instance" "jenkins_controller" {
  ami               = "ami-0c02fb55956c7d316" # Amazon Linux 2
  instance_type     = "t3.medium"
  availability_zone = data.aws_subnet.selected[0].availability_zone

  iam_instance_profile   = aws_iam_instance_profile.jenkins_agent.name
  vpc_security_group_ids = [aws_security_group.jenkins_controller.id]


  user_data = base64encode(<<-EOF
    #!/bin/bash
    yum update -y
    yum install -y java-17-amazon-corretto-devel wget
    
    # Install latest Jenkins LTS
    wget -O /etc/yum.repos.d/jenkins.repo https://pkg.jenkins.io/redhat-stable/jenkins.repo
    rpm --import https://pkg.jenkins.io/redhat-stable/jenkins.io-2023.key
    yum install -y jenkins
    
    # Configure Jenkins to use private IP
    #echo 'JENKINS_ARGS="--httpListenAddress=${aws_network_interface.jenkins_controller_private.private_ip}"' >> /etc/sysconfig/jenkins
    
    systemctl enable jenkins
    systemctl start jenkins
    EOF
  )

  tags = {
    Name = "jenkins-controller"
    Type = "jenkins-master"
  }
}

resource "aws_network_interface_attachment" "jenkins_controller_attach" {
  instance_id          = aws_instance.jenkins_controller.id
  network_interface_id = aws_network_interface.jenkins_controller_private.id
  device_index         = 1
}

resource "aws_security_group" "jenkins_controller" {
  name_prefix = "jenkins-controller-"
  description = "Security group for Jenkins controller"
  vpc_id      = data.aws_subnet.selected[0].vpc_id

  ingress {
    from_port   = 8080
    to_port     = 8080
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port       = 50000
    to_port         = 50000
    protocol        = "tcp"
    security_groups = [aws_security_group.jenkins_agent.id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "jenkins-controller-sg"
  }
}
