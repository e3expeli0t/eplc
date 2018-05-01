package main

import (
	"aplc/src/libapl/apllex"
	"aplc/src/libapl/Output"
	"io"
	"os"
)

func main() {
	Output.PrintStartMSG()

	if len(os.Args) <= 1 {
		Output.PrintErr("please supply file name")
	}

	args := os.Args[1:]
	file := args[0]

	var reader io.Reader
	reader, err := os.Open(file)

	if err != nil {
		Output.PrintErr("file: '" + file + "' don't exists")
	}
	var lexer = apllex.New(reader, file)
	Output.PrintLog(" <aplc> Parsing: '"+file+"' ")
	var tmp = lexer.Next()

	for tmp.Ttype != apllex.EOF {
		Output.PrintLog(tmp.Lexme)
		tmp = lexer.Next()
	}

}
