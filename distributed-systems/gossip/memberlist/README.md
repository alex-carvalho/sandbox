# HashiCorp memberlist Gossip Protocol PoC

This is a simple Proof of Concept (PoC) in Go that demonstrates cluster membership, node discovery, and failure detection using [HashiCorp's memberlist library](https://github.com/hashicorp/memberlist).

## How it Works

The nodes in the cluster use the **Gossip Protocol** to share membership information:
1. **Bootstrap / Join**: A new node starts up. If it is given seed addresses in the `-peers` flag, it contacts those nodes to join the cluster.
2. **State Gossip**: Nodes periodically select random peers and exchange membership lists.
3. **Failure Detection**: Nodes regularly ping each other. If a node does not respond within a timeout, it is marked as "suspicious". If it fails to refute the suspicion, it is marked as "dead" and this information is gossiped to the entire cluster.
4. **Graceful Leave**: When a node shuts down cleanly, it notifies the cluster so it is immediately removed from the active membership list without waiting for failure detection timeouts.

---

## How to Run

To run the PoC, you will need to open multiple terminal windows or run the processes in the background to see them interact.

### Step 1: Start the first node (seed node)
Run the first node on port `9001` (default):
```bash
go run main.go -name node1 -port 9001
```

### Step 2: Start a second node and join the first
In another terminal, start a second node on port `9002` and instruct it to join `node1` at `127.0.0.1:9001`:
```bash
go run main.go -name node2 -port 9002 -peers 127.0.0.1:9001
```

### Step 3: Start a third node and join the cluster
In a third terminal, start a third node on port `9003` and join `node2` at `127.0.0.1:9002`:
```bash
go run main.go -name node3 -port 9003 -peers 127.0.0.1:9002
```

---

## What to Observe

1. **Auto-Discovery**:
   As soon as nodes join, you will see color-coded notifications in the consoles. Even though `node3` only contacted `node2` (`127.0.0.1:9002`), it will automatically discover `node1` and join it via gossip!
   ```text
   [DISCOVER] [node2] joined cluster (addr: 127.0.0.1:9002)
   ```

   Every 5 seconds, each node will print its current active cluster members in a single line:
   ```text
   [MEMBERS] node1 (127.0.0.1:9001), node2 (127.0.0.1:9002), node3 (127.0.0.1:9003)
   ```

2. **Graceful Leave**:
   Stop `node3` by pressing `Ctrl+C`.
   - `node3` will log `Shutting down node gracefully...` and notify the cluster.
   - Instantly, other active nodes will output:
     ```text
     [DEAD] [node3] has left or failed
     ```
   - Their membership printout will immediately reflect 2 nodes.

3. **Node Failure Detection (Ungraceful Stop)**:
   Kill a node abruptly (e.g., using `kill -9` or stopping it in a way that bypasses graceful shutdown).
   - Because it did not send a leave message, other nodes will start suspecting the killed node is dead.
   - After a brief suspicion period, active nodes will print:
     ```text
     [DEAD] [<node>] has left or failed
     ```

---
