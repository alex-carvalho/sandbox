# TiKV 7.1.2 POC Setup

This repository contains a Docker Compose configuration for running a TiKV cluster for proof-of-concept and development purposes.

## Overview

TiKV is a distributed transactional key-value database. This setup includes:
- 1 PD (Placement Driver) node - manages and schedules the TiKV cluster
- 3 TiKV nodes - distributed storage nodes


## Quick Start

### Start the Cluster

```bash
docker-compose up -d
```

### Check Cluster Status

```bash
# Check if all containers are running
docker-compose ps

# Check PD cluster status
curl http://localhost:2379/pd/api/v1/health

# Check cluster members
curl http://localhost:2379/pd/api/v1/members

# Check cluster stores (TiKV nodes)
curl http://localhost:2379/pd/api/v1/stores
```

## API Endpoints

### PD API Endpoints

- Health check: `http://localhost:2379/pd/api/v1/health`
- Cluster status: `http://localhost:2379/pd/api/v1/status`
- Members: `http://localhost:2379/pd/api/v1/members`
- Stores: `http://localhost:2379/pd/api/v1/stores`
- Regions: `http://localhost:2379/pd/api/v1/regions`
- Config: `http://localhost:2379/pd/api/v1/config`


## Architecture

```
┌─────────────┐
│   PD Node   │  (Port 2379, 2380)
│  (pd0)      │  - Cluster management
└──────┬──────┘  - Timestamp allocation
       │         - Region scheduling
       │
   ┌───┴────────────────┐
   │                    │
┌──▼───┐  ┌────────┐  ┌▼────────┐
│TiKV0 │  │ TiKV1  │  │ TiKV2   │
│20160 │  │ 20161  │  │ 20162   │
└──────┘  └────────┘  └─────────┘
  Data storage nodes
```
