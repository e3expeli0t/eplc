package errors

import (
	"eplc/src/libepl/Output"
	"eplc/src/libepl/epllex"
)

func ParsingError(filename string, line uint, lineOffset uint, errorMsg string, currentline string, token epllex.Token) {
	descriptor := Output.ErrorDescriptor{ Fname: filename,  Line: line, LineOffset: lineOffset, CurrentLine: currentline, ErrorMSG: errorMsg,  Token: token.Lexme}
	//Output.PrintFatalErr(fmt.Sprintf("at %s:%d:%d: Syntax error: %s", filename, line, lineOffset, errorMsg))
	Output.ParserIntelligentError(descriptor)
}