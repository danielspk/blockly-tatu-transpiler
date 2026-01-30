package blockly

import "fmt"

// transpileVariableGet transpiles a variable getter block.
// Blockly: variables_get with VAR field
// Tatu: variable-name
func transpileVariableGet(block *Block) (string, error) {
	varName, err := validateField(block, "VAR")
	if err != nil {
		return "", err
	}

	return varName, nil
}

// transpileVariableSet transpiles a variable setter block.
// Blockly: variables_set with VAR field and VALUE
// Tatu: (set variable-name value)
func transpileVariableSet(block *Block) (string, error) {
	varName, err := validateField(block, "VAR")
	if err != nil {
		return "", err
	}

	value, err := transpileUnaryOp(block, "VALUE")
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("(set %s %s)", varName, value), nil
}

// transpileVariableCreate transpiles a variable creation block.
// Blockly: variables_create_tatu with VAR field and VALUE
// Tatu: (var variable-name initial-value)
func transpileVariableCreate(block *Block) (string, error) {
	varName, err := validateField(block, "VAR")
	if err != nil {
		return "", err
	}

	value, err := transpileUnaryOp(block, "VALUE")
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("(var %s %s)", varName, value), nil
}
