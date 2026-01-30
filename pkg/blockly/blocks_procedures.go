package blockly

import (
	"fmt"
	"strings"
)

// transpileProcedureDefNoReturn transpiles a procedure definition without return value.
// Blockly: procedures_defnoreturn with NAME field, mutation args, and STACK statement
// Tatu: (def name (param1 param2) (begin expr1 expr2 ...))
func transpileProcedureDefNoReturn(block *Block) (string, error) {
	return transpileProcedureDefinition(block, false, "procedures_defnoreturn")
}

// transpileProcedureDefReturn transpiles a procedure definition with return value.
// Blockly: procedures_defreturn with NAME field, mutation args, STACK statement, and RETURN value
// Tatu: (def name (param1 param2) (begin expr1 expr2 ... return value))
func transpileProcedureDefReturn(block *Block) (string, error) {
	return transpileProcedureDefinition(block, true, "procedures_defreturn")
}

// transpileProcedureDefinition transpiles a procedure definition.
func transpileProcedureDefinition(block *Block, hasReturnValue bool, blockTypeName string) (string, error) {
	var funcName string
	var params []string
	var statements []string
	var bodyCode string
	var err error

	if block.Mutation != nil && block.Mutation.Name != "" {
		funcName = block.Mutation.Name
	} else if nameField := block.Fields["NAME"]; nameField != nil {
		funcName = nameField.Value
	} else {
		return "", fmt.Errorf("%s missing NAME field", blockTypeName)
	}

	if block.Mutation != nil {
		params = block.Mutation.Args
	}

	stackStatement := block.Statements["STACK"]
	if stackStatement != nil && stackStatement.Block != nil {
		statements, err = transpileBlockSequence(stackStatement.Block)
		if err != nil {
			return "", err
		}
	}

	if hasReturnValue {
		finalReturnExpression := ""

		returnValue := block.Values["RETURN"]
		if returnValue == nil || returnValue.Block == nil {
			finalReturnExpression = "nil"
		} else {
			finalReturnExpression, err = transpileBlockInternal(returnValue.Block, true)
			if err != nil {
				return "", err
			}
		}

		if len(statements) == 0 {
			bodyCode = finalReturnExpression
		} else {
			allStatements := append(statements, finalReturnExpression)
			bodyCode = fmt.Sprintf("(begin\n%s)", joinStatements(allStatements))
		}
	} else {
		if len(statements) == 0 {
			bodyCode = "nil"
		} else if len(statements) == 1 {
			bodyCode = statements[0]
		} else {
			bodyCode = fmt.Sprintf("(begin\n%s)", joinStatements(statements))
		}
	}

	paramsStr := strings.Join(params, " ")

	return fmt.Sprintf("(def %s (%s)\n%s)", funcName, paramsStr, indent(bodyCode)), nil
}

// transpileProcedureCallNoReturn transpiles a procedure call without using the return value.
// Blockly: procedures_callnoreturn with mutation (name and args), and ARGn values
// Tatu: (function-name arg1 arg2 ...)
func transpileProcedureCallNoReturn(block *Block) (string, error) {
	return transpileProcedureCallReturn(block)
}

// transpileProcedureCallReturn transpiles a procedure call that returns a value.
// Blockly: procedures_callreturn with mutation (name and args), and ARGn values
// Tatu: (function-name arg1 arg2 ...)
func transpileProcedureCallReturn(block *Block) (string, error) {
	if block.Mutation == nil {
		return "", fmt.Errorf("missing mutation")
	}

	funcName := block.Mutation.Name
	argCount := len(block.Mutation.Args)

	args := make([]string, argCount)
	for i := 0; i < argCount; i++ {
		argName := fmt.Sprintf("ARG%d", i)
		argValue := block.Values[argName]

		if argValue == nil || argValue.Block == nil {
			args[i] = "nil"
		} else {
			argCode, err := transpileBlockInternal(argValue.Block, true)
			if err != nil {
				return "", err
			}

			args[i] = argCode
		}
	}

	if len(args) == 0 {
		return fmt.Sprintf("(%s)", funcName), nil
	}

	return fmt.Sprintf("(%s %s)", funcName, strings.Join(args, " ")), nil
}

// transpileProcedureIfReturn transpiles a conditional return statement inside a procedure.
// Blockly: procedures_ifreturn with CONDITION value and optional VALUE value
// Tatu: (if condition value nil) or (if condition nil nil) if no return value
func transpileProcedureIfReturn(block *Block) (string, error) {
	conditionValue := block.Values["CONDITION"]
	if conditionValue == nil {
		return "", fmt.Errorf("procedures_ifreturn missing CONDITION value")
	}

	condition, err := transpileBlockInternal(conditionValue.Block, true)
	if err != nil {
		return "", err
	}

	returnValue := block.Values["VALUE"]
	if returnValue != nil && returnValue.Block != nil {
		returnExpr, err := transpileBlockInternal(returnValue.Block, true)
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("(if %s %s)", condition, returnExpr), nil
	}

	return fmt.Sprintf("(if %s nil)", condition), nil
}
