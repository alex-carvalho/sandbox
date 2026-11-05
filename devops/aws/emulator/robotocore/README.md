# RobotoCore & The Local AWS Alternative

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
│          Docker Container / Rust Binary (:4566)          │
│                                                          │
│  ┌────────────────────────────────────────────────────┐  │
│  │             Axum + Tokio Gateway (Rust)            │  │
│  │  Request Router (Auth header, X-Amz-Target, Path)  │  │
│  │  Account ID Extractor (12-digit AWS Access Key)    │  │
│  └─────────────────────────┬──────────────────────────┘  │
│                            │                             │
│                            ▼                             │
│  ┌────────────────────────────────────────────────────┐  │
│  │          rotocore Native Rust Engine               │  │
│  │  • Per-service in-process Rust crates (160+ svcs)  │  │
│  │  • Native SigV4, IAM evaluation, SQS, Lambda       │  │
│  │  • Isolated in-memory state per Account & Region   │  │
│  └────────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────┘
```

* **Evolution** — Evolved from an early Python (Starlette + Moto) prototype into a standalone compiled **Rust** binary (Python is now only used for the `boto3`/`pytest` test suite).
* **Axum Gateway** — Parses `Authorization`, `X-Amz-Target`, path, or `Host` headers at native speed.
* **Multi-Account Isolation** — Treats any 12-digit `AWS_ACCESS_KEY_ID` as an isolated AWS Account ID.
* **`rotocore` Engine** — Clean-room Rust reimplementation of Moto's service layer running completely in-process.

---

<br>

## 4. State & Snapshots (`/_robotocore/state`)

* **Pure In-Memory by Default** — No continuous disk writes or embedded DB overhead.
* **On-Demand Checkpoints** — Free equivalent to LocalStack *Cloud Pods* for saving and resetting test state without re-running `terraform apply`:

```shell
# Save current in-memory AWS state
curl -X POST http://localhost:4566/_robotocore/state/save \
  -H "Content-Type: application/json" \
  -d '{"name": "after-terraform"}'

# Instantly restore / rollback to saved state
curl -X POST http://localhost:4566/_robotocore/state/load \
  -H "Content-Type: application/json" \
  -d '{"name": "after-terraform"}'
```

---

<br>

## 5. Live Demo (`Terraform` + `AWS CLI`)

* **Provision (`main.tf`)**: S3 Bucket + SQS Queue/DLQ + RDS Postgres + IAM Role/Policy
* **Manipulate (`AWS CLI`)**: Upload/download S3, send/receive SQS, inspect RDS, test `ENFORCE_IAM=1`

### 👉 **[Start Live Demo (`DEMO.md`)](./DEMO.md)**

---

<br>

## 6. Trade-offs

| ✅ Pros | ⚠️ Cons |
| :--- | :--- |
| **100% Free & MIT** — No paywalls or token lock-in | **Ephemeral RAM State** — Resets on container restart (supports snapshots) |
| **Rust Performance** — Compiled binary, no Python/JVM overhead | **No DB Sidecars** — Emulates RDS control plane & Data API, not a live Postgres/MySQL engine |
| **Broadest Coverage** — 160+ services / 10,100+ ops | **Policy Syntax on Create** — Stores raw policy JSON on `CreatePolicy`; evaluates at runtime |
| **Runtime IAM (`ENFORCE_IAM=1`)** — Included for free | **Younger Ecosystem** — Newer community than legacy tools |

---

<br>

## 7. Comparison: RobotoCore vs. Floci vs. MiniStack

| Dimension | **RobotoCore** | **Floci** | **MiniStack** |
| :--- | :--- | :--- | :--- |
| **Services** | **160+ services** (~10,100+ ops) | **133+ services** | **60+ services** |
| **Stack** | **Rust** (`Axum` + `rotocore` engine) | **Java** (Quarkus / Vert.x) | **Python 3.13** (ASGI / Hypercorn) |
| **DB & Sidecars** | In-memory control plane + RDS Data API | Launches real Docker sidecars (Postgres, MySQL, Redis, k3s) | Launches real Docker sidecars (Postgres, MySQL, Redis) |
| **Footprint** | Native Rust binary, fast boot & low RAM | **~24 ms startup**, ~13 MiB RAM | **< 2s startup**, ~30 MB RAM (`pip install ministack`) |
| **Persistence** | In-memory + Snapshots | Memory, Persistent, Hybrid, WAL | In-memory & Volume/S3 (`RDS_PERSIST=1`) |
| **License** | MIT | MIT | MIT |
