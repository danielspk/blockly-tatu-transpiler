package blockly

// transpileComment transpiles a comment block.
// Blockly: comment_tatu with TEXT field
// Tatu: ; comment text
func transpileComment(block *Block) (string, error) {
	text, err := validateFieldAllowEmpty(block, "TEXT")
	if err != nil {
		return "", err
	}

	return "; " + text, nil
}
