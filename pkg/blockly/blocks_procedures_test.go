package blockly

import (
	"strings"
	"testing"
)

func TestTranspileProcedureDefReturnWithParams(t *testing.T) {
	jsonData := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": [{
				"type": "procedures_defreturn",
				"id": "def1",
				"extraState": {"name": "double", "params": ["x"]},
				"inputs": {
					"RETURN": {"block": {"type": "math_arithmetic", "id": "math1", "fields": {"OP": "MULTIPLY"}, "inputs": {
						"A": {"block": {"type": "variables_get", "id": "var1", "fields": {"VAR": "x"}}},
						"B": {"block": {"type": "math_number", "id": "num1", "fields": {"NUM": 2}}}
					}}}
				}
			}]
		}
	}`

	result, err := Transpile([]byte(jsonData))
	if err != nil {
		t.Fatalf("Transpile failed: %v", err)
	}

	if !strings.Contains(result, "(def double") {
		t.Errorf("Expected function definition, got: %s", result)
	}

	if !strings.Contains(result, "(x)") {
		t.Errorf("Expected parameter x, got: %s", result)
	}

	if !strings.Contains(result, "* x 2") {
		t.Errorf("Expected multiplication, got: %s", result)
	}
}

func TestTranspileProcedureCallReturnWithArgs(t *testing.T) {
	jsonData := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": [{
				"type": "procedures_callreturn",
				"id": "call1",
				"extraState": {"name": "add", "params": ["a", "b"]},
				"inputs": {
					"ARG0": {"block": {"type": "math_number", "id": "num1", "fields": {"NUM": 5}}},
					"ARG1": {"block": {"type": "math_number", "id": "num2", "fields": {"NUM": 3}}}
				}
			}]
		}
	}`

	result, err := Transpile([]byte(jsonData))
	if err != nil {
		t.Fatalf("Transpile failed: %v", err)
	}

	if !strings.Contains(result, "(add 5 3)") {
		t.Errorf("Expected function call with args, got: %s", result)
	}
}

func TestTranspileProcedureIfReturn(t *testing.T) {
	jsonData := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": [{
				"type": "procedures_ifreturn",
				"id": "ifret1",
				"inputs": {
					"CONDITION": {"block": {"type": "logic_boolean", "id": "bool1", "fields": {"BOOL": "TRUE"}}},
					"VALUE": {"block": {"type": "math_number", "id": "num1", "fields": {"NUM": 42}}}
				}
			}]
		}
	}`

	result, err := Transpile([]byte(jsonData))
	if err != nil {
		t.Fatalf("Transpile failed: %v", err)
	}

	if !strings.Contains(result, "(if true 42)") {
		t.Errorf("Expected conditional return, got: %s", result)
	}
}

func TestTranspileProcedureCallNoReturn(t *testing.T) {
	jsonData := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": [{
				"type": "procedures_callnoreturn",
				"id": "call1",
				"extraState": {"name": "doSomething", "params": []}
			}]
		}
	}`

	result, err := Transpile([]byte(jsonData))
	if err != nil {
		t.Fatalf("Transpile failed: %v", err)
	}

	if !strings.Contains(result, "(doSomething)") {
		t.Errorf("Expected function call, got: %s", result)
	}
}

func TestTranspileProcedureDefNoReturnWithStatements(t *testing.T) {
	jsonData := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": [{
				"type": "procedures_defnoreturn",
				"id": "def1",
				"extraState": {"name": "greet", "params": ["name"]},
				"inputs": {
					"STACK": {"block": {
						"type": "text_print",
						"id": "print1",
						"inputs": {
							"TEXT": {"block": {"type": "variables_get", "id": "var1", "fields": {"VAR": "name"}}}
						}
					}}
				}
			}]
		}
	}`

	result, err := Transpile([]byte(jsonData))
	if err != nil {
		t.Fatalf("Transpile failed: %v", err)
	}

	if !strings.Contains(result, "(def greet") {
		t.Errorf("Expected function definition, got: %s", result)
	}

	if !strings.Contains(result, "(name)") {
		t.Errorf("Expected parameter, got: %s", result)
	}

	if !strings.Contains(result, "(print name)") {
		t.Errorf("Expected print statement, got: %s", result)
	}
}

// TestTranspileProcedureDefNoReturnWithMultipleStatements tests procedure with 2+ statements.
// This regression test ensures statements inside procedure bodies are not duplicated.
func TestTranspileProcedureDefNoReturnWithMultipleStatements(t *testing.T) {
	jsonData := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": [{
				"type": "procedures_defnoreturn",
				"id": "def1",
				"extraState": {"name": "multiPrint", "params": []},
				"inputs": {
					"STACK": {"block": {
						"type": "text_print",
						"id": "print1",
						"inputs": {"TEXT": {"block": {"type": "text", "id": "text1", "fields": {"TEXT": "first"}}}},
						"next": {
							"block": {
								"type": "text_print",
								"id": "print2",
								"inputs": {"TEXT": {"block": {"type": "text", "id": "text2", "fields": {"TEXT": "second"}}}},
								"next": {
									"block": {
										"type": "text_print",
										"id": "print3",
										"inputs": {"TEXT": {"block": {"type": "text", "id": "text3", "fields": {"TEXT": "third"}}}}
									}
								}
							}
						}
					}}
				}
			}]
		}
	}`

	result, err := Transpile([]byte(jsonData))
	if err != nil {
		t.Fatalf("Transpile failed: %v", err)
	}

	if !strings.Contains(result, "(block") {
		t.Errorf("Expected block for multiple statements, got: %s", result)
	}

	firstCount := strings.Count(result, `"first"`)
	secondCount := strings.Count(result, `"second"`)
	thirdCount := strings.Count(result, `"third"`)

	if firstCount != 1 {
		t.Errorf("Expected 'first' to appear exactly once, but appeared %d times in: %s", firstCount, result)
	}
	if secondCount != 1 {
		t.Errorf("Expected 'second' to appear exactly once, but appeared %d times in: %s", secondCount, result)
	}
	if thirdCount != 1 {
		t.Errorf("Expected 'third' to appear exactly once, but appeared %d times in: %s", thirdCount, result)
	}
}

// TestTranspileProcedureDefReturnWithMultipleStatements tests procedure with return and 2+ statements.
func TestTranspileProcedureDefReturnWithMultipleStatements(t *testing.T) {
	jsonData := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": [{
				"type": "procedures_defreturn",
				"id": "def1",
				"extraState": {"name": "calculate", "params": ["x"]},
				"inputs": {
					"STACK": {"block": {
						"type": "text_print",
						"id": "print1",
						"inputs": {"TEXT": {"block": {"type": "text", "id": "text1", "fields": {"TEXT": "calculating"}}}},
						"next": {
							"block": {
								"type": "text_print",
								"id": "print2",
								"inputs": {"TEXT": {"block": {"type": "text", "id": "text2", "fields": {"TEXT": "done"}}}}
							}
						}
					}},
					"RETURN": {"block": {"type": "variables_get", "id": "var1", "fields": {"VAR": "x"}}}
				}
			}]
		}
	}`

	result, err := Transpile([]byte(jsonData))
	if err != nil {
		t.Fatalf("Transpile failed: %v", err)
	}

	if !strings.Contains(result, "(block") {
		t.Errorf("Expected block for multiple statements, got: %s", result)
	}

	calculatingCount := strings.Count(result, `"calculating"`)
	doneCount := strings.Count(result, `"done"`)

	if calculatingCount != 1 {
		t.Errorf("Expected 'calculating' to appear exactly once, but appeared %d times in: %s", calculatingCount, result)
	}
	if doneCount != 1 {
		t.Errorf("Expected 'done' to appear exactly once, but appeared %d times in: %s", doneCount, result)
	}
}

func TestTranspileFunctionCallTatu(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
		want     string
	}{
		{
			name: "function call with no arguments",
			jsonData: `{
				"blocks": {
					"languageVersion": 0,
					"blocks": [
						{
							"type": "function_call_tatu",
							"id": "fn1",
							"fields": {"NAME": "myFunction"},
							"extraState": {"items": 0}
						}
					]
				}
			}`,
			want: "(myFunction)",
		},
		{
			name: "function call with one argument",
			jsonData: `{
				"blocks": {
					"languageVersion": 0,
					"blocks": [
						{
							"type": "function_call_tatu",
							"id": "fn1",
							"fields": {"NAME": "increment"},
							"extraState": {"items": 1},
							"inputs": {
								"ARG0": {
									"block": {
										"type": "math_number",
										"id": "num1",
										"fields": {"NUM": 5}
									}
								}
							}
						}
					]
				}
			}`,
			want: "(increment 5)",
		},
		{
			name: "function call with multiple arguments",
			jsonData: `{
				"blocks": {
					"languageVersion": 0,
					"blocks": [
						{
							"type": "function_call_tatu",
							"id": "fn1",
							"fields": {"NAME": "sum"},
							"extraState": {"items": 3},
							"inputs": {
								"ARG0": {
									"block": {
										"type": "math_number",
										"id": "num1",
										"fields": {"NUM": 10}
									}
								},
								"ARG1": {
									"block": {
										"type": "math_number",
										"id": "num2",
										"fields": {"NUM": 20}
									}
								},
								"ARG2": {
									"block": {
										"type": "math_number",
										"id": "num3",
										"fields": {"NUM": 30}
									}
								}
							}
						}
					]
				}
			}`,
			want: "(sum 10 20 30)",
		},
		{
			name: "function call with variable argument",
			jsonData: `{
				"blocks": {
					"languageVersion": 0,
					"blocks": [
						{
							"type": "function_call_tatu",
							"id": "fn1",
							"fields": {"NAME": "process"},
							"extraState": {"items": 1},
							"inputs": {
								"ARG0": {
									"block": {
										"type": "variables_get",
										"id": "var1",
										"fields": {"VAR": {"id": "x"}}
									}
								}
							}
						}
					]
				},
				"variables": [{"name": "data", "id": "x"}]
			}`,
			want: "(process data)",
		},
		{
			name: "nested function call",
			jsonData: `{
				"blocks": {
					"languageVersion": 0,
					"blocks": [
						{
							"type": "function_call_tatu",
							"id": "fn1",
							"fields": {"NAME": "outer"},
							"extraState": {"items": 1},
							"inputs": {
								"ARG0": {
									"block": {
										"type": "function_call_tatu",
										"id": "fn2",
										"fields": {"NAME": "inner"},
										"extraState": {"items": 1},
										"inputs": {
											"ARG0": {
												"block": {
													"type": "math_number",
													"id": "num1",
													"fields": {"NUM": 42}
												}
											}
										}
									}
								}
							}
						}
					]
				}
			}`,
			want: "(outer (inner 42))",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Transpile([]byte(tt.jsonData))
			if err != nil {
				t.Fatalf("Transpile failed: %v", err)
			}

			result = strings.TrimSpace(result)
			if !strings.Contains(result, tt.want) {
				t.Errorf("Expected result to contain:\n%s\n\nGot:\n%s", tt.want, result)
			}
		})
	}
}

func TestTranspileFunctionCallTatuErrors(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
		wantErr  string
	}{
		{
			name: "missing NAME field",
			jsonData: `{
				"blocks": {
					"languageVersion": 0,
					"blocks": [
						{
							"type": "function_call_tatu",
							"id": "fn1",
							"fields": {},
							"extraState": {"items": 0}
						}
					]
				}
			}`,
			wantErr: "NAME",
		},
		{
			name: "empty function name",
			jsonData: `{
				"blocks": {
					"languageVersion": 0,
					"blocks": [
						{
							"type": "function_call_tatu",
							"id": "fn1",
							"fields": {"NAME": ""},
							"extraState": {"items": 0}
						}
					]
				}
			}`,
			wantErr: "NAME",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Transpile([]byte(tt.jsonData))
			if err == nil {
				t.Fatal("Expected error but got nil")
			}

			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("Expected error containing '%s', got: %v", tt.wantErr, err)
			}
		})
	}
}
