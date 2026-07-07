# Write-Ahead Log (WAL) Proof of Concept

This is a simple PoC of WAL using Go. It demonstrates the fundamental concepts of how databases achieve durability and atomicity through append-only logging, sequential binary serialization, and crash recovery.

---

## Core Concepts

### 1. What is a WAL?
A Write-Ahead Log is an append-only log file on disk. Before any database operation (like insert, update, or delete) is committed and applied to the database's main storage files (or in-memory state), it is first written and flushed to the WAL. If the system crashes, the in-memory database can be fully reconstructed by replaying the log file from the beginning.

### 2. Key Terms Explained
* **LSN (Log Sequence Number)**: A 64-bit (`uint64`) unique, monotonically increasing identifier assigned to each log record. LSNs define a strict, absolute chronological order of mutations. They also prevent replaying duplicates (idempotency) and indicate replication progress.
* **OpType (Operation Type)**: Represents the operation being logged. In this PoC, we support two types: `SET` (value write/update) and `DELETE` (value removal).
* **KeyLen & ValLen**: Because keys and values have variable lengths (e.g., the key `"role"` is 4 bytes, but `"username"` is 8 bytes), we cannot predetermine their size on disk. We prefix them with length indicators:
  - **KeyLen** (2 bytes, `uint16`): The length of the key.
  - **ValLen** (4 bytes, `uint32`): The length of the value payload.
  
  During deserialization, the parser reads these lengths first, allocates the exact buffers, and reads that specific number of bytes from the file.

---

## Binary Record Format on Disk

Records are written to `wal.log` using a packed binary layout:

```
+------------------+------------------+------------------+-----------------+
| LSN (8 bytes)    | OpType (1 byte)  | KeyLen (2 bytes) | ValLen (4 bytes)|
+------------------+------------------+------------------+-----------------+
| Key (KeyLen B)   | Value (ValLen B) |
+------------------+------------------+
```

### Visualizing the Binary File
When running the demo, the WAL content is displayed as a Hex Dump:
```
00000000  00 00 00 00 00 00 00 01  00 00 08 00 00 00 04 75  |...............u|
00000010  73 65 72 6e 61 6d 65 61  6c 65 78                 |sernamealex|
```
- `00 00 00 00 00 00 00 01`: LSN = 1 (8 bytes)
- `00`: Op = SET (1 byte)
- `00 08`: KeyLen = 8 (2 bytes)
- `00 00 00 04`: ValLen = 4 (4 bytes)
- `75 73 65 72 6e 61 6d 65`: Key = `"username"` (8 bytes)
- `61 6c 65 78`: Value = `"alex"` (4 bytes)

---

## How to Run

1. Clone or navigate to the workspace directory.
2. Run the main demo:
   ```bash
   go run main.go
   ```
3. Run tests
   ```bash
   go test ./...
   ```

### Execution Flow & Output
The program goes through the following sequence:
1. **Initialize DB & WAL**: Starts a fresh `wal.log` file and an empty in-memory Key-Value store map.
2. **Write Operations**: Performs transactions (`SET username`, `SET role`, `DELETE role`). The operations are logged to the WAL and synced (`fsync`) before being applied to the in-memory map.
3. **Inspect Disk Log**: Outputs the raw hex dump of the WAL file.
4. **Simulate Crash**: Disposes of the in-memory map to simulate memory wipeout.
5. **Reopening & Replaying**: Opens the WAL file again, scans and parses all binary records, and replays them sequentially to rebuild the original database state.
6. **Verify State**: Confirms the recovered state matches the pre-crash state.
