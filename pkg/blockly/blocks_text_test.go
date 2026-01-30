package blockly

import (
	"strings"
	"testing"
)

func TestTranspileText(t *testing.T) {
	jsonData := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": [{
				"type": "text",
				"id": "text1",
				"fields": {"TEXT": "hello world"}
			}]
		}
	}`

	result, err := Transpile([]byte(jsonData))
	if err != nil {
		t.Fatalf("Transpile failed: %v", err)
	}

	if strings.TrimSpace(result) != `"hello world"` {
		t.Errorf("Expected \"hello world\", got: %s", result)
	}
}

func TestTranspileTextPrint(t *testing.T) {
	jsonData := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": [{
				"type": "text_print",
				"id": "print1",
				"inputs": {
					"TEXT": {"block": {"type": "text", "id": "text1", "fields": {"TEXT": "Hello"}}}
				}
			}]
		}
	}`

	result, err := Transpile([]byte(jsonData))
	if err != nil {
		t.Fatalf("Transpile failed: %v", err)
	}

	if strings.TrimSpace(result) != `(print "Hello")` {
		t.Errorf("Expected (print \"Hello\"), got: %s", result)
	}
}

func TestTranspileTextLength(t *testing.T) {
	jsonData := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": [{
				"type": "text_length",
				"id": "len1",
				"inputs": {
					"VALUE": {"block": {"type": "text", "id": "text1", "fields": {"TEXT": "hello"}}}
				}
			}]
		}
	}`

	result, err := Transpile([]byte(jsonData))
	if err != nil {
		t.Fatalf("Transpile failed: %v", err)
	}

	if strings.TrimSpace(result) != `(str:len "hello")` {
		t.Errorf("Expected (str:len \"hello\"), got: %s", result)
	}
}

func TestTranspileTextAppend(t *testing.T) {
	jsonData := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": [{
				"type": "text_append",
				"id": "append1",
				"fields": {"VAR": "message"},
				"inputs": {
					"TEXT": {"block": {"type": "text", "id": "text1", "fields": {"TEXT": " world"}}}
				}
			}]
		}
	}`

	result, err := Transpile([]byte(jsonData))
	if err != nil {
		t.Fatalf("Transpile failed: %v", err)
	}

	if !strings.Contains(result, "str:concat") {
		t.Errorf("Expected str:concat, got: %s", result)
	}
}

func TestTranspileTextJoin(t *testing.T) {
	jsonData := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": [{
				"type": "text_join",
				"id": "join1",
				"extraState": {"itemCount": 3},
				"inputs": {
					"ADD0": {"block": {"type": "text", "id": "text1", "fields": {"TEXT": "Hello"}}},
					"ADD1": {"block": {"type": "text", "id": "text2", "fields": {"TEXT": " "}}},
					"ADD2": {"block": {"type": "text", "id": "text3", "fields": {"TEXT": "World"}}}
				}
			}]
		}
	}`

	result, err := Transpile([]byte(jsonData))
	if err != nil {
		t.Fatalf("Transpile failed: %v", err)
	}

	if !strings.Contains(result, "str:concat") || !strings.Contains(result, "Hello") {
		t.Errorf("Expected str:concat with text, got: %s", result)
	}
}

func TestTranspileTextCharAt(t *testing.T) {
	jsonData := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": [{
				"type": "text_charAt",
				"id": "charAt1",
				"fields": {"WHERE": "FROM_START"},
				"inputs": {
					"VALUE": {"block": {"type": "text", "id": "text1", "fields": {"TEXT": "hello"}}},
					"AT": {"block": {"type": "math_number", "id": "num1", "fields": {"NUM": 0}}}
				}
			}]
		}
	}`

	result, err := Transpile([]byte(jsonData))
	if err != nil {
		t.Fatalf("Transpile failed: %v", err)
	}

	if !strings.Contains(result, "str:slice") {
		t.Errorf("Expected str:slice for charAt, got: %s", result)
	}
}

func TestTranspileTextIsEmpty(t *testing.T) {
	jsonData := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": [{
				"type": "text_isEmpty",
				"id": "empty1",
				"inputs": {
					"VALUE": {"block": {"type": "text", "id": "text1", "fields": {"TEXT": "hello"}}}
				}
			}]
		}
	}`

	result, err := Transpile([]byte(jsonData))
	if err != nil {
		t.Fatalf("Transpile failed: %v", err)
	}

	if !strings.Contains(result, "str:len") || !strings.Contains(result, "= ") {
		t.Errorf("Expected isEmpty check with str:len, got: %s", result)
	}
}

func TestTranspileTextGetSubstring(t *testing.T) {
	jsonData := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": [{
				"type": "text_getSubstring",
				"id": "substr1",
				"fields": {"WHERE1": "FROM_START", "WHERE2": "FROM_START"},
				"inputs": {
					"STRING": {"block": {"type": "text", "id": "text1", "fields": {"TEXT": "hello world"}}},
					"AT1": {"block": {"type": "math_number", "id": "num1", "fields": {"NUM": 0}}},
					"AT2": {"block": {"type": "math_number", "id": "num2", "fields": {"NUM": 5}}}
				}
			}]
		}
	}`

	result, err := Transpile([]byte(jsonData))
	if err != nil {
		t.Fatalf("Transpile failed: %v", err)
	}

	if !strings.Contains(result, "str:slice") {
		t.Errorf("Expected str:slice for substring, got: %s", result)
	}
}

func TestTranspileTextIndexOf(t *testing.T) {
	jsonData := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": [{
				"type": "text_indexOf",
				"id": "indexOf1",
				"fields": {"END": "FIRST"},
				"inputs": {
					"VALUE": {"block": {"type": "text", "id": "text1", "fields": {"TEXT": "hello world"}}},
					"FIND": {"block": {"type": "text", "id": "text2", "fields": {"TEXT": "world"}}}
				}
			}]
		}
	}`

	result, err := Transpile([]byte(jsonData))
	if err != nil {
		t.Fatalf("Transpile failed: %v", err)
	}

	if !strings.Contains(result, "str:index") {
		t.Errorf("Expected str:index, got: %s", result)
	}
}

func TestTranspileTextTrim(t *testing.T) {
	jsonData := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": [{
				"type": "text_trim",
				"id": "trim1",
				"fields": {"MODE": "BOTH"},
				"inputs": {
					"TEXT": {"block": {"type": "text", "id": "text1", "fields": {"TEXT": "  hello  "}}}
				}
			}]
		}
	}`

	result, err := Transpile([]byte(jsonData))
	if err != nil {
		t.Fatalf("Transpile failed: %v", err)
	}

	if !strings.Contains(result, "str:trim") {
		t.Errorf("Expected str:trim, got: %s", result)
	}
}

func TestTranspileTextChangeCase(t *testing.T) {
	jsonData := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": [{
				"type": "text_changeCase",
				"id": "case1",
				"fields": {"CASE": "UPPERCASE"},
				"inputs": {
					"TEXT": {"block": {"type": "text", "id": "text1", "fields": {"TEXT": "hello"}}}
				}
			}]
		}
	}`

	result, err := Transpile([]byte(jsonData))
	if err != nil {
		t.Fatalf("Transpile failed: %v", err)
	}

	if !strings.Contains(result, "str:upper") {
		t.Errorf("Expected str:upper for uppercase, got: %s", result)
	}
}
