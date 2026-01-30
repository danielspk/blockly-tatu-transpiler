package blockly

import (
	"fmt"
	"strconv"
	"strings"
)

// transpileText transpiles a text block.
// Blockly: text with TEXT field
// Tatu: "string value"
func transpileText(block *Block) (string, error) {
	text, err := validateFieldAllowEmpty(block, "TEXT")
	if err != nil {
		return "", err
	}

	return strconv.Quote(text), nil
}

// transpileTextPrint transpiles a print block.
// Blockly: text_print with TEXT value
// Tatu: (print value)
func transpileTextPrint(block *Block) (string, error) {
	value, err := transpileUnaryOp(block, "TEXT")
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("(print %s)", value), nil
}

// transpileTextJoin transpiles a text join block.
// Blockly: text_join with variable number of items
// Tatu: (str:concat (to-string item1) (to-string item2) ...)
func transpileTextJoin(block *Block) (string, error) {
	var items []string

	for i := 0; ; i++ {
		itemKey := fmt.Sprintf("ADD%d", i)
		itemValue := block.Values[itemKey]
		if itemValue == nil {
			break
		}

		item, err := transpileBlockInternal(itemValue.Block, true)
		if err != nil {
			return "", err
		}

		if !strings.HasPrefix(item, "\"") {
			item = fmt.Sprintf("(to-string %s)", item)
		}

		items = append(items, item)
	}

	if len(items) == 0 {
		return `""`, nil
	}
	if len(items) == 1 {
		return items[0], nil
	}

	return fmt.Sprintf("(str:concat %s)", strings.Join(items, " ")), nil
}

// transpileTextLength transpiles a text length block.
// Blockly: text_length with VALUE
// Tatu: (str:len value)
func transpileTextLength(block *Block) (string, error) {
	value, err := transpileUnaryOp(block, "VALUE")
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("(str:len %s)", value), nil
}

// transpileTextIsEmpty transpiles a text isEmpty check.
// Blockly: text_isEmpty with VALUE
// Tatu: (= (str:len value) 0)
func transpileTextIsEmpty(block *Block) (string, error) {
	value, err := transpileUnaryOp(block, "VALUE")
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("(= (str:len %s) 0)", value), nil
}

// transpileTextCharAt transpiles a get character at position block.
// Blockly: text_charAt with WHERE field and AT value
// Tatu: (str:slice text pos (+ pos 1))
func transpileTextCharAt(block *Block) (string, error) {
	where, err := validateField(block, "WHERE")
	if err != nil {
		return "", err
	}

	if where == "RANDOM" {
		return "", fmt.Errorf("text_charAt RANDOM not yet supported")
	}

	text, err := transpileUnaryOp(block, "VALUE")
	if err != nil {
		return "", err
	}

	switch where {
	case "FROM_START":
		atValue, err := validateValue(block, "AT")
		if err != nil {
			return "", err
		}
		at, err := transpileBlockInternal(atValue, true)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("(str:slice %s %s (+ %s 1))", text, at, at), nil

	case "FROM_END":
		atValue, err := validateValue(block, "AT")
		if err != nil {
			return "", err
		}
		at, err := transpileBlockInternal(atValue, true)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("(str:slice %s (- (str:len %s) %s) (- (str:len %s) (- %s 1)))", text, text, at, text, at), nil

	case "FIRST":
		return fmt.Sprintf("(str:slice %s 0 1)", text), nil

	case "LAST":
		return fmt.Sprintf("(str:slice %s (- (str:len %s) 1) (str:len %s))", text, text, text), nil

	default:
		return "", fmt.Errorf("unknown text_charAt WHERE: %s", where)
	}
}

// transpileTextGetSubstring transpiles a get substring block.
// Blockly: text_getSubstring with WHERE1/WHERE2 and AT1/AT2
// Tatu: (str:slice text start end)
func transpileTextGetSubstring(block *Block) (string, error) {
	where1, err := validateField(block, "WHERE1")
	if err != nil {
		return "", err
	}

	where2, err := validateField(block, "WHERE2")
	if err != nil {
		return "", err
	}

	text, err := transpileUnaryOp(block, "STRING")
	if err != nil {
		return "", err
	}

	start, err := transpileWhereClause(block, where1, "AT1", text, "str:len")
	if err != nil {
		return "", err
	}

	end, err := transpileWhereClause(block, where2, "AT2", text, "str:len")
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("(str:slice %s %s %s)", text, start, end), nil
}

// transpileTextAppend transpiles a text append block.
// Blockly: text_append with VAR and TEXT
// Tatu: (set var (str:concat var text))
func transpileTextAppend(block *Block) (string, error) {
	varName, err := validateField(block, "VAR")
	if err != nil {
		return "", err
	}

	text, err := transpileUnaryOp(block, "TEXT")
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("(set %s (str:concat %s %s))", varName, varName, text), nil
}

// transpileTextIndexOf transpiles a find text block.
// Blockly: text_indexOf with END, VALUE, FIND
// Tatu: (str:index value find)
func transpileTextIndexOf(block *Block) (string, error) {
	endField, err := validateField(block, "END")
	if err != nil {
		return "", err
	}

	if endField == "LAST" {
		return "", fmt.Errorf("text_indexOf LAST not supported")
	}

	value, find, err := transpileBinaryOp(block, "VALUE", "FIND")
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("(str:index %s %s)", value, find), nil
}

// transpileTextChangeCase transpiles a change case block.
// Blockly: text_changeCase with CASE (UPPERCASE, LOWERCASE, TITLECASE)
// Tatu: (str:upper text), (str:lower text)
func transpileTextChangeCase(block *Block) (string, error) {
	caseField, err := validateField(block, "CASE")
	if err != nil {
		return "", err
	}

	if caseField == "TITLECASE" {
		return "", fmt.Errorf("text_changeCase TITLECASE not yet supported")
	}

	text, err := transpileUnaryOp(block, "TEXT")
	if err != nil {
		return "", err
	}

	switch caseField {
	case "UPPERCASE":
		return fmt.Sprintf("(str:upper %s)", text), nil

	case "LOWERCASE":
		return fmt.Sprintf("(str:lower %s)", text), nil

	default:
		return "", fmt.Errorf("unknown text_changeCase CASE: %s", caseField)
	}
}

// transpileTextTrim transpiles a trim whitespace block.
// Blockly: text_trim with MODE (BOTH, LEFT, RIGHT)
// Tatu: (str:trim text) - only supports trimming both sides
func transpileTextTrim(block *Block) (string, error) {
	modeField, err := validateField(block, "MODE")
	if err != nil {
		return "", err
	}

	text, err := transpileUnaryOp(block, "TEXT")
	if err != nil {
		return "", err
	}

	if modeField != "BOTH" {
		return "", fmt.Errorf("text_trim only supports BOTH mode, got: %s", modeField)
	}

	return fmt.Sprintf("(str:trim %s)", text), nil
}
