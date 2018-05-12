package libepl

import (
	"eplc/src/libepl/eplccode"
	"io"
)
func Compile(source io.Reader, fname string) { 
	eplccode.GenerateAIR(source, fname)
}