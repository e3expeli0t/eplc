package ast

import (
	"fmt"
	"reflect"
)

func printer(n Node) bool {

	if n != nil {
		fmt.Println(reflect.TypeOf(n))
	}

	switch t := n.(type) {
	case String:
		fmt.Println(t.Value)
	case Number:
		fmt.Println(t.Value)
	case Boolean:
		if t.Val == BOOL_FALSE {
			fmt.Println("\tfalse")
		} else {
			fmt.Println("\ttrue")
		}
	case *BinaryMul:
		fmt.Print("*")
	case *BinaryDiv:
		fmt.Print("/")
	case *BinarySub:
		fmt.Print("-")
	case *BinaryAdd:
		fmt.Print("+")
	case *Ident:
		fmt.Println("\t"+t.Name)
	}
	return true
}

func PrintNode(n Node) {
	Travel(n, printer)
}