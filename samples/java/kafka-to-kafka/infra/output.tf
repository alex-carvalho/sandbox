
output "kafka_cluster_arn" {
  value = aws_msk_cluster.kafka.arn
}

data "aws_msk_bootstrap_brokers" "brokers" {
  cluster_arn = aws_msk_cluster.kafka.arn
}
