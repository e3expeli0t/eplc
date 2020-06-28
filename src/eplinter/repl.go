package main

import "eplc/src/eplinter/internal"

func main() {
	cmd := internal.NewCommandLine("REPL")
	cmd.Run()
}
