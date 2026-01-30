package blockly

import (
	"fmt"
	"strings"
)

const indentSpaces = 2

// validateField checks if a required field exists and returns its value.
func validateField(block *Block, fieldName string) (string, error) {
	if block.Fields == nil {
		return "", fmt.Errorf("block '%s' (ID: %s) has no fields defined, expected field '%s'",
			block.Type, block.ID, fieldName)
	}

	field, ok := block.Fields[fieldName]
	if !ok {
		return "", fmt.Errorf("block '%s' (ID: %s) missing required field '%s'",
			block.Type, block.ID, fieldName)
	}

	if field.Value == "" {
		return "", fmt.Errorf("block '%s' (ID: %s) field '%s' has empty value",
			block.Type, block.ID, fieldName)
	}

	return field.Value, nil
}

// validateFieldAllowEmpty checks if a required field exists and returns its value, allowing empty strings.
func validateFieldAllowEmpty(block *Block, fieldName string) (string, error) {
	if block.Fields == nil {
		return "", fmt.Errorf("block '%s' (ID: %s) has no fields defined, expected field '%s'",
			block.Type, block.ID, fieldName)
	}

	field, ok := block.Fields[fieldName]
	if !ok {
		return "", fmt.Errorf("block '%s' (ID: %s) missing required field '%s'",
			block.Type, block.ID, fieldName)
	}

	return field.Value, nil
}

// validateValue checks if a required value input exists and returns the block.
func validateValue(block *Block, valueName string) (*Block, error) {
	if block.Values == nil {
		return nil, fmt.Errorf("block '%s' (ID: %s) has no value inputs defined, expected input '%s'",
			block.Type, block.ID, valueName)
	}

	value, ok := block.Values[valueName]
	if !ok {
		return nil, fmt.Errorf("block '%s' (ID: %s) missing required value input '%s'",
			block.Type, block.ID, valueName)
	}

	if value.Block == nil {
		return nil, fmt.Errorf("block '%s' (ID: %s) value input '%s' has no connected block",
			block.Type, block.ID, valueName)
	}

	return value.Block, nil
}

// validateStatement checks if a required statement input exists and returns the block.
func validateStatement(block *Block, statementName string) (*Block, error) {
	if block.Statements == nil {
		return nil, fmt.Errorf("block '%s' (ID: %s) has no statement inputs defined, expected statement '%s'",
			block.Type, block.ID, statementName)
	}

	statement, ok := block.Statements[statementName]
	if !ok {
		return nil, fmt.Errorf("block '%s' (ID: %s) missing required statement input '%s'",
			block.Type, block.ID, statementName)
	}

	if statement.Block == nil {
		return nil, fmt.Errorf("block '%s' (ID: %s) statement input '%s' has no connected block",
			block.Type, block.ID, statementName)
	}

	return statement.Block, nil
}

// transpileUnaryOp is a helper for transpiling blocks with a single input.
func transpileUnaryOp(block *Block, inputName string) (string, error) {
	valueBlock, err := validateValue(block, inputName)
	if err != nil {
		return "", err
	}

	return transpileBlockInternal(valueBlock, true)
}

// transpileBinaryOp is a helper for transpiling blocks with two inputs.
func transpileBinaryOp(block *Block, inputA, inputB string) (string, string, error) {
	aBlock, err := validateValue(block, inputA)
	if err != nil {
		return "", "", err
	}

	bBlock, err := validateValue(block, inputB)
	if err != nil {
		return "", "", err
	}

	a, err := transpileBlockInternal(aBlock, true)
	if err != nil {
		return "", "", err
	}

	b, err := transpileBlockInternal(bBlock, true)
	if err != nil {
		return "", "", err
	}

	return a, b, nil
}

// indent adds indentation to each line of the input string.
func indent(s string) string {
	if s == "" {
		return s
	}

	indentStr := strings.Repeat(" ", indentSpaces)
	hasTrailingNewline := strings.HasSuffix(s, "\n")

	if hasTrailingNewline {
		s = s[:len(s)-1]
	}

	result := indentStr + strings.ReplaceAll(s, "\n", "\n"+indentStr)

	if hasTrailingNewline {
		result += "\n"
	}

	return result
}
