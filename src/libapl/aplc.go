package main

import (
	"./apllex"
	"./aplparse"
	"./Output"
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
	var parser = aplparse.Parser{lexer}
	var tmp = parser.Lexer.Next()

	for tmp.Ttype != apllex.EOF {
		Output.PrintLog(tmp.Lexme)
		tmp = parser.Lexer.Next()
	}
}
