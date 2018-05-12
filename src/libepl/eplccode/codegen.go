package eplccode

import (
	//"eplc/src/libepl/epllex"
	"eplc/src/libepl/eplparse"
	"io"
)

func GenerateAIR(source io.Reader, fname string) {
	parser := eplparse.New(source, fname)
	parser.Construct()
}