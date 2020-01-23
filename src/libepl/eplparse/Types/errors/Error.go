package errors

import (
	"eplc/src/libepl/Output"
	"eplc/src/libepl/epllex"
)

func UnresolvedTypeError(msg string, token epllex.Token, fname string, line string) {
	descriptor := Output.ErrorDescriptor{
		Fname:		 fname ,
		Line:        token.StartLine,
		LineOffset:  token.StartOffset,
		CurrentLine: line,
		ErrorMSG:    msg,
		Token:       token.Lexme,
	}
	Output.TypeIntelligentError(descriptor)
}
