package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/prometheus/prometheus/tsdb/index"
)

// realByteSlice implements index.ByteSlice over standard []byte
type realByteSlice []byte

func (b realByteSlice) Len() int {
	return len(b)
}

func (b realByteSlice) Range(start, end int) []byte {
	return b[start:end]
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <tsdb-file-or-directory-path>")
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
			// Check if file is a TSDB index (ends with .tsdb, .gz, etc.)
			ext := filepath.Ext(path)
			if ext != ".tsdb" && ext != ".gz" {
				return nil
			}

			fmt.Printf("\n##################################################\n")
			fmt.Printf("📄 Processing TSDB Index File: %s\n", path)
			fmt.Printf("##################################################\n")
			if parseErr := parseIndexFile(path); parseErr != nil {
				fmt.Printf("⚠️  Error parsing %s: %v\n", path, parseErr)
			}
			return nil
		})
		if err != nil {
			fmt.Printf("❌ Error walking directory: %v\n", err)
			os.Exit(1)
		}
	} else {
		if err := parseIndexFile(targetPath); err != nil {
			fmt.Printf("❌ Error parsing file: %v\n", err)
			os.Exit(1)
		}
	}
}

func parseIndexFile(filePath string) error {
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read TSDB file: %w", err)
	}

	// Decompress if gzip compressed
	if len(fileBytes) > 2 && fileBytes[0] == 0x1f && fileBytes[1] == 0x8b {
		fmt.Println("Decompressing gzip-compressed TSDB file...")
		gz, err := gzip.NewReader(bytes.NewReader(fileBytes))
		if err != nil {
			return fmt.Errorf("failed to create gzip reader: %w", err)
		}
		defer gz.Close()
		decompressedBytes, err := io.ReadAll(gz)
		if err != nil {
			return fmt.Errorf("failed to decompress TSDB file: %w", err)
		}
		fileBytes = decompressedBytes
	}

	// Loki TSDB index footer layout (last 68 bytes of uncompressed file)
	footerSize := 68
	if len(fileBytes) < footerSize {
		return fmt.Errorf("uncompressed file is too small to contain a Loki TSDB index footer")
	}

	footerOffset := len(fileBytes) - footerSize
	footerBytes := fileBytes[footerOffset:]

	// Extract offsets from Loki's footer
	symbolsEnd := binary.BigEndian.Uint64(footerBytes[0:8])
	seriesEnd := binary.BigEndian.Uint64(footerBytes[8:16])
	labelIndicesEnd := binary.BigEndian.Uint64(footerBytes[16:24])
	postingsEnd := binary.BigEndian.Uint64(footerBytes[24:32])
	labelIndicesTableEnd := binary.BigEndian.Uint64(footerBytes[32:40])
	postingsTableEnd := binary.BigEndian.Uint64(footerBytes[40:48])

	// Extract timestamps
	minTimeMs := binary.BigEndian.Uint64(footerBytes[48:56])
	maxTimeMs := binary.BigEndian.Uint64(footerBytes[56:64])

	fmt.Printf("==================================================\n")
	fmt.Printf("LOKI INDEX METADATA\n")
	fmt.Printf("==================================================\n")
	fmt.Printf("Uncompressed File Size: %d bytes\n", len(fileBytes))
	fmt.Printf("Index Time Range:       %s to %s\n\n",
		time.Unix(int64(minTimeMs/1000), 0).Format(time.RFC3339),
		time.Unix(int64(maxTimeMs/1000), 0).Format(time.RFC3339),
	)

	// Synthesize the 48-byte Prometheus Table of Contents (TOC)
	promSymbolsOffset := uint64(5) // Magic (4 bytes) + Version (1 byte)
	promSeriesOffset := symbolsEnd
	promPostingsOffset := seriesEnd
	promLabelIndicesOffset := postingsEnd
	promLabelIndicesTableOffset := labelIndicesEnd
	promPostingsTableOffset := labelIndicesTableEnd
	promIndexEndOffset := postingsTableEnd

	toc := make([]byte, 48)
	binary.BigEndian.PutUint64(toc[0:8], promSymbolsOffset)
	binary.BigEndian.PutUint64(toc[8:16], promSeriesOffset)
	binary.BigEndian.PutUint64(toc[16:24], promPostingsOffset)
	binary.BigEndian.PutUint64(toc[24:32], promLabelIndicesOffset)
	binary.BigEndian.PutUint64(toc[32:40], promLabelIndicesTableOffset)
	binary.BigEndian.PutUint64(toc[40:48], promPostingsTableOffset)

	// Compute CRC32-Castagnoli checksum for the TOC
	castagnoliTable := crc32.MakeTable(crc32.Castagnoli)
	checksum := crc32.Checksum(toc, castagnoliTable)

	// Assemble the 52-byte final Prometheus TOC block
	tocWithChecksum := make([]byte, 52)
	copy(tocWithChecksum, toc)
	binary.BigEndian.PutUint32(tocWithChecksum[48:52], checksum)

	// Recreate the standard Prometheus TSDB index file byte structure
	promIndexData := fileBytes[:promIndexEndOffset]
	synthesizedPromBytes := append(promIndexData, tocWithChecksum...)

	// Open Prometheus index reader to read Symbols and Postings
	r, err := index.NewReader(realByteSlice(synthesizedPromBytes), index.DecodePostingsRaw)
	if err != nil {
		return fmt.Errorf("failed to create Prometheus index reader: %w", err)
	}
	defer r.Close()

	// 1. Read Symbols
	fmt.Printf("==================================================\n")
	fmt.Printf("INDEX SYMBOL TABLE\n")
	fmt.Printf("==================================================\n")
	syms := r.Symbols()
	var symbols []string
	for syms.Next() {
		symbols = append(symbols, syms.At())
		fmt.Printf("  - %s\n", syms.At())
	}
	fmt.Printf("Total symbols count: %d\n\n", len(symbols))

	// 2. Read Postings (Label-Value associations)
	fmt.Printf("==================================================\n")
	fmt.Printf("INDEX POSTINGS (Label -> Value -> Series references)\n")
	fmt.Printf("==================================================\n")
	ctx := context.Background()
	labelNames, err := r.LabelNames(ctx)
	if err != nil {
		return fmt.Errorf("failed to read label names: %w", err)
	}

	// Keep track of all unique series reference IDs we find
	uniqueSeriesRefs := make(map[uint64]bool)

	for _, name := range labelNames {
		values, err := r.LabelValues(ctx, name, nil)
		if err != nil {
			continue
		}
		fmt.Printf("Label %q:\n", name)
		for _, val := range values {
			p, err := r.Postings(ctx, name, val)
			if err == nil {
				var refs []uint64
				for p.Next() {
					refVal := uint64(p.At())
					refs = append(refs, refVal)
					uniqueSeriesRefs[refVal] = true
				}
				fmt.Printf("  - %q -> Series Refs: %v\n", val, refs)
			}
		}
	}
	fmt.Println()

	// 3. Custom Parse Series Section
	fmt.Printf("==================================================\n")
	fmt.Printf("INDEX SERIES RECORDS\n")
	fmt.Printf("==================================================\n")

	// Sort unique series references for deterministic printout
	var sortedRefs []uint64
	for ref := range uniqueSeriesRefs {
		sortedRefs = append(sortedRefs, ref)
	}
	sort.Slice(sortedRefs, func(i, j int) bool {
		return sortedRefs[i] < sortedRefs[j]
	})

	for _, ref := range sortedRefs {
		fileOffset := ref * 16
		if fileOffset >= uint64(len(fileBytes)) {
			fmt.Printf("⚠️  Series ref %d offset %d out of file range.\n", ref, fileOffset)
			continue
		}

		recReader := bytes.NewReader(fileBytes[fileOffset:])

		// Read record length (uvarint)
		recLen, err := binary.ReadUvarint(recReader)
		if err != nil {
			fmt.Printf("⚠️  Failed to read length for series ref %d: %v\n", ref, err)
			continue
		}

		payloadBytes := make([]byte, recLen)
		if _, err := io.ReadFull(recReader, payloadBytes); err != nil {
			fmt.Printf("⚠️  Failed to read payload of length %d for series ref %d: %v\n", recLen, ref, err)
			continue
		}

		payloadBuf := bytes.NewReader(payloadBytes)

		// Read Fingerprint (8 bytes)
		var fingerprint uint64
		if err := binary.Read(payloadBuf, binary.BigEndian, &fingerprint); err != nil {
			continue
		}

		// Read Labels count (uvarint)
		labelsCount, _ := binary.ReadUvarint(payloadBuf)

		var seriesLabels []string
		for l := uint64(0); l < labelsCount; l++ {
			kIdx, _ := binary.ReadUvarint(payloadBuf)
			vIdx, _ := binary.ReadUvarint(payloadBuf)

			var key, val string
			if kIdx < uint64(len(symbols)) {
				key = symbols[kIdx]
			} else {
				key = fmt.Sprintf("<invalid_idx_%d>", kIdx)
			}
			if vIdx < uint64(len(symbols)) {
				val = symbols[vIdx]
			} else {
				val = fmt.Sprintf("<invalid_idx_%d>", vIdx)
			}
			seriesLabels = append(seriesLabels, fmt.Sprintf("%s=%q", key, val))
		}

		// Read Chunks list count (uvarint)
		chunksCount, _ := binary.ReadUvarint(payloadBuf)

		fmt.Printf("Series %d (File Offset: %d):\n", ref, fileOffset)
		fmt.Printf("  Fingerprint:  %016x (%d)\n", fingerprint, fingerprint)
		fmt.Printf("  Labels:       %v\n", seriesLabels)
		fmt.Printf("  Chunks Count: %d\n", chunksCount)

		// The remaining bytes contain Loki chunk metadata
		remBytes := make([]byte, payloadBuf.Len())
		payloadBuf.Read(remBytes)
		if len(remBytes) > 0 {
			fmt.Printf("  Raw Chunk metadata block (hex): %x\n", remBytes)
		}
		fmt.Println()
	}
	fmt.Printf("==================================================\n")
	fmt.Printf("✅ Finished reading TSDB index.\n")
	return nil
}
