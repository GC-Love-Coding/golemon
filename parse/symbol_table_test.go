package parse

import (
	"fmt"
	"testing"
)

func TestInsert(t *testing.T) {
	size := 8192
	symTable := NewSymbolTable()
	inputs := make([]string, 0, size)
	symbols := make([]*Symbol, 0, size)

	for i := 0; i < size; i++ {
		inputs = append(inputs, fmt.Sprintf("name = %d", i))
	}

	for _, str := range inputs {
		sym := symTable.Insert(str)
		symbols = append(symbols, sym)
	}

	for i := 0; i < len(inputs); i++ {
		if inputs[i] != symbols[i].name {
			t.Errorf("Expect name: %s, actual name: %s\n", inputs[i], symbols[i].name)
		}
	}

	if symTable.Len() != size {
		t.Errorf("Expect size: %d, actual size: %d\n", size, symTable.Len())
	}
}

func generateLowercaseAlphabets() []string {
	var lower []string

	for i := 'a'; i <= 'z'; i++ {
		lower = append(lower, string(i))
	}

	return lower
}

func generateuppercaseAlphabets() []string {
	var upper []string

	for i := 'A'; i <= 'Z'; i++ {
		upper = append(upper, string(i))
	}

	return upper
}

func TestSortedSymbols(t *testing.T) {
	var inputs []string
	var expected []string
	symTable := NewSymbolTable()

	inputs = append(inputs, generateLowercaseAlphabets()...)
	inputs = append(inputs, generateuppercaseAlphabets()...)
	expected = append(expected, generateuppercaseAlphabets()...)
	expected = append(expected, generateLowercaseAlphabets()...)

	for _, str := range inputs {
		symTable.Insert(str)
	}

	sortedSymbols := symTable.SortedSymbols()

	for i, symbol := range sortedSymbols {
		if name := symbol.Name(); name != expected[i] {
			t.Errorf("Expect name: `%s`, actual: `%s`", expected[i], name)
		}
	}
}
