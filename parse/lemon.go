package parse

import (
	"fmt"
	"os"

	"github.com/golemon/util"
)

// The state vector for the entire parser generator is recorded as
// follows. (LEMON uses no global variables and makes little use of
// static variables. Fields in the following structure can be thought
// of as begin global variables in the program.)

// char* stacksize

// preccounter:
type Lemon struct {
	sortedState *[]State // Table of states sorted by state number
	rule        []Rule   // List of all rules
	nstate      int      // Number of states
	nrule       int      // Number of rules
	nsymbol     int      // Number of terminal and nonterminal symbols
	nterminal   int      // Number of terminal symbols
	errSym      *Symbol  // The error symbol
	name        string   // Name of the generated parser
	arg         string   // Declaration of the 3th argument to parser
	tokenType   string   // Type of terminal symbols in the parser stack
	varType     string   // The default type of non-terminal symbols
	start       string   // Name of the start symbol for the grammar
	include     string   // Code to put at the start of the C file
	includeLn   int      // Line nunmber for start of include code
	errorCode   string   // Code to execyte when an error is seen
	errorLn     int      // Line number for start of error code
	failure     string   // Code to execute on parser failure
	failureLn   int      // Line number for start of failure code
	accept      string   // Code to execute when the parser accepts
	acceptLn    int      // Line number for the start of accept code
	extraCode   string   // Code appended to the generated file
	extraCodeLn int      // Line number for the start of the extra code
	overflow    string   // Code to execute on a stack overflow
	overflowLn  int      // Line number for start of overflow code
	tokenDest   string   // Code to execute to destroy token data
	tokenDestLn int      // Line number for token destroyer code
	varDest     string   // Code for the default non-terminal destructor code
	varDestLn   int      // Line number for default non-term destructor code
	infile      string   // Name of the input file
	outfile     string   // Name of the current output file
	tokenPrefix string   // A prefix added to token names in the .h file
	nconflict   int      // Number of parsing conflicts
	tableSize   int      // Size of the parse table
	basisFlag   bool     // Print only basis configurations
	argv0       string   // Name of the program
	runeBuf     *RuneBuffer
	lineno      int
}

func NewLemon(infile string, outfile string) *Lemon {
	fp, err := os.Open(infile)

	if err != nil {
		errorf(infile, 0, "Fail to open: "+infile)
	}

	return &Lemon{
		lineno:  1,
		infile:  infile,
		outfile: outfile,
		runeBuf: NewRuneBuffer(fp),
	}
}

func (lemon *Lemon) Parse() {
	ps := NewParserState(lemon)
	filename := lemon.infile
	runeBuf := lemon.runeBuf
	token := NewToken()

	for {
		curRune := runeBuf.GetRune()

		// TODO: verify outside this for loop.
		if curRune == EOF {
			break
		}

		// Keep track of the line number.
		if curRune == NewLine {
			lemon.lineno++
		}

		// Skip the space and newline is also a space.
		// So newline check is done before.
		if util.IsSpace(curRune) {
			continue
		}

		// Skip comment.
		if curRune == '/' {
			lemon.skipComment()
			continue
		}

		ps.startTokLineno = lemon.lineno
		token.AppendRune(curRune)

		// TODO: `'` and `"`
		// TODO: check content between ''
		if curRune == '\'' {
			for curRune = runeBuf.GetRune(); curRune != EOF && curRune != '\''; curRune = runeBuf.GetRune() {
				token.AppendRune(curRune)
			}

			if curRune == EOF {
				ps.errorCnt++
				errorf(filename, lemon.lineno, "String starting on this line is not terminated before the end of the file.")
			} else {
				token.AppendRune(curRune)
			}
		} else if curRune == '<' {
			for curRune = runeBuf.GetRune(); curRune != EOF && curRune != '>'; curRune = runeBuf.GetRune() {
				token.AppendRune(curRune)
			}

			if curRune == EOF {
				ps.errorCnt++
				errorf(filename, lemon.lineno, "Type specifier `<type>` on this line is not terminated before the end of the file.")
			} else {
				token.AppendRune(curRune)
			}
		} else if curRune == '{' {
			level := 1
			for curRune = runeBuf.GetRune(); curRune != EOF && (level > 1 || curRune != '}'); curRune = runeBuf.GetRune() {
				token.AppendRune(curRune)

				if curRune == NewLine {
					lemon.lineno++
				} else if curRune == '}' {
					level--
				} else if curRune == '{' {
					level++
				} else if curRune == '/' {
					lemon.skipComment()
				} else if curRune == '\'' || curRune == '"' {
					var prevRune rune
					expRune := curRune

					for curRune = runeBuf.GetRune(); curRune != EOF && (curRune != expRune || prevRune == '\\'); curRune = runeBuf.GetRune() {
						token.AppendRune(curRune)

						if curRune == NewLine {
							lemon.lineno++
						}

						if prevRune == '\\' {
							prevRune = 0
						} else {
							prevRune = curRune
						}
					}

					if curRune != EOF {
						token.AppendRune(curRune)
					}
				}
			}

			if curRune == EOF {
				ps.errorCnt++
				errorf(lemon.infile, lemon.lineno, "Code starting on this line is not terminated before the end of the file.")
			} else {
				token.AppendRune(curRune)
			}
		} else if util.IsAlphaNum(curRune) {
			for curRune = runeBuf.GetRune(); curRune != EOF && (util.IsAlphaNum(curRune) || curRune == '_'); curRune = runeBuf.GetRune() {
				token.AppendRune(curRune)
			}

			if curRune != EOF {
				runeBuf.UngetRune(curRune)
			}
		}

		ps.parseOneToken(token)
		token.Reset()
	}

	ps.Reprint()
	ps.PrintFirstSets()
}

// Skip over space.
func (lemon *Lemon) skipSpace() {
	var r rune
	lineno := 0
	runeBuf := lemon.runeBuf

	for r = runeBuf.GetRune(); util.IsSpace(r); r = runeBuf.GetRune() {
		if r == NewLine {
			lineno++
		}
	}

	runeBuf.UngetRune(r)
	lemon.lineno += lineno
}

// Skip over comments.
// skipComment is called after reading a '/'.
func (lemon *Lemon) skipComment() {
	runeBuf := lemon.runeBuf
	curRune := runeBuf.GetRune()

	if curRune == '/' {
		for curRune != EOF {
			if curRune == NewLine {
				lemon.lineno++
				return
			}

			curRune = runeBuf.GetRune()
		}

		errorf(lemon.infile, lemon.lineno, "EOF inside comment.")
		return
	}

	if curRune != '*' {
		errorf(lemon.infile, lemon.lineno, "Illegal comment: `%q`", curRune)
	}

	var prevRune rune

	for curRune = runeBuf.GetRune(); curRune != EOF && (prevRune != '*' || curRune != '/'); curRune = runeBuf.GetRune() {
		if curRune == NewLine {
			lemon.lineno++
		}

		prevRune = curRune
	}

	if curRune == EOF {
		errorf(lemon.infile, lemon.lineno, "EOF inside comment.")
	}
}

// Return the number of rules defined in `.y` file.
// Note that, given the `expr : expr '+' expr | expr '-' expr`
// The number of rules will be 2.
func (lemon *Lemon) RuleCount() int {
	return lemon.nrule
}

// Write out error comment.
func errorf(filename string, lineno int, format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	fmt.Fprintf(os.Stderr, ": %v:%v\n", filename, lineno)
	os.Exit(1)
}

// Get the input file name.
func (lemon *Lemon) InputFile() string {
	return lemon.infile
}

// Get the output file name.
func (lemon *Lemon) OutputFile() string {
	return lemon.outfile
}
