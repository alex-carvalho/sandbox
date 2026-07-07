package wal

import (
	"bytes"
	"os"
	"testing"
)

func TestWAL_WriteAndRecover(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "wal_test_*.log")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	// 1. Open the WAL and write a record
	w1, _, err := Open(tmpFile.Name())
	if err != nil {
		t.Fatalf("failed to open WAL: %v", err)
	}

	w1.Write(OpSet, "theme", []byte("dark"))
	w1.Close()

	// 2. Reopen the WAL to trigger recovery and verify the content
	w2, recovered, _ := Open(tmpFile.Name())
	defer w2.Close()

	if len(recovered) != 1 {
		t.Fatalf("expected 1 recovered record, got %d", len(recovered))
	}

	record := recovered[0]
	if record.LSN != 1 {
		t.Errorf("expected recovered LSN 1, got %d", record.LSN)
	}
	if record.Op != OpSet {
		t.Errorf("expected recovered Op OpSet, got %v", record.Op)
	}
	if record.Key != "theme" {
		t.Errorf("expected recovered key 'theme', got '%s'", record.Key)
	}
	if !bytes.Equal(record.Value, []byte("dark")) {
		t.Errorf("expected recovered value 'dark', got '%s'", string(record.Value))
	}

	// Verify next LSN has been incremented correctly
	if w2.nextLSN != 2 {
		t.Errorf("expected next LSN to be 2, got %d", w2.nextLSN)
	}
}
