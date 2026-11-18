## Slide: Title Slide
**Title:** Grafana LGTM

- End-to-end observability: **Loki, Grafana, Tempo, Mimir**

### Why Traditional Monitoring Falls Short

- Siloed tools for logs, metrics, and traces slow down root cause analysis
- Manual timestamp & ID matching increases **MTTR**
- Proprietary APM agents create vendor lock-in
- **Solution:** OpenTelemetry standard + unified **LGTM** correlation

---

## Slide: Grafana LGTM Stack Overview
**Title:** Grafana LGTM Stack Overview

- **Loki:** Cost-effective log aggregation ("Prometheus for logs")
- **Grafana:** Single pane of glass for dashboards & correlations
- **Tempo:** Trace storage without heavy indexing + auto **Span Metrics**
- **Mimir:** Horizontally scalable Prometheus TSDB with **Exemplars**

- **Grafana Alloy:** Unified OTel collector (replaces Promtail & Grafana Agent)

---

## Slide: Deep Dive — Grafana Loki Architecture
**Title:** Loki Architecture Deep Dive: Why It’s Fast & Cheap

- **Core Idea:** Index only **metadata labels** (`service_name`, `namespace`), never full-text log lines
- **Write Path:**
  - **Distributor** validates & hashes streams → routes to **Ingesters**
  - **Ingester** buffers chunks in memory → flushes **TSDB Index + Compressed Chunks** to Object Storage (S3/GCS/FS)
- **Read Path:**
  - **Query Frontend** splits & caches queries → **Queriers** scan TSDB index + chunks in parallel
- **Deployment Modes:** Monolithic (Single Binary) → Simple Scalable → Microservices

---

## Slide: Why Mimir Replaces Prometheus in LGTM
**Title:** Grafana Mimir vs. Prometheus: Scaling Metrics

- **100% Prometheus-Compatible:** Same PromQL, Alerting Rules & TSDB engine
- **Why move beyond standalone Prometheus?**
  - **Long-Term Storage:** Stores TSDB blocks in **Object Storage (S3/GCS)** instead of single-node disks
  - **Horizontal Scale & HA:** Microservices architecture (*Distributor, Ingester, Store-Gateway, Querier, Compactor*)
  - **Multi-Tenancy:** Native tenant isolation via `X-Scope-OrgID`
  - **Push-Based Ingestion:** Receives `remote_write` (with **Exemplars**) from **Grafana Alloy**, which handles OTLP & `/metrics` scraping

---

## Slide: Same Operational DNA — Loki vs. Tempo vs. Mimir
**Title:** One Shared Architecture Across Logs, Traces & Metrics

- **Common Core:** Stateless compute + **Object Storage (S3/GCS)** + Hash Ring (no ZooKeeper/Elasticsearch)
- **Shared Write/Read Pipeline:** `Distributor → Ingester → Compactor` & `Query Frontend → Querier`

| Backend | Signal | How It Indexes & Stores | Unique Component |
| :--- | :--- | :--- | :--- |
| **Loki** | **Logs** | Indexes **labels only** (TSDB); raw logs stored as compressed chunks | — |
| **Tempo** | **Traces** | **Zero index** for `trace_id` lookups; **Apache Parquet** for TraceQL | **`metrics-generator`** (derives RED metrics & Service Graphs) |
| **Mimir** | **Metrics** | **Prometheus TSDB** (inverted index + compressed samples) | **`store-gateway`** (caches S3 TSDB block indexes) |

> **Speaker Note:** *"The big advantage of LGTM is that Loki, Tempo, and Mimir share the exact same operational DNA: stateless compute decoupled from cheap Object Storage — only the block format inside changes (log chunks in Loki, Parquet traces in Tempo, and TSDB metric blocks in Mimir)."*

---

## Slide: Cross-Signal Correlation
**Title:** Navigating Signals: Metrics ↔ Traces ↔ Logs

- **Metrics → Traces:** Click **Exemplars** on Mimir charts to open Tempo traces
- **Traces → Logs:** *"Logs for this span"* in Tempo queries Loki by `trace_id`
- **Logs → Traces:** Loki `derivedFields` extract `trace_id` into Tempo links
- **Traces → Metrics:** Tempo `metrics-generator` builds **RED metrics** & **Service Graphs**

---

## Slide: Live Demo Walkthrough
**Title:** Live Demo: Troubleshooting Real Scenarios
Demo
---

## Slide: AI in Grafana Cloud
**Title:** AI in Grafana Cloud (Managed SaaS)

- **AI Assistant:** Natural language → PromQL/LogQL/TraceQL, dashboards & 30+ sources
- **AI Investigations:** Multi-agent root cause analysis across metrics, logs, traces & SQL
- **Agent Observability:** LLM token/cost tracking, agent/MCP replays & guardrails

---

## Slide: Open-Source AI with Grafana LGTM
**Title:** Open-Source AI Capabilities (Self-Hosted)

- **`mcp-grafana`:** Official MCP server connecting Claude/Cursor/Copilot to self-hosted LGTM

---

## Slide: Key Takeaways
**Title:** Key Takeaways

- **Open Standards (`OTel` + `Alloy`):** One collector, zero proprietary agent lock-in
- **LGTM Efficiency:** Loki, Tempo & Mimir share one operational model on Object Storage
- **True Correlation:** Metrics, logs, traces & DB queries linked in one click
- **AI-Ready:** Works seamlessly with both Grafana Cloud AI and OSS MCP/LLM tools
