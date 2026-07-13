# Gossip Protocol POC

This project implements a Go-based **Gossip Protocol Proof of Concept (POC)**. It demonstrates decentralized, eventual-consistency state replication (Last-Write-Wins Key-Value database) and failure detection/membership tracking (based on heartbeat counters and suspected/dead status timers).

It comes with an interactive **Web Dashboard** featuring a custom HTML5 canvas network visualizer that draws live gossip message propagation, dynamic cluster controls, and real-time logs.

---

## Codebase Architecture

The project has zero third-party dependencies and uses only the Go standard library:

1. **[packet.go](./packet.go)**:
   Defines the network message structures (`GossipMessage`, `Member`, `Value`). Uses JSON serialization over UDP.
2. **[membership.go](./membership.go)**:
   Implements the thread-safe `MembershipTable`. Handles peer list merging (alive/suspected/dead updates), self-status refutation (if another node reports this node is suspected/dead, it increments its heartbeat and refutes), and failure timeouts.
3. **[node.go](./node.go)**:
   The core gossip worker. Spawns UDP sockets and runs three independent goroutines per node:
   - **UDP Listener**: Listens for incoming gossip and merges state & membership.
   - **Gossip Loop**: Periodically increments heartbeat, selects $k$ random peers (fanout), and writes UDP packets.
   - **Failure Detector Loop**: Periodically scans known peers and transitions stale peers to `SUSPECTED` or `DEAD`.
4. **[simulator.go](./simulator.go)**:
   Manages a cluster of in-memory nodes. Exposes controls to crash/kill nodes, revive them (with restored KV state), dynamically add nodes, and inject key-value updates. Pipes internal node events to a unified channel.
5. **[web.go](./web.go)**:
   Runs an HTTP server. Emits SSE (Server-Sent Events) to stream real-time node events to the client and exposes REST endpoints for dashboard commands.
6. **[index.html](./index.html)**:
   The UI dashboard. Features a responsive grid, glowing glassmorphism theme, custom HTML5 Canvas rendering of the circular node layout, particle-based message animations, log filters, and interactive simulation controls.
7. **[main.go](./main.go)**:
   CLI entry point. Parses configuration flags and launches either the visual simulation or a standalone interactive terminal-based node.

---

## Verification Results

### Automated Unit Tests
We wrote comprehensive unit tests verifying membership table updates, failure detection timeout steps, self-status refutation, and loopback state convergence:

```bash
go test -v ./...
```

**Output:**
```
=== RUN   TestNewMembershipTable
--- PASS: TestNewMembershipTable (0.00s)
=== RUN   TestMergeMembership
--- PASS: TestMergeMembership (0.00s)
=== RUN   TestRefutation
--- PASS: TestRefutation (0.00s)
=== RUN   TestDetectFailures
--- PASS: TestDetectFailures (0.00s)
=== RUN   TestStateConvergence
--- PASS: TestStateConvergence (0.15s)
PASS
ok  	simple-gossip	0.558s
```

---

## Running the Live Simulation

To run the interactive cluster simulation and view the web dashboard:

1. **Start the simulation**:
   ```bash
   go run . --http-port 8080 --cluster-size 5
   ```
2. **View the Web Dashboard**:
   Open your browser to: **[http://localhost:8080](http://localhost:8080)**

#### What to test:
- **Inject State**: Under **Inject State Update**, select a port (e.g. `127.0.0.1:9001`), set key `hello` to `world`, and click **Inject Update**.
  - *Observation*: Watch the indigo gossip particles fly across the canvas graph and see the logs light up in cyan. Check the other node cards: they will rapidly synchronize and display `hello = world`.
- **Kill a Node**: Click **Kill Node** on any card (e.g. `127.0.0.1:9003`).
  - *Observation*: The node becomes gray. The other nodes will keep sending gossips, notice that it is unresponsive, mark it as `SUSPECTED` (yellow outline/badges), and eventually transition it to `DEAD` (red outline/badges) when the heartbeat timeout expires.
- **Revive a Node**: Click **Revive Node** on the killed node's card.
  - *Observation*: The node boots back up, resumes communication, gets the latest keys it missed during its downtime, and refutes the old suspicion state, going back to green `ALIVE` status.
- **Add a Dynamic Node**: Under **Join Dynamic Node**, the next sequential port (e.g., `127.0.0.1:9006`) is automatically calculated and pre-populated. Simply click **Add Node**.
  - *Observation*: A new node circle appears in the topology graph. It links to a bootstrap peer, joins the cluster, and automatically receives the shared replicated state and membership list.

---

## Running Standalone Terminal Nodes with Dashboard Visualizer

You can run standalone nodes and interconnect them manually from different terminals. **You can also enable the interactive Web Dashboard on Node 1** to visualize the multi-terminal nodes and watch messages propagate between separate processes in real-time!

1. **Terminal 1**: Start Node 1 on UDP port `9001` and host the visualizer dashboard on HTTP port `8080`:
   ```bash
   go run . -mode node -port 9001 -ui -http-port 8080
   ```
   *Dashboard Link*: Open your browser to **[http://localhost:8080](http://localhost:8080)**. Initially, it will show only Node 1.

2. **Terminal 2**: Start Node 2 on UDP port `9002`, bootstrap it with Node 1, and report its telemetry events to Node 1's visualizer:
   ```bash
   go run . -mode node -port 9002 -peers 127.0.0.1:9001 -visualizer http://localhost:8080
   ```
   *Result*: Node 2 connects to Node 1. Look at Node 1's browser dashboard: Node 2 dynamically pops up on the graph. Because of the `-visualizer` flag, Node 2's packet sends and receives will now also animate on the dashboard!

3. **Terminal 3**: Start Node 3 on UDP port `9003`, bootstrap it with Node 2, and report its telemetry events to Node 1's visualizer:
   ```bash
   go run . -mode node -port 9003 -peers 127.0.0.1:9002 -visualizer http://localhost:8080
   ```
   *Result*: Node 3 connects to Node 2, discovers Node 1 indirectly via gossip, and begins communicating. Its activity is reported to Node 1's visualizer, animating the entire cluster flow!

Type commands directly in any terminal session:
- `set score 42` (on Node 1, and watch it replicate to Node 2 and Node 3!)
- `get score` (on Node 3 to check if it has the value)
- `members` (to inspect the local cluster membership view)
- `exit` (to stop the node)

## Simulator Video

[![Simulator Video](http://img.youtube.com/vi/vVfyOetC0EY/0.jpg)](http://www.youtube.com/watch?v=vVfyOetC0EY "Simulator Video")
