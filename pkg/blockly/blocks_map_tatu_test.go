package blockly

import (
	"strings"
	"testing"
)

func TestTranspileMapCreate(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
		expected string
		wantErr  bool
	}{
		{
			name: "empty map",
			jsonData: `{
				"blocks": {
					"languageVersion": 0,
					"blocks": [{
						"type": "map_create_with_tatu",
						"id": "map1",
						"extraState": {"itemCount": 0},
						"inputs": {}
					}]
				}
			}`,
			expected: `(map)`,
			wantErr:  false,
		},
		{
			name: "map with single pair",
			jsonData: `{
				"blocks": {
					"languageVersion": 0,
					"blocks": [{
						"type": "map_create_with_tatu",
						"id": "map1",
						"extraState": {"itemCount": 1},
						"inputs": {
							"PAIR0": {
								"block": {
									"type": "map_pair_tatu",
									"id": "pair1",
									"inputs": {
										"KEY": {"block": {"type": "text", "id": "key1", "fields": {"TEXT": "name"}}},
										"VALUE": {"block": {"type": "text", "id": "val1", "fields": {"TEXT": "John"}}}
									}
								}
							}
						}
					}]
				}
			}`,
			expected: `(map "name" "John")`,
			wantErr:  false,
		},
		{
			name: "map with multiple pairs",
			jsonData: `{
				"blocks": {
					"languageVersion": 0,
					"blocks": [{
						"type": "map_create_with_tatu",
						"id": "map1",
						"extraState": {"itemCount": 3},
						"inputs": {
							"PAIR0": {
								"block": {
									"type": "map_pair_tatu",
									"id": "pair1",
									"inputs": {
										"KEY": {"block": {"type": "text", "id": "key1", "fields": {"TEXT": "name"}}},
										"VALUE": {"block": {"type": "text", "id": "val1", "fields": {"TEXT": "John"}}}
									}
								}
							},
							"PAIR1": {
								"block": {
									"type": "map_pair_tatu",
									"id": "pair2",
									"inputs": {
										"KEY": {"block": {"type": "text", "id": "key2", "fields": {"TEXT": "age"}}},
										"VALUE": {"block": {"type": "math_number", "id": "val2", "fields": {"NUM": 30}}}
									}
								}
							},
							"PAIR2": {
								"block": {
									"type": "map_pair_tatu",
									"id": "pair3",
									"inputs": {
										"KEY": {"block": {"type": "text", "id": "key3", "fields": {"TEXT": "active"}}},
										"VALUE": {"block": {"type": "logic_boolean", "id": "val3", "fields": {"BOOL": "TRUE"}}}
									}
								}
							}
						}
					}]
				}
			}`,
			expected: `(map "name" "John" "age" 30 "active" true)`,
			wantErr:  false,
		},
		{
			name: "map with numeric keys",
			jsonData: `{
				"blocks": {
					"languageVersion": 0,
					"blocks": [{
						"type": "map_create_with_tatu",
						"id": "map1",
						"extraState": {"itemCount": 2},
						"inputs": {
							"PAIR0": {
								"block": {
									"type": "map_pair_tatu",
									"id": "pair1",
									"inputs": {
										"KEY": {"block": {"type": "math_number", "id": "key1", "fields": {"NUM": 1}}},
										"VALUE": {"block": {"type": "text", "id": "val1", "fields": {"TEXT": "first"}}}
									}
								}
							},
							"PAIR1": {
								"block": {
									"type": "map_pair_tatu",
									"id": "pair2",
									"inputs": {
										"KEY": {"block": {"type": "math_number", "id": "key2", "fields": {"NUM": 2}}},
										"VALUE": {"block": {"type": "text", "id": "val2", "fields": {"TEXT": "second"}}}
									}
								}
							}
						}
					}]
				}
			}`,
			expected: `(map 1 "first" 2 "second")`,
			wantErr:  false,
		},
		{
			name: "map with nested expressions as values",
			jsonData: `{
				"blocks": {
					"languageVersion": 0,
					"blocks": [{
						"type": "map_create_with_tatu",
						"id": "map1",
						"extraState": {"itemCount": 1},
						"inputs": {
							"PAIR0": {
								"block": {
									"type": "map_pair_tatu",
									"id": "pair1",
									"inputs": {
										"KEY": {"block": {"type": "text", "id": "key1", "fields": {"TEXT": "sum"}}},
										"VALUE": {
											"block": {
												"type": "math_arithmetic",
												"id": "val1",
												"fields": {"OP": "ADD"},
												"inputs": {
													"A": {"block": {"type": "math_number", "id": "num1", "fields": {"NUM": 5}}},
													"B": {"block": {"type": "math_number", "id": "num2", "fields": {"NUM": 3}}}
												}
											}
										}
									}
								}
							}
						}
					}]
				}
			}`,
			expected: `(map "sum" (+ 5 3))`,
			wantErr:  false,
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

func TestTranspileMapCreateSequence(t *testing.T) {
	jsonData := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": [
				{
					"type": "variables_set",
					"id": "set1",
					"fields": {"VAR": "user"},
					"inputs": {
						"VALUE": {
							"block": {
								"type": "map_create_with_tatu",
								"id": "map1",
								"extraState": {"itemCount": 2},
								"inputs": {
									"PAIR0": {
										"block": {
											"type": "map_pair_tatu",
											"id": "pair1",
											"inputs": {
												"KEY": {"block": {"type": "text", "id": "key1", "fields": {"TEXT": "name"}}},
												"VALUE": {"block": {"type": "text", "id": "val1", "fields": {"TEXT": "Alice"}}}
											}
										}
									},
									"PAIR1": {
										"block": {
											"type": "map_pair_tatu",
											"id": "pair2",
											"inputs": {
												"KEY": {"block": {"type": "text", "id": "key2", "fields": {"TEXT": "age"}}},
												"VALUE": {"block": {"type": "math_number", "id": "val2", "fields": {"NUM": 25}}}
											}
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

	if !strings.Contains(result, `(set user (map "name" "Alice" "age" 25))`) {
		t.Errorf("Expected map assignment, got: %s", result)
	}
}
