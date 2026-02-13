package blockly

import (
	"strings"
	"testing"
)

// TestTranspileBreak tests break statement in while loop.
func TestTranspileBreak(t *testing.T) {
	jsonData := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": [
				{
					"type": "controls_whileUntil",
					"id": "while1",
					"fields": {
						"MODE": "WHILE"
					},
					"inputs": {
						"BOOL": {
							"block": {
								"type": "logic_boolean",
								"id": "bool1",
								"fields": {
									"BOOL": "TRUE"
								}
							}
						},
						"DO": {
							"block": {
								"type": "controls_flow_statements",
								"id": "break1",
								"fields": {
									"FLOW": "BREAK"
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

	if !strings.Contains(result, "var __break") {
		t.Errorf("Expected '__break' variable declaration, got: %s", result)
	}

	if !strings.Contains(result, "set __break true") {
		t.Errorf("Expected 'set __break true', got: %s", result)
	}

	if !strings.Contains(result, "(and true (not __break))") {
		t.Errorf("Expected condition to check break flag, got: %s", result)
	}
}

// TestTranspileContinue tests continue statement in while loop.
func TestTranspileContinue(t *testing.T) {
	jsonData := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": [
				{
					"type": "controls_whileUntil",
					"id": "while1",
					"fields": {
						"MODE": "WHILE"
					},
					"inputs": {
						"BOOL": {
							"block": {
								"type": "logic_boolean",
								"id": "bool1",
								"fields": {
									"BOOL": "TRUE"
								}
							}
						},
						"DO": {
							"block": {
								"type": "controls_flow_statements",
								"id": "continue1",
								"fields": {
									"FLOW": "CONTINUE"
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

	if !strings.Contains(result, "var __continue") {
		t.Errorf("Expected '__continue' variable declaration, got: %s", result)
	}

	if !strings.Contains(result, "set __continue true") {
		t.Errorf("Expected 'set __continue true', got: %s", result)
	}
}

// TestTranspileBreakInFor tests break statement in for loop.
func TestTranspileBreakInFor(t *testing.T) {
	jsonData := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": [
				{
					"type": "controls_for",
					"id": "for1",
					"fields": {
						"VAR": "i"
					},
					"inputs": {
						"FROM": {
							"block": {
								"type": "math_number",
								"id": "num1",
								"fields": {"NUM": 1}
							}
						},
						"TO": {
							"block": {
								"type": "math_number",
								"id": "num2",
								"fields": {"NUM": 10}
							}
						},
						"BY": {
							"block": {
								"type": "math_number",
								"id": "num3",
								"fields": {"NUM": 1}
							}
						},
						"DO": {
							"block": {
								"type": "controls_if",
								"id": "if1",
								"inputs": {
									"IF0": {
										"block": {
											"type": "logic_compare",
											"id": "cmp1",
											"fields": {"OP": "GT"},
											"inputs": {
												"A": {
													"block": {
														"type": "variables_get",
														"id": "var1",
														"fields": {"VAR": "i"}
													}
												},
												"B": {
													"block": {
														"type": "math_number",
														"id": "num4",
														"fields": {"NUM": 5}
													}
												}
											}
										}
									},
									"DO0": {
										"block": {
											"type": "controls_flow_statements",
											"id": "break1",
											"fields": {"FLOW": "BREAK"}
										}
									}
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

	if !strings.Contains(result, "var __break") {
		t.Errorf("Expected '__break' variable, got: %s", result)
	}

	if !strings.Contains(result, "(and (<= i 10) (not __break))") {
		t.Errorf("Expected for condition to check break flag, got: %s", result)
	}
}

// TestTranspileIfWithElseIf tests if/elseif/else structure.
func TestTranspileIfWithElseIf(t *testing.T) {
	jsonData := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": [{
				"type": "controls_if",
				"id": "if1",
				"inputs": {
					"IF0": {"block": {"type": "logic_boolean", "id": "bool1", "fields": {"BOOL": "TRUE"}}},
					"DO0": {"block": {"type": "text_print", "id": "print1", "inputs": {"TEXT": {"block": {"type": "text", "id": "text1", "fields": {"TEXT": "first"}}}}}},
					"IF1": {"block": {"type": "logic_boolean", "id": "bool2", "fields": {"BOOL": "FALSE"}}},
					"DO1": {"block": {"type": "text_print", "id": "print2", "inputs": {"TEXT": {"block": {"type": "text", "id": "text2", "fields": {"TEXT": "second"}}}}}},
					"ELSE": {"block": {"type": "text_print", "id": "print3", "inputs": {"TEXT": {"block": {"type": "text", "id": "text3", "fields": {"TEXT": "else"}}}}}}
				}
			}]
		}
	}`

	result, err := Transpile([]byte(jsonData))
	if err != nil {
		t.Fatalf("Transpile failed: %v", err)
	}

	if !strings.Contains(result, "(if true") {
		t.Errorf("Expected first if condition, got: %s", result)
	}
	if !strings.Contains(result, "(if false") {
		t.Errorf("Expected second if condition (elseif), got: %s", result)
	}
	if !strings.Contains(result, `"first"`) {
		t.Errorf("Expected first branch, got: %s", result)
	}
	if !strings.Contains(result, `"second"`) {
		t.Errorf("Expected second branch (elseif), got: %s", result)
	}
	if !strings.Contains(result, `"else"`) {
		t.Errorf("Expected else branch, got: %s", result)
	}
}

// TestTranspileIfWithMultipleStatements tests if block with 2+ statements in then branch.
// This regression test ensures statements inside control blocks are not duplicated.
func TestTranspileIfWithMultipleStatements(t *testing.T) {
	jsonData := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": [{
				"type": "controls_if",
				"id": "if1",
				"inputs": {
					"IF0": {"block": {"type": "logic_boolean", "id": "bool1", "fields": {"BOOL": "TRUE"}}},
					"DO0": {"block": {
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

	if !strings.Contains(result, "(begin") {
		t.Errorf("Expected begin for multiple statements, got: %s", result)
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

	lines := strings.Split(result, "\n")
	var printLines []string
	for _, line := range lines {
		if strings.Contains(line, "(print") {
			printLines = append(printLines, line)
		}
	}

	if len(printLines) != 3 {
		t.Errorf("Expected 3 print statements, found %d in: %s", len(printLines), result)
	}

	if len(printLines) >= 2 {
		firstIndent := len(printLines[0]) - len(strings.TrimLeft(printLines[0], " \t"))
		for i, line := range printLines[1:] {
			indent := len(line) - len(strings.TrimLeft(line, " \t"))
			if indent != firstIndent {
				t.Errorf("Print statement %d has different indentation (%d) than first (%d) in: %s", i+2, indent, firstIndent, result)
			}
		}
	}
}

// TestTranspileIfWithMultipleStatementsInElse tests if block with 2+ statements in else branch.
func TestTranspileIfWithMultipleStatementsInElse(t *testing.T) {
	jsonData := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": [{
				"type": "controls_if",
				"id": "if1",
				"inputs": {
					"IF0": {"block": {"type": "logic_boolean", "id": "bool1", "fields": {"BOOL": "FALSE"}}},
					"DO0": {"block": {"type": "text_print", "id": "print1", "inputs": {"TEXT": {"block": {"type": "text", "id": "text1", "fields": {"TEXT": "then"}}}}}},
					"ELSE": {"block": {
						"type": "text_print",
						"id": "print2",
						"inputs": {"TEXT": {"block": {"type": "text", "id": "text2", "fields": {"TEXT": "else1"}}}},
						"next": {
							"block": {
								"type": "text_print",
								"id": "print3",
								"inputs": {"TEXT": {"block": {"type": "text", "id": "text3", "fields": {"TEXT": "else2"}}}}
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

	else1Count := strings.Count(result, `"else1"`)
	else2Count := strings.Count(result, `"else2"`)

	if else1Count != 1 {
		t.Errorf("Expected 'else1' to appear exactly once, but appeared %d times in: %s", else1Count, result)
	}
	if else2Count != 1 {
		t.Errorf("Expected 'else2' to appear exactly once, but appeared %d times in: %s", else2Count, result)
	}
}

// TestTranspileWhileWithMultipleStatements tests while loop with 2+ statements in body.
func TestTranspileWhileWithMultipleStatements(t *testing.T) {
	jsonData := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": [{
				"type": "controls_whileUntil",
				"id": "while1",
				"fields": {"MODE": "WHILE"},
				"inputs": {
					"BOOL": {"block": {"type": "logic_boolean", "id": "bool1", "fields": {"BOOL": "TRUE"}}},
					"DO": {"block": {
						"type": "text_print",
						"id": "print1",
						"inputs": {"TEXT": {"block": {"type": "text", "id": "text1", "fields": {"TEXT": "loop1"}}}},
						"next": {
							"block": {
								"type": "text_print",
								"id": "print2",
								"inputs": {"TEXT": {"block": {"type": "text", "id": "text2", "fields": {"TEXT": "loop2"}}}}
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

	if !strings.Contains(result, "(begin") {
		t.Errorf("Expected begin for multiple statements, got: %s", result)
	}

	loop1Count := strings.Count(result, `"loop1"`)
	loop2Count := strings.Count(result, `"loop2"`)

	if loop1Count != 1 {
		t.Errorf("Expected 'loop1' to appear exactly once, but appeared %d times in: %s", loop1Count, result)
	}
	if loop2Count != 1 {
		t.Errorf("Expected 'loop2' to appear exactly once, but appeared %d times in: %s", loop2Count, result)
	}
}

// TestTranspileForWithMultipleStatements tests for loop with 2+ statements in body.
func TestTranspileForWithMultipleStatements(t *testing.T) {
	jsonData := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": [{
				"type": "controls_for",
				"id": "for1",
				"fields": {"VAR": "i"},
				"inputs": {
					"FROM": {"block": {"type": "math_number", "id": "num1", "fields": {"NUM": 1}}},
					"TO": {"block": {"type": "math_number", "id": "num2", "fields": {"NUM": 5}}},
					"BY": {"block": {"type": "math_number", "id": "num3", "fields": {"NUM": 1}}},
					"DO": {"block": {
						"type": "text_print",
						"id": "print1",
						"inputs": {"TEXT": {"block": {"type": "variables_get", "id": "var1", "fields": {"VAR": "i"}}}},
						"next": {
							"block": {
								"type": "text_print",
								"id": "print2",
								"inputs": {"TEXT": {"block": {"type": "text", "id": "text1", "fields": {"TEXT": "iteration"}}}}
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

	if !strings.Contains(result, "(begin") {
		t.Errorf("Expected begin for multiple statements, got: %s", result)
	}

	iterationCount := strings.Count(result, `"iteration"`)
	if iterationCount != 1 {
		t.Errorf("Expected 'iteration' to appear exactly once, but appeared %d times in: %s", iterationCount, result)
	}
}

// TestTranspileRepeatWithMultipleStatements tests repeat loop with 2+ statements in body.
func TestTranspileRepeatWithMultipleStatements(t *testing.T) {
	jsonData := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": [{
				"type": "controls_repeat_ext",
				"id": "repeat1",
				"inputs": {
					"TIMES": {"block": {"type": "math_number", "id": "num1", "fields": {"NUM": 3}}},
					"DO": {"block": {
						"type": "text_print",
						"id": "print1",
						"inputs": {"TEXT": {"block": {"type": "text", "id": "text1", "fields": {"TEXT": "repeat1"}}}},
						"next": {
							"block": {
								"type": "text_print",
								"id": "print2",
								"inputs": {"TEXT": {"block": {"type": "text", "id": "text2", "fields": {"TEXT": "repeat2"}}}}
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

	if !strings.Contains(result, "(begin") {
		t.Errorf("Expected begin for multiple statements, got: %s", result)
	}

	repeat1Count := strings.Count(result, `"repeat1"`)
	repeat2Count := strings.Count(result, `"repeat2"`)

	if repeat1Count != 1 {
		t.Errorf("Expected 'repeat1' to appear exactly once, but appeared %d times in: %s", repeat1Count, result)
	}
	if repeat2Count != 1 {
		t.Errorf("Expected 'repeat2' to appear exactly once, but appeared %d times in: %s", repeat2Count, result)
	}
}

// TestTranspileForEachWithMultipleStatements tests forEach loop with 2+ statements in body.
func TestTranspileForEachWithMultipleStatements(t *testing.T) {
	jsonData := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": [{
				"type": "controls_forEach",
				"id": "foreach1",
				"fields": {"VAR": "item"},
				"inputs": {
					"LIST": {"block": {"type": "lists_create_with", "id": "list1", "extraState": {"itemCount": 2}, "inputs": {
						"ADD0": {"block": {"type": "math_number", "id": "num1", "fields": {"NUM": 1}}},
						"ADD1": {"block": {"type": "math_number", "id": "num2", "fields": {"NUM": 2}}}
					}}},
					"DO": {"block": {
						"type": "text_print",
						"id": "print1",
						"inputs": {"TEXT": {"block": {"type": "variables_get", "id": "var1", "fields": {"VAR": "item"}}}},
						"next": {
							"block": {
								"type": "text_print",
								"id": "print2",
								"inputs": {"TEXT": {"block": {"type": "text", "id": "text1", "fields": {"TEXT": "processed"}}}}
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

	if !strings.Contains(result, "(begin") {
		t.Errorf("Expected begin for multiple statements, got: %s", result)
	}

	processedCount := strings.Count(result, `"processed"`)
	if processedCount != 1 {
		t.Errorf("Expected 'processed' to appear exactly once, but appeared %d times in: %s", processedCount, result)
	}
}

// TestTranspileIfWithComments tests that comments don't cause unnecessary begin wrappers.
func TestTranspileIfWithComments(t *testing.T) {
	jsonData := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": [{
				"type": "controls_if",
				"id": "if1",
				"extraState": {"hasElse": true},
				"inputs": {
					"IF0": {"block": {"type": "logic_boolean", "id": "bool1", "fields": {"BOOL": "TRUE"}}},
					"DO0": {"block": {
						"type": "comment_tatu",
						"id": "cmt1",
						"fields": {"TEXT": "then branch"},
						"next": {"block": {
							"type": "variables_set",
							"id": "set1",
							"fields": {"VAR": {"id": "var1"}},
							"inputs": {"VALUE": {"block": {"type": "logic_boolean", "id": "bool2", "fields": {"BOOL": "TRUE"}}}}
						}}
					}},
					"ELSE": {"block": {
						"type": "comment_tatu",
						"id": "cmt2",
						"fields": {"TEXT": "else branch"},
						"next": {"block": {
							"type": "text_print",
							"id": "print1",
							"inputs": {"TEXT": {"block": {"type": "text", "id": "text1", "fields": {"TEXT": "ok"}}}}
						}}
					}}
				}
			}]
		},
		"variables": [{"name": "x", "id": "var1"}]
	}`

	result, err := Transpile([]byte(jsonData))
	if err != nil {
		t.Fatalf("Transpile failed: %v", err)
	}

	if strings.Contains(result, "(begin") {
		t.Errorf("Expected no begin wrapper for comment + single statement, got: %s", result)
	}

	if !strings.Contains(result, "; then branch") {
		t.Errorf("Expected comment 'then branch', got: %s", result)
	}

	if !strings.Contains(result, "; else branch") {
		t.Errorf("Expected comment 'else branch', got: %s", result)
	}

	if !strings.Contains(result, "(if true") {
		t.Errorf("Expected if statement, got: %s", result)
	}
}

// TestTranspileIfWithoutElse tests if block without else branch.
func TestTranspileIfWithoutElse(t *testing.T) {
	jsonData := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": [{
				"type": "controls_if",
				"id": "if1",
				"inputs": {
					"IF0": {"block": {"type": "logic_boolean", "id": "bool1", "fields": {"BOOL": "TRUE"}}},
					"DO0": {"block": {"type": "text_print", "id": "print1", "inputs": {"TEXT": {"block": {"type": "text", "id": "text1", "fields": {"TEXT": "then"}}}}}}
				}
			}]
		}
	}`

	result, err := Transpile([]byte(jsonData))
	if err != nil {
		t.Fatalf("Transpile failed: %v", err)
	}

	if !strings.Contains(result, "(if true") {
		t.Errorf("Expected if statement, got: %s", result)
	}
	if !strings.Contains(result, `"then"`) {
		t.Errorf("Expected then branch, got: %s", result)
	}

	if strings.Contains(result, "nil") {
		t.Errorf("Expected no nil (else should be optional), got: %s", result)
	}

	openParens := strings.Count(result, "(")
	closeParens := strings.Count(result, ")")
	if openParens != closeParens {
		t.Errorf("Unbalanced parentheses: %d open vs %d close in: %s", openParens, closeParens, result)
	}
}

// TestTranspileExpression tests expression_tatu block.
func TestTranspileExpression(t *testing.T) {
	jsonData := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": [
				{
					"type": "variables_set",
					"id": "set1",
					"fields": {
						"VAR": "x"
					},
					"inputs": {
						"VALUE": {
							"block": {
								"type": "math_number",
								"id": "num1",
								"fields": {
									"NUM": "42"
								}
							}
						}
					},
					"next": {
						"block": {
							"type": "expression_tatu",
							"id": "expr1",
							"inputs": {
								"VALUE": {
									"block": {
										"type": "variables_get",
										"id": "var1",
										"fields": {
											"VAR": "x"
										}
									}
								}
							}
						}
					}
				}
			],
			"variables": [
				{"name": "x", "id": "var_x"}
			]
		}
	}`

	result, err := Transpile([]byte(jsonData))
	if err != nil {
		t.Fatalf("Transpile failed: %v", err)
	}

	if !strings.Contains(result, "(set x 42)") {
		t.Errorf("Expected set statement, got: %s", result)
	}

	lines := strings.Split(strings.TrimSpace(result), "\n")
	if len(lines) < 2 {
		t.Errorf("Expected at least 2 lines (set + expression), got: %s", result)
	}

	lastLine := strings.TrimSpace(lines[len(lines)-1])
	if lastLine != "x" {
		t.Errorf("Expected last line to be 'x', got: %s", lastLine)
	}
}
