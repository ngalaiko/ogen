package jsonschema

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTypeArraySupport(t *testing.T) {
	parser := NewParser(Settings{})

	tests := []struct {
		name     string
		input    *RawSchema
		expected *Schema
		wantErr  bool
	}{
		// Single types (backwards compatibility)
		{
			name:     "single string type",
			input:    &RawSchema{Type: RawType{"string"}},
			expected: &Schema{Type: String},
		},
		{
			name:     "single integer type",
			input:    &RawSchema{Type: RawType{"integer"}},
			expected: &Schema{Type: Integer},
		},
		{
			name:     "single number type",
			input:    &RawSchema{Type: RawType{"number"}},
			expected: &Schema{Type: Number},
		},
		{
			name:     "single boolean type",
			input:    &RawSchema{Type: RawType{"boolean"}},
			expected: &Schema{Type: Boolean},
		},
		{
			name:     "single array type",
			input:    &RawSchema{Type: RawType{"array"}},
			expected: &Schema{Type: Array},
		},
		{
			name:     "single object type",
			input:    &RawSchema{Type: RawType{"object"}},
			expected: &Schema{Type: Object},
		},
		{
			name:     "single null type",
			input:    &RawSchema{Type: RawType{"null"}},
			expected: &Schema{Type: Null},
		},
		{
			name:     "empty type",
			input:    &RawSchema{Type: RawType{}},
			expected: &Schema{Type: Empty},
		},

		// OpenAPI 3.1 nullable types (type + null)
		{
			name:     "nullable string - OpenAPI 3.1 style",
			input:    &RawSchema{Type: RawType{"string", "null"}},
			expected: &Schema{Type: String, Nullable: true},
		},
		{
			name:     "nullable integer - OpenAPI 3.1 style",
			input:    &RawSchema{Type: RawType{"integer", "null"}},
			expected: &Schema{Type: Integer, Nullable: true},
		},
		{
			name:     "nullable number - OpenAPI 3.1 style",
			input:    &RawSchema{Type: RawType{"number", "null"}},
			expected: &Schema{Type: Number, Nullable: true},
		},
		{
			name:     "nullable boolean - OpenAPI 3.1 style",
			input:    &RawSchema{Type: RawType{"boolean", "null"}},
			expected: &Schema{Type: Boolean, Nullable: true},
		},

		// Union types (converted to oneOf)
		{
			name:  "union: string or number",
			input: &RawSchema{Type: RawType{"string", "number"}},
			expected: &Schema{
				Type: Empty,
				OneOf: []*Schema{
					{Type: String},
					{Type: Number},
				},
			},
		},
		{
			name:  "union: string or integer",
			input: &RawSchema{Type: RawType{"string", "integer"}},
			expected: &Schema{
				Type: Empty,
				OneOf: []*Schema{
					{Type: String},
					{Type: Integer},
				},
			},
		},
		{
			name:  "complex union: string, integer, boolean",
			input: &RawSchema{Type: RawType{"string", "integer", "boolean"}},
			expected: &Schema{
				Type: Empty,
				OneOf: []*Schema{
					{Type: String},
					{Type: Integer},
					{Type: Boolean},
				},
			},
		},

		// Union types with nullable (type + type + null)
		{
			name:  "nullable union: string, number, or null",
			input: &RawSchema{Type: RawType{"string", "number", "null"}},
			expected: &Schema{
				Type:     Empty,
				Nullable: true,
				OneOf: []*Schema{
					{Type: String},
					{Type: Number},
				},
			},
		},
		{
			name:  "nullable complex union: string, integer, boolean, or null",
			input: &RawSchema{Type: RawType{"string", "integer", "boolean", "null"}},
			expected: &Schema{
				Type:     Empty,
				Nullable: true,
				OneOf: []*Schema{
					{Type: String},
					{Type: Integer},
					{Type: Boolean},
				},
			},
		},

		// Edge cases and validation
		{
			name:    "invalid type in array",
			input:   &RawSchema{Type: RawType{"string", "invalid"}},
			wantErr: true,
		},
		{
			name:    "another invalid type",
			input:   &RawSchema{Type: RawType{"integer", "unknown", "null"}},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parser.Parse(tt.input, testCtx())
			
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			
			require.NoError(t, err)

			// Zero out location pointers for comparison
			result.Pointer = tt.expected.Pointer
			if result.OneOf != nil {
				for i := range result.OneOf {
					result.OneOf[i].Pointer = tt.expected.OneOf[i].Pointer
				}
			}

			require.Equal(t, tt.expected, result)
		})
	}
}