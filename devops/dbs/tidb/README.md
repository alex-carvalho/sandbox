## TiDB —

### Why TiDB was created / Problem it solves
- Background: Modern online services need both transactional (OLTP) and real-time analytical (OLAP) capabilities. Traditional relational databases provide strong transactional guarantees but don’t scale horizontally easily. NoSQL or NewSQL systems often trade away SQL compatibility, strong transactions, or easy adoption.
- Goal: TiDB was designed to provide a MySQL-compatible, distributed HTAP (Hybrid Transactional/Analytical Processing) database that supports horizontal scalability, strong consistency, and familiar SQL while minimizing operational complexity.
- Problems solved:
	- Seamless horizontal scaling for transactions without sharding application code.
	- Strong consistency (serializable or snapshot isolation depending on configuration) across distributed nodes.
	- Single logical SQL interface (MySQL-compatible) for both OLTP and near real-time OLAP workloads.
	- Operational flexibility for cloud-native and commodity hardware deployments.
	- Reduces need for separate systems (MySQL + data warehouse) by enabling mixed workloads.

### High-level architecture
- Layered, modular components:
	- **TiDB Server (stateless SQL layer)**
		- Acts like a MySQL-compatible SQL gateway.
		- Parses SQL, plans queries, coordinates distributed execution, and returns results.
		- Stateless design enables easy horizontal scaling and rolling upgrades.
	- **Placement Driver (PD)**
		- Cluster metadata manager and distributed scheduling control plane.
		- Maintains cluster topology, global timestamp allocation (for distributed transactions using TSO — Timestamp Oracle), and placement of data (Region leaders and peers).
	- **TiKV (distributed transactional storage)**
		- Key-value storage built on RocksDB.
		- Provides distributed transactions with MVCC (Multi-Version Concurrency Control).
		- Stores data as Regions (ranges of keys); Regions are replicated using the Raft consensus protocol for fault tolerance.
	- **TiFlash (columnar engine for analytics, optional)**
		- Asynchronously replicates data from TiKV to a columnar store optimized for analytical queries.
		- Allows TiDB to serve both OLTP (via TiKV) and OLAP (via TiFlash) using the same SQL layer.
	- **Binlog / TiCDC (change data capture)**
		- Tools and components for replication, backup, and CDC streaming to other systems.
- Important control flows:
	- Transaction coordination: TiDB requests timestamps from PD, coordinates read/write via TiKV, and commits using 2-phase commit with Raft-backed Region leaders.
	- Query execution: TiDB pushes compute down where possible (TiKV coprocessor) and coordinates distributed execution across Regions; TiFlash is used for columnar scans and aggregation-heavy queries.

### Principal features
- MySQL protocol and dialect compatibility
	- Applications speaking MySQL can generally connect without changes.
- Distributed transactions with strong consistency
	- MVCC and support for ACID semantics across nodes.
	- Timestamp Oracle (TSO) provides globally-ordered timestamps for transactions.
- Horizontal scaling
	- Stateless SQL nodes and independently scalable storage nodes.
	- Automatic Region splitting and rebalancing.
- High availability and fault tolerance
	- Raft replication per Region; PD orchestrates failover and rebalancing.
- HTAP capabilities
	- TiFlash provides columnar, MPP-like reads for analytics while TiKV serves transactional workloads.
- Online schema change / DDL
	- Non-blocking schema changes and compatibility with MySQL’s DDL patterns.
- Rich ecosystem
	- TiCDC for change data capture, BR for backups, Lightning for fast import, Dashboard for monitoring and management.
- Performance features
	- Coprocessor push-down to TiKV (filtering and partial aggregation at storage).
	- Vectorized execution and cost-based optimizer improvements.
- Security & enterprise features
	- Role-based access control, encryption-at-rest/in-transit (Enterprise features may vary by distribution).
- Cloud-native and operator support
	- Kubernetes operators (TiDB Operator) for running and managing clusters in K8s.

### How TiDB scales
- Horizontal scalability model:
	- Stateless TiDB servers scale out linearly for connection concurrency and query planning capacity.
	- TiKV nodes scale storage capacity and IO; Regions (key ranges) are split automatically when they become large or hot, and rebalanced across nodes.
	- Raft replication ensures availability as nodes fail; PD reassigns leaders and rebalances Regions.
	- TiFlash can be scaled independently to meet analytical throughput needs.
- Data distribution and balancing:
	- Data is partitioned into Regions (typical default ~96MB) and distributed across TiKV nodes.
	- PD monitors Region size, IO, and hot spots; it issues scheduling tasks (split Region, move peer, change leader) to maintain balance and performance.
- Concurrency and throughput:
	- MVCC enables many concurrent reads without blocking writers; write throughput is bounded by Raft/IO and network.
	- Coprocessor and TiFlash reduce network and CPU overhead by pushing computation close to data.
- Elastic operations:
	- Adding nodes triggers PD-driven rebalancing; removing nodes triggers region peer relocation.
	- Rolling upgrades supported because TiDB servers are stateless.

### References & next steps
- Official docs and case studies are on the TiDB (PingCAP) website — check the latest architecture diagrams, version-specific features (TiKV, TiFlash), and real-user case studies.
- If you want, I can:
	- Produce a one-page printable resume version (PDF/Markdown).
	- Expand any section into a detailed technical explainer (e.g., transaction protocol with sequence diagram, PD internals, TiKV storage layout).
	- Create comparison matrix vs MySQL, CockroachDB, YugabyteDB, and Snowflake for HTAP decisions.



### TiDB Comparison Matrix

A short, high-level comparison of TiDB against MySQL, CockroachDB, YugabyteDB, and Snowflake for HTAP/OLTP decisions.

| Product | Compatibility | Transactions | Scaling model | HTAP | Primary use-case |
|---|---|---|---|---|---|
| TiDB | MySQL protocol and dialect compatibility | Distributed ACID (MVCC + 2PC + Raft) | Horizontal; TiDB stateless, TiKV scale-out storage | Yes (TiFlash) | MySQL-compatible distributed OLTP + HTAP |
| MySQL (single-node / Group Replication) | Native | ACID on single node; Group Replication adds HA but limited scaling | Vertical or manual sharding (ProxySharding, Vitess) | No (requires separate analytical system) | Traditional OLTP workloads on single node or managed clusters |
| CockroachDB | PostgreSQL wire protocol (dialect differs) | Distributed ACID with Raft-based replication | Horizontal, geo-distribution focused | Limited; primarily OLTP | Distributed SQL for global apps with strong consistency |
| YugabyteDB | PostgreSQL-compatible and Cassandra-like interfaces | Distributed ACID (MVCC + consensus) | Horizontal, multi-shard | Limited; OLTP focus with analytical integrations | Cloud-native distributed SQL as a PostgreSQL-compatible alternative |
| Snowflake | ANSI SQL (analytical DW) | Not a general-purpose transactional DB (supports ACID within queries) | Elastic compute separated from storage | No (analytical data warehouse) | Large-scale analytics and data warehousing |


## TiDB Transaction Protocol — Deep Dive

This document explains TiDB's transaction protocol and the components involved: PD (TSO), MVCC, 2PC, Raft, and how TiKV and TiDB coordinate to provide distributed ACID transactions.

## Contract
- Inputs: SQL transactional workload (reads/writes) from clients.
- Outputs: Durable, consistent commits or rollbacks across Regions.
- Success: Transaction commits with isolation guarantees; failures return clear errors.

## Components
- PD (Placement Driver): issues global timestamps (TSO) and stores metadata about Regions and topology.
- TiDB server: coordinates transactions and executes SQL plans.
- TiKV: stores key-value data with MVCC and handles Raft replication per Region.
- Raft: consensus protocol used by TiKV to replicate Region peers.

## Flow: a simplified write transaction
1. Begin transaction: client issues BEGIN; TiDB contacts PD to get a start timestamp (start_ts) via TSO.
2. Execution: reads and writes are performed against TiKV using start_ts for snapshot reads; writes are buffered at TiDB.
3. Prewrite (Phase 1 of 2PC): TiDB sends prewrite requests to the Region leaders for each affected key. TiKV writes a Lock and a provisional value with start_ts.
4. Commit timestamp (commit_ts): TiDB asks PD for a commit timestamp (via TSO) after prewrites succeed.
5. Commit (Phase 2 of 2PC): TiDB sends commit requests to the same Region leaders. TiKV converts provisional values to committed versions with commit_ts and releases locks.
6. Cleanup: any locks from failed prewrites are garbage-collected or resolved by TiKV's resolve-lock mechanisms.

## Read behavior and MVCC
- Snapshot reads use start_ts to return the latest committed version whose commit_ts <= start_ts.
- For strong consistency and repeatable reads, TiDB uses pessimistic or optimistic transaction modes.

## Optimizations and failure handling
- One-phase commit (1PC): If a transaction touches keys in a single Region, TiDB may use a fast path to commit in one round.
- Async Commit: optimization to reduce latency by allowing commit_ts to be visible earlier under certain safety checks.
- Pessimistic transactions: acquire locks earlier to avoid conflicts (useful for high-contention workloads).
- Resolve locks: background processes (and clients) can resolve stale locks caused by failed coordinators.

## Raft and durability
- Each Region is a Raft group. Writes flow through the Region leader and are replicated by Raft to peers before being acknowledged.
- TiKV persists to RocksDB and Raft logs to ensure durability.

## Edge cases
- PD unavailability: cluster cannot obtain new TSOs; Io remains possible for existing timestamped reads but commits requiring new TS will fail.
- Split/Merge Region during transaction: transaction uses Region metadata and leader routing; PD updates guide retries.
- Network partition: Raft ensures safety by requiring quorum; partitions may reduce availability until quorum is restored.

## Further reading
- TiDB docs: Transaction Model and MVCC
- TiKV docs: Raftstore and storage internals


Diagram (text):
- Client -> TiDB: BEGIN
- TiDB -> PD: get start_ts
- TiDB -> TiKV(RegionA leader): prewrite keyA
- TiDB -> TiKV(RegionB leader): prewrite keyB
- TiDB -> PD: get commit_ts
- TiDB -> TiKV(RegionA leader): commit keyA
- TiDB -> TiKV(RegionB leader): commit keyB
