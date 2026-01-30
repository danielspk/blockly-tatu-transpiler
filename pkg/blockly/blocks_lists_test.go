package blockly

import (
	"strings"
	"testing"
)

func TestTranspileListsCreate(t *testing.T) {
	jsonData := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": [{
				"type": "lists_create_with",
				"id": "list1",
				"inputs": {
					"ADD0": {"block": {"type": "math_number", "id": "num1", "fields": {"NUM": 1}}},
					"ADD1": {"block": {"type": "math_number", "id": "num2", "fields": {"NUM": 2}}},
					"ADD2": {"block": {"type": "math_number", "id": "num3", "fields": {"NUM": 3}}}
				}
			}]
		}
	}`

	result, err := Transpile([]byte(jsonData))
	if err != nil {
		t.Fatalf("Transpile failed: %v", err)
	}

	if strings.TrimSpace(result) != "(vector 1 2 3)" {
		t.Errorf("Expected (vector 1 2 3), got: %s", result)
	}
}

func TestTranspileListsLength(t *testing.T) {
	jsonData := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": [{
				"type": "lists_length",
				"id": "len1",
				"inputs": {
					"VALUE": {"block": {"type": "variables_get", "id": "var1", "fields": {"VAR": "myList"}}}
				}
			}]
		}
	}`

	result, err := Transpile([]byte(jsonData))
	if err != nil {
		t.Fatalf("Transpile failed: %v", err)
	}

	if !strings.Contains(result, "vec:len") {
		t.Errorf("Expected vec:len, got: %s", result)
	}
}

func TestTranspileListsSort(t *testing.T) {
	jsonData := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": [{
				"type": "lists_sort",
				"id": "sort1",
				"fields": {"TYPE": "NUMERIC", "DIRECTION": "ASCENDING"},
				"inputs": {
					"LIST": {"block": {"type": "variables_get", "id": "var1", "fields": {"VAR": "myList"}}}
				}
			}]
		}
	}`

	result, err := Transpile([]byte(jsonData))
	if err != nil {
		t.Fatalf("Transpile failed: %v", err)
	}

	if !strings.Contains(result, "vec:sort") {
		t.Errorf("Expected vec:sort, got: %s", result)
	}
}

func TestTranspileListsIsEmpty(t *testing.T) {
	jsonData := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": [{
				"type": "lists_isEmpty",
				"id": "empty1",
				"inputs": {
					"VALUE": {"block": {"type": "lists_create_with", "id": "list1", "extraState": {"itemCount": 0}}}
				}
			}]
		}
	}`

	result, err := Transpile([]byte(jsonData))
	if err != nil {
		t.Fatalf("Transpile failed: %v", err)
	}

	if !strings.Contains(result, "vec:len") && !strings.Contains(result, "= ") {
		t.Errorf("Expected isEmpty check, got: %s", result)
	}
}

func TestTranspileListsGetIndex(t *testing.T) {
	jsonData := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": [{
				"type": "lists_getIndex",
				"id": "get1",
				"fields": {"MODE": "GET", "WHERE": "FROM_START"},
				"inputs": {
					"VALUE": {"block": {"type": "lists_create_with", "id": "list1", "extraState": {"itemCount": 3}, "inputs": {
						"ADD0": {"block": {"type": "math_number", "id": "num1", "fields": {"NUM": 1}}},
						"ADD1": {"block": {"type": "math_number", "id": "num2", "fields": {"NUM": 2}}},
						"ADD2": {"block": {"type": "math_number", "id": "num3", "fields": {"NUM": 3}}}
					}}},
					"AT": {"block": {"type": "math_number", "id": "num4", "fields": {"NUM": 0}}}
				}
			}]
		}
	}`

	result, err := Transpile([]byte(jsonData))
	if err != nil {
		t.Fatalf("Transpile failed: %v", err)
	}

	if !strings.Contains(result, "vec:get") {
		t.Errorf("Expected vec:get, got: %s", result)
	}
}

func TestTranspileListsSetIndex(t *testing.T) {
	jsonData := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": [{
				"type": "lists_setIndex",
				"id": "set1",
				"fields": {"MODE": "SET", "WHERE": "FROM_START"},
				"inputs": {
					"LIST": {"block": {"type": "variables_get", "id": "var1", "fields": {"VAR": "myList"}}},
					"AT": {"block": {"type": "math_number", "id": "num1", "fields": {"NUM": 0}}},
					"TO": {"block": {"type": "math_number", "id": "num2", "fields": {"NUM": 100}}}
				}
			}]
		}
	}`

	result, err := Transpile([]byte(jsonData))
	if err != nil {
		t.Fatalf("Transpile failed: %v", err)
	}

	if !strings.Contains(result, "vec:set") {
		t.Errorf("Expected vec:set, got: %s", result)
	}
}

func TestTranspileListsIndexOf(t *testing.T) {
	jsonData := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": [{
				"type": "lists_indexOf",
				"id": "indexOf1",
				"fields": {"END": "FIRST"},
				"inputs": {
					"VALUE": {"block": {"type": "variables_get", "id": "var1", "fields": {"VAR": "myList"}}},
					"FIND": {"block": {"type": "math_number", "id": "num1", "fields": {"NUM": 5}}}
				}
			}]
		}
	}`

	result, err := Transpile([]byte(jsonData))
	if err != nil {
		t.Fatalf("Transpile failed: %v", err)
	}

	if !strings.Contains(result, "vec:find") {
		t.Errorf("Expected vec:find, got: %s", result)
	}
}

func TestTranspileListsGetSublist(t *testing.T) {
	jsonData := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": [{
				"type": "lists_getSublist",
				"id": "sublist1",
				"fields": {"WHERE1": "FROM_START", "WHERE2": "FROM_START"},
				"inputs": {
					"LIST": {"block": {"type": "variables_get", "id": "var1", "fields": {"VAR": "myList"}}},
					"AT1": {"block": {"type": "math_number", "id": "num1", "fields": {"NUM": 1}}},
					"AT2": {"block": {"type": "math_number", "id": "num2", "fields": {"NUM": 3}}}
				}
			}]
		}
	}`

	result, err := Transpile([]byte(jsonData))
	if err != nil {
		t.Fatalf("Transpile failed: %v", err)
	}

	if !strings.Contains(result, "vec:slice") {
		t.Errorf("Expected vec:slice, got: %s", result)
	}
}

func TestTranspileListsReverse(t *testing.T) {
	jsonData := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": [{
				"type": "lists_reverse",
				"id": "reverse1",
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

	result, err := Transpile([]byte(jsonData))
	if err != nil {
		t.Fatalf("Transpile failed: %v", err)
	}

	if !strings.Contains(result, "vec:reverse") {
		t.Errorf("Expected vec:reverse, got: %s", result)
	}
}

func TestTranspileListsSplit(t *testing.T) {
	jsonData := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": [{
				"type": "lists_split",
				"id": "split1",
				"fields": {"MODE": "SPLIT"},
				"inputs": {
					"INPUT": {"block": {"type": "text", "id": "text1", "fields": {"TEXT": "a,b,c"}}},
					"DELIM": {"block": {"type": "text", "id": "text2", "fields": {"TEXT": ","}}}
				}
			}]
		}
	}`

	result, err := Transpile([]byte(jsonData))
	if err != nil {
		t.Fatalf("Transpile failed: %v", err)
	}

	if !strings.Contains(result, "str:split") {
		t.Errorf("Expected str:split, got: %s", result)
	}
}

func TestTranspileListsRepeat(t *testing.T) {
	jsonData := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": [{
				"type": "lists_repeat",
				"id": "repeat1",
				"inputs": {
					"ITEM": {"block": {"type": "math_number", "id": "num1", "fields": {"NUM": 5}}},
					"NUM": {"block": {"type": "math_number", "id": "num2", "fields": {"NUM": 3}}}
				}
			}]
		}
	}`

	result, err := Transpile([]byte(jsonData))
	if err != nil {
		t.Fatalf("Transpile failed: %v", err)
	}

	if !strings.Contains(result, "vector") {
		t.Errorf("Expected vector creation, got: %s", result)
	}
}
