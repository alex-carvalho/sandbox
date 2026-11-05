# RobotoCore & The Local AWS Landscape

> **100% Free, MIT-Licensed Local AWS Digital Twin (Written in Rust)**
> **160 Services** · **10,100+ API Operations (~96% coverage)** · **Port `4566`**

---

<br>

## 1. What is RobotoCore?

* **Drop-in Endpoint (`localhost:4566`)** — Works with AWS CLI, SDKs, Terraform, CDK, and Packer.
* **From Moto to Rust** — Created by Jack Danger (Moto maintainer). Started as a Python/Starlette + Moto bridge and was **rewritten 100% in Rust (`rotocore`)**.
* **Zero Friction** — No signup, no auth tokens, no telemetry, $0 cloud cost.

---

<br>

## 2. Why It Was Created

### Catalyst: LocalStack Community Sunset (*March 23, 2026*)
* LocalStack archived its open-source Community Edition and required auth tokens / commercial tiers.

### Problems Solved
* **Zero Surprise Bills** — No forgotten dev resources in real AWS.
* **Instant Feedback** — Millisecond local provisioning vs. slow cloud deploys.
* **No Local AWS Keys** — Eliminates real IAM credentials on developer laptops.
* **Isolated State** — No naming collisions in shared dev accounts.
* **AI Agent Sandbox** — Safe offline environment for coding agents (Claude, Cursor, Antigravity).

---

<br>

## 3. Architecture (Rust: `Axum` + `rotocore`)

```text
┌──────────────────────────────────────────────────────────┐
│                 Rust Server (port 4566)                  │
│                                                          │
│  ┌────────────────────────────────────────────────────┐  │
│  │        Request Router & Axum Gateway (src/core/)   │  │
│  │  • Protocol layer (query / json / rest-xml / ec2)  │  │
│  │  • Account ID Extractor (12-digit AWS Access Key)  │  │
│  └──────────────┬──────────────────────┬──────────────┘  │
│                 │                      │                 │
│                 ▼                      ▼                 │
│  ┌─────────────────────────┐  ┌───────────────────────┐  │
│  │  Core Service Crates    │  │    rotocore Crate     │  │
│  │  (s3, sqs, sts, lambda, │  │  ("Moto in Rust" —    │  │
│  │   dynamodb, ...)        │  │   139 service mods)   │  │
│  └──────────────┬──────────┘  └───────────┬───────────┘  │
│                 │                         │              │
│                 ▼                         ▼              │
│      In-Memory Stores (Isolated per Account & Region)    │
└──────────────────────────────────────────────────────────┘
```

* **Pure Rust Binary** — The old Python server and vendored Moto backend were completely replaced by `src/core/` + `rotocore/`.
* **Axum Gateway** — Parses `Authorization`, `X-Amz-Target`, URL patterns, sign scopes, and S3 virtual hosts.
* **Multi-Account Isolation** — Treats any 12-digit `AWS_ACCESS_KEY_ID` as an isolated AWS Account ID.
* **Strict `501 NotImplemented` Contract** — Unimplemented operations return `501 NotImplemented` instead of silent fake successes.

---

<br>

## 4. Platform Features (`/_robotocore/*`)

Verified endpoints implemented in [`src/core/server.rs`](https://github.com/robotocore/robotocore/blob/main/src/core/server.rs):

* **Chaos Engineering (`/_robotocore/chaos/rules`)** — Inject latency (`latency_ms`), throttling (`429`), or HTTP faults (`500`/`503`) by service, operation regex, region, account, and probability.
* **State Snapshots & Export (`/_robotocore/state/*`)** — Save (`save`), restore (`load`), list (`snapshots`), export/import JSON (`export`/`import`), or wipe (`reset`) in-memory state.
* **API Audit Log & Error Tracking (`/_robotocore/audit`, `/usage/errors`)** — In-memory ring buffer of every API call, status code, and error.
* **Service Posture & Runtimes (`/_robotocore/services`, `/runtimes`)** — Query all 160 registered services and available Lambda runtimes.
* **Multi-Account Isolation** — Any 12-digit `AWS_ACCESS_KEY_ID` maps to its own isolated AWS Account state.

```shell
# Chaos: Inject 503 error on S3 ListBuckets, then clear rules
curl -X POST http://localhost:4566/_robotocore/chaos/rules \
  -d '{"service":"s3","operation":"ListBuckets","error_code":"ServiceUnavailableException","status_code":503}'
curl -X POST http://localhost:4566/_robotocore/chaos/rules/clear

# State: Save, list, and restore named snapshots
curl -X POST http://localhost:4566/_robotocore/state/save -d '{"name": "after-tf"}'
curl -s http://localhost:4566/_robotocore/state/snapshots
curl -X POST http://localhost:4566/_robotocore/state/load -d '{"name": "after-tf"}'

# Introspection: Audit log, usage errors, and service catalog
curl -s http://localhost:4566/_robotocore/audit
curl -s http://localhost:4566/_robotocore/services
```

---

<br>

## 5. Live Demo (`Terraform` + `AWS CLI`)

* **Provision (`main.tf`)**: S3 Bucket + SQS Queue/DLQ + RDS Postgres + IAM Role/Policy
* **Manipulate (`AWS CLI` & `/_robotocore/*`)**: S3, SQS, RDS, IAM/STS, Chaos Injection, and Snapshots

### 👉 **[Start Live Demo (`DEMO.md`)](./DEMO.md)**

---

<br>

## 6. Trade-offs

| ✅ Pros | ⚠️ Cons |
| :--- | :--- |
| **100% Free & MIT** — No paywalls or token lock-in | **Ephemeral RAM State** — Resets on container restart (unless snapshotted/exported) |
| **Rust Performance** — Compiled binary, zero Python/JVM overhead | **No DB Sidecars** — Emulates RDS control plane & Data API, not a live Postgres/MySQL engine |
| **Broadest Coverage** — 160 services registered in-process | **IAM Control Plane Only** — Stores IAM policies/roles & STS sessions, but does not enforce IAM permissions or validate policy syntax in the Rust gateway yet |
| **Built-in Chaos & Snapshots** — Free `/_robotocore/chaos/*` & `/state/*` APIs | **Younger Ecosystem** — Active Rust port with evolving parity |

---

<br>

## 7. Comparison: RobotoCore vs. Floci vs. MiniStack

| Dimension | **RobotoCore** | **Floci** | **MiniStack** |
| :--- | :--- | :--- | :--- |
| **Services** | **160 services** (~96% ops) | **133+ services** | **60+ services** |
| **Stack** | **Native Rust** (`Axum` + `rotocore` engine) | **Java + GraalVM Native** (Quarkus / Vert.x) | **Python 3.13** (ASGI / Hypercorn) |
| **DB & Sidecars** | In-memory control plane + RDS Data API | Launches real Docker sidecars (Postgres, MySQL, Redis, k3s) | Launches real Docker sidecars (Postgres, MySQL, Redis) |
| **Footprint** | Native Rust binary, fast boot & low RAM | **~24 ms startup**, ~13 MiB RAM | **< 2s startup**, ~30 MB RAM (`pip install ministack`) |
| **Persistence** | In-memory + Snapshots (`save`/`load`/`export`) | Memory, Persistent, Hybrid, WAL | In-memory & Volume/S3 (`RDS_PERSIST=1`) |
| **Chaos / Fault Injection** | **Yes** (`/_robotocore/chaos/rules` API) | **No** (Open feature request `#4623`) | **No** |
| **License** | MIT | MIT | MIT |
