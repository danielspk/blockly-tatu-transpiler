package blockly

import (
	"strings"
	"testing"
)

func TestTranspileMathNumber(t *testing.T) {
	jsonData := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": [{
				"type": "math_number",
				"id": "num1",
				"fields": {"NUM": 42}
			}]
		}
	}`

	result, err := Transpile([]byte(jsonData))
	if err != nil {
		t.Fatalf("Transpile failed: %v", err)
	}

	if strings.TrimSpace(result) != "42" {
		t.Errorf("Expected 42, got: %s", result)
	}
}

func TestTranspileMathArithmetic(t *testing.T) {
	jsonData := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": [{
				"type": "math_arithmetic",
				"id": "add1",
				"fields": {"OP": "ADD"},
				"inputs": {
					"A": {"block": {"type": "math_number", "id": "num1", "fields": {"NUM": 5}}},
					"B": {"block": {"type": "math_number", "id": "num2", "fields": {"NUM": 3}}}
				}
			}]
		}
	}`

	result, err := Transpile([]byte(jsonData))
	if err != nil {
		t.Fatalf("Transpile failed: %v", err)
	}

	if strings.TrimSpace(result) != "(+ 5 3)" {
		t.Errorf("Expected (+ 5 3), got: %s", result)
	}
}

func TestTranspileMathSingle(t *testing.T) {
	jsonData := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": [{
				"type": "math_single",
				"id": "sqrt1",
				"fields": {"OP": "ROOT"},
				"inputs": {
					"NUM": {"block": {"type": "math_number", "id": "num1", "fields": {"NUM": 16}}}
				}
			}]
		}
	}`

	result, err := Transpile([]byte(jsonData))
	if err != nil {
		t.Fatalf("Transpile failed: %v", err)
	}

	if !strings.Contains(result, "math:sqrt") {
		t.Errorf("Expected math:sqrt, got: %s", result)
	}
}

func TestTranspileMathConstant(t *testing.T) {
	jsonData := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": [{
				"type": "math_constant",
				"id": "pi1",
				"fields": {"CONSTANT": "PI"}
			}]
		}
	}`

	result, err := Transpile([]byte(jsonData))
	if err != nil {
		t.Fatalf("Transpile failed: %v", err)
	}

	if strings.TrimSpace(result) != "(math:pi)" {
		t.Errorf("Expected (math:pi), got: %s", result)
	}
}

func TestTranspileMathConstrain(t *testing.T) {
	jsonData := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": [{
				"type": "math_constrain",
				"id": "constrain1",
				"inputs": {
					"VALUE": {"block": {"type": "math_number", "id": "num1", "fields": {"NUM": 50}}},
					"LOW": {"block": {"type": "math_number", "id": "num2", "fields": {"NUM": 0}}},
					"HIGH": {"block": {"type": "math_number", "id": "num3", "fields": {"NUM": 100}}}
				}
			}]
		}
	}`

	result, err := Transpile([]byte(jsonData))
	if err != nil {
		t.Fatalf("Transpile failed: %v", err)
	}

	if !strings.Contains(result, "math:min") && !strings.Contains(result, "math:max") {
		t.Errorf("Expected constrain expression, got: %s", result)
	}
}

func TestTranspileMathTrig(t *testing.T) {
	tests := []struct {
		name     string
		op       string
		expected string
	}{
		{"sine", "SIN", "math:sin"},
		{"cosine", "COS", "math:cos"},
		{"tangent", "TAN", "math:tan"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData := `{
				"blocks": {
					"languageVersion": 0,
					"blocks": [{
						"type": "math_trig",
						"id": "trig1",
						"fields": {"OP": "` + tt.op + `"},
						"inputs": {
							"NUM": {"block": {"type": "math_number", "id": "num1", "fields": {"NUM": 45}}}
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

func TestTranspileMathModulo(t *testing.T) {
	jsonData := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": [{
				"type": "math_modulo",
				"id": "mod1",
				"inputs": {
					"DIVIDEND": {"block": {"type": "math_number", "id": "num1", "fields": {"NUM": 10}}},
					"DIVISOR": {"block": {"type": "math_number", "id": "num2", "fields": {"NUM": 3}}}
				}
			}]
		}
	}`

	result, err := Transpile([]byte(jsonData))
	if err != nil {
		t.Fatalf("Transpile failed: %v", err)
	}

	if !strings.Contains(result, "%") {
		t.Errorf("Expected modulo operator, got: %s", result)
	}
}

func TestTranspileMathRandomInt(t *testing.T) {
	jsonData := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": [{
				"type": "math_random_int",
				"id": "rand1",
				"inputs": {
					"FROM": {"block": {"type": "math_number", "id": "num1", "fields": {"NUM": 1}}},
					"TO": {"block": {"type": "math_number", "id": "num2", "fields": {"NUM": 100}}}
				}
			}]
		}
	}`

	result, err := Transpile([]byte(jsonData))
	if err != nil {
		t.Fatalf("Transpile failed: %v", err)
	}

	if !strings.Contains(result, "math:rand") {
		t.Errorf("Expected math:rand, got: %s", result)
	}
}

func TestTranspileMathRandomFloat(t *testing.T) {
	jsonData := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": [{
				"type": "math_random_float",
				"id": "rand1"
			}]
		}
	}`

	result, err := Transpile([]byte(jsonData))
	if err != nil {
		t.Fatalf("Transpile failed: %v", err)
	}

	if !strings.Contains(result, "math:rand") {
		t.Errorf("Expected math:rand, got: %s", result)
	}
}

func TestTranspileMathOnList(t *testing.T) {
	jsonData := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": [{
				"type": "math_on_list",
				"id": "onlist1",
				"fields": {"OP": "SUM"},
				"inputs": {
					"LIST": {"block": {"type": "lists_create_with", "id": "list1", "extraState": {"itemCount": 3}, "inputs": {
						"ADD0": {"block": {"type": "math_number", "id": "num1", "fields": {"NUM": 1}}},
						"ADD1": {"block": {"type": "math_number", "id": "num2", "fields": {"NUM": 2}}},
						"ADD2": {"block": {"type": "math_number", "id": "num3", "fields": {"NUM": 3}}}
					}}}
				}
			}]
		}
	}`

	_, err := Transpile([]byte(jsonData))
	if err == nil {
		t.Fatalf("Expected error for unsupported math_on_list operation")
	}

	if !strings.Contains(err.Error(), "not supported") {
		t.Errorf("Expected 'not supported' error, got: %v", err)
	}
}

func TestTranspileMathRound(t *testing.T) {
	jsonData := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": [{
				"type": "math_round",
				"id": "round1",
				"fields": {"OP": "ROUND"},
				"inputs": {
					"NUM": {"block": {"type": "math_number", "id": "num1", "fields": {"NUM": 3.5}}}
				}
			}]
		}
	}`

	result, err := Transpile([]byte(jsonData))
	if err != nil {
		t.Fatalf("Transpile failed: %v", err)
	}

	if !strings.Contains(result, "math:round") {
		t.Errorf("Expected math:round, got: %s", result)
	}
}

func TestTranspileMathChange(t *testing.T) {
	jsonData := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": [{
				"type": "math_change",
				"id": "change1",
				"fields": {"VAR": "counter"},
				"inputs": {
					"DELTA": {"block": {"type": "math_number", "id": "num1", "fields": {"NUM": 1}}}
				}
			}]
		}
	}`

	result, err := Transpile([]byte(jsonData))
	if err != nil {
		t.Fatalf("Transpile failed: %v", err)
	}

	if !strings.Contains(result, "(set counter (+ counter 1))") {
		t.Errorf("Expected math change expression, got: %s", result)
	}
}

func TestTranspileMathNumberProperty(t *testing.T) {
	tests := []struct {
		name     string
		property string
		contains string
	}{
		{"even", "EVEN", "= (% "},
		{"odd", "ODD", "(not (= (% "},
		{"positive", "POSITIVE", "> "},
		{"negative", "NEGATIVE", "< "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData := `{
				"blocks": {
					"languageVersion": 0,
					"blocks": [{
						"type": "math_number_property",
						"id": "prop1",
						"fields": {"PROPERTY": "` + tt.property + `"},
						"inputs": {
							"NUMBER_TO_CHECK": {"block": {"type": "math_number", "id": "num1", "fields": {"NUM": 10}}}
						}
					}]
				}
			}`

			result, err := Transpile([]byte(jsonData))
			if err != nil {
				t.Fatalf("Transpile failed: %v", err)
			}

			if !strings.Contains(result, tt.contains) {
				t.Errorf("Expected %s in result, got: %s", tt.contains, result)
			}
		})
	}
}
