package errors

import (
	"eplc/src/libepl/Output"
	"eplc/src/libepl/epllex"
)

type InternalParserError struct {
	ErrCount uint
}

//todo: support position marker (marks the token in specific position)
func (ipe *InternalParserError) ParsingError(filename string, line uint, lineOffset uint, errorMsg string, currentline string, token epllex.Token) {
	descriptor := Output.ErrorDescriptor{
		Fname: filename,
		Line: line,
		LineOffset: lineOffset,
		CurrentLine: currentline,
		ErrorMSG: errorMsg,
		Token: token.Lexme,
	}
	//Output.PrintFatalErr(fmt.Sprintf("at %s:%d:%d: Syntax error: %s", filename, line, lineOffset, errorMsg))
	Output.ParserIntelligentError(descriptor)
	ipe.ErrCount++
}

func (ipe *InternalParserError) IsValidFile() {
	if ipe.ErrCount != 0 {
		Output.PrintFatalErr("Parser: Aborting due to previous errors")
	}
}
