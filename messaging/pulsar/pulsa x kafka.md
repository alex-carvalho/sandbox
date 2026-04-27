# Pulsar vs Kafka

## Overview
:contentReference[oaicite:0]{index=0} and :contentReference[oaicite:1]{index=1} are two of the most popular distributed messaging and event streaming platforms. Both are designed for high-throughput, scalable, and fault-tolerant data pipelines, but they differ significantly in architecture and operational model.

---

## Origins

- **Apache Pulsar**
  - Created at Yahoo (2016)
  - Open-sourced under the Apache Software Foundation (2018)

- **Apache Kafka**
  - Created at LinkedIn (2011)
  - Open-sourced shortly after and widely adopted across industries

---

## Architecture

- **Pulsar**
  - Separates compute and storage
  - Uses brokers for serving traffic
  - Uses Apache BookKeeper for durable storage
  - Enables independent scaling of storage and compute

- **Kafka**
  - Combines compute and storage in brokers
  - Each broker stores and serves data
  - Scaling requires rebalancing partitions across brokers

---

## Key Differences

### 1. Storage Model
- **Pulsar**: Segment-based storage via BookKeeper (tiered by design)
- **Kafka**: Log-based storage directly on broker disks

### 2. Scalability
- **Pulsar**: Scales storage and compute independently
- **Kafka**: Scales by adding brokers (tightly coupled)

### 3. Multi-Tenancy
- **Pulsar**: Built-in multi-tenancy with isolation
- **Kafka**: Limited native support (often handled externally)

### 4. Message Retention & Replay
- **Pulsar**: Flexible retention policies and easy replay
- **Kafka**: Strong replay support, retention based on time/size

### 5. Messaging Models
- **Pulsar**:
  - Queue (exclusive, shared, failover)
  - Pub/Sub
  - Streaming

- **Kafka**:
  - Primarily pub/sub and streaming
  - Queue semantics require additional handling

### 6. Geo-Replication
- **Pulsar**: Native and simpler to configure
- **Kafka**: Requires tools like MirrorMaker

---

## Performance Considerations

- **Kafka**
  - Extremely mature and optimized
  - Strong ecosystem and tooling
  - Predictable performance

- **Pulsar**
  - Competitive performance
  - Better for use cases needing flexible scaling
  - Slightly more complex architecture

---

## When to Choose Each

### Choose Pulsar if:
- You need multi-tenancy out of the box
- You want separation of storage and compute
- You need built-in geo-replication
- You want flexible messaging patterns

### Choose Kafka if:
- You want a mature and widely adopted ecosystem
- Your team already has Kafka expertise
- You prefer simpler architecture (fewer moving parts)
- You rely on Kafka-native tools (e.g., Kafka Streams)

---

## Summary

| Feature                | Pulsar                         | Kafka                          |
|----------------------|--------------------------------|--------------------------------|
| Architecture         | Decoupled                      | Coupled                        |
| Multi-Tenancy        | Native                         | Limited                        |
| Geo-Replication      | Built-in                       | External tools                 |
| Scalability          | Independent scaling            | Broker-based scaling           |
| Ecosystem            | Growing                        | Very mature                    |

---

## Conclusion

Both Apache Pulsar and Apache Kafka are powerful platforms for event-driven systems. The right choice depends on your architectural needs, operational preferences, and ecosystem requirements.
