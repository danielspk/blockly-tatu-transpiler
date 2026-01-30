package blockly

import (
	"strings"
	"testing"
)

func TestTranspileInclude(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
		expected string
		wantErr  bool
	}{
		{
			name: "include with simple path",
			jsonData: `{
				"blocks": {
					"languageVersion": 0,
					"blocks": [{
						"type": "include_file_tatu",
						"id": "include1",
						"fields": {"PATH": "lib/utils.tatu"}
					}]
				}
			}`,
			expected: `(include "lib/utils.tatu")`,
			wantErr:  false,
		},
		{
			name: "include with nested path",
			jsonData: `{
				"blocks": {
					"languageVersion": 0,
					"blocks": [{
						"type": "include_file_tatu",
						"id": "include2",
						"fields": {"PATH": "lib/math/geometry.tatu"}
					}]
				}
			}`,
			expected: `(include "lib/math/geometry.tatu")`,
			wantErr:  false,
		},
		{
			name: "include with absolute path",
			jsonData: `{
				"blocks": {
					"languageVersion": 0,
					"blocks": [{
						"type": "include_file_tatu",
						"id": "include3",
						"fields": {"PATH": "/usr/local/tatu/stdlib.tatu"}
					}]
				}
			}`,
			expected: `(include "/usr/local/tatu/stdlib.tatu")`,
			wantErr:  false,
		},
		{
			name: "include missing PATH field",
			jsonData: `{
				"blocks": {
					"languageVersion": 0,
					"blocks": [{
						"type": "include_file_tatu",
						"id": "include4",
						"fields": {}
					}]
				}
			}`,
			expected: "",
			wantErr:  true,
		},
		{
			name: "include with empty PATH",
			jsonData: `{
				"blocks": {
					"languageVersion": 0,
					"blocks": [{
						"type": "include_file_tatu",
						"id": "include5",
						"fields": {"PATH": ""}
					}]
				}
			}`,
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Transpile([]byte(tt.jsonData))

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Transpile failed: %v", err)
			}

			if !strings.Contains(result, tt.expected) {
				t.Errorf("Expected to contain %q, got: %s", tt.expected, result)
			}
		})
	}
}

func TestTranspileIncludeSequence(t *testing.T) {
	jsonData := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": [
				{
					"type": "include_file_tatu",
					"id": "include1",
					"fields": {"PATH": "lib/math.tatu"}
				},
				{
					"type": "include_file_tatu",
					"id": "include2",
					"fields": {"PATH": "lib/string.tatu"}
				}
			]
		}
	}`

	result, err := Transpile([]byte(jsonData))
	if err != nil {
		t.Fatalf("Transpile failed: %v", err)
	}

	if !strings.Contains(result, `(include "lib/math.tatu")`) {
		t.Errorf("Expected first include statement, got: %s", result)
	}

	if !strings.Contains(result, `(include "lib/string.tatu")`) {
		t.Errorf("Expected second include statement, got: %s", result)
	}
}
