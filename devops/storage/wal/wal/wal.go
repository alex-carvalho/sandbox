package wal

import (
	"encoding/binary"
	"io"
	"os"
)

type OpType uint8

const (
	OpSet    OpType = 0
	OpDelete OpType = 1
)

func (op OpType) String() string {
	switch op {
	case OpSet:
		return "SET"
	case OpDelete:
		return "DELETE"
	default:
		return "UNKNOWN"
	}
}

// Record represents a single entry in the Write-Ahead Log.
type Record struct {
	// LSN (Log Sequence Number) is a unique, monotonically increasing ID
	LSN   uint64
	// Op is the operation type (e.g., SET or DELETE).
	Op    OpType
	// Key is the database key affected by this operation.
	Key   string
	// Value is the raw byte slice payload associated with the key.
	Value []byte
}

type WAL struct {
	file    *os.File
	nextLSN uint64
}

func Open(path string) (*WAL, []Record, error) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return nil, nil, err
	}

	w := &WAL{file: file}
	records, err := w.recover()
	if err != nil {
		file.Close()
		return nil, nil, err
	}

	return w, records, nil
}

func (w *WAL) Write(op OpType, key string, value []byte) (uint64, error) {
	lsn := w.nextLSN
	w.nextLSN++

	// 1. Write Header:
	// - LSN (8 bytes): Sequence order.
	// - Op (1 byte): Operation code.
	// - KeyLen (2 bytes): Prefix length for the variable-length Key.
	// - ValLen (4 bytes): Prefix length for the variable-length Value.
	// Specifying lengths in the header allows the parser to know exactly
	// how many bytes to read for each payload field.
	if err := binary.Write(w.file, binary.BigEndian, lsn); err != nil {
		return 0, err
	}
	if err := binary.Write(w.file, binary.BigEndian, uint8(op)); err != nil {
		return 0, err
	}
	if err := binary.Write(w.file, binary.BigEndian, uint16(len(key))); err != nil {
		return 0, err
	}
	if err := binary.Write(w.file, binary.BigEndian, uint32(len(value))); err != nil {
		return 0, err
	}

	// 2. Write Payload: Key + Value bytes directly matching lengths specified above.
	if _, err := w.file.Write([]byte(key)); err != nil {
		return 0, err
	}
	if _, err := w.file.Write(value); err != nil {
		return 0, err
	}

	// 3. Flush to disk (fsync) to guarantee persistence.
	if err := w.file.Sync(); err != nil {
		return 0, err
	}

	return lsn, nil
}

func (w *WAL) Close() error {
	return w.file.Close()
}

func (w *WAL) recover() ([]Record, error) {
	// Seek to the start of the file to scan all logs.
	if _, err := w.file.Seek(0, 0); err != nil {
		return nil, err
	}

	var records []Record
	for {
		var lsn uint64
		var opVal uint8
		var keyLen uint16
		var valLen uint32

		// Read LSN. If EOF is hit here, we reached the end of the log file.
		err := binary.Read(w.file, binary.BigEndian, &lsn)
		if err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				break
			}
			return nil, err
		}

		// Read the rest of the metadata header.
		if err := binary.Read(w.file, binary.BigEndian, &opVal); err != nil {
			return nil, err
		}
		if err := binary.Read(w.file, binary.BigEndian, &keyLen); err != nil {
			return nil, err
		}
		if err := binary.Read(w.file, binary.BigEndian, &valLen); err != nil {
			return nil, err
		}

		// Allocate and read the key and value buffers based on length fields.
		keyBuf := make([]byte, keyLen)
		if _, err := io.ReadFull(w.file, keyBuf); err != nil {
			return nil, err
		}

		valBuf := make([]byte, valLen)
		if _, err := io.ReadFull(w.file, valBuf); err != nil {
			return nil, err
		}

		records = append(records, Record{
			LSN:   lsn,
			Op:    OpType(opVal),
			Key:   string(keyBuf),
			Value: valBuf,
		})
		w.nextLSN = lsn + 1
	}

	if w.nextLSN == 0 {
		w.nextLSN = 1
	}

	return records, nil
}
