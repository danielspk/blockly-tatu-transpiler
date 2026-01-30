package blockly

import (
	"strings"
	"testing"
)

// TestErrorEmptyJSON tests error handling for empty JSON data.
func TestErrorEmptyJSON(t *testing.T) {
	_, err := Transpile([]byte(""))
	if err == nil {
		t.Fatal("Expected error for empty JSON, got nil")
	}
	if !strings.Contains(err.Error(), "empty JSON data") {
		t.Errorf("Expected 'empty JSON data' error, got: %v", err)
	}
}

// TestErrorInvalidJSON tests error handling for invalid JSON syntax.
func TestErrorInvalidJSON(t *testing.T) {
	invalidJSON := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": [
				{"type": "logic_boolean" MISSING_COMMA "id": "bool1"}
			]
		}
	}`

	_, err := Transpile([]byte(invalidJSON))
	if err == nil {
		t.Fatal("Expected error for invalid JSON, got nil")
	}
	if !strings.Contains(err.Error(), "JSON syntax error") {
		t.Errorf("Expected 'JSON syntax error', got: %v", err)
	}
}

// TestErrorMissingBlocksArray tests error handling when blocks array is missing.
func TestErrorMissingBlocksArray(t *testing.T) {
	jsonData := `{
		"blocks": {
			"languageVersion": 0
		}
	}`

	_, err := Transpile([]byte(jsonData))
	if err == nil {
		t.Fatal("Expected error for missing blocks array, got nil")
	}
	if !strings.Contains(err.Error(), "missing 'blocks.blocks' array") {
		t.Errorf("Expected 'missing blocks array' error, got: %v", err)
	}
}

// TestErrorEmptyBlockType tests error handling for block with empty type.
func TestErrorEmptyBlockType(t *testing.T) {
	jsonData := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": [
				{
					"type": "",
					"id": "block1"
				}
			]
		}
	}`

	_, err := Transpile([]byte(jsonData))
	if err == nil {
		t.Fatal("Expected error for empty block type, got nil")
	}
	if !strings.Contains(err.Error(), "empty type") {
		t.Errorf("Expected 'empty type' error, got: %v", err)
	}
}

// TestErrorJSONTypeError tests handling of JSON type mismatches.
func TestErrorJSONTypeError(t *testing.T) {
	// NUM field should be a number but is given as a nested object
	jsonData := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": [
				{
					"type": "math_number",
					"id": "num1",
					"fields": {
						"NUM": {"invalid": "object"}
					}
				}
			]
		}
	}`

	_, err := Transpile([]byte(jsonData))
	// This might succeed or fail depending on how we handle field conversion
	// The test verifies we don't panic
	_ = err
}

// TestVariableMapping tests that variable IDs are correctly mapped to variable names.
func TestVariableMapping(t *testing.T) {
	jsonData := []byte(`{
  "blocks": {
    "blocks": [
      {
        "fields": {
          "VAR": {
            "id": "tOk#yCD[kF$q` + "`" + `k42w[(0"
          }
        },
        "id": "uT8s~m,hMlS0U+|nC_}a",
        "inputs": {
          "VALUE": {
            "block": {
              "extraState": "<mutation arguments=\"1\" function=\"str:upper\"></mutation>",
              "fields": {
                "NAME": "str:len"
              },
              "id": "u0C%tigvku~5+1GLDprd",
              "inputs": {
                "ARG0": {
                  "block": {
                    "fields": {
                      "TEXT": "hola"
                    },
                    "id": "-eT1?LEd|o?p#{XMV803",
                    "type": "text"
                  }
                }
              },
              "type": "string_call_tatu"
            }
          }
        },
        "next": {
          "block": {
            "id": "=4?RVZajoxi.BkGff}]H",
            "inputs": {
              "TEXT": {
                "block": {
                  "fields": {
                    "VAR": {
                      "id": "tOk#yCD[kF$q` + "`" + `k42w[(0"
                    }
                  },
                  "id": "G` + "`" + `Zey@j#Uu[YzH1SkUg[",
                  "type": "variables_get"
                }
              }
            },
            "type": "text_print"
          }
        },
        "type": "variables_set",
        "x": 150,
        "y": 30
      }
    ],
    "languageVersion": 0
  },
  "variables": [
    {
      "id": "tOk#yCD[kF$q` + "`" + `k42w[(0",
      "name": "nombre"
    },
    {
      "id": "NVo=^Lsg!M;BwlgFfxNi",
      "name": "item"
    }
  ]
}`)

	result, err := Transpile(jsonData)
	if err != nil {
		t.Fatalf("Transpile failed: %v", err)
	}

	if !strings.Contains(result, "nombre") {
		t.Errorf("Expected 'nombre' in output, got: %s", result)
	}

	if strings.Contains(result, "tOk#yCD[kF") {
		t.Errorf("Expected variable name, not ID. Got: %s", result)
	}
}

// TestNumberFormatting tests that numbers are formatted correctly without scientific notation.
func TestNumberFormatting(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
		want     string
		notWant  string
	}{
		{
			name: "large integer",
			jsonData: `{
				"blocks": {
					"languageVersion": 0,
					"blocks": [
						{
							"type": "math_number",
							"id": "num1",
							"fields": {"NUM": 1870000}
						}
					]
				}
			}`,
			want:    "1870000",
			notWant: "1.87e+06",
		},
		{
			name: "small decimal",
			jsonData: `{
				"blocks": {
					"languageVersion": 0,
					"blocks": [
						{
							"type": "math_number",
							"id": "num1",
							"fields": {"NUM": 3.14}
						}
					]
				}
			}`,
			want:    "3.14",
			notWant: "3.140000",
		},
		{
			name: "integer with .0",
			jsonData: `{
				"blocks": {
					"languageVersion": 0,
					"blocks": [
						{
							"type": "math_number",
							"id": "num1",
							"fields": {"NUM": 42.0}
						}
					]
				}
			}`,
			want:    "42",
			notWant: "42.0",
		},
		{
			name: "very large number",
			jsonData: `{
				"blocks": {
					"languageVersion": 0,
					"blocks": [
						{
							"type": "math_number",
							"id": "num1",
							"fields": {"NUM": 9999999999}
						}
					]
				}
			}`,
			want:    "9999999999",
			notWant: "e+",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Transpile([]byte(tt.jsonData))
			if err != nil {
				t.Fatalf("Transpile failed: %v", err)
			}

			if !strings.Contains(result, tt.want) {
				t.Errorf("Expected result to contain '%s', got: %s", tt.want, result)
			}

			if strings.Contains(result, tt.notWant) {
				t.Errorf("Expected result NOT to contain '%s', got: %s", tt.notWant, result)
			}
		})
	}
}

func TestParseParamsAsObjects(t *testing.T) {
	jsonData := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": [{
				"type": "procedures_defnoreturn",
				"id": "proc1",
				"extraState": {
					"params": [
						{"name": "arg1", "id": "id1"},
						{"name": "arg2", "id": "id2"}
					]
				},
				"fields": {"NAME": "myFunc"},
				"inputs": {
					"STACK": {
						"block": {
							"type": "text_print",
							"id": "print1",
							"inputs": {
								"TEXT": {
									"block": {"type": "text", "id": "text1", "fields": {"TEXT": "hello"}}
								}
							}
						}
					}
				}
			}]
		}
	}`

	result, err := Transpile([]byte(jsonData))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "(def myFunc (arg1 arg2)") {
		t.Errorf("expected function with params arg1 arg2, got: %s", result)
	}
}
