package error

import (
	"fmt"
)

func ErrorMsg(filename string, lineno int, format string, a ...interface{}) {
	fmt.Printf("%s:%d => %s\n", filename, lineno, fmt.Sprintf(format, a...))
}
