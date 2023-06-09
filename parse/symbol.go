package parse

import (
	"hash/fnv"

	"github.com/golemon/util"
)

type SymbolType int
type SymbolAssoc int

const (
	Terminal SymbolType = iota
	NonTerminal
)

const (
	Left SymbolAssoc = iota
	Right
	None
	Unknown
)

// Symbols (terminals and nonterminals) of the grammar are stored in the following.
// TODO: fix data type
type Symbol struct {
	name       string      // Name of the symbol
	index      int         // Index number for this symbol
	symType    SymbolType  // Symbols are all either TERMINALS or NTs
	rule       *Rule       // Linked list of rules of this (if an NT)
	precedence int         // Precedence if defined (-1 otherwise)
	assoc      SymbolAssoc // Associativity if predcence is defined
	firstset   util.IntSet // First-set for all rules of this symbol
	nullable   bool        // True if NT and can generate an empty string
	datatype   string      // The data type of information held by this object. Only used if type==NONTERMINAL
	dtnum      int         // The data type number. In the parser, the value stack is a union. The .yy%d element of this union is the correct data type for this object
}

func NewSymbol(name string) *Symbol {
	symbol := &Symbol{
		name:     name,
		firstset: make(util.IntSet),
		nullable: false,
	}

	if util.IsUpper(name) {
		symbol.symType = Terminal
	} else {
		symbol.symType = NonTerminal
	}

	return symbol
}

func (symbol Symbol) Equal(other Symbol) bool {
	return symbol.name == other.name
}

func (symbol Symbol) Hash() uint32 {
	h := fnv.New32a()
	h.Write([]byte(symbol.name))

	return h.Sum32()
}

func (symbol *Symbol) Name() string {
	return symbol.name
}

func (symbol *Symbol) IsNullable() bool {
	return symbol.nullable
}

func (symbol *Symbol) IsTerminal() bool {
	return symbol.symType == Terminal
}
