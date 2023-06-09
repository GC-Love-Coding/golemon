package parse

import (
	"fmt"
	"strings"

	"github.com/golemon/util"
)

// TODO: right hand side non-terminal must occur in rule left hand side once or defined as %type.

// https://www.ibm.com/docs/en/aix/7.2?topic=information-yacc-grammar-file-declarations
// Terminal (or token) names can be declared using the %token declaration
// Nonterminal names can be declared using the %type declaration.
// Keyword 		Description
// %left		Identifies tokens that are left-associative with other tokens.
// %nonassoc	Identifies tokens that are not associative with other tokens.
// %right		Identifies tokens that are right-associative with other tokens.
// %start		Identifies a nonterminal name for the start symbol.
// %token		Identifies the token names that the yacc command accepts. Declares all token names in the declarations section.
// %type		Identifies the type of nonterminals. Type-checking is performed when this construct is present.
// %union		Identifies the yacc value stack as the union of the various type of values desired. By default, the values returned are integers. The effect of this construct is to provide the declaration of YYSTYPE directly from the input.
//
// All of the tokens on the same line have the same precedence level and associativity.
// The lines appear in the file in order of increasing precedence or binding strength.
// For example, the following describes the precedence and associativity of the four arithmetic operators:
// %token [<Tag>] Name [Number] [Name [Number]]...

// Precedence in the grammar rules
// https://www.ibm.com/docs/en/zos/2.3.0?topic=section-precedence-in-grammar-rules

// flex & bison
// Every token in a grammar can have a precedence and an associativity assigned by a precedence declaration.
// Every rule can also have a precedence and an associativity, which is taken from a %prec clause in the rule or,
// failing that, the rightmost token in the rule that has a precedence assigned.
// Whenever there is a shift/reduce conflict, bison compares the precedence of the token that might be shifted
// to that of the rule that might be reduced. It shifts if the token’s precedence is higher or reduces if the
// rule’s precedence is higher. If both have the same precedence, bison checks the associativity.
// If they are left associative, it reduces; if they are right associative, it shifts; and if they are nonassociative,
// bison generates an error.
// Typical Uses of Precedence
// Although you can in theory use precedence to resolve any kind of shift/reduce conflict, you should use precedence only for a few well-understood situations and rewrite the grammar otherwise. Precedence declarations were designed to handle expression gram- mars, with large numbers of rules like this:
type FsmState int

const (
	WaitPercentSign FsmState = iota
	WaitOpenBrace
	WaitKwDefOrRule1
	WaitKwDefOrRule2
	WaitOptTagOrOpenBrace
	WaitSymbolAfterKeyword

	WaitRuleLhsSymbol
	WaitColon
	WaitRuleRhsSymbol
	WaitPrecedence
	WaitPrecedenceTerm
	WaitSymbolAfterPrec

	WaitSubRoutine1
	WaitSubRoutine2
)

type Keyword int

const (
	KwUnknown Keyword = iota
	KwType
	KwToken
	KwLeft
	KwRight
	KwNonassoc
	KwStart
	KwPrec
	KwUnion
)

// TODO: case sensitivity
var ReservedKeywords = map[Keyword]string{
	KwType:     "TYPE",
	KwToken:    "TOKEN",
	KwLeft:     "LEFT",
	KwRight:    "RIGHT",
	KwNonassoc: "NONASSOC",
	KwStart:    "START",
	KwPrec:     "PREC",
	KwUnion:    "UNION",
}

// The state of the parser.
type ParserState struct {
	gp              *Lemon   // The owner of this parser state
	errorCnt        int      // Number of errors so far
	currToken       *Token   // Text of current token
	curState        FsmState // Current state of the parser
	importCode      []rune   // Import code inside %{ %}
	unionCode       string   // Union type definition
	unionCodeLineno int      // Union code line number
	datatype        string   // %type definition
	startTokLineno  int      // Start token line number
	prevKeyword     Keyword  // Previous keyword
	prevTag         string
	// lhs            *Symbol     // Left-hand side of current rule
	// nrhs           int         // Number of right-hand side symbols seen
	// rhs            []*Symbol   // RHS symbols
	prevRule    *Rule       // Previous rule parsed.
	declKeyword string      // Keyword of a declaration
	declArgSlot *string     // Where the declaration argument should be put
	declLnSlot  *int        // Where the declaration line number is put
	declAssoc   SymbolAssoc // Assign this association to decl arguments
	precCounter int         // Assign this precedence to decl arguments
	firstRule   *Rule       // Pointer to first rule in the grammar
	lastRule    *Rule       // Pointer to the most recently parsed rule
	subroutine  *strings.Builder
	symTable    *SymbolTable
}

func stateToString(state FsmState) string {
	switch state {
	case WaitPercentSign:
		return "Wait %"
	case WaitOpenBrace:
		return "Wait {"
	case WaitKwDefOrRule1:
		return "WaitKwDefOrRule1"
	case WaitKwDefOrRule2:
		return "WaitKwDefOrRule2"
	case WaitRuleLhsSymbol:
		return "Wait for rule left hand side"
	case WaitOptTagOrOpenBrace:
		return "Wait < or Lhs symbol"
	case WaitSymbolAfterKeyword:
		return "Wait symbol after keyword"
	case WaitColon:
		return "Wait `:`"
	case WaitRuleRhsSymbol:
		return "Wait rule rhs symbol"
	case WaitPrecedence:
		return "Wait `prec`"
	case WaitPrecedenceTerm:
		return "Wait precedence terminal"
	case WaitSymbolAfterPrec:
		return "Wait symbol after precedence"
	case WaitSubRoutine1:
		return "Wait subroutine1"
	case WaitSubRoutine2:
		return "Wait subroutine2"

	default:
		return "Not implemented"
	}
}

func NewParserState(gp *Lemon) *ParserState {
	return &ParserState{
		gp:          gp,
		subroutine:  &strings.Builder{},
		prevKeyword: KwUnknown,
		symTable:    NewSymbolTable(),
	}
}

func (ps *ParserState) appendRule(rule *Rule) {
	if ps.firstRule == nil {
		ps.firstRule = rule
	} else {
		ps.lastRule.next = rule
	}

	ps.lastRule = rule
	ps.prevRule = rule
	ps.gp.nrule++
}

func (ps *ParserState) parseOneToken(token *Token) {
	runeCount := token.RuneCount()
	tokenStr := token.String()
	upperStr := strings.ToUpper(tokenStr)
	fstRune := token.FirstRune()
	lstRune := token.LastRune()
	startLineno := ps.startTokLineno
	filename := ps.gp.InputFile()
	symTable := ps.symTable

	fmt.Println("State=", stateToString(ps.curState), "[Token=", tokenStr, "]")

	if tokenStr == "'('" {
		fmt.Println()
	}

	switch ps.curState {
	case WaitPercentSign:
		if fstRune != '%' {
			ps.errorCnt++
			errorf(filename, startLineno, "Declaration must start with `%{` and end with `%}`. Find: `%s`", tokenStr)
		} else {
			ps.curState = WaitOpenBrace
		}

	case WaitOpenBrace:
		// TODO: bug here `%   {` will be allowed.
		// Precondition:
		// 1. First rune must be `{`.
		// 2. Last rune must be `}`.
		// 3. Second last rune must be `%`.
		if fstRune != '{' || lstRune != '}' || token.NthRune(runeCount-2) != '%' {
			ps.errorCnt++
			errorf(filename, startLineno, "Declaration must start with `%{` and end with `%}`. Find: `%s`", tokenStr)
		} else {
			// Ignore first `{` and last `}`.
			ps.importCode = make([]rune, runeCount-2)
			copy(ps.importCode, token.Buffer()[1:runeCount])
			ps.curState = WaitKwDefOrRule1
		}

	case WaitKwDefOrRule1:
		if fstRune != '%' {
			ps.errorCnt++
			errorf(filename, startLineno, "Expect `%keyword` to declare keyword or `%%` to start rule definition. Find: `%s`", tokenStr)
		} else {
			ps.curState = WaitKwDefOrRule2
		}

	case WaitKwDefOrRule2:
		if fstRune == '%' {
			ps.curState = WaitRuleLhsSymbol
		} else {
			for k, v := range ReservedKeywords {
				if v == upperStr {
					ps.prevKeyword = k
					break
				}
			}

			if ps.prevKeyword == KwUnknown {
				ps.errorCnt++
				errorf(filename, startLineno, "Expect `%keyword` to declare keyword or `%%` to start rule definition. Find: `%s`", tokenStr)
			} else {
				ps.curState = WaitOptTagOrOpenBrace
			}
		}

	case WaitOptTagOrOpenBrace:
		// 1. %union {}
		// 2. %type [<tag>] non-terminal
		// 3. %token [<tag>] terminal
		// 4. %left [<tag>] terminal
		// 5. %right [<tag>] terminal
		// 6. %nonassoc [<tag>] terminal
		// 7. %start [<tag>] non-terminal
		if ps.prevKeyword == KwUnion {
			if fstRune != '{' || lstRune != '}' {
				ps.errorCnt++
				errorf(filename, startLineno, "Expect `{}` after `%union`: `%s`", tokenStr)
			} else {
				// TODO: union declared once?
				if len(ps.unionCode) > 0 {
					ps.errorCnt++
					errorf(filename, startLineno, "Multiple `%union` definitions are found. Previous definition is at: %d", ps.unionCodeLineno)
				} else {
					ps.unionCode = tokenStr
					ps.unionCodeLineno = startLineno
					ps.curState = WaitKwDefOrRule1
				}
			}
		} else if fstRune == '<' {
			if ps.prevKeyword == KwUnion {
				ps.errorCnt++
				errorf(filename, startLineno, "Tag specifier `<>` can't follow after `%union`: `%s`", tokenStr)
			} else {
				ps.prevTag = string(token.Buffer()[1:runeCount])
				ps.curState = WaitSymbolAfterKeyword
			}
		} else {
			ps.defineSymbol(tokenStr)
			ps.curState = WaitSymbolAfterKeyword
		}

	case WaitSymbolAfterKeyword:
		if fstRune == '%' {
			// We need to clear previous keyword and tag.
			ps.prevKeyword = KwUnknown
			ps.prevTag = ""
			ps.curState = WaitKwDefOrRule2
		} else {
			ps.defineSymbol(tokenStr)
		}

	case WaitRuleLhsSymbol:
		// TODO: may need to check the existence of symbol
		// At least, one rule is defined.
		if fstRune == '%' {
			if ps.gp.RuleCount() == 0 {
				ps.errorCnt++
				errorf(filename, startLineno, "Unexpected `%%`, at least 1 rule must be defined.")
			} else {
				ps.curState = WaitSubRoutine1
			}
		} else if !util.IsLower(tokenStr) {
			ps.errorCnt++
			errorf(filename, startLineno, "For rule definition, left hand side symbol must be non-terminal: `%s`.", tokenStr)
		} else {
			symbol := symTable.Insert(tokenStr)
			rule := NewRule(symbol, startLineno)
			ps.appendRule(rule)

			ps.curState = WaitColon
		}

	case WaitColon:
		if fstRune != ':' {
			ps.errorCnt++
			errorf(filename, startLineno, "Expect `:` after non-terminal: `%s`", tokenStr)
		} else {
			ps.curState = WaitRuleRhsSymbol
		}

	case WaitRuleRhsSymbol:
		prevRule := ps.prevRule
		// For grammar `lhs: | expr;`.
		// There is no symbol on the right hand side for the first rule.
		// This means `lhs` symbol coule be nullable.
		if fstRune == '|' {
			count := prevRule.GetRhsSymbolCount()
			symbol := prevRule.GetLhsSymbol()

			if count == 0 {
				// TODO: how about `symbol := | | {}`. This is a warnning in bison.
				if symbol.nullable {
					ps.errorCnt++
					errorf(filename, startLineno, "Find multiple empty expression for: `%s`", symbol.Name())
				}
				symbol.nullable = true
			} else {
				// TODO: previous rule may needs default action code.
				rule := NewRule(symbol, startLineno)
				ps.appendRule(rule)
			}
		} else if fstRune == '{' {
			// TODO: check {}{}
			// Grammar like: `expr: {}` is ok.
			prevRule.SetCodeAndLine(tokenStr, startLineno)
		} else if fstRune == '%' {
			ps.curState = WaitPrecedence
		} else if fstRune == ';' {
			// End of this rule.
			ps.prevRule = nil
			ps.curState = WaitRuleLhsSymbol
		} else {
			// ps.errorCnt++
			// errorf(filename, startLineno, "Expect right hand symbol but find `%s`", tokenStr)
			else if symbol, ok := symTable.Get(tokenStr); ok {
				ps.prevRule.AppendRhsSymbol(symbol)
			} 
		}

	case WaitPrecedence:
		if upperStr != ReservedKeywords[KwPrec] {
			ps.errorCnt++
			errorf(filename, startLineno, "Expect `%%prec`. Find: `%s`.", &tokenStr)
		} else {
			ps.curState = WaitPrecedenceTerm
		}

	case WaitPrecedenceTerm:
		if symbol, ok := symTable.Get(tokenStr); !ok || symbol == nil {
			ps.errorCnt++
			errorf(filename, startLineno, "Terminal after `%%prec` must be defined: `%s`.", tokenStr)
		} else {
			ps.prevRule.precSym = symbol
			ps.curState = WaitSymbolAfterPrec
		}

	case WaitSymbolAfterPrec:
		switch fstRune {
		case '{':
			// TODO: include `{}`?
			ps.prevRule.code = tokenStr
			fallthrough
		case '|':
			// TODO: default action.
			fallthrough
		case ';':
			ps.prevRule = nil
			ps.curState = WaitRuleLhsSymbol
		default:
			errorf(filename, startLineno, "Expect `|` or `{` or `;` after `%%prec term`: `%s`", tokenStr)
		}

	case WaitSubRoutine1:
		if fstRune != '%' {
			ps.errorCnt++
			errorf(filename, startLineno, "Expect `%%` after `%%` before subroutine: `%s`", tokenStr)
		} else {
			ps.curState = WaitSubRoutine2
		}

	case WaitSubRoutine2:
		ps.subroutine.WriteString(tokenStr)
	}
}

// Define a symbol based on previous keyword.
func (ps *ParserState) defineSymbol(symName string) *Symbol {
	kw := ps.prevKeyword
	startLineno := ps.startTokLineno
	filename := ps.gp.InputFile()
	symTable := ps.symTable

	// TODO: is it ok insert before `errorf`
	symbol := symTable.Insert(symName)

	if len(ps.prevTag) > 0 {
		symbol.datatype = ps.prevTag
	}

	switch kw {
	case KwType:
		if !util.IsLower(symName) {
			errorf(filename, startLineno, "Non-terminal must be lower case: `%s`", symName)
		} else {
			symbol.symType = NonTerminal
		}

	case KwToken:
		// '+' or NUMBER.
		if !util.IsUpper(symName) && !util.IsStringLiteral(symName) {
			errorf(filename, startLineno, "Terminal must be upper case or string literal: `%s`", symName)
		} else {
			symbol.symType = Terminal
		}

	case KwLeft, KwRight, KwNonassoc:
		// Must be terminal.
		if !util.IsUpper(symName) {
			errorf(filename, startLineno, "%s must followed by terminal: `%s`", ReservedKeywords[kw], symName)
		} else {
			symbol.symType = Terminal
			switch kw {
			case KwLeft:
				symbol.assoc = Left
			case KwRight:
				symbol.assoc = Right
			case KwNonassoc:
				symbol.assoc = None
			}
			symbol.precedence = ps.precCounter
			ps.precCounter++
		}

	case KwStart, KwPrec:
		fmt.Println("Not implement KwStart")
	}

	return symbol
}

// Duplicate the input file without comments and without actions on rules.
func (ps *ParserState) Reprint() {
	maxLen := 10
	sortedSymbols := ps.symTable.SortedSymbols()

	fmt.Printf("// Reprint of input file \"%s\".\n// Symbols:\n", ps.gp.InputFile())

	for _, sym := range sortedSymbols {
		nameLen := len(sym.Name())

		if nameLen > maxLen {
			maxLen = nameLen
		}
	}

	ncolumns := 76 / (maxLen + 5)

	if ncolumns < 1 {
		ncolumns = 1
	}

	skip := (len(sortedSymbols) + ncolumns - 1) / ncolumns
	for i := 0; i < skip; i++ {
		fmt.Printf("//")
		for j := i; j < len(sortedSymbols); j += skip {
			fmt.Printf(" %3d %-*.*s", j, maxLen, maxLen, sortedSymbols[j].Name())
		}
		fmt.Println()
	}

	for rule := ps.firstRule; rule != nil; rule = rule.next {
		fmt.Println(rule.String())
	}
}

// Those rules which have a precedence symbol coded in the input
// grammar using the "[symbol]" construct will already have the
// rp->precsym field filled. Other rules take as their precedence
// symbol the first RHS symbol with a defined precedence. If there
// are not RHS symbols with a defined precedence, the precedece
// symbol field is left blank.
func (ps *ParserState) updateRulePrecedences() {
	for rp := ps.firstRule; rp != nil; rp = rp.next {
		if rp.precSym == nil {
			continue
		}

		rp.updatePrecedence()
	}
}

// For each terminal t
//   Nullable(t) = false
// For each non-terminal N
//   Nullable(N) = is there a production N ::= ε(epsilon)
// Repeat
//   For each production N ::= x1x2x3...xn
//     If Nullable(xi) for all of xi then set Nullable(N) to true
// Util nothing new becomes Nullable
func (ps *ParserState) computeNullableSets() {
	changed := true

	for changed {
		changed = false

		for rp := ps.firstRule; rp != nil; rp = rp.next {
			symbol := rp.GetLhsSymbol()

			if symbol.IsNullable() {
				continue
			}

			changed = changed || rp.computeNullable()
		}
	}
}

func (ps *ParserState) computeFirstSets() {
	ps.computeNullableSets()

	for rule := ps.firstRule; rule != nil; rule = rule.next {
		rule.computeFirstSet()
	}
}

// Print the first sets.
func (ps *ParserState) PrintFirstSets() {
	ps.computeFirstSets()
	ps.symTable.PrintFirstSets()
}
