/*
*	Copyright (C) 2018-2020 Elia Ariaz
*
*	This program is free software: you can redistribute it and/or modify
*	it under the terms of the GNU General Public License as published by
*	the Free Software Foundation, either version 3 of the License, or
*	(at your option) any later version.
*
*	This program is distributed in the hope that it will be useful,
*	but WITHOUT ANY WARRANTY; without even the implied warranty of
*	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
*	GNU General Public License for more details.
*
*	You should have received a copy of the GNU General Public License
*	along with this program. If not, see <http://www.gnu.org/licenses/>.
 */

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