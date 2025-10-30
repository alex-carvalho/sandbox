# Filesystem ingestion using kafka connector

## Overview
Apache NiFi pipeline that automatically reads CSV files from `/data/input`, parses them, and sends to Kafka.

## Quick Start

```bash
# Start services
docker-compose up -d

# Check logs
docker-compose logs -f nifi

# Access UIs
# NiFi: https://localhost:8443/nifi (admin/adminadminadmin)
# Kafka UI: http://localhost:8090
```

## Pipeline Flow
1. **GetFile** - Reads CSV files from `/data/input`
2. **PublishKafka** - Sends data to Kafka topic `csv-data`

## Add CSV Files
Place CSV files in `./my-data/` directory - they'll be automatically processed.

## Verify
```bash
# Check Kafka topic
docker exec -it kafka kafka-console-consumer --bootstrap-server localhost:9092 --topic csv-data --from-beginning
```
