package Errors

import (
	"fmt"
	"aplc/src/libapl/Output"
)

//TokenError prints error msg with precise info about the token that cause the error
func TokenError(line uint, lineOffset uint, token rune, filename string) {
	Output.PrintErr("Lexical", filename+":"+fmt.Sprint(line)+":"+fmt.Sprint(lineOffset)+": Could't resolve Token '"+string(token)+"'")
}

//Lexical prints lexical error msg
func Lexical(msg ...interface{}) {
	Output.PrintErr("Lexical", msg...)
}