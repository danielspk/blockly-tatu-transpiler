package blockly

import (
	"fmt"
	"strings"
)

// transpileMapPair transpiles a map pair block.
// Blockly: map_pair_tatu with KEY and VALUE
// Tatu: (key val)
func transpileMapPair(block *Block) (string, error) {
	key, val, err := transpileBinaryOp(block, "KEY", "VALUE")
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s %s", key, val), nil
}

// transpileMapCreate transpiles a map creation block.
// Blockly: map_create_with_tatu with PAIR0, PAIR1, PAIR2, ... (map_pair_tatu blocks)
// Tatu: (map key1 val1 key2 val2 ...)
func transpileMapCreate(block *Block) (string, error) {
	var pairs []string

	for i := 0; ; i++ {
		pairInput := fmt.Sprintf("PAIR%d", i)
		pairValue := block.Values[pairInput]

		if pairValue == nil {
			break
		}

		if pairValue.Block == nil {
			return "", fmt.Errorf("map_create_with_tatu PAIR%d is empty", i)
		}

		if pairValue.Block.Type != "map_pair_tatu" {
			return "", fmt.Errorf("map_create_with_tatu PAIR%d must be a map_pair_tatu block, got %s", i, pairValue.Block.Type)
		}

		pair, err := transpileMapPair(pairValue.Block)
		if err != nil {
			return "", fmt.Errorf("map_create_with_tatu pair %d error: %w", i, err)
		}

		parts := strings.SplitN(pair, " ", 2)
		if len(parts) == 2 {
			pairs = append(pairs, parts[0], parts[1])
		}
	}

	if len(pairs) == 0 {
		return "(map)", nil
	}

	return fmt.Sprintf("(map %s)", strings.Join(pairs, " ")), nil
}
