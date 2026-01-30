package blockly

import (
	"fmt"
)

// transpileInclude transpiles an include block.
// Blockly: include_file_tatu with PATH field
// Tatu: (include "path/to/file.tatu")
func transpileInclude(block *Block) (string, error) {
	path, err := validateField(block, "PATH")
	if err != nil {
		return "", err
	}

	if path == "" {
		return "", fmt.Errorf("include_file_tatu block has empty PATH")
	}

	return fmt.Sprintf("(include \"%s\")", path), nil
}
