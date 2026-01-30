package blockly

// Workspace represents a Blockly workspace containing blocks.
type Workspace struct {
	Blocks []*Block
}

// Block represents a single Blockly block.
type Block struct {
	Type       string
	ID         string
	Disabled   bool
	Comment    string
	Fields     map[string]*Field
	Values     map[string]*Value
	Statements map[string]*Statement
	Mutation   *Mutation
	Next       *Block
}

// Field represents a Blockly field (dropdown, text input, etc.).
type Field struct {
	Value string
}

// Value represents a Blockly value input (contains nested blocks).
type Value struct {
	Block *Block
}

// Statement represents a Blockly statement input (sequence of blocks).
type Statement struct {
	Block *Block
}

// Mutation represents a Blockly mutation tag (used for procedures, if-else, etc.).
type Mutation struct {
	Name string
	Args []string
}
