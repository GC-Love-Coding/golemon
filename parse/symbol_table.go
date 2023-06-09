package parse

import (
	"fmt"
	"sort"
	"strings"
)

const TableSize = 1024

type SymbolTable struct {
	hasNewInsert   bool
	numTerminal    int
	numNonTerminal int
	sortedSymbols  []*Symbol
	symbols        map[string]*Symbol
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		sortedSymbols: make([]*Symbol, 0, TableSize),
		symbols:       make(map[string]*Symbol, TableSize),
	}
}

// This function inserts the symbol with the given name into the symbol table.
// If a symbol with the same name already exists, the existing symbol will be returned.
func (symTable *SymbolTable) Insert(name string) *Symbol {
	if newSym, ok := symTable.Get(name); ok {
		return newSym
	}

	newSym := NewSymbol(name)
	symTable.symbols[name] = newSym
	symTable.sortedSymbols = append(symTable.sortedSymbols, newSym)
	symTable.hasNewInsert = true

	switch newSym.symType {
	case Terminal:
		symTable.numTerminal++
	case NonTerminal:
		symTable.numNonTerminal++
	}

	return newSym
}

func (symTable *SymbolTable) Get(name string) (*Symbol, bool) {
	symbol, ok := symTable.symbols[name]

	return symbol, ok
}

// Get the number of symbols in the symbol table.
func (symTable *SymbolTable) Len() int {
	return len(symTable.symbols)
}

// Get the number of terminals in the symbol table.
func (symTable *SymbolTable) TerminalCount() int {
	return symTable.numTerminal
}

// Get the number of non-terminals in the symbol table.
func (symTable *SymbolTable) NonTerminalCount() int {
	return symTable.numNonTerminal
}

func (symTable *SymbolTable) SortedSymbols() []*Symbol {
	size := symTable.Len()

	if !symTable.hasNewInsert {
		return symTable.sortedSymbols[:size]
	}

	sort.Slice(symTable.sortedSymbols[:size], func(i, j int) bool {
		return symTable.sortedSymbols[i].name < symTable.sortedSymbols[j].name
	})

	// TODO: this may be not necessary.
	// TODO: keep the original order may keep the parser table smaller?
	for i, symbol := range symTable.sortedSymbols {
		symbol.index = i
	}

	symTable.hasNewInsert = false

	return symTable.sortedSymbols[:size]
}

func (symTable *SymbolTable) PrintFirstSets() {
	sortedSymbols := symTable.SortedSymbols()

	for _, symbol := range sortedSymbols {
		fmt.Printf("%s => { ", symbol.Name())
		symNames := make([]string, 0, symbol.firstset.Len())

		for index := range symbol.firstset {
			symNames = append(symNames, sortedSymbols[index].name)
		}

		if symbol.IsNullable() {
			symNames = append(symNames, "ε")
		}

		fmt.Printf(strings.Join(symNames, ", "))
		fmt.Println(" }")
	}
}

// Each symbol in symbol table must either be:
// 1. Terminal
// 2. Non-terminal
//   2.1 Defined with %type
//   2.2
func (symTable *SymbolTable) verifyAllSymbols() {

}
