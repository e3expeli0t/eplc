package errors

import (
	"eplc/src/libepl/epllex"
	"eplc/src/libio"
)

func UnresolvedTypeError(msg string, token epllex.Token, fname string, line string) {
	descriptor := libio.ErrorDescriptor{
		Fname:		 fname ,
		Line:        token.StartLine,
		LineOffset:  token.StartOffset,
		CurrentLine: line,
		ErrorMSG:    msg,
		Token:       token.Lexme,
	}
	libio.TypeIntelligentError(descriptor)
}
