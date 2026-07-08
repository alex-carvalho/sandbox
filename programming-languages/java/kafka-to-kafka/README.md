# kafka-to-kafka
Java 25 + Gradle Kafka to Kafka app 

Features:
- Consumes from a Kafka topic and republishes to another
- Dockerfile for containerized build/run
- Unit test using Testcontainers Kafka
- Structured logging with Logback

Run locally (requires Java 25 + Gradle):
```
./gradlew run
```

Build docker (optional):
```
docker build -t kafka-to-kakfa .
docker run --env KAFKA_BOOTSTRAP_SERVERS=host:9092 kafka-to-kafka
```

Tests (requires Docker for Testcontainers):
```
./gradlew test
```
