package blockly

import (
	"fmt"
	"strings"
)

// transpileListsCreate transpiles a list creation block.
// Blockly: lists_create_with with ADD0, ADD1, ADD2, ... values
// Tatu: (vector item1 item2 ...)
func transpileListsCreate(block *Block) (string, error) {
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

		items = append(items, item)
	}

	if len(items) == 0 {
		return "(vector)", nil
	}

	return fmt.Sprintf("(vector %s)", strings.Join(items, " ")), nil
}

// transpileListsLength transpiles a list length block.
// Blockly: lists_length with VALUE
// Tatu: (vec:len list)
func transpileListsLength(block *Block) (string, error) {
	value, err := transpileUnaryOp(block, "VALUE")
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("(vec:len %s)", value), nil
}

// transpileListsIsEmpty transpiles a list isEmpty check.
// Blockly: lists_isEmpty with VALUE
// Tatu: (= (vec:len list) 0)
func transpileListsIsEmpty(block *Block) (string, error) {
	value, err := transpileUnaryOp(block, "VALUE")
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("(= (vec:len %s) 0)", value), nil
}

// transpileListsGetIndex transpiles a get item from list block.
// Blockly: lists_getIndex with MODE, WHERE, AT values
// Tatu: (vec:get list index)
func transpileListsGetIndex(block *Block) (string, error) {
	modeField, err := validateField(block, "MODE")
	if err != nil {
		return "", err
	}

	whereField, err := validateField(block, "WHERE")
	if err != nil {
		return "", err
	}

	list, err := transpileUnaryOp(block, "VALUE")
	if err != nil {
		return "", err
	}

	if modeField != "GET" {
		return "", fmt.Errorf("lists_getIndex mode %s not yet supported", modeField)
	}

	index, err := transpileWhereClause(block, whereField, "AT", list, "vec:len")
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("(vec:get %s %s)", list, index), nil
}

// transpileListsSetIndex transpiles a set item in list block.
// Blockly: lists_setIndex with MODE, WHERE, AT, TO values
// Tatu: (vec:set list index value)
func transpileListsSetIndex(block *Block) (string, error) {
	modeField, err := validateField(block, "MODE")
	if err != nil {
		return "", err
	}

	whereField, err := validateField(block, "WHERE")
	if err != nil {
		return "", err
	}

	list, to, err := transpileBinaryOp(block, "LIST", "TO")
	if err != nil {
		return "", err
	}

	if modeField != "SET" {
		return "", fmt.Errorf("lists_setIndex mode %s not yet supported", modeField)
	}

	index, err := transpileWhereClause(block, whereField, "AT", list, "vec:len")
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("(vec:set %s %s %s)", list, index, to), nil
}

// transpileListsIndexOf transpiles a find item in list block.
// Blockly: lists_indexOf with MODE, VALUE, FIND
// Tatu: (vec:find list item)
func transpileListsIndexOf(block *Block) (string, error) {
	_, err := validateField(block, "END")
	if err != nil {
		return "", err
	}

	value, find, err := transpileBinaryOp(block, "VALUE", "FIND")
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("(vec:find %s %s)", value, find), nil
}

// transpileListsGetSublist transpiles a get sublist block.
// Blockly: lists_getSublist with WHERE1, WHERE2, AT1, AT2
// Tatu: (vec:slice list start end)
func transpileListsGetSublist(block *Block) (string, error) {
	where1Field, err := validateField(block, "WHERE1")
	if err != nil {
		return "", err
	}

	where2Field, err := validateField(block, "WHERE2")
	if err != nil {
		return "", err
	}

	list, err := transpileUnaryOp(block, "LIST")
	if err != nil {
		return "", err
	}

	start, err := transpileWhereClause(block, where1Field, "AT1", list, "vec:len")
	if err != nil {
		return "", err
	}

	// For WHERE2, we need special handling for FROM_START case (add 1)
	var end string
	switch where2Field {
	case "FROM_START":
		at2Value, err := validateValue(block, "AT2")
		if err != nil {
			return "", err
		}
		at, err := transpileBlockInternal(at2Value, true)
		if err != nil {
			return "", err
		}
		end = fmt.Sprintf("(+ %s 1)", at)

	case "FROM_END":
		at2Value, err := validateValue(block, "AT2")
		if err != nil {
			return "", err
		}
		at, err := transpileBlockInternal(at2Value, true)
		if err != nil {
			return "", err
		}
		end = fmt.Sprintf("(- (vec:len %s) (- %s 1))", list, at)

	case "LAST":
		end = fmt.Sprintf("(vec:len %s)", list)

	default:
		return "", fmt.Errorf("unknown WHERE2: %s", where2Field)
	}

	return fmt.Sprintf("(vec:slice %s %s %s)", list, start, end), nil
}

// transpileListsSplit transpiles a split/join block.
// Blockly: lists_split with MODE (SPLIT/JOIN)
// Tatu: (str:split text delim) or (str:join list delim)
func transpileListsSplit(block *Block) (string, error) {
	modeField, err := validateField(block, "MODE")
	if err != nil {
		return "", err
	}

	input, delim, err := transpileBinaryOp(block, "INPUT", "DELIM")
	if err != nil {
		return "", err
	}

	switch modeField {
	case "SPLIT":
		return fmt.Sprintf("(str:split %s %s)", input, delim), nil

	case "JOIN":
		return fmt.Sprintf("(str:join %s %s)", input, delim), nil

	default:
		return "", fmt.Errorf("unknown lists_split MODE: %s", modeField)
	}
}

// transpileListsSort transpiles a sort list block.
// Blockly: lists_sort with TYPE, DIRECTION
// Tatu: (vec:sort list) or (vec:reverse (vec:sort list))
func transpileListsSort(block *Block) (string, error) {
	_, err := validateField(block, "TYPE")
	if err != nil {
		return "", err
	}

	direction, err := validateField(block, "DIRECTION")
	if err != nil {
		return "", err
	}

	list, err := transpileUnaryOp(block, "LIST")
	if err != nil {
		return "", err
	}

	sortExpr := fmt.Sprintf("(vec:sort %s)", list)

	if direction == "DESCENDING" {
		return fmt.Sprintf("(vec:reverse %s)", sortExpr), nil
	}

	return sortExpr, nil
}

// transpileListsRepeat transpiles a list repeat block.
// Blockly: lists_repeat with ITEM and NUM values
// Tatu: Create a vector by repeating an item N times using a loop
func transpileListsRepeat(block *Block) (string, error) {
	item, num, err := transpileBinaryOp(block, "ITEM", "NUM")
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(`(begin
  (var __result (vector))
  (var __i 0)
  (while (< __i %s)
    (begin
      (set __result (vec:push __result %s))
      (set __i (+ __i 1))))
  __result)`, num, item), nil
}

// transpileListsReverse transpiles a lists reverse block.
// Blockly: lists_reverse with LIST value
// Tatu: (vec:reverse list)
func transpileListsReverse(block *Block) (string, error) {
	list, err := transpileUnaryOp(block, "LIST")
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("(vec:reverse %s)", list), nil
}
