package properties

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func Parse(r io.Reader) (map[string]string, error) {
	props := make(map[string]string)
	scanner := bufio.NewScanner(r)

	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		parts := strings.SplitN(trimmed, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid property at line %d: %s", lineNum, trimmed)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if key == "" {
			return nil, fmt.Errorf("empty key at line %d", lineNum)
		}

		props[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading properties: %w", err)
	}

	return props, nil
}

// Merge combines multiple property maps into one (later maps override earlier ones)
func Merge(maps ...map[string]string) map[string]string {
	result := make(map[string]string)

	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}

	return result
}
