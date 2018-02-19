package Errors

import (
	"../../Output"
)

func TokenError(msg ...interface{}) {
	Output.PrintErr("Lexical", msg...)
	//panic("Lexical analysis error")
}

func Lexical(msg ...interface{}) {
	Output.PrintErr("Lexical", msg...)
}
