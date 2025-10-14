# Full TiDB Stack with Docker Compose (v8.5)

This file describes how to run a larger TiDB test environment locally using `docker-compose.full.yml`. The stack includes:
- PD (3 instances)
- TiKV (3 instances)
- TiDB (1 instance)
- TiFlash (1 instance)
- TiCDC (1 instance)
- Prometheus + Grafana for monitoring (basic)


Prerequisites
- Docker Desktop with enough resources (at least 8–12 GB RAM allocated; more is better)
- docker and docker-compose available

Bring up the full stack

```bash
docker compose -f docker-compose.full.yml up -d
```

Check services

```bash
# PD health
curl http://127.0.0.1:2379/pd/api/v1/health

# TiDB status page
open http://127.0.0.1:10080

# Grafana
open http://127.0.0.1:3000
# default grafana credentials are admin/admin (change on first login)
```

Connect to TiDB

```bash
mysql -h 127.0.0.1 -P 4000 -u root
```


Caveats and resource guidance
- TiKV and TiFlash are memory- and disk-intensive. On macOS allocate at least 8–12GB RAM to Docker Desktop; 16GB is preferable.
- On low-resource machines, reduce TiKV/TiFlash count to 1 each to save RAM/CPU.
