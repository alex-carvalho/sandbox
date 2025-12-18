package properties

import (
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    map[string]string
		wantErr bool
	}{
		{
			name: "simple properties",
			input: `
key1=value1
key2=value2
`,
			want: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
			wantErr: false,
		},
		{
			name: "with comments and empty lines",
			input: `
# This is a comment
key1=value1

# Another comment
key2=value2
`,
			want: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
			wantErr: false,
		},
		{
			name: "with whitespace",
			input: `
  key1  =  value1  
key2=value2
`,
			want: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
			wantErr: false,
		},
		{
			name:    "invalid line",
			input:   "invalid_line_without_equals",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "empty key",
			input:   "=value",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(strings.NewReader(tt.input))
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !mapsEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMerge(t *testing.T) {
	map1 := map[string]string{"key1": "value1", "key2": "value2"}
	map2 := map[string]string{"key2": "overridden", "key3": "value3"}
	map3 := map[string]string{"key1": "new_value1"}

	result := Merge(map1, map2, map3)

	expected := map[string]string{
		"key1": "new_value1",
		"key2": "overridden",
		"key3": "value3",
	}

	if !mapsEqual(result, expected) {
		t.Errorf("Merge() = %v, want %v", result, expected)
	}
}

func mapsEqual(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if bv, ok := b[k]; !ok || bv != v {
			return false
		}
	}
	return true
}
