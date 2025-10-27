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
    yum install -y java-17-amazon-corretto-devel git openssh-server openssh-clients

    # Create a unique node name for this instance using the provided Terraform variable
    # The instance-id from metadata is appended to make the name unique across fleet
    export NODE_NAME="${var.NODE_NAME}-$(curl -s http://169.254.169.254/latest/meta-data/instance-id)"

    # Configure SSH so the Jenkins controller can connect (controller holds the private key)
    AGENT_USER=ec2-user
    PUBKEY='${var.agent_ssh_pubkey}'

    # Ensure sshd is enabled
    systemctl enable sshd || true
    systemctl restart sshd || true

    # Create .ssh and add controller public key for the ec2-user (idempotent, do not remove existing keys)
    SSH_DIR="/home/$${AGENT_USER}/.ssh"
    AUTH_FILE="$${SSH_DIR}/authorized_keys"
    mkdir -p "$${SSH_DIR}"
    chmod 700 "$${SSH_DIR}"
    touch "$${AUTH_FILE}"
    # Append the public key only if it's not already present
    if ! grep -qxF "$PUBKEY" "$${AUTH_FILE}"; then
      echo "$PUBKEY" >> "$${AUTH_FILE}"
    fi
    chmod 600 "$${AUTH_FILE}"
    chown -R $${AGENT_USER}:$${AGENT_USER} "$${SSH_DIR}"

    # Note: we do not start the inbound JNLP agent here because Jenkins will connect via SSH.
    # If you later want to start the agent via JNLP/WebSocket, uncomment and adjust the lines below.
    # Install latest Jenkins agent JAR (optional)
    # wget -O /opt/agent.jar http://${aws_network_interface.jenkins_controller_private.private_ip}:8080/jnlpJars/agent.jar

    # Example JNLP start (commented):
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

  ingress {
    from_port       = 22
    to_port         = 22
    protocol        = "tcp"
    security_groups = [aws_security_group.jenkins_controller.id]
    description     = "Allow SSH from Jenkins controller"
  }

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

resource "aws_autoscaling_group" "jenkins_agents_asg" {
  name                      = "jenkins-agents-asg"
  max_size                  = "5"
  min_size                  = "0"
  desired_capacity          = "0"
  vpc_zone_identifier       = [
    data.aws_subnet.selected[0].id,
    data.aws_subnet.selected[1].id,
    data.aws_subnet.selected[2].id,
  ]

  launch_template {
    id      = aws_launch_template.jenkins_agent.id
    version = aws_launch_template.jenkins_agent.latest_version
  }

  health_check_type         = "EC2"
  health_check_grace_period = 300

  tag {
    key                 = "Name"
    value               = "jenkins-agent"
    propagate_at_launch = true
  }

  lifecycle {
    create_before_destroy = true
  }
}