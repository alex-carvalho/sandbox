# Running TiDB locally with Docker (macOS / zsh)

This guide shows a minimal TiDB cluster (PD, TiKV, TiDB) using Docker Compose. It's intended for development and experimentation only.

Start the cluster

```bash
docker compose up -d
```

Notes:
- The provided `docker-compose.yml` uses PingCAP images
- Ports exposed:
  - 4000: MySQL-compatible SQL port
  - 2379: PD client
  - 20160: TiKV gRPC
  - 10080: TiDB status/http

Connect with the MySQL client

```bash
# using mysql client
mysql -h 127.0.0.1 -P 4000 -u root

# or using Docker exec into tidb container
docker compose exec tidb mysql -uroot -h 127.0.0.1 -P 4000
```

Check status (example)

```bash
# query PD
curl http://127.0.0.1:2379/pd/api/v1/health

# TiDB status page
open http://127.0.0.1:10080
```
