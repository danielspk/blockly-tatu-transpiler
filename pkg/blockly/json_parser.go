package blockly

import (
	"encoding/json"
	"fmt"
	"strings"
)

// JSONWorkspace represents the root JSON structure from Blockly.
type JSONWorkspace struct {
	Blocks struct {
		LanguageVersion int         `json:"languageVersion"`
		Blocks          []JSONBlock `json:"blocks"`
	} `json:"blocks"`
	Variables []JSONVariable `json:"variables,omitempty"`
}

// JSONVariable represents a variable definition in Blockly.
type JSONVariable struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// JSONBlock represents a block in Blockly JSON format.
type JSONBlock struct {
	Type       string               `json:"type"`
	ID         string               `json:"id"`
	X          float64              `json:"x,omitempty"`
	Y          float64              `json:"y,omitempty"`
	Fields     map[string]any       `json:"fields,omitempty"`
	Inputs     map[string]JSONInput `json:"inputs,omitempty"`
	Next       *JSONBlockWrapper    `json:"next,omitempty"`
	ExtraState *JSONExtraState      `json:"extraState,omitempty"`
}

// JSONInput represents an input connection.
type JSONInput struct {
	Block *JSONBlock `json:"block,omitempty"`
}

// JSONBlockWrapper wraps a block for next/parent references.
type JSONBlockWrapper struct {
	Block *JSONBlock `json:"block,omitempty"`
}

// JSONExtraState represents extra state data (mutations).
type JSONExtraState struct {
	Name   string   `json:"name,omitempty"`
	Params []string `json:"params,omitempty"`
}

// UnmarshalJSON workaround for handle `extraState` as string or object.
func (e *JSONExtraState) UnmarshalJSON(data []byte) error {
	var str string

	if err := json.Unmarshal(data, &str); err == nil {
		return nil
	}

	type Alias JSONExtraState

	return json.Unmarshal(data, (*Alias)(e))
}

// parseJSON parses Blockly JSON into a Workspace structure.
func parseJSON(jsonData []byte) (*Workspace, error) {
	if len(jsonData) == 0 {
		return nil, fmt.Errorf("empty JSON data provided")
	}

	var jsonWorkspace JSONWorkspace

	if err := json.Unmarshal(jsonData, &jsonWorkspace); err != nil {
		if syntaxErr, ok := err.(*json.SyntaxError); ok {
			return nil, fmt.Errorf("JSON syntax error at byte offset %d: %w", syntaxErr.Offset, err)
		}

		if typeErr, ok := err.(*json.UnmarshalTypeError); ok {
			return nil, fmt.Errorf("JSON type error: expected %s but got %s for field '%s': %w",
				typeErr.Type, typeErr.Value, typeErr.Field, err)
		}

		return nil, fmt.Errorf("failed to parse Blockly JSON workspace: %w (ensure JSON is from Blockly.serialization.workspaces.save())", err)
	}

	if jsonWorkspace.Blocks.Blocks == nil {
		return nil, fmt.Errorf("invalid Blockly JSON: missing 'blocks.blocks' array (ensure JSON is from Blockly.serialization.workspaces.save())")
	}

	varMap := make(map[string]string)
	for _, v := range jsonWorkspace.Variables {
		varMap[v.ID] = v.Name
	}

	workspace := &Workspace{
		Blocks: make([]*Block, 0, len(jsonWorkspace.Blocks.Blocks)),
	}

	for i := range jsonWorkspace.Blocks.Blocks {
		block, err := convertJSONBlock(&jsonWorkspace.Blocks.Blocks[i], varMap)
		if err != nil {
			return nil, fmt.Errorf("failed to convert block at index %d: %w", i, err)
		}

		workspace.Blocks = append(workspace.Blocks, block)
	}

	return workspace, nil
}

// convertJSONBlock converts a JSONBlock to a Block.
func convertJSONBlock(jsonBlock *JSONBlock, varMap map[string]string) (*Block, error) {
	if jsonBlock == nil {
		return nil, nil
	}

	if jsonBlock.Type == "" {
		return nil, fmt.Errorf("block with ID '%s' has empty type", jsonBlock.ID)
	}

	block := &Block{
		Type:       jsonBlock.Type,
		ID:         jsonBlock.ID,
		Fields:     make(map[string]*Field),
		Values:     make(map[string]*Value),
		Statements: make(map[string]*Statement),
	}

	for name, value := range jsonBlock.Fields {
		fieldValue := ""

		switch v := value.(type) {
		case string:
			fieldValue = v
		case float64:
			if v == float64(int64(v)) {
				fieldValue = fmt.Sprintf("%.0f", v)
			} else {
				fieldValue = fmt.Sprintf("%f", v)
				fieldValue = strings.TrimRight(strings.TrimRight(fieldValue, "0"), ".")
			}
		case int:
			fieldValue = fmt.Sprintf("%d", v)
		case bool:
			if v {
				fieldValue = "TRUE"
			} else {
				fieldValue = "FALSE"
			}
		case map[string]any:
			if varName, ok := v["name"].(string); ok {
				fieldValue = varName
			} else if id, ok := v["id"].(string); ok {
				if mappedName, found := varMap[id]; found {
					fieldValue = mappedName
				} else {
					fieldValue = id
				}
			}
		}

		block.Fields[name] = &Field{
			Value: fieldValue,
		}
	}

	for name, input := range jsonBlock.Inputs {

		if input.Block != nil {
			isStatement := name == "DO" || name == "STACK" ||
				name == "DO0" || name == "DO1" || name == "DO2"

			if name == "ELSE" && jsonBlock.Type != "logic_ternary" {
				isStatement = true
			}

			convertedBlock, err := convertJSONBlock(input.Block, varMap)
			if err != nil {
				return nil, fmt.Errorf("failed to convert input '%s' for block type '%s' (ID: %s): %w",
					name, jsonBlock.Type, jsonBlock.ID, err)
			}

			if isStatement {
				block.Statements[name] = &Statement{
					Block: convertedBlock,
				}
			} else {
				block.Values[name] = &Value{
					Block: convertedBlock,
				}
			}
		}
	}

	if jsonBlock.ExtraState != nil {
		mutation := &Mutation{
			Name: jsonBlock.ExtraState.Name,
			Args: jsonBlock.ExtraState.Params,
		}

		block.Mutation = mutation
	}

	if jsonBlock.Next != nil && jsonBlock.Next.Block != nil {
		nextBlock, err := convertJSONBlock(jsonBlock.Next.Block, varMap)
		if err != nil {
			return nil, fmt.Errorf("failed to convert next block for block type '%s' (ID: %s): %w",
				jsonBlock.Type, jsonBlock.ID, err)
		}

		block.Next = nextBlock
	}

	return block, nil
}
