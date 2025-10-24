output "fleet_id" {
  description = "EC2 Fleet ID for Jenkins plugin configuration"
  value       = aws_ec2_fleet.jenkins_fleet.id
}

output "fleet_arn" {
  description = "EC2 Fleet ARN"
  value       = "arn:aws:ec2:${var.region}:${data.aws_caller_identity.current.account_id}:fleet/${aws_ec2_fleet.jenkins_fleet.id}"
}

data "aws_caller_identity" "current" {}

output "launch_template_id" {
  description = "Launch template ID"
  value       = aws_launch_template.jenkins_agent.id
}

output "jenkins_controller_private_ip" {
  description = "Jenkins controller private IP"
  value       = aws_network_interface.jenkins_controller_private.private_ip
}

output "jenkins_controller_public_ip" {
  description = "Jenkins controller public IP"
  value       = aws_instance.jenkins_controller.public_ip
}

output "jenkins_private_url" {
  description = "Jenkins private URL for agents"
  value       = "http://${aws_network_interface.jenkins_controller_private.private_ip}:8080"
}

output "jenkins_public_url" {
  description = "Jenkins public URL"
  value       = "http://${aws_instance.jenkins_controller.public_ip}:8080"
}

output "vpc_id" {
  description = "VPC ID"
  value       = data.aws_subnet.selected[0].vpc_id
}

output "subnet_ids" {
  description = "Subnet IDs used"
  value       = data.aws_subnet.selected[*].id
}