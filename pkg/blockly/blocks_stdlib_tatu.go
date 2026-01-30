package blockly

import (
	"fmt"
	"strings"
)

// transpileFunctionCall transpiles a generic function call (stdlib or custom).
// Blockly: function_call with NAME field and ARG0, ARG1, ... ARGn inputs
// Tatu: (function-name arg1 arg2 ...)
func transpileFunctionCall(block *Block) (string, error) {
	functionName, err := validateField(block, "NAME")
	if err != nil {
		return "", err
	}

	var args []string

	for i := 0; ; i++ {
		argKey := fmt.Sprintf("ARG%d", i)
		argValue := block.Values[argKey]

		if argValue == nil || argValue.Block == nil {
			break
		}

		argExpr, err := transpileBlockInternal(argValue.Block, true)
		if err != nil {
			return "", err
		}
		args = append(args, argExpr)
	}

	if len(args) == 0 {
		return fmt.Sprintf("(%s)", functionName), nil
	}

	return fmt.Sprintf("(%s %s)", functionName, strings.Join(args, " ")), nil
}
