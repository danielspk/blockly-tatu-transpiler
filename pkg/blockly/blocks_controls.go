package blockly

import (
	"fmt"
	"strings"
)

// transpileIf transpiles an if/else block with support for multiple elseif branches.
// Blockly: controls_if with IF0, DO0, IF1, DO1, ..., and optional ELSE
// Tatu: Nested (if condition then-block else-block) structures
func transpileIf(block *Block) (string, error) {
	var conditions []string
	var thenBranches []string
	var thenCode string
	var elseCode string

	ifIndex := 0

	for {
		ifKey := fmt.Sprintf("IF%d", ifIndex)
		doKey := fmt.Sprintf("DO%d", ifIndex)

		ifValue := block.Values[ifKey]
		doStatement := block.Statements[doKey]

		if ifValue == nil || doStatement == nil {
			break
		}

		cond, err := transpileBlockInternal(ifValue.Block, true)
		if err != nil {
			return "", fmt.Errorf("error transpiling %s: %w", ifKey, err)
		}

		conditions = append(conditions, cond)

		thenStatements, err := transpileBlockSequence(doStatement.Block)
		if err != nil {
			return "", fmt.Errorf("error transpiling %s: %w", doKey, err)
		}

		thenCode = getBodyCode(thenStatements)
		thenBranches = append(thenBranches, thenCode)
		ifIndex++
	}

	if len(conditions) == 0 {
		return "", fmt.Errorf("controls_if missing IF0 or DO0")
	}

	elseStatement := block.Statements["ELSE"]

	if elseStatement != nil {
		elseStatements, err := transpileBlockSequence(elseStatement.Block)
		if err != nil {
			return "", err
		}

		elseCode = getBodyCode(elseStatements)
	} else {
		elseCode = "nil"
	}

	result := elseCode

	for i := len(conditions) - 1; i >= 0; i-- {
		result = fmt.Sprintf("(if %s\n%s\n%s)", conditions[i], indent(thenBranches[i]), indent(result))
	}

	return result, nil
}

// transpileRepeat transpiles a repeat block.
// Blockly: controls_repeat_ext with TIMES value
// Tatu: (var i 0) (while (< i n) (begin ... (set i (+ i 1))))
func transpileRepeat(block *Block) (string, error) {
	var bodyCode string

	timesValue := block.Values["TIMES"]
	doStatement := block.Statements["DO"]

	if timesValue == nil || doStatement == nil {
		return "", fmt.Errorf("controls_repeat_ext missing TIMES or DO")
	}

	times, err := transpileBlockInternal(timesValue.Block, true)
	if err != nil {
		return "", err
	}

	doStatements, err := transpileBlockSequence(doStatement.Block)
	if err != nil {
		return "", err
	}

	bodyCode = getBodyCode(doStatements)
	loopVar := "__i"

	whileBody := fmt.Sprintf("(begin\n%s\n%s)",
		indent(bodyCode),
		indent(fmt.Sprintf("(set %s (+ %s 1))", loopVar, loopVar)))

	whileLoop := fmt.Sprintf("(while (< %s %s)\n%s)",
		loopVar, times, indent(whileBody))

	return fmt.Sprintf("(begin\n%s\n%s)",
		indent(fmt.Sprintf("(var %s 0)", loopVar)),
		indent(whileLoop)), nil
}

// transpileWhileUntil transpiles a while/until block.
// Blockly: controls_whileUntil with MODE field (WHILE/UNTIL)
// Tatu: (while cond body) or (while (not cond) body)
func transpileWhileUntil(block *Block) (string, error) {
	var bodyCode string

	mode, err := validateField(block, "MODE")
	if err != nil {
		return "", err
	}

	boolBlock, err := validateValue(block, "BOOL")
	if err != nil {
		return "", err
	}

	doBlock, err := validateStatement(block, "DO")
	if err != nil {
		return "", err
	}

	cond, err := transpileBlockInternal(boolBlock, true)
	if err != nil {
		return "", err
	}

	hasBreak := hasFlowControl(doBlock, "BREAK")
	hasContinue := hasFlowControl(doBlock, "CONTINUE")
	doStatements, err := transpileBlockSequence(doBlock)
	if err != nil {
		return "", err
	}

	bodyCode = getBodyCode(doStatements)
	bodyCode = wrapWithFlowControl(bodyCode, hasBreak, hasContinue)

	if mode == "UNTIL" {
		cond = fmt.Sprintf("(if %s false true)", cond)
	}

	cond = modifyLoopCondition(cond, hasBreak)

	loopCode := fmt.Sprintf("(while %s\n%s)", cond, indent(bodyCode))

	return addFlowControlVars(loopCode, hasBreak, hasContinue), nil
}

// transpileFor transpiles a for loop block.
// Blockly: controls_for with VAR field, FROM, TO, BY values
// Tatu: (for (var i from) (<= i to) body (set i (+ i by)))
func transpileFor(block *Block) (string, error) {
	var bodyCode string

	varName, err := validateField(block, "VAR")
	if err != nil {
		return "", err
	}

	fromValue, err := validateValue(block, "FROM")
	if err != nil {
		return "", err
	}

	toValue, err := validateValue(block, "TO")
	if err != nil {
		return "", err
	}

	byValue, err := validateValue(block, "BY")
	if err != nil {
		return "", err
	}

	doStatement, err := validateStatement(block, "DO")
	if err != nil {
		return "", err
	}

	from, err := transpileBlockInternal(fromValue, true)
	if err != nil {
		return "", err
	}

	to, err := transpileBlockInternal(toValue, true)
	if err != nil {
		return "", err
	}

	by, err := transpileBlockInternal(byValue, true)
	if err != nil {
		return "", err
	}

	hasBreak := hasFlowControl(doStatement, "BREAK")
	hasContinue := hasFlowControl(doStatement, "CONTINUE")
	doStatements, err := transpileBlockSequence(doStatement)
	if err != nil {
		return "", err
	}

	bodyCode = getBodyCode(doStatements)
	bodyCode = wrapWithFlowControl(bodyCode, hasBreak, hasContinue)

	cond := fmt.Sprintf("(<= %s %s)", varName, to)
	cond = modifyLoopCondition(cond, hasBreak)

	loopCode := fmt.Sprintf("(for (var %s %s)\n%s\n%s\n%s)",
		varName, from, indent(cond), indent(bodyCode), indent(fmt.Sprintf("(set %s (+ %s %s))", varName, varName, by)))

	return addFlowControlVars(loopCode, hasBreak, hasContinue), nil
}

// transpileForEach transpiles a forEach loop block.
// Blockly: controls_forEach with VAR, LIST
// Tatu: Loop with vec:len and vec:get
func transpileForEach(block *Block) (string, error) {
	var bodyCode string

	varName, err := validateField(block, "VAR")
	if err != nil {
		return "", err
	}

	listValue, err := validateValue(block, "LIST")
	if err != nil {
		return "", err
	}

	doStatement, err := validateStatement(block, "DO")
	if err != nil {
		return "", err
	}

	list, err := transpileBlockInternal(listValue, true)
	if err != nil {
		return "", err
	}

	doStatements, err := transpileBlockSequence(doStatement)
	if err != nil {
		return "", err
	}

	bodyCode = getBodyCode(doStatements)
	indexVar := "__i_foreach"

	forBody := fmt.Sprintf("(begin\n%s\n%s\n%s)",
		indent(fmt.Sprintf("(set %s (vec:get %s %s))", varName, list, indexVar)),
		indent(bodyCode),
		indent(fmt.Sprintf("(set %s (+ %s 1))", indexVar, indexVar)))

	forLoop := fmt.Sprintf("(for (var %s 0)\n%s\n%s)",
		indexVar,
		indent(fmt.Sprintf("(< %s (vec:len %s))", indexVar, list)),
		indent(forBody))

	return fmt.Sprintf("(begin\n%s\n%s)",
		indent(fmt.Sprintf("(var %s %s)", varName, "nil")),
		indent(forLoop)), nil
}

// transpileFlowStatement transpiles break/continue statements.
// Blockly: controls_flow_statements with FLOW field (BREAK/CONTINUE)
// Tatu: Sets a flag variable to control loop execution
func transpileFlowStatement(block *Block) (string, error) {
	flowType, err := validateField(block, "FLOW")
	if err != nil {
		return "", err
	}

	switch flowType {
	case "BREAK":
		return "(set __break true)", nil

	case "CONTINUE":
		return "(set __continue true)", nil

	default:
		return "", fmt.Errorf("unknown flow control type: %s", flowType)
	}
}

// hasFlowControl checks if a block or its children contain break/continue statements.
func hasFlowControl(block *Block, flowType string) bool {
	if block == nil {
		return false
	}

	if block.Type == "controls_flow_statements" {
		if flowField := block.Fields["FLOW"]; flowField != nil {
			if flowType == "" || flowField.Value == flowType {
				return true
			}
		}
	}

	for _, value := range block.Values {
		if hasFlowControl(value.Block, flowType) {
			return true
		}
	}

	for _, stmt := range block.Statements {
		if hasFlowControl(stmt.Block, flowType) {
			return true
		}
	}

	return hasFlowControl(block.Next, flowType)
}

// wrapWithFlowControl wraps loop body code with break/continue logic.
func wrapWithFlowControl(bodyCode string, hasBreak, hasContinue bool) string {
	if !hasBreak && !hasContinue {
		return bodyCode
	}

	if hasContinue && !hasBreak {
		return fmt.Sprintf("(begin\n%s\n%s)",
			indent("(set __continue false)"),
			indent(fmt.Sprintf("(if (not __continue)\n%s)", indent(bodyCode))))
	}

	if hasBreak && !hasContinue {
		return fmt.Sprintf("(if (not __break)\n%s)", indent(bodyCode))
	}

	return fmt.Sprintf("(begin\n%s\n%s)",
		indent("(set __continue false)"),
		indent(fmt.Sprintf("(if (and (not __break) (not __continue))\n%s)", indent(bodyCode))))
}

// addFlowControlVars adds variable declarations for break/continue flags.
func addFlowControlVars(code string, hasBreak, hasContinue bool) string {
	if !hasBreak && !hasContinue {
		return code
	}

	var vars string

	if hasBreak {
		vars += indent("(var __break false)") + "\n"
	}
	if hasContinue {
		vars += indent("(var __continue false)") + "\n"
	}

	return fmt.Sprintf("(begin\n%s%s)", vars, code)
}

// modifyLoopCondition modifies a loop condition to check for break flag.
func modifyLoopCondition(condition string, hasBreak bool) string {
	if !hasBreak {
		return condition
	}

	return fmt.Sprintf("(and %s (not __break))", condition)
}

// joinStatements joins multiple statements with proper indentation.
func joinStatements(statements []string) string {
	result := ""

	for i, stmt := range statements {
		if i > 0 {
			result += "\n"
		}

		result += indent(stmt)
	}

	return result
}

// transpileExpression transpiles an expression statement block.
// Blockly: expression_tatu with VALUE input
// Tatu: simply returns the transpiled expression
func transpileExpression(block *Block) (string, error) {
	return transpileUnaryOp(block, "VALUE")
}

// getBodyCode wraps statements in begin expression if needed.
func getBodyCode(doStatements []string) string {
	if len(doStatements) == 0 {
		return "nil"
	}

	firstIsComment := strings.HasPrefix(strings.TrimSpace(doStatements[0]), ";")
	statementsToWrap := doStatements
	if firstIsComment {
		statementsToWrap = doStatements[1:]
	}

	if len(statementsToWrap) <= 1 {
		return strings.Join(doStatements, "\n")
	}

	beginBody := ""

	for i, stmt := range statementsToWrap {
		if i > 0 {
			beginBody += "\n"
		}

		beginBody += indent(stmt)
	}

	result := fmt.Sprintf("(begin\n%s)", beginBody)

	if firstIsComment {
		result = doStatements[0] + "\n" + result
	}

	return result
}
