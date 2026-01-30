// Package blockly provides a transpiler from Blockly XML to Tatu source code.
package blockly

import (
	"fmt"
	"strings"
)

// Transpile converts Blockly workspace (JSON format) to Tatu source code.
// Note: function definitions will be added at startup.
func Transpile(jsonData []byte) (string, error) {
	workspace, err := parseJSON(jsonData)
	if err != nil {
		return "", err
	}

	var result strings.Builder
	var functionDefs []string
	var mainBlocks []string

	for _, block := range workspace.Blocks {
		code, err := transpileBlockInternal(block, true)
		if err != nil {
			return "", err
		}

		if code == "" {
			continue
		}

		if block.Type == "procedures_defnoreturn" || block.Type == "procedures_defreturn" {
			functionDefs = append(functionDefs, code)
		} else {
			mainBlocks = append(mainBlocks, code)
		}
	}

	for i, def := range functionDefs {
		result.WriteString(def)
		result.WriteString("\n")

		if i < len(functionDefs)-1 {
			result.WriteString("\n")
		}
	}

	if len(functionDefs) > 0 && len(mainBlocks) > 0 {
		result.WriteString("\n")
	}

	for _, mainCode := range mainBlocks {
		result.WriteString(mainCode)
		result.WriteString("\n")
	}

	return result.String(), nil
}

// transpileBlockInternal routes a block to its specific transpiler.
// If processNext is true, it will recursively process the Next block chain.
func transpileBlockInternal(block *Block, processNext bool) (string, error) {
	if block == nil {
		return "", nil
	}

	if block.Disabled {
		if processNext && block.Next != nil {
			return transpileBlockInternal(block.Next, true)
		}

		return "", nil
	}

	var result string
	var err error

	switch {
	// logic blocks
	case block.Type == "controls_if":
		result, err = transpileIf(block)
	case block.Type == "logic_boolean":
		result, err = transpileBoolean(block)
	case block.Type == "logic_null":
		result = "nil"
	case block.Type == "logic_compare":
		result, err = transpileCompare(block)
	case block.Type == "logic_operation":
		result, err = transpileOperation(block)
	case block.Type == "logic_negate":
		result, err = transpileNegate(block)
	case block.Type == "logic_ternary":
		result, err = transpileTernary(block)

	// loop blocks
	case block.Type == "controls_repeat_ext":
		result, err = transpileRepeat(block)
	case block.Type == "controls_whileUntil":
		result, err = transpileWhileUntil(block)
	case block.Type == "controls_for":
		result, err = transpileFor(block)
	case block.Type == "controls_forEach":
		result, err = transpileForEach(block)
	case block.Type == "controls_flow_statements":
		result, err = transpileFlowStatement(block)
	case block.Type == "expression_tatu":
		result, err = transpileExpression(block)

	// math blocks
	case block.Type == "math_number":
		result, err = transpileMathNumber(block)
	case block.Type == "math_arithmetic":
		result, err = transpileMathArithmetic(block)
	case block.Type == "math_single":
		result, err = transpileMathSingle(block)
	case block.Type == "math_trig":
		result, err = transpileMathTrig(block)
	case block.Type == "math_constant":
		result, err = transpileMathConstant(block)
	case block.Type == "math_round":
		result, err = transpileMathRound(block)
	case block.Type == "math_modulo":
		result, err = transpileMathModulo(block)
	case block.Type == "math_random_int":
		result, err = transpileMathRandomInt(block)
	case block.Type == "math_constrain":
		result, err = transpileMathConstrain(block)
	case block.Type == "math_random_float":
		result, err = transpileMathRandomFloat(block)
	case block.Type == "math_on_list":
		result, err = transpileMathOnList(block)
	case block.Type == "math_number_property":
		result, err = transpileMathNumberProperty(block)
	case block.Type == "math_change":
		result, err = transpileMathChange(block)

	// variables
	case block.Type == "variables_get":
		result, err = transpileVariableGet(block)
	case block.Type == "variables_set":
		result, err = transpileVariableSet(block)
	case block.Type == "variables_create_tatu":
		result, err = transpileVariableCreate(block)

	// text blocks
	case block.Type == "text":
		result, err = transpileText(block)
	case block.Type == "text_print":
		result, err = transpileTextPrint(block)
	case block.Type == "text_join":
		result, err = transpileTextJoin(block)
	case block.Type == "text_length":
		result, err = transpileTextLength(block)
	case block.Type == "text_isEmpty":
		result, err = transpileTextIsEmpty(block)
	case block.Type == "text_charAt":
		result, err = transpileTextCharAt(block)
	case block.Type == "text_getSubstring":
		result, err = transpileTextGetSubstring(block)
	case block.Type == "text_append":
		result, err = transpileTextAppend(block)
	case block.Type == "text_indexOf":
		result, err = transpileTextIndexOf(block)
	case block.Type == "text_changeCase":
		result, err = transpileTextChangeCase(block)
	case block.Type == "text_trim":
		result, err = transpileTextTrim(block)

	// list blocks
	case block.Type == "lists_create_with":
		result, err = transpileListsCreate(block)
	case block.Type == "lists_length":
		result, err = transpileListsLength(block)
	case block.Type == "lists_isEmpty":
		result, err = transpileListsIsEmpty(block)
	case block.Type == "lists_getIndex":
		result, err = transpileListsGetIndex(block)
	case block.Type == "lists_setIndex":
		result, err = transpileListsSetIndex(block)
	case block.Type == "lists_indexOf":
		result, err = transpileListsIndexOf(block)
	case block.Type == "lists_getSublist":
		result, err = transpileListsGetSublist(block)
	case block.Type == "lists_split":
		result, err = transpileListsSplit(block)
	case block.Type == "lists_sort":
		result, err = transpileListsSort(block)
	case block.Type == "lists_repeat":
		result, err = transpileListsRepeat(block)
	case block.Type == "lists_reverse":
		result, err = transpileListsReverse(block)

	// procedure blocks
	case block.Type == "procedures_defnoreturn":
		result, err = transpileProcedureDefNoReturn(block)
	case block.Type == "procedures_defreturn":
		result, err = transpileProcedureDefReturn(block)
	case block.Type == "procedures_callnoreturn":
		result, err = transpileProcedureCallNoReturn(block)
	case block.Type == "procedures_callreturn":
		result, err = transpileProcedureCallReturn(block)
	case block.Type == "procedures_ifreturn":
		result, err = transpileProcedureIfReturn(block)

	// map blocks - Tatu
	case block.Type == "map_create_with_tatu":
		result, err = transpileMapCreate(block)

	// comment block - Tatu
	case block.Type == "comment_tatu":
		result, err = transpileComment(block)

	// include block - Tatu
	case block.Type == "include_file_tatu":
		result, err = transpileInclude(block)

	// function call - Tatu
	case strings.HasSuffix(block.Type, "_call_tatu"):
		result, err = transpileFunctionCall(block)

	default:
		return "", fmt.Errorf("unsupported block type '%s' (ID: %s)", block.Type, block.ID)
	}

	if err != nil {
		return "", fmt.Errorf("transpilation error for block type '%s' (ID: %s): %w", block.Type, block.ID, err)
	}

	if block.Comment != "" {
		commentPrefix := ""

		for _, line := range strings.Split(block.Comment, "\n") {
			commentPrefix += "; " + line + "\n"
		}

		result = commentPrefix + result
	}

	if processNext && block.Next != nil {
		nextCode, err := transpileBlockInternal(block.Next, true)
		if err != nil {
			return "", err
		}

		if nextCode != "" {
			result = result + "\n" + nextCode
		}
	}

	return result, nil
}

// transpileBlockSequence transpiles a sequence of blocks.
func transpileBlockSequence(block *Block) ([]string, error) {
	var statements []string

	current := block

	for current != nil {
		code, err := transpileBlockInternal(current, false)
		if err != nil {
			return nil, err
		}
		if code != "" {
			statements = append(statements, code)
		}

		current = current.Next
	}

	return statements, nil
}

// transpileWhereClause transpiles a WHERE clause to an index expression.
// Used for blocks that have WHERE field (FROM_START, FROM_END, FIRST, LAST, RANDOM).
// Returns the transpiled index expression for use in Tatu code.
//
// Parameters:
//   - block: The block containing the WHERE field and optional AT value
//   - whereValue: The WHERE field value (FROM_START, FROM_END, FIRST, LAST, RANDOM)
//   - atFieldName: The name of the field containing the AT value (e.g., "AT", "AT1", "AT2")
//   - collectionExpr: The transpiled expression for the collection (e.g., list variable)
//   - lenFunc: The length function to use (e.g., "vec:len", "str:len")
func transpileWhereClause(block *Block, whereValue, atFieldName, collectionExpr, lenFunc string) (string, error) {
	switch whereValue {
	case "FROM_START":
		atValue, err := validateValue(block, atFieldName)
		if err != nil {
			return "", err
		}
		return transpileBlockInternal(atValue, true)

	case "FROM_END":
		atValue, err := validateValue(block, atFieldName)
		if err != nil {
			return "", err
		}
		at, err := transpileBlockInternal(atValue, true)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("(- (%s %s) %s)", lenFunc, collectionExpr, at), nil

	case "FIRST":
		return "0", nil

	case "LAST":
		return fmt.Sprintf("(- (%s %s) 1)", lenFunc, collectionExpr), nil

	case "RANDOM":
		return "", fmt.Errorf("RANDOM index not yet supported")

	default:
		return "", fmt.Errorf("unknown WHERE value: %s", whereValue)
	}
}
