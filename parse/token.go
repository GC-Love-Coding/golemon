package parse

import "fmt"

const MaxBufferSize = (1 << 20)

// A token is represented by a continuous slice of runes.
// TODO: use bytes.Buffer or []rune
type Token struct {
	buffer []rune
}

func NewToken() *Token {
	return &Token{
		buffer: make([]rune, 0, MaxBufferSize),
	}
}

// Append a rune to current token.
// Example:
// Current token contains "abc" and after AppendRune('d'), token will become "abcd"
func (token *Token) AppendRune(r rune) {
	token.buffer = append(token.buffer, r)
}

func (token *Token) NthRune(n int) rune {
	if n < 0 || n >= len(token.buffer) {
		panic(fmt.Sprintf("Index out of bound: %d.", n))
	}

	return token.buffer[n]
}

// Get the first rune.
func (token *Token) FirstRune() rune {
	return token.NthRune(0)
}

// Get the last rune.
func (token *Token) LastRune() rune {
	return token.NthRune(token.RuneCount() - 1)
}

// Reset token only sets the cursor to 0 so as to reuse the buffer.
func (token *Token) Reset() {
	token.buffer = token.buffer[:0]
}

func (token *Token) String() string {
	return string(token.buffer[:token.RuneCount()])
}

// Return the number of runes inside token.
// Note that this may differ the number of bytes inside token.
func (token *Token) RuneCount() int {
	return len(token.buffer)
}

// Get the underlying buffer.
func (token *Token) Buffer() []rune {
	return token.buffer
}
