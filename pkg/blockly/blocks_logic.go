package blockly

import "fmt"

// transpileBoolean transpiles a boolean block.
// Blockly: logic_boolean with field BOOL (TRUE/FALSE)
// Tatu: true / false
func transpileBoolean(block *Block) (string, error) {
	value, err := validateField(block, "BOOL")
	if err != nil {
		return "", err
	}

	if value == "TRUE" {
		return "true", nil
	}

	return "false", nil
}

// transpileCompare transpiles a comparison block.
// Blockly: logic_compare with field OP (EQ, NEQ, LT, LTE, GT, GTE)
// Tatu: (= a b), (< a b), (<= a b), (> a b), (>= a b), (not (= a b))
func transpileCompare(block *Block) (string, error) {
	op, err := validateField(block, "OP")
	if err != nil {
		return "", err
	}

	a, b, err := transpileBinaryOp(block, "A", "B")
	if err != nil {
		return "", err
	}

	var operator string
	switch op {
	case "EQ":
		operator = "="
	case "NEQ":
		return fmt.Sprintf("(not (= %s %s))", a, b), nil
	case "LT":
		operator = "<"
	case "LTE":
		operator = "<="
	case "GT":
		operator = ">"
	case "GTE":
		operator = ">="
	default:
		return "", fmt.Errorf("block '%s' (ID: %s) has unknown comparison operator: %s", block.Type, block.ID, op)
	}

	return fmt.Sprintf("(%s %s %s)", operator, a, b), nil
}

// transpileOperation transpiles a logical operation block.
// Blockly: logic_operation with field OP (AND/OR)
// Tatu: (and a b), (or a b)
func transpileOperation(block *Block) (string, error) {
	opValue, err := validateField(block, "OP")
	if err != nil {
		return "", err
	}

	a, b, err := transpileBinaryOp(block, "A", "B")
	if err != nil {
		return "", err
	}

	var op string
	switch opValue {
	case "AND":
		op = "and"
	case "OR":
		op = "or"
	default:
		return "", fmt.Errorf("block '%s' (ID: %s) has unknown logical operator: %s", block.Type, block.ID, opValue)
	}

	return fmt.Sprintf("(%s %s %s)", op, a, b), nil
}

// transpileNegate transpiles a NOT block.
// Blockly: logic_negate
// Tatu: (not x)
func transpileNegate(block *Block) (string, error) {
	value, err := transpileUnaryOp(block, "BOOL")
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("(not %s)", value), nil
}

// transpileTernary transpiles a ternary operator block.
// Blockly: logic_ternary (condition ? then : else)
// Tatu: (if condition then else)
func transpileTernary(block *Block) (string, error) {
	ifValue := block.Values["IF"]
	thenValue := block.Values["THEN"]
	elseValue := block.Values["ELSE"]

	if ifValue == nil || thenValue == nil || elseValue == nil {
		return "", fmt.Errorf("logic_ternary missing IF, THEN, or ELSE value")
	}

	cond, err := transpileBlockInternal(ifValue.Block, true)
	if err != nil {
		return "", err
	}

	thenCode, err := transpileBlockInternal(thenValue.Block, true)
	if err != nil {
		return "", err
	}

	elseCode, err := transpileBlockInternal(elseValue.Block, true)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("(if %s %s %s)", cond, thenCode, elseCode), nil
}
