resource "aws_msk_cluster" "kafka" {
  cluster_name  = var.cluster_name
  kafka_version = "3.9.x"

  number_of_broker_nodes = 2  

  broker_node_group_info {
    instance_type   = "kafka.t3.small"
    client_subnets  = [data.aws_subnets.default.ids[0], data.aws_subnets.default.ids[1]]
    security_groups = [aws_security_group.msk_sg.id]

    storage_info {
      ebs_storage_info {
        volume_size = var.ebs_volume_gb
      }
    }
  }

  encryption_info {
    encryption_in_transit {
      client_broker = "TLS_PLAINTEXT"
      in_cluster    = true
    }
  }

  logging_info {
    broker_logs {
      cloudwatch_logs {
        enabled = false
      }
      firehose {
        enabled = false
      }
      s3 {
        enabled = false
      }
    }
  }

  tags = {
    ManagedBy   = "terraform"
  }
}


resource "aws_security_group" "msk_sg" {
  name        = "${var.cluster_name}-sg"
  description = "Allow Kafka client traffic"
  vpc_id      = data.aws_vpc.default.id

  ingress {
    description = "Kafka port"
    from_port   = 9092
    to_port     = 9092
    protocol    = "tcp"
    cidr_blocks = [var.allowed_cidr]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "${var.cluster_name}-sg"
  }
}