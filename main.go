package main

import (
	"./Output"
	"./applex"
	"io"
	"os"
)

func main() {
	var reader io.Reader
	reader, err := os.Open("/tmp/test.apl")

	if err != nil {
		println("Shit", err.Error())
		os.Exit(0)
	}

	var lexer = apllex.New(reader)

	tmp := lexer.Next()

	for tmp.Ttype != apllex.EOF {
		Output.PrintLog(tmp.Lexme, "is type", tmp.Ttype)
		tmp = lexer.Next()
	}
}
