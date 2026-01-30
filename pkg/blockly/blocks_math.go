package blockly

import "fmt"

// transpileMathNumber transpiles a number block.
// Blockly: math_number with NUM field
// Tatu: 42, 3.14, etc.
func transpileMathNumber(block *Block) (string, error) {
	return validateField(block, "NUM")
}

// transpileMathArithmetic transpiles an arithmetic operation block.
// Blockly: math_arithmetic with OP field (ADD, MINUS, MULTIPLY, DIVIDE, POWER)
// Tatu: (+ a b), (- a b), (* a b), (/ a b), (math:pow a b)
func transpileMathArithmetic(block *Block) (string, error) {
	opField, err := validateField(block, "OP")
	if err != nil {
		return "", err
	}

	a, b, err := transpileBinaryOp(block, "A", "B")
	if err != nil {
		return "", err
	}

	var op string
	switch opField {
	case "ADD":
		op = "+"
	case "MINUS":
		op = "-"
	case "MULTIPLY":
		op = "*"
	case "DIVIDE":
		op = "/"
	case "POWER":
		return fmt.Sprintf("(math:pow %s %s)", a, b), nil
	default:
		return "", fmt.Errorf("unknown arithmetic operator: %s", opField)
	}

	return fmt.Sprintf("(%s %s %s)", op, a, b), nil
}

// transpileMathSingle transpiles a single-operand math function.
// Blockly: math_single with OP field (ROOT, ABS, NEG, LN, LOG10, EXP, POW10)
// Tatu: (math:sqrt x), (math:abs x), (- x), (math:log x), (math:exp x), etc.
func transpileMathSingle(block *Block) (string, error) {
	opField, err := validateField(block, "OP")
	if err != nil {
		return "", err
	}

	num, err := transpileUnaryOp(block, "NUM")
	if err != nil {
		return "", err
	}

	switch opField {
	case "ROOT":
		return fmt.Sprintf("(math:sqrt %s)", num), nil
	case "ABS":
		return fmt.Sprintf("(math:abs %s)", num), nil
	case "NEG":
		return fmt.Sprintf("(- %s)", num), nil
	case "LN":
		return fmt.Sprintf("(math:log %s)", num), nil
	case "LOG10":
		return fmt.Sprintf("(/ (math:log %s) (math:log 10))", num), nil
	case "EXP":
		return fmt.Sprintf("(math:exp %s)", num), nil
	case "POW10":
		return fmt.Sprintf("(math:pow 10 %s)", num), nil
	default:
		return "", fmt.Errorf("unknown math_single operator: %s", opField)
	}
}

// transpileMathTrig transpiles a trigonometric function.
// Blockly: math_trig with OP field (SIN, COS, TAN, ASIN, ACOS, ATAN)
// Tatu: (math:sin x), (math:cos x), (math:tan x), etc.
func transpileMathTrig(block *Block) (string, error) {
	opField, err := validateField(block, "OP")
	if err != nil {
		return "", err
	}

	// Validate operation is supported before processing values
	var fn string
	switch opField {
	case "SIN":
		fn = "math:sin"
	case "COS":
		fn = "math:cos"
	case "TAN":
		fn = "math:tan"
	case "ASIN":
		return "", fmt.Errorf("asin not yet supported in Tatu stdlib")
	case "ACOS":
		return "", fmt.Errorf("acos not yet supported in Tatu stdlib")
	case "ATAN":
		return "", fmt.Errorf("atan not yet supported in Tatu stdlib")
	default:
		return "", fmt.Errorf("unknown trig function: %s", opField)
	}

	num, err := transpileUnaryOp(block, "NUM")
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("(%s %s)", fn, num), nil
}

// transpileMathConstant transpiles a math constant.
// Blockly: math_constant with CONSTANT field (PI, E, GOLDEN_RATIO, SQRT2, SQRT1_2, INFINITY)
// Tatu: (math:pi), (math:e), calculated values
func transpileMathConstant(block *Block) (string, error) {
	constField, err := validateField(block, "CONSTANT")
	if err != nil {
		return "", err
	}

	switch constField {
	case "PI":
		return "(math:pi)", nil
	case "E":
		return "(math:e)", nil
	case "GOLDEN_RATIO":
		return "1.61803398875", nil
	case "SQRT2":
		return "(math:sqrt 2)", nil
	case "SQRT1_2":
		return "(/ 1 (math:sqrt 2))", nil
	case "INFINITY":
		return "9999999999", nil // Tatu doesn't have infinity, use large number
	default:
		return "", fmt.Errorf("unknown math constant: %s", constField)
	}
}

// transpileMathRound transpiles a rounding function.
// Blockly: math_round with OP field (ROUND, ROUNDUP, ROUNDDOWN)
// Tatu: (math:round x), (math:ceil x), (math:floor x)
func transpileMathRound(block *Block) (string, error) {
	opField, err := validateField(block, "OP")
	if err != nil {
		return "", err
	}

	num, err := transpileUnaryOp(block, "NUM")
	if err != nil {
		return "", err
	}

	var fn string
	switch opField {
	case "ROUND":
		fn = "math:round"
	case "ROUNDUP":
		fn = "math:ceil"
	case "ROUNDDOWN":
		fn = "math:floor"
	default:
		return "", fmt.Errorf("unknown rounding function: %s", opField)
	}

	return fmt.Sprintf("(%s %s)", fn, num), nil
}

// transpileMathModulo transpiles a modulo operation.
// Blockly: math_modulo
// Tatu: (% dividend divisor)
func transpileMathModulo(block *Block) (string, error) {
	dividend, divisor, err := transpileBinaryOp(block, "DIVIDEND", "DIVISOR")
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("(%% %s %s)", dividend, divisor), nil
}

// transpileMathRandomInt transpiles a random integer block.
// Blockly: math_random_int with FROM and TO values
// Tatu: (math:rand from to)
func transpileMathRandomInt(block *Block) (string, error) {
	from, to, err := transpileBinaryOp(block, "FROM", "TO")
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("(math:rand %s %s)", from, to), nil
}

// transpileMathConstrain transpiles a constrain/clamp block.
// Blockly: math_constrain with VALUE, LOW, HIGH
// Tatu: (math:min (math:max value low) high)
func transpileMathConstrain(block *Block) (string, error) {
	value, err := transpileUnaryOp(block, "VALUE")
	if err != nil {
		return "", err
	}

	low, err := transpileUnaryOp(block, "LOW")
	if err != nil {
		return "", err
	}

	high, err := transpileUnaryOp(block, "HIGH")
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("(math:min (math:max %s %s) %s)", value, low, high), nil
}

// transpileMathRandomFloat transpiles a random float block.
func transpileMathRandomFloat(block *Block) (string, error) {
	// generate a random float between 0 and 1
	// since Tatu's math:rand returns integers, we divide by a large number
	return "(/ (math:rand 0 1000000) 1000000.0)", nil
}

// transpileMathOnList transpiles aggregate operations on lists.
// Blockly: math_on_list with OP field (SUM, MIN, MAX, AVERAGE, MEDIAN, MODE, STD_DEV, RANDOM)
// Tatu: Various vec: functions
func transpileMathOnList(block *Block) (string, error) {
	opField, err := validateField(block, "OP")
	if err != nil {
		return "", err
	}

	switch opField {
	case "SUM":
		return "", fmt.Errorf("math_on_list SUM operation not supported in Tatu stdlib")
	case "MIN":
		return "", fmt.Errorf("math_on_list MIN operation not supported in Tatu stdlib")
	case "MAX":
		return "", fmt.Errorf("math_on_list MAX operation not supported in Tatu stdlib")
	case "AVERAGE":
		return "", fmt.Errorf("math_on_list AVERAGE operation not supported in Tatu stdlib")
	case "MEDIAN":
		return "", fmt.Errorf("math_on_list MEDIAN operation not supported in Tatu stdlib")
	case "MODE":
		return "", fmt.Errorf("math_on_list MODE operation not supported in Tatu stdlib")
	case "STD_DEV":
		return "", fmt.Errorf("math_on_list STD_DEV operation not supported in Tatu stdlib")
	case "RANDOM":
		return "", fmt.Errorf("math_on_list RANDOM operation not supported in Tatu stdlib")
	default:
		return "", fmt.Errorf("unknown math_on_list operation: %s", opField)
	}
}

// transpileMathNumberProperty transpiles number property checks.
// Blockly: math_number_property with PROPERTY field and DIVISOR (for divisible by)
// Tatu: Various comparison expressions
func transpileMathNumberProperty(block *Block) (string, error) {
	propertyField, err := validateField(block, "PROPERTY")
	if err != nil {
		return "", err
	}

	num, err := transpileUnaryOp(block, "NUMBER_TO_CHECK")
	if err != nil {
		return "", err
	}

	switch propertyField {
	case "EVEN":
		return fmt.Sprintf("(= (%% %s 2) 0)", num), nil
	case "ODD":
		return fmt.Sprintf("(not (= (%% %s 2) 0))", num), nil
	case "PRIME":
		return fmt.Sprintf("(math:is-prime %s)", num), nil
	case "WHOLE":
		return fmt.Sprintf("(= %s (math:floor %s))", num, num), nil
	case "POSITIVE":
		return fmt.Sprintf("(> %s 0)", num), nil
	case "NEGATIVE":
		return fmt.Sprintf("(< %s 0)", num), nil
	case "DIVISIBLE_BY":
		divisor, err := transpileUnaryOp(block, "DIVISOR")
		if err != nil {
			return "", err
		}

		// (= (% n divisor) 0)
		return fmt.Sprintf("(= (%% %s %s) 0)", num, divisor), nil
	default:
		return "", fmt.Errorf("unknown number property: %s", propertyField)
	}
}

// transpileMathChange transpiles a math change block.
// Blockly: math_change with VAR field and DELTA value
// Tatu: (set var (+ var delta))
func transpileMathChange(block *Block) (string, error) {
	varName, err := validateField(block, "VAR")
	if err != nil {
		return "", err
	}

	delta, err := transpileUnaryOp(block, "DELTA")
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("(set %s (+ %s %s))", varName, varName, delta), nil
}
