package blockly

import (
	"strings"
	"testing"
)

func TestTranspileVariablesGetSet(t *testing.T) {
	jsonData := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": [{
				"type": "variables_set",
				"id": "set1",
				"fields": {"VAR": "myVar"},
				"inputs": {
					"VALUE": {"block": {"type": "math_number", "id": "num1", "fields": {"NUM": 42}}}
				},
				"next": {
					"block": {
						"type": "variables_get",
						"id": "get1",
						"fields": {"VAR": "myVar"}
					}
				}
			}]
		}
	}`

	result, err := Transpile([]byte(jsonData))
	if err != nil {
		t.Fatalf("Transpile failed: %v", err)
	}

	if !strings.Contains(result, "(set myVar 42)") {
		t.Errorf("Expected set statement, got: %s", result)
	}

	if !strings.Contains(result, "myVar") {
		t.Errorf("Expected variable get, got: %s", result)
	}
}

func TestTranspileVariablesWithMapping(t *testing.T) {
	jsonData := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": [{
				"type": "variables_set",
				"id": "set1",
				"fields": {"VAR": {"id": "xyz123"}},
				"inputs": {
					"VALUE": {"block": {"type": "text", "id": "text1", "fields": {"TEXT": "hello"}}}
				}
			}]
		},
		"variables": [
			{"id": "xyz123", "name": "userName"}
		]
	}`

	result, err := Transpile([]byte(jsonData))
	if err != nil {
		t.Fatalf("Transpile failed: %v", err)
	}

	if !strings.Contains(result, "userName") {
		t.Errorf("Expected mapped variable name 'userName', got: %s", result)
	}

	if strings.Contains(result, "xyz123") {
		t.Errorf("Should not contain raw ID 'xyz123', got: %s", result)
	}
}

func TestTranspileVariableCreate(t *testing.T) {
	jsonData := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": [{
				"type": "variables_create_tatu",
				"id": "create1",
				"fields": {"VAR": "counter"},
				"inputs": {
					"VALUE": {"block": {"type": "math_number", "id": "num1", "fields": {"NUM": 0}}}
				},
				"next": {
					"block": {
						"type": "variables_set",
						"id": "set1",
						"fields": {"VAR": "counter"},
						"inputs": {
							"VALUE": {"block": {"type": "math_number", "id": "num2", "fields": {"NUM": 5}}}
						}
					}
				}
			}]
		}
	}`

	result, err := Transpile([]byte(jsonData))
	if err != nil {
		t.Fatalf("Transpile failed: %v", err)
	}

	if !strings.Contains(result, "(var counter 0)") {
		t.Errorf("Expected variable creation with (var counter 0), got: %s", result)
	}

	if !strings.Contains(result, "(set counter 5)") {
		t.Errorf("Expected set statement, got: %s", result)
	}
}

func TestTranspileVariableCreateWithExpression(t *testing.T) {
	jsonData := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": [{
				"type": "variables_create_tatu",
				"id": "create1",
				"fields": {"VAR": "result"},
				"inputs": {
					"VALUE": {
						"block": {
							"type": "math_arithmetic",
							"id": "math1",
							"fields": {"OP": "ADD"},
							"inputs": {
								"A": {"block": {"type": "math_number", "id": "num1", "fields": {"NUM": 10}}},
								"B": {"block": {"type": "math_number", "id": "num2", "fields": {"NUM": 20}}}
							}
						}
					}
				}
			}]
		}
	}`

	result, err := Transpile([]byte(jsonData))
	if err != nil {
		t.Fatalf("Transpile failed: %v", err)
	}

	if !strings.Contains(result, "(var result (+ 10 20))") {
		t.Errorf("Expected variable creation with expression, got: %s", result)
	}
}

func TestTranspileNoAutomaticVariableDeclarations(t *testing.T) {
	jsonData := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": [{
				"type": "variables_set",
				"id": "set1",
				"fields": {"VAR": "envVar"},
				"inputs": {
					"VALUE": {"block": {"type": "math_number", "id": "num1", "fields": {"NUM": 100}}}
				}
			}]
		},
		"variables": [
			{"id": "var1", "name": "envVar"}
		]
	}`

	result, err := Transpile([]byte(jsonData))
	if err != nil {
		t.Fatalf("Transpile failed: %v", err)
	}

	if strings.Contains(result, "(var envVar nil)") {
		t.Errorf("Should not auto-generate variable declarations, got: %s", result)
	}

	if !strings.Contains(result, "(set envVar 100)") {
		t.Errorf("Expected set statement, got: %s", result)
	}
}
