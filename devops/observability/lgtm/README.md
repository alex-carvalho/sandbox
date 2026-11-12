# Grafana LGTM Stack POC (Kind + Helm + Alloy + Polyglot Loan Microservices + PostgreSQL)

A complete local Proof of Concept for the **Grafana LGTM** observability stack (**L**oki, **G**rafana, **T**empo, **M**imir) collected via **Grafana Alloy** on a local **Kind** Kubernetes cluster, featuring a 3-tier polyglot architecture with **PostgreSQL**:

- **`loan-frontend` (Node.js / Express)**: Interactive Web Portal & API Gateway (`http://localhost:8080`)
- **`loan-backend` (Go)**: Loan Orchestrator & PostgreSQL Persistence (`SELECT`, `INSERT`, `UPDATE`, slow analytical queries, and SQL constraint errors)
- **`loan-underwriting` (Python / FastAPI)**: Fraud Blacklist verification (PostgreSQL), External Credit Bureau client, and Actuarial Risk Engine
- **`postgres` (PostgreSQL 16)**: Relational database with full SQL span and query-duration instrumentation

See [DESIGN.md](./DESIGN.md) for the full architecture and correlation design.

---

## Quickstart

### Deploy or Update the Stack

```bash
make all
```

If your Kind cluster and LGTM stack are already running and you only want to rebuild and redeploy the new apps and updated dashboard:

```bash
make build-load-apps deploy-apps deploy-lgtm status
```

---

## Live Presentation Scenarios (Web UI or API)

Open **http://localhost:8080** to trigger any scenario with one click, or pass specific user data:

| Presentation Scenario | How to Trigger in Form / API | What It Demonstrates in Grafana LGTM |
| :--- | :--- | :--- |
| **Normal Multi-Hop Flow** | Score `760`, Amount `$250,000` | Distributed trace across **Node.js $\rightarrow$ Go $\rightarrow$ Python $\rightarrow$ PostgreSQL** (`SELECT`, `INSERT`, `UPDATE`) |
| **Slow PostgreSQL Query (~900ms)** | Click **Demo: Slow PostgreSQL Query** or set Amount `> $900,000` | Long `postgres.slow_portfolio_exposure_analysis` span (`SELECT pg_sleep(0.85)...`) + Exemplar on the **PostgreSQL Query Latency** panel |
| **Fraud Blacklist Block (`403`)** | Click **Demo: Fraud Blacklist Error** or include `"Fraud"` / `"Blacklist"` in Applicant Name | `SELECT` against `fraud_blacklist` table in PostgreSQL followed by `403 FRAUD_BLOCKED` error span & `SECURITY ALERT` log in Loki |
| **Credit Bureau Timeout (`504`)** | Click **Demo: Credit Bureau Timeout** or set Credit Score `< 450` | Upstream timeout in `loan-underwriting` (`credit-bureau.external_lookup`) propagating `504` back through Go and Node.js |
| **PostgreSQL SQL Constraint Error (`500`)** | Click **Demo: PostgreSQL Constraint Error** or set Amount `= 666` | Real PostgreSQL `chk_valid_loan_amount` constraint violation in `loan-backend` with SQL statement & error recorded in Tempo and Loki |
