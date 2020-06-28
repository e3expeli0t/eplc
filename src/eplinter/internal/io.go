package internal

import (
	"eplc/src/libepl/Output/color"
	"fmt"
)

func (cl *CommandLine) Print(input...interface{}) {
	fmt.Print(cl.prompt)
	fmt.Print(input...)
}

func (cl *CommandLine) Println(input...interface{}) {
	fmt.Print(cl.prompt)
	fmt.Println(input...)
}
func (cl *CommandLine) PrintError(msg...interface{}) {
	fmt.Print(color.Blink(color.GLightRed("!!")))
	fmt.Print(fmt.Sprintf(color.GLightRed("Error:{}:{} "), cl.lex.Filename, cl.lex.Line))
	fmt.Println(msg...)
}