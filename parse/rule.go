package parse

import (
	"bytes"
)

// Each production rule in the grammar is stored in the following structure.
type Rule struct {
	lhs        *Symbol   // Left-hand side of the rule
	ruleLineno int       // Line number for the rule
	nrhs       int       // Number of RHS symbols
	rhs        []*Symbol // The RHS symbols
	line       int       // Line number at which code begins
	code       string    // The code executed when this rule is reduced
	precSym    *Symbol   // Precedence symbol for this rule
	index      int       // An index number for this rule
	canReduce  bool      // True if this rule is ever reduced
	nextlhs    *Rule     // Next rule with the same LHS
	next       *Rule     // Next rule in the global list
}

func NewRule(symbol *Symbol, ruleLineno int) *Rule {
	rule := &Rule{
		lhs:        symbol,
		ruleLineno: ruleLineno,
		rhs:        make([]*Symbol, 0, 32),
		nextlhs:    symbol.rule,
	}

	// TODO: elegent?
	symbol.rule = rule

	return rule
}

// Get rule's left hand side symbol.
func (rule *Rule) GetLhsSymbol() *Symbol {
	return rule.lhs
}

func (rule *Rule) AppendRhsSymbol(symbol *Symbol) {
	rule.rhs = append(rule.rhs, symbol)
	rule.nrhs++
}

func (rule *Rule) GetRhsSymbolCount() int {
	return rule.nrhs
}

func (rule *Rule) SetCodeAndLine(code string, line int) {
	rule.code = code
	rule.line = line
}

// Get a string representation of this rule.
// This is `LHS: RHS.`. The righ hand side may be empty if
// the symbol on left hand side could be nullable.
func (rule *Rule) String() string {
	var buf bytes.Buffer
	buf.WriteString(rule.lhs.Name())
	buf.WriteRune(':')

	if rule.nrhs > 0 {
		buf.WriteString(rule.rhs[0].Name())
	}

	for i := 1; i < rule.nrhs; i++ {
		buf.WriteString(" " + rule.rhs[i].Name())
	}

	buf.WriteRune('.')

	if rule.precSym != nil {
		buf.WriteRune('[')
		buf.WriteString(rule.precSym.Name())
		buf.WriteRune(']')
	}

	return buf.String()
}

func (rule *Rule) updatePrecedence() {
	for i := 0; i < rule.nrhs; i++ {
		if rule.rhs[i].precedence >= 0 {
			rule.precSym = rule.rhs[i]
			break
		}
	}
}

func (rule *Rule) computeNullable() bool {
	for i := 0; i < rule.nrhs; i++ {
		if !rule.rhs[i].nullable {
			return false
		}
	}

	rule.lhs.nullable = true

	return true
}

func (rule *Rule) computeFirstSet() bool {
	changed := false
	lhsSymbol := rule.lhs

	for i := 0; i < rule.nrhs; i++ {
		rhsSymbol := rule.rhs[i]
		index := rhsSymbol.index

		if rhsSymbol.IsTerminal() {
			if !lhsSymbol.firstset.Includes(index) {
				changed = true
				lhsSymbol.firstset.Add(index)
			}
		} else if lhsSymbol == rhsSymbol {
			if !lhsSymbol.IsNullable() {
				break
			}
		} else {
			oldLen := lhsSymbol.firstset.Len()
			lhsSymbol.firstset.AddSet(rhsSymbol.firstset)

			if lhsSymbol.firstset.Len() != oldLen {
				changed = true
			}

			if !rhsSymbol.IsNullable() {
				break
			}
		}
	}

	return changed
}
