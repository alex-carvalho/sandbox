package main

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/example/go-wal-poc/wal"
)

func main() {
	logPath := "wal.log"

	// Clean up previous runs
	os.Remove(logPath)
	defer os.Remove(logPath)

	fmt.Println("=== STEP 1: Initialize DB & WAL ===")
	db := make(map[string]string)

	w, _, err := wal.Open(logPath)
	if err != nil {
		fmt.Println("Error opening WAL:", err)
		return
	}

	// Perform operations (Logged to WAL first, then applied to memory)
	fmt.Println("\n=== STEP 2: Write Operations (Logged to WAL first) ===")
	write := func(op wal.OpType, key string, val string) {
		fmt.Printf("Writing: %s %s = %s\n", op, key, val)

		lsn, err := w.Write(op, key, []byte(val))
		if err != nil {
			fmt.Println("WAL write error:", err)
			return
		}

		switch op {
			case wal.OpSet:
				db[key] = val
			case wal.OpDelete:
				delete(db, key)
		}
		fmt.Printf("  LSN: %d | Memory DB State: %v\n", lsn, db)
	}

	write(wal.OpSet, "username1", "john")
	write(wal.OpSet, "username2", "alex")
	write(wal.OpDelete, "username1", "")
	w.Close()

	// Show raw WAL content on disk (using Hex Dump to visualize the binary layout)
	fmt.Println("\n=== STEP 3: Raw WAL File Content on Disk (Hex Dump) ===")
	rawBytes, _ := os.ReadFile(logPath)
	fmt.Print(hex.Dump(rawBytes))

	// Simulate system crash
	fmt.Println("\n=== STEP 4: Simulate Sudden Crash (Memory wiped) ===")
	db = nil
	fmt.Println("Memory DB is now:", db)

	// Reopen WAL and recover database state
	fmt.Println("\n=== STEP 5: Reopening WAL & Replaying Logs ===")
	db = make(map[string]string)
	w2, recoveredRecords, err := wal.Open(logPath)
	if err != nil {
		fmt.Println("Error opening WAL for recovery:", err)
		return
	}
	defer w2.Close()

	for _, record := range recoveredRecords {
		switch record.Op {
			case wal.OpSet:
				db[record.Key] = string(record.Value)
			case wal.OpDelete:
				delete(db, record.Key)
		}

		fmt.Printf("  Replayed LSN %d: %s %s = %s | Current DB: %v\n",
			record.LSN, record.Op, record.Key, string(record.Value), db)
	}

	fmt.Println("\n=== Final Recovered DB State ===")
	fmt.Println(db)
}
