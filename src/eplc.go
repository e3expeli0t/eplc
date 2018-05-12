package main

import (
	"eplc/src/libepl"
	"eplc/src/libepl/Output"
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

	Output.PrintLog(" <aplc> Parsing: '"+file+"' ")
	libepl.Compile(reader, file)
}
