package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	"github.com/prometheus/prometheus/tsdb/record"
	"github.com/prometheus/prometheus/tsdb/wlog"
)

// Loki Series Record
type LokiSeries struct {
	Fingerprint uint64
	Labels      []Label
}

type Label struct {
	Name  string
	Value string
}

// Loki Logs Record
type LokiLogs struct {
	UserID      string
	Timestamp   time.Time
	Fingerprint uint64
	SeriesRef   uint64
	Entries     []LokiEntry
}

type LokiEntry struct {
	Timestamp      time.Time
	Line           string
	StructuredMeta []Label
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <wal-file-or-directory-path>")
		fmt.Println("Example: go run main.go ../loki-data/wal/")
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
			baseName := filepath.Base(path)
			if baseName[0] == '.' {
				return nil
			}
			// Skip temporary checkpoint files/folders
			if strings.Contains(path, ".tmp") {
				return nil
			}
			// Only process files that have numeric names (WAL segments like 00000000)
			if !isNumeric(baseName) {
				return nil
			}

			fmt.Printf("\n##################################################\n")
			fmt.Printf("📄 Processing WAL Segment: %s\n", path)
			fmt.Printf("##################################################\n")
			if parseErr := parseWalFile(path); parseErr != nil {
				fmt.Printf("⚠️  Error parsing %s: %v\n", path, parseErr)
			}
			return nil
		})
		if err != nil {
			fmt.Printf("❌ Error walking directory: %v\n", err)
			os.Exit(1)
		}
	} else {
		if err := parseWalFile(targetPath); err != nil {
			fmt.Printf("❌ Error parsing file: %v\n", err)
			os.Exit(1)
		}
	}
}

func isNumeric(s string) bool {
	if len(s) == 0 {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func parseWalFile(filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open WAL file: %w", err)
	}
	defer f.Close()

	fmt.Printf("📖 Reading and Parsing Loki WAL file: %s\n\n", filePath)

	r := wlog.NewReader(f)
	recordCount := 0

	for r.Next() {
		recordCount++
		rec := r.Record()
		recType := rec[0]

		fmt.Printf("--------------------------------------------------\n")
		fmt.Printf("Record #%d | Type: %s (%d) | Size: %d bytes\n", recordCount, typeToString(recType), recType, len(rec))
		fmt.Printf("--------------------------------------------------\n")

		switch recType {
		case 1: // Series/Stream Record
			seriesList, userID, err := parseLokiSeriesRecord(rec)
			if err == nil {
				fmt.Printf("Tenant/UserID: %q\n", userID)
				fmt.Println("Registered Streams:")
				for _, s := range seriesList {
					fmt.Printf("  - Fingerprint: %016x (%d)\n", s.Fingerprint, s.Fingerprint)
					fmt.Println("    Labels:")
					for _, l := range s.Labels {
						fmt.Printf("      %s=%q\n", l.Name, l.Value)
					}
				}
			} else {
				fmt.Printf("  [Failed to parse Loki Series: %v]\n", err)
				printFallbackStrings(rec)
			}

		case 5: // Loki Logs Record
			logs, err := parseLokiLogsRecord(rec)
			if err == nil {
				fmt.Printf("Tenant/UserID: %q\n", logs.UserID)
				fmt.Printf("Batch Reference Time: %s\n", logs.Timestamp.Format(time.RFC3339Nano))
				fmt.Printf("Stream Fingerprint: %016x\n", logs.Fingerprint)
				fmt.Printf("Series Reference ID: %d\n", logs.SeriesRef)
				fmt.Printf("Log Entries (%d):\n", len(logs.Entries))
				for i, entry := range logs.Entries {
					fmt.Printf("  [%d] Time: %s\n", i+1, entry.Timestamp.Format(time.RFC3339Nano))
					fmt.Printf("      Line: %q\n", entry.Line)
					if len(entry.StructuredMeta) > 0 {
						fmt.Println("      Structured Metadata:")
						for _, m := range entry.StructuredMeta {
							fmt.Printf("        %s=%q\n", m.Name, m.Value)
						}
					}
				}
			} else {
				fmt.Printf("  [Failed to parse Loki Logs: %v]\n", err)
				printFallbackStrings(rec)
			}

		default:
			printFallbackStrings(rec)
		}
		fmt.Println()
	}

	if err := r.Err(); err != nil {
		return fmt.Errorf("error encountered during WAL reading: %w", err)
	}

	fmt.Printf("✅ Finished reading WAL. Total records: %d\n", recordCount)
	return nil
}

// Parse Loki Type 1 Record (Series metadata)
func parseLokiSeriesRecord(b []byte) ([]LokiSeries, string, error) {
	if len(b) < 2 {
		return nil, "", fmt.Errorf("record too short")
	}

	offset := 1 // skip type byte
	userID, err := readString(b, &offset)
	if err != nil {
		return nil, "", err
	}

	seriesCount, err := readUvarint(b, &offset)
	if err != nil {
		return nil, "", err
	}

	var seriesList []LokiSeries
	for i := uint64(0); i < seriesCount; i++ {
		if offset+8 > len(b) {
			return nil, "", fmt.Errorf("unexpected end of series fingerprint")
		}
		fingerprint := binary.BigEndian.Uint64(b[offset : offset+8])
		offset += 8

		labelCount, err := readUvarint(b, &offset)
		if err != nil {
			return nil, "", err
		}

		var labels []Label
		for j := uint64(0); j < labelCount; j++ {
			name, err := readString(b, &offset)
			if err != nil {
				return nil, "", err
			}
			val, err := readString(b, &offset)
			if err != nil {
				return nil, "", err
			}
			labels = append(labels, Label{Name: name, Value: val})
		}

		seriesList = append(seriesList, LokiSeries{
			Fingerprint: fingerprint,
			Labels:      labels,
		})
	}

	return seriesList, userID, nil
}

// Parse Loki Type 5 Record (Log Entries)
func parseLokiLogsRecord(b []byte) (*LokiLogs, error) {
	if len(b) < 26 {
		return nil, fmt.Errorf("record too short")
	}

	offset := 1 // skip type byte
	userID, err := readString(b, &offset)
	if err != nil {
		return nil, err
	}

	if offset+8 > len(b) {
		return nil, fmt.Errorf("unexpected end of batch timestamp")
	}
	batchTimeNanos := int64(binary.BigEndian.Uint64(b[offset : offset+8]))
	batchTime := time.Unix(0, batchTimeNanos)
	offset += 8

	if offset+8 > len(b) {
		return nil, fmt.Errorf("unexpected end of stream fingerprint")
	}
	fingerprint := binary.BigEndian.Uint64(b[offset : offset+8])
	offset += 8

	if offset+8 > len(b) {
		return nil, fmt.Errorf("unexpected end of series ref")
	}
	seriesRef := binary.BigEndian.Uint64(b[offset : offset+8])
	offset += 8

	entryCount, err := readUvarint(b, &offset)
	if err != nil {
		return nil, err
	}

	var entries []LokiEntry
	for i := uint64(0); i < entryCount; i++ {
		// Read timestamp delta
		delta, err := readUvarint(b, &offset)
		if err != nil {
			return nil, err
		}
		entryTime := time.Unix(0, batchTimeNanos+int64(delta))

		// Read line length (uvarint)
		lineLen, err := readUvarint(b, &offset)
		if err != nil {
			return nil, err
		}

		if offset+int(lineLen) > len(b) {
			return nil, fmt.Errorf("unexpected end of line bytes")
		}
		line := string(b[offset : offset+int(lineLen)])
		offset += int(lineLen)

		// Read structured metadata count
		metaCount, err := readUvarint(b, &offset)
		if err != nil {
			return nil, err
		}

		var structuredMeta []Label
		for j := uint64(0); j < metaCount; j++ {
			name, err := readString(b, &offset)
			if err != nil {
				return nil, err
			}
			val, err := readString(b, &offset)
			if err != nil {
				return nil, err
			}
			structuredMeta = append(structuredMeta, Label{Name: name, Value: val})
		}

		entries = append(entries, LokiEntry{
			Timestamp:      entryTime,
			Line:           line,
			StructuredMeta: structuredMeta,
		})
	}

	return &LokiLogs{
		UserID:      userID,
		Timestamp:   batchTime,
		Fingerprint: fingerprint,
		SeriesRef:   seriesRef,
		Entries:     entries,
	}, nil
}

// Read uvarint from byte slice and advance offset
func readUvarint(b []byte, offset *int) (uint64, error) {
	if *offset >= len(b) {
		return 0, fmt.Errorf("read beyond record boundary")
	}
	val, n := binary.Uvarint(b[*offset:])
	if n <= 0 {
		return 0, fmt.Errorf("invalid varint encoding")
	}
	*offset += n
	return val, nil
}

// Read length-prefixed string from byte slice and advance offset
func readString(b []byte, offset *int) (string, error) {
	length, err := readUvarint(b, offset)
	if err != nil {
		return "", err
	}
	if *offset+int(length) > len(b) {
		return "", fmt.Errorf("string length exceeds record size")
	}
	str := string(b[*offset : *offset+int(length)])
	*offset += int(length)
	return str, nil
}

// Convert type constants to string representation
func typeToString(t byte) string {
	switch record.Type(t) {
	case record.Series:
		return "Loki Series/Stream"
	case record.Samples:
		return "Samples"
	case record.Tombstones:
		return "Tombstones"
	case record.Metadata:
		return "Metadata"
	case 5:
		return "Loki Logs Record"
	default:
		return "Unknown"
	}
}

func printFallbackStrings(rec []byte) {
	printable := extractPrintableStrings(rec, 4)
	if len(printable) > 0 {
		fmt.Println("Printable Strings / Log Content (Fallback):")
		for _, str := range printable {
			fmt.Printf("  %s\n", str)
		}
	} else {
		fmt.Println("  [No readable strings found]")
	}
}

// Helper to extract sequences of printable ASCII characters
func extractPrintableStrings(data []byte, minLength int) []string {
	var result []string
	var current []rune

	for _, b := range data {
		r := rune(b)
		if unicode.IsPrint(r) && b >= 32 && b <= 126 {
			current = append(current, r)
		} else {
			if len(current) >= minLength {
				result = append(result, string(current))
			}
			current = nil
		}
	}
	if len(current) >= minLength {
		result = append(result, string(current))
	}
	return result
}
