package jsonschema

import (
	"encoding/json"
	"testing"

	"github.com/go-faster/yaml"
	"github.com/stretchr/testify/require"
)

func TestRawType_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		rawType  RawType
		expected string
	}{
		{
			name:     "single type",
			rawType:  RawType{"string"},
			expected: `"string"`,
		},
		{
			name:     "empty type",
			rawType:  RawType{},
			expected: `[]`,
		},
		{
			name:     "multiple types",
			rawType:  RawType{"string", "null"},
			expected: `["string","null"]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.rawType)
			require.NoError(t, err)
			require.Equal(t, tt.expected, string(data))
		})
	}
}

func TestRawType_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected RawType
		wantErr  bool
	}{
		{
			name:     "single type",
			input:    `"string"`,
			expected: RawType{"string"},
		},
		{
			name:     "multiple types",
			input:    `["string","null"]`,
			expected: RawType{"string", "null"},
		},
		{
			name:     "union type",
			input:    `["string","integer"]`,
			expected: RawType{"string", "integer"},
		},
		{
			name:     "single element array",
			input:    `["string"]`,
			expected: RawType{"string"},
		},
		{
			name:    "invalid type",
			input:   `123`,
			wantErr: true,
		},
		{
			name:    "invalid object",
			input:   `{"type": "string"}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var rawType RawType
			err := json.Unmarshal([]byte(tt.input), &rawType)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.expected, rawType)
		})
	}
}

func TestRawType_MarshalYAML(t *testing.T) {
	tests := []struct {
		name     string
		rawType  RawType
		expected string
	}{
		{
			name:     "single type",
			rawType:  RawType{"string"},
			expected: "string\n",
		},
		{
			name:     "multiple types",
			rawType:  RawType{"string", "null"},
			expected: "- string\n- \"null\"\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := yaml.Marshal(tt.rawType)
			require.NoError(t, err)
			require.Equal(t, tt.expected, string(data))
		})
	}
}

func TestRawType_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected RawType
		wantErr  bool
	}{
		{
			name:     "single type",
			input:    "string",
			expected: RawType{"string"},
		},
		{
			name:     "multiple types",
			input:    "- string\n- \"null\"",
			expected: RawType{"string", "null"},
		},
		{
			name:     "single element array",
			input:    "- string",
			expected: RawType{"string"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var rawType RawType
			err := yaml.Unmarshal([]byte(tt.input), &rawType)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.expected, rawType)
		})
	}
}
