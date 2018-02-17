package main

import
(
	"./apllex"
	"./aplparse"
	"./Output"
	"io"
	"os"
)
/*
todo: add errors in the lexer to the fallowing operations: StringLiterals matching ,MultiLines comments (will be added soon)
*/
func main() {
	var reader io.Reader
	reader, err := os.Open("/tmp/test.apl")

	if err != nil {
		println("Shit", err.Error())
		os.Exit(0)
	}

	var lexer = apllex.New(reader)
	var parser = aplparse.Parser{lexer}
	Output.PrintLog(&parser, &lexer)
}
