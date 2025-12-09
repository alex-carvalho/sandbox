# Elastic Stack Observability Architecture

## Overview

This is a proof-of-concept (POC) implementation of a complete observability stack using the Elastic ecosystem. It captures logs, metrics, traces, and APM data from a Java Spring Boot application running in Kubernetes, providing comprehensive visibility into application performance and behavior.

---

## System Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         KUBERNETES CLUSTER (KinD)                          │
└─────────────────────────────────────────────────────────────────────────────┘

                              ┌──────────────────────┐
                              │   Java Application   │
                              │  (elastic-stack-     │
                              │   demo:8080)         │
                              └──────────┬───────────┘
                                         │
                    ┌────────────────────┼────────────────────┐
                    │                    │                    │
                    ▼                    ▼                    ▼
            ┌──────────────┐    ┌──────────────┐    ┌──────────────┐
            │    Logs      │    │   Metrics    │    │   Traces &   │
            │  (stdout)    │    │  (Micrometer)│    │     APM      │
            │  JSON format │    │ (Prometheus) │    │              │
            └──────┬───────┘    └──────┬───────┘    └──────┬───────┘
                   │                   │                   │
                   │ Collected by      │ Scraped by        │ Sent to
                   │ filestream input  │ Prometheus        │ APM Server
                   │                   │                   │
                   ▼                   ▼                   ▼
        ┌─────────────────────────────────────────────────────────┐
        │          Elastic Agent (DaemonSet)                      │
        │  ┌─────────────────────────────────────────────────┐   │
        │  │  filestream-default: Reads container logs       │   │
        │  │  - Input: /var/log/containers/*.log            │   │
        │  │  - Parser: NDJSON (flattens JSON logs)         │   │
        │  │  - Output: logs-kubernetes.container_logs      │   │
        │  └─────────────────────────────────────────────────┘   │
        │  ┌─────────────────────────────────────────────────┐   │
        │  │  filestream-monitoring: Agent self-monitoring   │   │
        │  │  - Input: Agent internal logs                  │   │
        │  │  - Output: logs-elastic_agent.*                │   │
        │  └─────────────────────────────────────────────────┘   │
        │  ┌─────────────────────────────────────────────────┐   │
        │  │  APM Input: Receives APM data                   │   │
        │  │  - Input: Localhost:8200                        │   │
        │  │  - Output: apm-* indices                        │   │
        │  └─────────────────────────────────────────────────┘   │
        └────────────────────┬────────────────────────────────────┘
                             │
                             ▼
        ┌─────────────────────────────────────────────────────────┐
        │       Elasticsearch (Single Node - 1Gi Storage)         │
        │  ┌─────────────────────────────────────────────────┐   │
        │  │ Data Streams & Indices:                         │   │
        │  │ - logs-kubernetes.container_logs-*             │   │
        │  │ - logs-elastic_agent.*-*                       │   │
        │  │ - apm-*-* (Traces, Transactions)               │   │
        │  │ - metrics-* (System metrics)                   │   │
        │  └─────────────────────────────────────────────────┘   │
        └────────────────────┬────────────────────────────────────┘
                             │
                ┌────────────┼────────────┐
                │            │            │
                ▼            ▼            ▼
        ┌──────────────┐ ┌─────────┐ ┌──────────────┐
        │   Kibana     │ │ Fleet   │ │  APM Server  │
        │   (UI & Viz) │ │ Server  │ │              │
        │   :5601      │ │         │ │   :8200      │
        └──────────────┘ └─────────┘ └──────────────┘
```

---

## Component Details

### 1. **Java Application (Spring Boot)**

**Purpose**: Business logic application that generates logs, metrics, and traces

**Location**: `java-app/src/main/java/com/elastic/test/`

**Key Features**:
- **Logback Logging**: Configured with JSON output for easy parsing
  - Console appender outputs JSON format
  - Different profiles: `local`, `docker`, `kubernetes`
  - Log levels configurable per package
  
- **Micrometer Metrics**: 
  - Request counters
  - Error tracking
  - Response time measurements
  - Prometheus endpoint exposure

- **APM Instrumentation**:
  - Elastic APM Java Agent integrated
  - Automatic transaction tracing
  - Error capture and reporting

**Configuration Files**:
- `src/main/resources/application.yml` - Spring configuration
- `src/main/resources/logback-spring.xml` - Logging configuration
- `build.gradle` - Dependencies (logstash encoder, APM agent, Micrometer)

**Endpoints**:
- `GET /api/hello?name=<name>` - Test endpoint with request/response logging
- `GET /api/health-check` - Health status
- `POST /api/test-error` - Simulate error logging
- `POST /api/test-exception` - Test exception handling
- `GET /api/slow-endpoint` - Test slow request tracing

---

### 2. **Elasticsearch**

**Purpose**: Centralized data store for logs, metrics, and traces

**Configuration**: `terraform/k8s/elasticsearch.yaml`

**Specifications**:
- Version: 8.11.0
- Single node cluster (for POC)
- Storage: 1Gi persistent volume (expandable)
- Hosted in `elastic` namespace

**Data Stored**:
```
logs-kubernetes.container_logs-default-*        # Container logs from all pods
logs-elastic_agent.filebeat-*                   # Agent internal logs
logs-elastic_agent.fleet_server-*               # Fleet server logs
logs-elastic_agent.apm_server-*                 # APM server logs
apm-*-*                                         # APM transactions & traces
metrics-*                                        # System metrics
```

**Index Strategy**: Data streams with automatic rollover and retention policies

---

### 3. **Kibana**

**Purpose**: Web UI for visualization and data exploration

**Configuration**: `terraform/k8s/kibana.yaml`

**Features**:
- Unified interface for logs, metrics, and APM
- Fleet management for agent policies
- Data discovery and ad-hoc querying
- Dashboard creation and visualization

**Access**:
```bash
kubectl port-forward svc/kibana-sample-kb-http 5601:5601 -n elastic
# http://localhost:5601
# Username: elastic
# Password: kubectl get secret elasticsearch-sample-es-elastic-user -n elastic -o jsonpath='{.data.elastic}' | base64 -d
```

---

### 4. **Elastic Agent (DaemonSet)**

**Purpose**: Unified data collector running on every node

**Configuration**: `terraform/k8s/elastic-agent.yaml`

**Mode**: Fleet (policy-based configuration through Kibana)

**Policy**: `eck-agent` defined in Kibana configuration

**Input Configurations**:

#### 4.1 **filestream-default** (Container Log Collection)
```yaml
Type: filestream
Input: /var/log/containers/*${kubernetes.container.id}.log
Parser: NDJSON (NEW - parses JSON logs)
Output: logs-kubernetes.container_logs
```

**How it works**:
1. Reads raw container logs from Kubernetes container runtime logs
2. Applies NDJSON parser to extract JSON fields
3. Enriches logs with Kubernetes metadata (pod name, namespace, container name, etc.)
4. Ships to Elasticsearch as `logs-kubernetes.container_logs-*` data streams

**Key Enhancement**: The NDJSON parser (added in latest config) converts multi-line JSON log entries into properly structured documents:
```json
Input (raw log line):
{"@timestamp":"2025-12-09T14:15:18.846609197Z","message":"Successfully processed","logger_name":"com.elastic.test.controller.DemoController",...}

Output (parsed and indexed):
{
  "@timestamp": "2025-12-09T14:15:18.846609197Z",
  "message": "Successfully processed hello request for: NewTest1",
  "logger_name": "com.elastic.test.controller.DemoController",
  "log.level": "DEBUG",
  "kubernetes": {
    "pod.name": "java-app-cd65c8b9f-dsf67",
    "container.name": "java-app",
    "namespace": "elastic"
  },
  ...
}
```

#### 4.2 **filestream-monitoring** (Agent Self-Monitoring)
```yaml
Type: filestream
Input: Agent internal state logs
Output: logs-elastic_agent-*
```

#### 4.3 **apm Input** (APM Server)
```yaml
Type: apm
Listen: 0.0.0.0:8200
Output: apm-* indices
```

**Volume Mounts**:
- `/var/log` - Host logs
- `/var/lib/docker/containers` - Docker container logs
- `/var/log/containers` - Kubernetes container logs
- `/var/log/pods` - Kubernetes pod logs

---

### 5. **Fleet Server**

**Purpose**: Central management server for Elastic Agents

**Configuration**: `terraform/k8s/fleet-server.yaml`

**Responsibilities**:
- Manages agent policies and enrollment
- Distributes integration packages
- Handles agent communication
- Provides central control plane

---

### 6. **APM Server**

**Purpose**: Receives and processes application performance monitoring data

**Runs as part of**: Elastic Agent (embedded apm input)

**Receives**:
- Transactions (request traces)
- Spans (detailed operation timing)
- Errors
- Metrics

**Sends to**: Elasticsearch as `apm-*` indices

---

## Data Flow Diagram

```
┌──────────────────────────────────────────────────────────────────────────┐
│                          DATA COLLECTION FLOWS                           │
└──────────────────────────────────────────────────────────────────────────┘

LOGS COLLECTION FLOW:
═════════════════════

Java App
  │
  ├─ stdout (JSON formatted by Logback)
  │         │
  │         └─→ Container Runtime
  │              │
  │              └─→ /var/log/containers/java-app-*.log
  │                   │
  │                   └─→ Elastic Agent (filestream-default)
  │                        │
  │                        ├─ Read raw log lines
  │                        ├─ Parse NDJSON (flatten JSON)
  │                        ├─ Enrich with K8s metadata
  │                        │
  │                        └─→ Elasticsearch
  │                            logs-kubernetes.container_logs-*
  │
  └─→ Queryable in Kibana
      - Search by message, logger_name, log.level
      - Filter by pod, namespace, container
      - Correlate with traces via trace.id


METRICS COLLECTION FLOW:
═════════════════════════

Java App
  │
  ├─ Micrometer metrics (in-memory)
  │  │
  │  └─→ /actuator/metrics/prometheus endpoint
  │      │
  │      └─→ Prometheus (scrape configured)
  │          │
  │          └─→ Elasticsearch (Prometheus exporter)
  │              │
  │              └─→ metrics-*
  │
  └─→ Queryable in Kibana
      - Request rates
      - Error counts
      - Response times (percentiles)


APM/TRACES COLLECTION FLOW:
════════════════════════════

Java App
  │
  ├─ Elastic APM Agent (javaagent)
  │  │
  │  ├─ Intercepts HTTP requests → Transactions
  │  ├─ Intercepts method calls → Spans
  │  ├─ Catches exceptions → Errors
  │  │
  │  └─→ HTTP POST to APM Server (apm.elastic.svc:8200)
  │      │
  │      └─→ Elastic Agent (apm input)
  │          │
  │          └─→ Elasticsearch
  │              apm-*-transaction-*
  │              apm-*-error-*
  │              apm-*-metrics-*
  │
  └─→ Queryable in Kibana
      - Full request traces
      - Service map
      - Error tracking
      - Performance analysis

```

---

## Component Connections Map

```
┌─────────────────────────────────────────────────────────────────┐
│                    CONNECTIONS & PROTOCOLS                      │
└─────────────────────────────────────────────────────────────────┘

Java Application
    │
    ├─ :8080/actuator/metrics/prometheus
    │  └─→ [HTTP] Prometheus scraper
    │
    ├─ :8200 (APM Agent embedded)
    │  └─→ [HTTP/gRPC] APM Server @ apm.elastic.svc:8200
    │
    └─ stdout
       └─→ Kubernetes container runtime
          └─→ /var/log/containers/

Elastic Agent DaemonSet (Kubernetes)
    │
    ├─ filestream input (filestream-default)
    │  └─→ Reads: /var/log/containers/*${kubernetes.container.id}.log
    │
    ├─ apm input (embedded APM Server)
    │  ├─ Listens: 0.0.0.0:8200
    │  └─ Receives: apm-server.elastic.svc:8200 requests from Java apps
    │
    └─→ [HTTP] :9200 → elasticsearch-sample-es-http.elastic.svc
       └─→ Index logs, APM traces, metrics

Fleet Server
    │
    ├─ Communicates with: Kibana (policy management)
    ├─ Communicates with: Elastic Agents (policy distribution)
    └─→ Hostname: fleet-server-agent-http.elastic.svc:8220

Kibana
    │
    ├─→ [HTTPS] elasticsearch-sample-es-http.elastic.svc:9200
    ├─→ [HTTPS] fleet-server-agent-http.elastic.svc:8220
    └─ UI accessible at :5601 (after port-forward)

Elasticsearch
    │
    └─→ Persistent storage: PVC elasticsearch-data (1Gi)
```

---

## Kubernetes Namespaces & Services

```
NAMESPACE: elastic
═════════════════

Deployments:
  ├─ java-app (1 replica)
  │  └─ Pod: java-app-<hash>-<pod-id>
  │     └─ Container: java-app:latest
  │        ├─ Port: 8080 (HTTP)
  │        ├─ Env: SPRING_PROFILES_ACTIVE=kubernetes
  │        ├─ Env: ELASTIC_APM_SERVER_URLS=http://apm.elastic.svc:8200
  │        └─ Probes: liveness & readiness
  │
  └─ kibana-sample
     └─ Pod: kibana-sample-<hash>
        ├─ Port: 5601 (HTTP)
        └─ Reference: elasticsearch-sample-es-http

StatefulSets:
  └─ elasticsearch-sample
     └─ Pod: elasticsearch-sample-0
        ├─ Port: 9200 (HTTPS API)
        ├─ Port: 9300 (node communication)
        └─ Storage: 1Gi

DaemonSets:
  └─ elastic-agent
     └─ Pod per node: elastic-agent-<hash>
        ├─ Mount: /var/log (readOnly)
        ├─ Mount: /var/lib/docker/containers (readOnly)
        ├─ Mount: /var/log/containers (readOnly)
        └─ Mode: fleet (policy-managed)

Services:
  ├─ java-app (ClusterIP:8080) → java-app pod
  ├─ apm (ClusterIP:8200) → elastic-agent pod
  ├─ elasticsearch-sample-es-http (ClusterIP:9200) → elasticsearch pod
  ├─ kibana-sample-kb-http (ClusterIP:5601) → kibana pod
  └─ fleet-server-agent-http (ClusterIP:8220) → fleet-server
```

---

## Configuration Management

```
Terraform (IaC)
═══════════════

├─ main.tf
│  └─ Define Kubernetes cluster configuration
│
├─ providers.tf
│  └─ Configure Kubernetes provider & ECK operator
│
├─ versions.tf
│  └─ Specify Terraform & provider versions
│
├─ variables.tf
│  └─ Input variables for cluster customization
│
├─ outputs.tf
│  └─ Output cluster details
│
└─ k8s/ (Kubernetes manifests)
   ├─ elasticsearch.yaml
   │  └─ Elasticsearch CRD resource
   │
   ├─ kibana.yaml
   │  └─ Kibana CRD + Fleet Server policies
   │     ├─ system-1 (system integration)
   │     ├─ apm-1 (APM server)
   │     └─ kubernetes-1 (container log collection)
   │
   ├─ fleet-server.yaml
   │  └─ Fleet Server agent deployment
   │
   └─ elastic-agent.yaml
      └─ Elastic Agent DaemonSet
         └─ Mounted volumes for log collection

Java Application (Docker)
═════════════════════════

├─ Dockerfile
│  ├─ Base: eclipse-temurin:25-jdk
│  ├─ Copy: application JAR & logback config
│  ├─ Download: Elastic APM agent
│  └─ Command: java -javaagent:elastic-apm-agent.jar ...
│
├─ build.gradle
│  ├─ Dependencies:
│  │  ├─ spring-boot:web
│  │  ├─ spring-boot:actuator
│  │  ├─ micrometer-prometheus
│  │  ├─ logstash-logback-encoder
│  │  └─ elastic-apm-agent
│  │
│  └─ Tasks:
│     └─ downloadApmAgent
│
├─ src/main/resources/
│  ├─ application.yml
│  │  ├─ Server port: 8080
│  │  ├─ Actuator endpoints
│  │  ├─ Micrometer config
│  │  ├─ Logging config
│  │  └─ APM settings
│  │
│  └─ logback-spring.xml
│     ├─ Profiles:
│     │  ├─ local: Plain text console (colorized)
│     │  ├─ docker: JSON console (LogstashEncoder)
│     │  └─ kubernetes: JSON console (LogstashEncoder)
│     │
│     └─ Custom fields in JSON
│        ├─ service: elastic-stack-demo
│        └─ environment: (from spring.profiles.active)
│
└─ src/main/java/com/elastic/test/
   ├─ ElasticStackDemoApplication.java (Bootstrap)
   ├─ MetricsConfiguration.java (Metric definitions)
   └─ controller/DemoController.java (REST endpoints)
      ├─ Logging on each request
      ├─ Metric increments
      ├─ Timer tracking
      └─ Error simulation endpoints
```

---

## Deployment Sequence

```
DEPLOYMENT ORDER & DEPENDENCIES:
═════════════════════════════════

1. Infrastructure (terraform apply)
   │
   ├─→ Kubernetes Cluster (KinD)
   │
   ├─→ Elasticsearch
   │   └─ Waits for: ECK operator
   │
   ├─→ Kibana + Fleet Server
   │   ├─ Depends on: Elasticsearch running
   │   └─ Configures: Fleet policies & integrations
   │
   └─→ Elastic Agent
       ├─ Depends on: Fleet Server ready
       ├─ Depends on: Kibana policies defined
       └─ Policy: eck-agent

2. Java Application (java-app build & deploy)
   │
   ├─→ Build: gradle clean build → JAR artifact
   │
   ├─→ Docker: docker build → java-app:latest image
   │
   └─→ Deploy: kubectl apply -f deployment.yaml
       ├─ Depends on: Elasticsearch running
       ├─ Depends on: APM Server reachable (apm.elastic.svc:8200)
       └─ Sends logs to: stdout (collected by agent)

3. Verification
   │
   ├─→ Check pod status: kubectl get pods -n elastic
   │
   ├─→ Check logs: kubectl logs -n elastic -l app=java-app
   │
   ├─→ Test endpoint: kubectl exec -it <pod> -- curl http://java-app:8080/api/hello
   │
   ├─→ Port forward: kubectl port-forward svc/kibana-sample-kb-http 5601:5601
   │
   └─→ Verify in Kibana: Check logs-kubernetes.container_logs-* index
```

---

## Log Flow in Detail

```
DETAILED LOG JOURNEY FROM APP TO KIBANA:
═════════════════════════════════════════

Step 1: LOG GENERATION (Java Application)
──────────────────────────────────────────
Application code:
  logger.info("Received request for hello endpoint with name: {}", name);

Logback processes:
  ├─ Thread context: http-nio-8080-exec-1
  ├─ Level: INFO
  ├─ Logger name: com.elastic.test.controller.DemoController
  └─ Formatted as JSON by LogstashEncoder:
     {
       "@timestamp": "2025-12-09T14:15:18.846609197Z",
       "@version": "1",
       "message": "Received request for hello endpoint with name: NewTest1",
       "logger_name": "com.elastic.test.controller.DemoController",
       "thread_name": "http-nio-8080-exec-2",
       "level": "INFO",
       "level_value": 20000,
       "transaction.id": "98090104b66ca9b1",
       "trace.id": "92e8a2419488f3dfbe587940b29b7767",
       "service": "elastic-stack-demo",
       "environment": "kubernetes"
     }


Step 2: LOG ROUTING (Container Runtime)
────────────────────────────────────────
Java app writes to: STDOUT
                    │
                    └─→ Kubernetes container runtime captures
                        │
                        └─→ Writes to: /var/log/containers/
                            File name: <pod-name>_<namespace>_<container>-<container-id>.log


Step 3: LOG COLLECTION (Elastic Agent - filestream input)
──────────────────────────────────────────────────────────
Elastic Agent watches:
  └─ /var/log/containers/*${kubernetes.container.id}.log

For each java-app-*.log file:
  ├─ Read line (raw log entry)
  │
  ├─ Parse with NDJSON parser
  │  └─ Extract: message, logger_name, level, trace.id, etc.
  │
  ├─ Enrich with Kubernetes metadata
  │  ├─ kubernetes.pod.name: java-app-<hash>
  │  ├─ kubernetes.pod.namespace: elastic
  │  ├─ kubernetes.container.name: java-app
  │  ├─ kubernetes.node.name: <node-name>
  │  └─ host.name: <node-hostname>
  │
  └─ Transform into structured document:
     {
       "@timestamp": "2025-12-09T14:15:18.846Z",
       "message": "Received request for hello endpoint with name: NewTest1",
       "logger_name": "com.elastic.test.controller.DemoController",
       "log.level": "info",
       "log.origin": { "file.name": "DemoController.java", "file.line": 26 },
       "service": { "name": "elastic-stack-demo" },
       "kubernetes": {
         "pod": { "name": "java-app-cd65c8b9f-dsf67", "namespace": "elastic" },
         "container": { "name": "java-app", "id": "1605877f-eabd-41de..." },
         "node": { "name": "elastic-control-plane" }
       },
       "host": { "name": "elastic-control-plane" },
       "agent": { "id": "...", "name": "elastic-agent", "type": "filebeat" },
       "ecs": { "version": "8.10.0" },
       "data_stream": {
         "namespace": "default",
         "type": "logs",
         "dataset": "kubernetes.container_logs"
       }
     }


Step 4: TRANSMISSION (Elasticsearch)
─────────────────────────────────────
Elastic Agent sends document via:
  └─ HTTP POST → https://elasticsearch-sample-es-http.elastic.svc:9200
     │
     └─→ Endpoint: /_bulk
         │
         └─→ Index: .ds-logs-kubernetes.container_logs-default-YYYY.MM.DD-000001
             (Auto-created by data stream template)


Step 5: INDEXING (Elasticsearch)
──────────────────────────────────
Elasticsearch:
  ├─ Receives bulk request
  ├─ Parses JSON documents
  ├─ Applies index mapping (field types)
  ├─ Inverts text (full-text search capability)
  └─ Stores document in shard


Step 6: QUERYING (Kibana)
──────────────────────────
User in Kibana:
  └─ Creates data view for "logs-kubernetes.container_logs-*"
     │
     └─→ Can search:
         ├─ By message text: "hello endpoint"
         ├─ By log level: log.level:INFO
         ├─ By pod: kubernetes.pod.name:"java-app*"
         ├─ By namespace: kubernetes.pod.namespace:elastic
         └─ By service: service.name:"elastic-stack-demo"

Results include:
  ├─ All original log fields
  ├─ Kubernetes enrichment
  ├─ Host information
  ├─ Trace correlation (trace.id)
  └─ Timestamp for timeline view
```

---

## Observability Pillars

### 1. **Logging** 📋
- **Source**: Java application logger (Logback)
- **Format**: Structured JSON (LogstashEncoder)
- **Collection**: Elastic Agent filestream input
- **Storage**: `logs-kubernetes.container_logs-*` indices
- **Query**: Message content, log level, logger name, Kubernetes metadata
- **Retention**: Configurable via ILM (Index Lifecycle Management)

### 2. **Metrics** 📊
- **Source**: Spring Boot Actuator + Micrometer
- **Types**:
  - Request counts
  - Error counts
  - Response time distributions (p50, p95, p99)
- **Endpoints**: `/actuator/metrics/prometheus`
- **Scraping**: Prometheus (separate component, can be integrated)
- **Storage**: Metrics indices in Elasticsearch
- **Visualization**: Dashboards in Kibana

### 3. **Traces & APM** 🔍
- **Source**: Elastic APM Java Agent
- **Captures**:
  - HTTP transactions (request path, method, status, duration)
  - Span details (method calls, database queries)
  - Exception details with stack traces
- **Transmission**: Sent to APM Server (embedded in Elastic Agent)
- **Storage**: `apm-*` indices in Elasticsearch
- **Visualization**: 
  - Service map (dependencies)
  - Trace timeline (request flow)
  - Error tracking
  - Performance insights

---

## Key Technologies & Versions

```
Component                Version    Purpose
════════════════════════════════════════════════════════════════
Elasticsearch            8.11.0     Data storage & search engine
Kibana                   8.11.0     Visualization & exploration UI
Elastic Agent            8.11.0     Unified data collection
Fleet Server             8.11.0     Agent orchestration
APM Server               8.11.0     Trace aggregation (embedded)
ECK Operator             2.x        Kubernetes automation

Java/Spring Boot         25/4.0     Application framework
Logback                  1.x        Logging framework
LogstashEncoder          7.4        JSON log encoding
Micrometer               Core       Metrics collection
Elastic APM Agent        1.55.1     Performance monitoring

Terraform                1.x        Infrastructure as Code
Kubernetes (KinD)        Latest     Cluster orchestration
Docker                   Latest     Container runtime
```

---

## Troubleshooting Guide

### Logs Not Appearing?
1. **Check Java app is running**:
   ```bash
   kubectl logs -n elastic -l app=java-app
   ```

2. **Check Elastic Agent is collecting**:
   ```bash
   kubectl logs -n elastic -l agent.k8s.elastic.co/name=elastic-agent
   ```

3. **Verify Elasticsearch is receiving**:
   ```bash
   kubectl port-forward -n elastic svc/elasticsearch-sample-es-http 9200:9200
   ELASTIC_PASSWORD=$(kubectl get secret -n elastic elasticsearch-sample-es-elastic-user -o jsonpath='{.data.elastic}' | base64 -d)
   curl -k -u elastic:$ELASTIC_PASSWORD https://localhost:9200/.ds-logs-kubernetes.container_logs-*/_count
   ```

### Traces Not Showing?
1. **Check APM Server is reachable**:
   ```bash
   kubectl logs -n elastic java-app-* | grep "APM Server"
   ```

2. **Verify APM configuration**:
   ```bash
   kubectl describe pod -n elastic java-app-* | grep -A 10 "Environment"
   ```

3. **Test APM endpoint**:
   ```bash
   kubectl exec -it -n elastic java-app-* -- curl -v http://apm.elastic.svc:8200
   ```
