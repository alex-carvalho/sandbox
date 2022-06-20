from diagrams import Cluster, Diagram
from diagrams.aws.compute import EC2,ECS, EKS, Lambda
from diagrams.aws.database import RDS, Redshift, ElastiCache, Aurora
from diagrams.aws.network import ELB, Route53
from diagrams.aws.integration import SQS
from diagrams.aws.migration import DMS
from diagrams.aws.storage import S3
from diagrams.aws.security import WAF

from diagrams.onprem.analytics import Spark
from diagrams.onprem.compute import Server
from diagrams.onprem.database import PostgreSQL
from diagrams.onprem.inmemory import Redis
from diagrams.onprem.aggregator import Fluentd
from diagrams.onprem.monitoring import Grafana, Prometheus
from diagrams.onprem.network import Nginx
from diagrams.onprem.queue import Kafka

from diagrams.generic.device import Mobile

with Diagram("Simple Web Service", show=False):
  
    with Cluster("Service Cluster"):
        websvc = [
             EC2("web"),
             EC2("web"),
             EC2("web")]
    
    ELB("ELB") >> websvc >> RDS("userdb")



with Diagram("Complex workflow", show=False):
    dns = Route53("DNS")
    lb = ELB("ELB")
    redis = ElastiCache("redis")
  
    waf = WAF("waf")
    user = Mobile("user")

    with Cluster("DW"):
        dms = DMS("DMS")
        s3 = S3("S3 dw")
        dms >> s3
        Redshift("redshift") >> s3

    with Cluster("Services"):
        svc_group = [EKS ("web1"),
                     EKS ("web2"),
                     EKS ("web3")]

    with Cluster("Cluster Aurora"):
        db_primary = Aurora("userdb")
        db_replica = Aurora("userdb read")
        db_primary - db_replica

    
    user >> dns >> waf >> lb >> svc_group
    svc_group >> db_primary
    svc_group >> redis
    dms >> db_replica
   


with Diagram("onprem Service", show=False):
    ingress = Nginx("ingress")

    metrics = Prometheus("metric")
    metrics << Grafana("monitoring")

    with Cluster("Service Cluster"):
        grpcsvc = [
            Server("grpc1"),
            Server("grpc2"),
            Server("grpc3")]

    with Cluster("Sessions HA"):
        primary = Redis("session")
        primary - Redis("replica") << metrics
        grpcsvc >> primary

    with Cluster("Database HA"):
        primary = PostgreSQL("users")
        primary - PostgreSQL("replica") << metrics
        grpcsvc >> primary

    aggregator = Fluentd("logging")
    aggregator >> Kafka("stream") >> Spark("analytics")

    ingress >> grpcsvc >> aggregator