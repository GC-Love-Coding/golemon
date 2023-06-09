package parse

import (
	"fmt"
	"testing"
)

func TestString(t *testing.T) {
	expect := "abc"
	token := NewToken()

	for _, s := range expect {
		token.AppendRune(rune(s))
	}

	if actual := token.String(); actual != expect {
		t.Errorf("Expect: %s, got: %s helloworld\n", expect, actual)
	}
}

func TestCap(t *testing.T) {
	token := NewToken()

	for i := 0; i < 1024; i++ {
		if i == 256 {
			fmt.Println("hello")
		}

		token.AppendRune(rune(i))
	}
}
