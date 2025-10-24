resource "aws_iam_role" "jenkins_agent" {
  name_prefix = "jenkins-agent-"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
        Service = "ec2.amazonaws.com"
      }
    }]
  })
}

resource "aws_iam_role_policy_attachment" "jenkins_agent_ssm" {
  role       = aws_iam_role.jenkins_agent.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore"
}

resource "aws_iam_role_policy" "jenkins_ec2_fleet" {
  name = "jenkins-ec2-fleet-policy"
  role = aws_iam_role.jenkins_agent.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "ec2:DescribeFleets",
          "ec2:DescribeFleetInstances",
          "ec2:ModifyFleet",
          "ec2:DescribeInstances",
          "ec2:DescribeInstanceTypes",
          "ec2:DescribeInstanceStatus",
          "ec2:DescribeAutoScalingGroups",
          "ec2:UpdateAutoScalingGroup",
          "ec2:TerminateInstances",
          "ec2:RunInstances",
          "ec2:CreateTags",
          "ec2:DescribeTags",
          "ec2:DescribeSubnets",
          "ec2:DescribeSecurityGroups",
          "ec2:DescribeLaunchTemplates",
          "ec2:DescribeLaunchTemplateVersions",
          "autoscaling:DescribeAutoScalingGroups",

          "ec2:DescribeSpotFleetRequests",
          "ec2:DescribeSpotFleetInstances",
          "ec2:ModifySpotFleetRequest"
        ]
        Resource = "*"
      }
    ]
  })
}

resource "aws_iam_instance_profile" "jenkins_agent" {
  name_prefix = "jenkins-agent-"
  role        = aws_iam_role.jenkins_agent.name
}