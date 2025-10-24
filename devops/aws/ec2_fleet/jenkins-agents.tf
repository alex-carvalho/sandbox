resource "aws_launch_template" "jenkins_agent" {
  name_prefix   = "jenkins-agent-"
  image_id      = "ami-0c02fb55956c7d316" # Amazon Linux 2
  instance_type = "t3.small"

  key_name = "ops"

  vpc_security_group_ids = [aws_security_group.jenkins_agent.id]
  iam_instance_profile {
    name = aws_iam_instance_profile.jenkins_agent.name
  }

  user_data = base64encode(<<-EOF
    #!/bin/bash
    set -euo pipefail
    yum update -y
    yum install -y java-17-amazon-corretto-devel git

    # Create a unique node name for this instance using the provided Terraform variable
    # The instance-id from metadata is appended to make the name unique across fleet
    export NODE_NAME="${var.NODE_NAME}-$(curl -s http://169.254.169.254/latest/meta-data/instance-id)"

    # Install latest Jenkins agent JAR
    # wget -O /opt/agent.jar http://${aws_network_interface.jenkins_controller_private.private_ip}:8080/jnlpJars/agent.jar

    # Start Jenkins agent using WebSocket (avoids using -secret in userdata).
    # Run at the end of userdata in background so cloud-init finishes while agent keeps running.
    # nohup java -jar /opt/agent.jar -workDir /tmp -url http://${aws_network_interface.jenkins_controller_private.private_ip}:8080 -webSocket -name "$NODE_NAME" &>/var/log/jenkins-agent.log &
    EOF
  )

  tag_specifications {
    resource_type = "instance"
    tags = {
      Name = "jenkins-agent"
      Type = "jenkins-worker"
    }
  }
}

resource "aws_ec2_fleet" "jenkins_fleet" {
  launch_template_config {
    launch_template_specification {
      launch_template_id = aws_launch_template.jenkins_agent.id
      version            = aws_launch_template.jenkins_agent.latest_version
    }

    override {
      subnet_id = data.aws_subnet.selected[0].id
    }
    override {
      subnet_id = data.aws_subnet.selected[1].id
    }
    override {
      subnet_id = data.aws_subnet.selected[2].id
    }
  }

  target_capacity_specification {
    default_target_capacity_type = "on-demand"
    total_target_capacity        = var.fleet_capacity
  }

  type = "maintain"

  tags = {
    Name = "jenkins-agent-fleet"
  }
}

resource "aws_security_group" "jenkins_agent" {
  name_prefix = "jenkins-agent-"
  description = "Security group for Jenkins agents"
  vpc_id      = data.aws_subnet.selected[0].vpc_id

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "jenkins-agent-sg"
  }
}