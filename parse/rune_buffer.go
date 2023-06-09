package parse

import (
	"bufio"
	"io"
)

const (
	EOF     = -1
	NewLine = '\n'
)

type RuneBuffer struct {
	reader   *bufio.Reader // A pointer to buffer reader
	peekRune rune          // Peek rune
}

func NewRuneBuffer(rd io.Reader) *RuneBuffer {
	return &RuneBuffer{bufio.NewReader(rd), 0}
}

func (runeBuf *RuneBuffer) GetRune() rune {
	var r rune
	peekRune := runeBuf.peekRune

	if peekRune != 0 {
		if peekRune == EOF {
			return EOF
		}

		r = peekRune
		runeBuf.peekRune = 0

		return r
	}

	r, n, err := runeBuf.reader.ReadRune()

	if n == 0 {
		return EOF
	}

	if err != nil {
		panic("Read error: " + err.Error())
	}

	return r
}

func (runeBuf *RuneBuffer) UngetRune(c rune) {
	if runeBuf.peekRune != 0 {
		panic("UngetRune - 2nd unget")
	}

	runeBuf.peekRune = c
}
