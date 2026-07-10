# Minimal Heartbeat Gossip POC

This code contains a highly simplified, single-file implementation of the Gossip Protocol. Focusing purely on **decentralized cluster membership and heartbeat failure detection**.

---

## Code Design

- **State Representation**: Each node tracks known peers as an **`(Epoch, Heartbeat)` tuple**. 
  - **`Epoch`**: Monotonically increasing startup timestamp in seconds. It ensures that when a crashed node restarts with a clean slate, its new messages have a higher epoch, bypassing epoch regression.
  - **`Heartbeat`**: A simple sequential counter starting at `1` and incrementing monotonically.
- **Heartbeat Dissemination**: Every second (`-gossip-ms` interval), the node increments its own heartbeat, picks a random active peer from its list, and sends its entire peer version table over UDP (as JSON).
- **Membership Merge**: When a node receives a gosip packet, it compares the incoming version with its local cache. It updates its records if `incoming.Epoch > local.Epoch` OR `(incoming.Epoch == local.Epoch && incoming.Heartbeat > local.Heartbeat)`. If it sees an update for a failed peer, it triggers a recovery.
- **Failure Detector**: Every second, the node checks if any peer has not advanced its version for more than the timeout (`-fail-ms`). If so, it prints a `[DEAD]` warning.

---

## Launching Nodes Individually

Open multiple terminal windows to run standalone nodes and watch them interact:

### 1. Terminal 1 (Node 1)
Start the first node on port `9001`:
```bash
go run main.go -port 9001
```

### 2. Terminal 2 (Node 2)
Start the second node on port `9002` and bootstrap it with Node 1:
```bash
go run main.go -port 9002 -peers 127.0.0.1:9001
```
*Result*: Node 2 will connect to Node 1. You will see both nodes outputting `[DISCOVER]` and periodic `[GOSSIP]` heartbeat advances.

### 3. Terminal 3 (Node 3)
Start the third node on port `9003` and bootstrap it with Node 2:
```bash
go run main.go -port 9003 -peers 127.0.0.1:9002
```
*Result*: Node 3 will connect to Node 2, discover Node 1 *indirectly via gossip propagation*, and automatically begin exchanging heartbeats directly with Node 1 as well!

---

## Observing Failure Detection & Recovery

1. **Kill Node 2**: In Terminal 2, press `Ctrl+C` to stop Node 2.
2. **Watch the logs**: Within 5 seconds, Terminals 1 and 3 will print:
   ```
   [DEAD] Peer 127.0.0.1:9002 has failed (no heartbeat update for 5s)
   ```
3. **Revive Node 2**: Restart Node 2 in Terminal 2:
   ```bash
   go run minimal/main.go -port 9002 -peers 127.0.0.1:9001
   ```
4. **Observe Recovery**: Node 1 and Node 3 will automatically detect the new heartbeat activity, print a `[RECOVER]` message, and resume normal tracking!
