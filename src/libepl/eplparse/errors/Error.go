package errors

import (
	"eplc/src/libepl/Output"
	"fmt"
)

func ParsingError(filename string, line uint, lineOffset uint, errorMsg string) {
	Output.PrintFatalErr("<Syntatic>", fmt.Sprintf("%s:%d:%d: Syntax error: %s", filename, line, lineOffset, errorMsg))
}