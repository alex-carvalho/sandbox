package main

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/golang/snappy"
)

// ChunkMetadata represents the parsed JSON from the snappy-compressed metadata block
type ChunkMetadata struct {
	Fingerprint uint64            `json:"fingerprint"`
	UserID      string            `json:"userID"`
	From        int64             `json:"from"`
	Through     int64             `json:"through"`
	Metric      map[string]string `json:"metric"`
	Encoding    int               `json:"encoding"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <chunk-file-or-directory-path>")
		fmt.Println("Example: go run main.go ../loki-data/chunks/")
		os.Exit(1)
	}

	targetPath := os.Args[1]
	info, err := os.Stat(targetPath)
	if err != nil {
		fmt.Printf("❌ Failed to access path %q: %v\n", targetPath, err)
		os.Exit(1)
	}

	if info.IsDir() {
		fmt.Printf("📂 Walking directory recursively: %s\n\n", targetPath)
		err = filepath.Walk(targetPath, func(path string, fileInfo os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			// Skip directories
			if fileInfo.IsDir() {
				return nil
			}
			// Skip hidden files
			if filepath.Base(path)[0] == '.' {
				return nil
			}
			// Skip index and TSDB database files
			if filepath.Ext(path) == ".tsdb" || filepath.Ext(path) == ".gz" || filepath.Base(path) == "name" {
				return nil
			}
			if bytes.Contains([]byte(path), []byte("/index/")) {
				return nil
			}

			fmt.Printf("\n##################################################\n")
			fmt.Printf("📄 Processing File: %s\n", path)
			fmt.Printf("##################################################\n")
			if parseErr := parseChunkFile(path); parseErr != nil {
				fmt.Printf("⚠️  Error parsing %s: %v\n", path, parseErr)
			}
			return nil
		})
		if err != nil {
			fmt.Printf("❌ Error walking directory: %v\n", err)
			os.Exit(1)
		}
	} else {
		if err := parseChunkFile(targetPath); err != nil {
			fmt.Printf("❌ Error parsing file: %v\n", err)
			os.Exit(1)
		}
	}
}

func parseChunkFile(filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open chunk file: %w", err)
	}
	defer f.Close()

	// 1. Read L_meta (4 bytes, big-endian)
	var lMeta uint32
	if err := binary.Read(f, binary.BigEndian, &lMeta); err != nil {
		return fmt.Errorf("failed to read metadata length (L_meta): %w", err)
	}

	if lMeta < 4 {
		return fmt.Errorf("invalid L_meta length %d (too small)", lMeta)
	}

	// 2. Read Snappy metadata stream (L_meta - 4 bytes)
	snappyLen := lMeta - 4
	snappyBytes := make([]byte, snappyLen)
	if _, err := io.ReadFull(f, snappyBytes); err != nil {
		return fmt.Errorf("failed to read Snappy metadata block: %w", err)
	}

	// Decompress Snappy metadata block
	sr := snappy.NewReader(bytes.NewReader(snappyBytes))
	decodedMetaJSON, err := io.ReadAll(sr)
	if err != nil {
		return fmt.Errorf("failed to decompress Snappy metadata block: %w", err)
	}

	var meta ChunkMetadata
	if err := json.Unmarshal(decodedMetaJSON, &meta); err != nil {
		return fmt.Errorf("failed to parse metadata JSON: %w (raw JSON: %s)", err, string(decodedMetaJSON))
	}

	// Print Metadata Header
	fmt.Printf("==================================================\n")
	fmt.Printf("CHUNK METADATA\n")
	fmt.Printf("==================================================\n")
	fmt.Printf("Tenant/UserID:    %s\n", meta.UserID)
	fmt.Printf("Fingerprint:      %016x (%d)\n", meta.Fingerprint, meta.Fingerprint)
	fmt.Printf("Time Range:       %s to %s\n", 
		time.Unix(meta.From, 0).Format(time.RFC3339),
		time.Unix(meta.Through, 0).Format(time.RFC3339),
	)
	fmt.Printf("Encoding:         %d (0x%02x)\n", meta.Encoding, meta.Encoding)
	fmt.Println("Labels:")
	for k, v := range meta.Metric {
		fmt.Printf("  %s=%q\n", k, v)
	}
	fmt.Println()

	// 3. Read L_data (4 bytes, big-endian)
	var lData uint32
	if err := binary.Read(f, binary.BigEndian, &lData); err != nil {
		return fmt.Errorf("failed to read data block length (L_data): %w", err)
	}

	// 4. Read Data block (L_data bytes)
	dataBytes := make([]byte, lData)
	if _, err := io.ReadFull(f, dataBytes); err != nil {
		return fmt.Errorf("failed to read data block: %w", err)
	}

	// 5. Read Block Index Offset (last 8 bytes of data block)
	if len(dataBytes) < 8 {
		return fmt.Errorf("data block is too small (%d bytes) to contain footer index", len(dataBytes))
	}
	idxOffsetBytes := dataBytes[len(dataBytes)-8:]
	idxOffset := binary.BigEndian.Uint64(idxOffsetBytes)

	if idxOffset > uint64(len(dataBytes)-8) {
		return fmt.Errorf("invalid block index offset: %d", idxOffset)
	}

	indexBytes := dataBytes[idxOffset : len(dataBytes)-8]
	r := bytes.NewReader(indexBytes)
	numBlocks, err := binary.ReadUvarint(r)
	if err != nil {
		return fmt.Errorf("failed to parse block count from index: %w", err)
	}

	fmt.Printf("==================================================\n")
	fmt.Printf("DATA BLOCKS INDEX (Count: %d)\n", numBlocks)
	fmt.Printf("==================================================\n")

	// Read all blocks metadata from the index
	type BlockMeta struct {
		Entries uint64
		MinTime int64
		MaxTime int64
		Offset  uint64
		Length  uint64
	}
	var blocks []BlockMeta
	for i := uint64(0); i < numBlocks; i++ {
		entries, _ := binary.ReadUvarint(r)
		minTime, _ := binary.ReadUvarint(r)
		maxTime, _ := binary.ReadUvarint(r)
		offset, _ := binary.ReadUvarint(r)
		length, _ := binary.ReadUvarint(r)

		bm := BlockMeta{
			Entries: entries,
			MinTime: int64(minTime / 2),
			MaxTime: int64(maxTime / 2),
			Offset:  offset,
			Length:  length,
		}
		blocks = append(blocks, bm)

		fmt.Printf("Block #%d:\n", i+1)
		fmt.Printf("  Log Entries count: %d\n", bm.Entries)
		fmt.Printf("  Time Range:        %s to %s\n",
			time.Unix(0, bm.MinTime).Format(time.RFC3339Nano),
			time.Unix(0, bm.MaxTime).Format(time.RFC3339Nano),
		)
		fmt.Printf("  File Offset:       %d\n", bm.Offset)
		fmt.Printf("  Compressed Size:   %d bytes\n", bm.Length)
	}
	fmt.Println()

	if len(blocks) == 0 {
		fmt.Println("No log blocks to display.")
		return nil
	}

	// Decompress and Parse Block 1: Labels Pool
	// Starts at offset 7, and goes until the offset of the first log block (excluding block 1's checksum)
	block1Start := uint64(7)
	block2Offset := blocks[0].Offset
	if block2Offset < 11 {
		return fmt.Errorf("invalid block 2 offset in index: %d", block2Offset)
	}
	block1Len := block2Offset - 4 - block1Start
	block1Bytes := dataBytes[block1Start : block1Start+block1Len]

	gz1, err := gzip.NewReader(bytes.NewReader(block1Bytes))
	if err != nil {
		return fmt.Errorf("failed to decompress labels pool block: %w", err)
	}
	gz1.Multistream(false)
	decLabelsBytes, err := io.ReadAll(gz1)
	if err != nil {
		gz1.Close()
		return fmt.Errorf("failed to read labels pool stream: %w", err)
	}
	gz1.Close()

	// Parse labels pool into a slice of strings
	var labelsPool []string
	poolReader := bytes.NewReader(decLabelsBytes)
	for poolReader.Len() > 0 {
		strLen, err := binary.ReadUvarint(poolReader)
		if err != nil {
			break
		}
		strBuf := make([]byte, strLen)
		if _, err := io.ReadFull(poolReader, strBuf); err != nil {
			break
		}
		labelsPool = append(labelsPool, string(strBuf))
	}

	fmt.Printf("==================================================\n")
	fmt.Printf("LOG ENTRIES\n")
	fmt.Printf("==================================================\n")

	// Decompress and display log entries from all log data blocks
	for bIdx, b := range blocks {
		// Gzip stream starts at block offset and runs to the end of data block
		blockBytes := dataBytes[b.Offset:]
		gz, err := gzip.NewReader(bytes.NewReader(blockBytes))
		if err != nil {
			fmt.Printf("❌ Failed to decode gzip stream for Block #%d: %v\n", bIdx+1, err)
			continue
		}
		gz.Multistream(false)
		decEntriesBytes, err := io.ReadAll(gz)
		gz.Close()
		if err != nil {
			fmt.Printf("❌ Failed to read entries stream for Block #%d: %v\n", bIdx+1, err)
			continue
		}

		s2Reader := bytes.NewReader(decEntriesBytes)
		entryIdx := 1
		for s2Reader.Len() > 0 {
			tsZigZag, err := binary.ReadUvarint(s2Reader)
			if err != nil {
				break
			}
			tsNS := int64(tsZigZag / 2)
			entryTime := time.Unix(0, tsNS)

			lineLen, _ := binary.ReadUvarint(s2Reader)
			lineBytes := make([]byte, lineLen)
			io.ReadFull(s2Reader, lineBytes)

			metaCount, _ := binary.ReadUvarint(s2Reader)
			var structuredMeta []string
			if metaCount > 0 {
				uvarints := make([]uint64, metaCount)
				for m := uint64(0); m < metaCount; m++ {
					metaIdx, _ := binary.ReadUvarint(s2Reader)
					uvarints[m] = metaIdx
				}
				numPairs := uvarints[0]
				for i := uint64(0); i < numPairs; i++ {
					if 1+2*i+1 < uint64(len(uvarints)) {
						keyIdx := uvarints[1+2*i]
						valIdx := uvarints[2+2*i]
						var key, val string
						if keyIdx < uint64(len(labelsPool)) {
							key = labelsPool[keyIdx]
						} else {
							key = fmt.Sprintf("<invalid_idx_%d>", keyIdx)
						}
						if valIdx < uint64(len(labelsPool)) {
							val = labelsPool[valIdx]
						} else {
							val = fmt.Sprintf("<invalid_idx_%d>", valIdx)
						}
						structuredMeta = append(structuredMeta, fmt.Sprintf("%s=%s", key, val))
					}
				}
			}

			// Format and print log line
			metaStr := ""
			if len(structuredMeta) > 0 {
				metaStr = " {"
				for i, sm := range structuredMeta {
					if i > 0 {
						metaStr += ", "
					}
					metaStr += sm
				}
				metaStr += "}"
			}
			fmt.Printf("[%s] %s%s\n", entryTime.Format("2006-01-02 15:04:05.000000"), string(lineBytes), metaStr)
			entryIdx++
		}
	}
	fmt.Printf("==================================================\n")
	fmt.Printf("✅ Finished reading chunk. Total blocks parsed: %d\n", numBlocks)
	return nil
}
