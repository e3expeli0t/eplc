package main

import (
	"eplc/src/libepl"
	"eplc/src/libepl/Output"
//	"eplc/src/libepl/epllex"
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	var quiet bool
	var debug bool
	var Optimize int
	var DestFile string

	if len(os.Args) <= 1 {
		Output.PrintErr("please supply file name")
	}

	args := os.Args[1:]
	file := args[0]

	flag.BoolVar(&quiet, "quiet", true, "Prints start massge and print program compilation status (the defualt is false)")
	flag.BoolVar(&debug, "debug", false, "Prints debugging information and adds symbol table to the binary (the defualt is false)")
	flag.IntVar(&Optimize, "opt", 3, "Optimization level (the defualt is 3")
	flag.StringVar(&DestFile, "o", "a.out", "Name the target binary (defualt is a.out)")

	flag.Usage = func() {
		fmt.Printf("Usage: eplc <filename> [flags]:")

		flag.PrintDefaults()
	}

	flag.Parse()

	if !quiet {
		Output.PrintStartMSG()
	}

	var reader io.Reader
	reader, err := os.Open(file)

	if err != nil {
		Output.PrintErr("file: '" + file + "' don't exists")
	}

	if !quiet {
		Output.PrintLog(" <aplc> Parsing: '" + file + "' ")
	}
/*
	lx := epllex.New(reader, "name.epl")
	var tmp = lx.Next()

	for tmp.Ttype != epllex.EOF {
		Output.PrintLog(tmp.Lexme)
		tmp = lx.Next()
	}
	*/
	libepl.Compile(reader, file)
}
