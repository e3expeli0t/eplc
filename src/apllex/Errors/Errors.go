package Errors

import "fmt"

func TokenError(msg ...interface{}) {
	fmt.Println("Aplc_runtime: <Lexical>: ", msg)
	//panic("Lexical analysis error")
}
