package blockly

import (
	"strings"
	"testing"
)

func TestTranspileComment(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
		want     string
	}{
		{
			name: "simple comment",
			jsonData: `{
				"blocks": {
					"languageVersion": 0,
					"blocks": [
						{
							"type": "comment_tatu",
							"id": "comment1",
							"fields": {"TEXT": "This is a comment"}
						}
					]
				}
			}`,
			want: "; This is a comment",
		},
		{
			name: "comment before variable",
			jsonData: `{
				"blocks": {
					"languageVersion": 0,
					"blocks": [
						{
							"type": "comment_tatu",
							"id": "comment1",
							"fields": {"TEXT": "Define counter variable"},
							"next": {
								"block": {
									"type": "variables_create_tatu",
									"id": "var1",
									"fields": {"VAR": {"id": "x"}},
									"inputs": {
										"VALUE": {
											"block": {
												"type": "math_number",
												"id": "num1",
												"fields": {"NUM": 0}
											}
										}
									}
								}
							}
						}
					]
				},
				"variables": [{"name": "counter", "id": "x"}]
			}`,
			want: "; Define counter variable\n(var counter 0)",
		},
		{
			name: "empty comment",
			jsonData: `{
				"blocks": {
					"languageVersion": 0,
					"blocks": [
						{
							"type": "comment_tatu",
							"id": "comment1",
							"fields": {"TEXT": ""}
						}
					]
				}
			}`,
			want: ";",
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

func TestTranspileCommentMissingField(t *testing.T) {
	jsonData := `{
		"blocks": {
			"languageVersion": 0,
			"blocks": [
				{
					"type": "comment_tatu",
					"id": "comment1",
					"fields": {}
				}
			]
		}
	}`

	_, err := Transpile([]byte(jsonData))
	if err == nil {
		t.Fatal("Expected error for missing TEXT field, got nil")
	}

	if !strings.Contains(err.Error(), "TEXT") {
		t.Errorf("Expected error about missing TEXT field, got: %v", err)
	}
}
