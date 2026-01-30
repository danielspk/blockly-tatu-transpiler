package blockly

import (
	"strings"
	"testing"
)

func TestTranspileBasicOK(t *testing.T) {
	jsonData := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": [
				{
					"type": "text_print",
					"id": "print1",
					"inputs": {
						"TEXT": {
							"block": {
								"type": "text",
								"id": "text1",
								"fields": {
									"TEXT": "Hello"
								}
							}
						}
					}
				}
			]
		}
	}`

	result, err := Transpile([]byte(jsonData))
	if err != nil {
		t.Fatalf("Transpile failed: %v", err)
	}

	expected := `(print "Hello")`
	if strings.TrimSpace(result) != expected {
		t.Errorf("Expected %q, got %q", expected, strings.TrimSpace(result))
	}
}

func TestTranspileEmptyWorkspace(t *testing.T) {
	jsonData := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": []
		}
	}`

	result, err := Transpile([]byte(jsonData))
	if err != nil {
		t.Fatalf("Transpile failed: %v", err)
	}

	if strings.TrimSpace(result) != "" {
		t.Errorf("Expected empty result, got %q", result)
	}
}

func TestTranspileInvalidJSON(t *testing.T) {
	jsonData := `{invalid json}`

	_, err := Transpile([]byte(jsonData))
	if err == nil {
		t.Fatal("Expected error for invalid JSON, got nil")
	}
}

func TestTranspileUnsupportedBlockType(t *testing.T) {
	jsonData := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": [
				{
					"type": "unknown_block_type",
					"id": "block1"
				}
			]
		}
	}`

	_, err := Transpile([]byte(jsonData))
	if err == nil {
		t.Fatal("Expected error for unsupported block type, got nil")
	}

	if !strings.Contains(err.Error(), "unsupported block type") {
		t.Errorf("Expected 'unsupported block type' error, got: %v", err)
	}
}

func TestTranspileMissingRequiredField(t *testing.T) {
	jsonData := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": [
				{
					"type": "logic_boolean",
					"id": "bool1"
				}
			]
		}
	}`

	_, err := Transpile([]byte(jsonData))
	if err == nil {
		t.Fatal("Expected error for missing field, got nil")
	}

	if !strings.Contains(err.Error(), "BOOL") {
		t.Errorf("Expected error to mention 'BOOL' field, got: %v", err)
	}
}

func TestTranspileMissingRequiredValue(t *testing.T) {
	jsonData := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": [
				{
					"type": "logic_compare",
					"id": "cmp1",
					"fields": {
						"OP": "EQ"
					},
					"inputs": {
						"B": {
							"block": {
								"type": "math_number",
								"id": "num1",
								"fields": {
									"NUM": 5
								}
							}
						}
					}
				}
			]
		}
	}`

	_, err := Transpile([]byte(jsonData))
	if err == nil {
		t.Fatal("Expected error for missing value input, got nil")
	}

	if !strings.Contains(err.Error(), "A") {
		t.Errorf("Expected error to mention value 'A', got: %v", err)
	}
}

func TestTranspileInvalidOperator(t *testing.T) {
	jsonData := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": [
				{
					"type": "logic_compare",
					"id": "cmp1",
					"fields": {
						"OP": "INVALID_OP"
					},
					"inputs": {
						"A": {
							"block": {
								"type": "math_number",
								"id": "num1",
								"fields": {"NUM": 5}
							}
						},
						"B": {
							"block": {
								"type": "math_number",
								"id": "num2",
								"fields": {"NUM": 10}
							}
						}
					}
				}
			]
		}
	}`

	_, err := Transpile([]byte(jsonData))
	if err == nil {
		t.Fatal("Expected error for invalid operator, got nil")
	}

	if !strings.Contains(err.Error(), "INVALID_OP") {
		t.Errorf("Expected error to mention 'INVALID_OP', got: %v", err)
	}
}

func TestTranspileNestedBlockError(t *testing.T) {
	jsonData := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": [
				{
					"type": "text_print",
					"id": "print1",
					"inputs": {
						"TEXT": {
							"block": {
								"type": "logic_boolean",
								"id": "bool1"
							}
						}
					}
				}
			]
		}
	}`

	_, err := Transpile([]byte(jsonData))
	if err == nil {
		t.Fatal("Expected error for nested block with missing field, got nil")
	}

	if !strings.Contains(err.Error(), "bool1") || !strings.Contains(err.Error(), "BOOL") {
		t.Errorf("Expected error to mention nested block 'bool1' and field 'BOOL', got: %v", err)
	}
}

func TestTranspileBlockWithComment(t *testing.T) {
	jsonData := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": [{
				"type": "variables_set",
				"id": "test1",
				"icons": {
					"comment": {
						"text": "Line 1\nLine 2"
					}
				},
				"fields": {"VAR": {"id": "varX", "name": "x"}},
				"inputs": {
					"VALUE": {
						"block": {"type": "math_number", "id": "num1", "fields": {"NUM": 5}}
					}
				}
			}]
		},
		"variables": [{"id": "varX", "name": "x"}]
	}`

	result, err := Transpile([]byte(jsonData))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result = strings.TrimSpace(result)
	expected := "; Line 1\n; Line 2\n(set x 5)"

	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}
