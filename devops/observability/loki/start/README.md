# Starting on Grafana Loki

---

## How Loki Works (Concepts)

Grafana Loki is a horizontally scalable, highly available, multi-tenant log aggregation system inspired by Prometheus. Unlike traditional log search databases (like Elasticsearch), Loki takes a unique approach:

* **Metadata-Only Indexing**: Loki **does not** index the raw text of the logs. Instead, it only indexes the **labels** (metadata) associated with a log stream (e.g., `job="bash-test"`, `env="dev"`).
* **Streams**: A stream is a set of log entries sharing the exact same labels.
* **Chunks**: Raw log messages are compressed and stored as chunks in object storage (or local disk in this dev setup), grouped by stream.
* **LogQL**: Loki's query language, which looks and behaves like PromQL. You filter streams by labels, then apply line filters (e.g., `|= "error"`) to search the raw log text dynamically at query time.

---

## Starting the Environment
```bash
docker compose up -d


curl http://localhost:3101/ready
```

---

## Pushing Logs via HTTP API

Loki exposes a push endpoint at `POST /loki/api/v1/push`. 
Log entries must be formatted as JSON containing streams and list of values (where each value is a `[timestamp_ns_string, log_line_string]` array).

We have provided a helper script `./push-log.sh` to construct the payload and send it:

```bash
# Run with default log message
./push-log.sh

# Run with custom job, level, and message
./push-log.sh "myapp" "ERROR" "Database connection timed out!"
```
---

## Querying Logs via HTTP API

Loki exposes a query range endpoint at `GET /loki/api/v1/query_range` supporting **LogQL**.

We have provided a helper script `./query-log.sh` to query Loki:

```bash

# Run default query (fetches logs for job="bash-test")
./query-log.sh

# Run custom LogQL query
./query-log.sh '{job="myapp"} |~ "(?i)error"'
```


Loki use wal to make sure not loose any log data. How to see wal content:
```bash
cd wal-reader && go mod tidy && go run main.go ../loki-data/wal/
```

Read chunks via `chunk-reader`:
```bash
cd chunk-reader && go mod tidy && go run main.go ../loki-data/chunks/
```

Read index via `index-reader`:
```bash
cd index-reader && go mod tidy && go run main.go ../loki-data/chunks/index/
```

TODO: structured metadata
