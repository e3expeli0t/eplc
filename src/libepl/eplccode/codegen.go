package eplccode

import (
	"eplc/src/libepl/epllex"
	"eplc/src/libepl/eplparse"
	"io"
)

func GenerateAIR(source io.Reader, fname string) {
	lexer := epllex.New(source, fname)
	parser := eplparse.New(lexer)
	parser.Construct()
}