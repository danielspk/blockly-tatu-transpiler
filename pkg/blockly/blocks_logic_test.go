package blockly

import (
	"strings"
	"testing"
)

func TestTranspileLogicBoolean(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected string
	}{
		{"true value", "TRUE", "true"},
		{"false value", "FALSE", "false"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData := `{
				"blocks": {
					"languageVersion": 0,
					"blocks": [{
						"type": "logic_boolean",
						"id": "bool1",
						"fields": {"BOOL": "` + tt.value + `"}
					}]
				}
			}`

			result, err := Transpile([]byte(jsonData))
			if err != nil {
				t.Fatalf("Transpile failed: %v", err)
			}

			if !strings.Contains(result, tt.expected) {
				t.Errorf("Expected %s, got: %s", tt.expected, result)
			}
		})
	}
}

func TestTranspileLogicOperation(t *testing.T) {
	tests := []struct {
		name     string
		op       string
		expected string
	}{
		{"AND operation", "AND", "(and true false)"},
		{"OR operation", "OR", "(or true false)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData := `{
				"blocks": {
					"languageVersion": 0,
					"blocks": [{
						"type": "logic_operation",
						"id": "op1",
						"fields": {"OP": "` + tt.op + `"},
						"inputs": {
							"A": {"block": {"type": "logic_boolean", "id": "bool1", "fields": {"BOOL": "TRUE"}}},
							"B": {"block": {"type": "logic_boolean", "id": "bool2", "fields": {"BOOL": "FALSE"}}}
						}
					}]
				}
			}`

			result, err := Transpile([]byte(jsonData))
			if err != nil {
				t.Fatalf("Transpile failed: %v", err)
			}

			if !strings.Contains(result, tt.expected) {
				t.Errorf("Expected %s, got: %s", tt.expected, result)
			}
		})
	}
}

func TestTranspileLogicCompare(t *testing.T) {
	jsonData := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": [{
				"type": "logic_compare",
				"id": "cmp1",
				"fields": {"OP": "LT"},
				"inputs": {
					"A": {"block": {"type": "math_number", "id": "num1", "fields": {"NUM": 5}}},
					"B": {"block": {"type": "math_number", "id": "num2", "fields": {"NUM": 10}}}
				}
			}]
		}
	}`

	result, err := Transpile([]byte(jsonData))
	if err != nil {
		t.Fatalf("Transpile failed: %v", err)
	}

	if strings.TrimSpace(result) != "(< 5 10)" {
		t.Errorf("Expected (< 5 10), got: %s", result)
	}
}

func TestTranspileLogicNegate(t *testing.T) {
	jsonData := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": [{
				"type": "logic_negate",
				"id": "neg1",
				"inputs": {
					"BOOL": {"block": {"type": "logic_boolean", "id": "bool1", "fields": {"BOOL": "TRUE"}}}
				}
			}]
		}
	}`

	result, err := Transpile([]byte(jsonData))
	if err != nil {
		t.Fatalf("Transpile failed: %v", err)
	}

	if !strings.Contains(result, "(if true false true)") {
		t.Errorf("Expected negation expression, got: %s", result)
	}
}

func TestTranspileLogicTernary(t *testing.T) {
	jsonData := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": [{
				"type": "logic_ternary",
				"id": "ternary1",
				"inputs": {
					"IF": {"block": {"type": "logic_boolean", "id": "bool1", "fields": {"BOOL": "TRUE"}}},
					"THEN": {"block": {"type": "text", "id": "text1", "fields": {"TEXT": "yes"}}},
					"ELSE": {"block": {"type": "text", "id": "text2", "fields": {"TEXT": "no"}}}
				}
			}]
		}
	}`

	result, err := Transpile([]byte(jsonData))
	if err != nil {
		t.Fatalf("Transpile failed: %v", err)
	}

	if !strings.Contains(result, "(if true") || !strings.Contains(result, `"yes"`) || !strings.Contains(result, `"no"`) {
		t.Errorf("Expected ternary expression, got: %s", result)
	}
}

func TestTranspileLogicNull(t *testing.T) {
	jsonData := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": [{
				"type": "logic_null",
				"id": "null1"
			}]
		}
	}`

	result, err := Transpile([]byte(jsonData))
	if err != nil {
		t.Fatalf("Transpile failed: %v", err)
	}

	if !strings.Contains(result, "nil") {
		t.Errorf("Expected nil, got: %s", result)
	}
}
