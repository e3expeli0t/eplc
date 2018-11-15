/*
*	eplc
*	Copyright (C) 2018 eplc core team
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

package eplccode

import (
	"eplc/src/libepl/epllex"
	"eplc/src/libepl/eplparse"
	"eplc/src/libepl/eplparse/ast"

	"io"
	"fmt"
)

var label uint = 0

/*
	GenerateAIR generates AIR (AVM IR) for optimization and machine 
	code generation by AVM
*/
func GenerateAIR(source io.Reader, fname string) {
	lexer := epllex.New(source, fname)
	parser := eplparse.New(lexer)
	genImport(parser.ParseImport())
}

func genImport(node ast.Import) {
	for _, i := range node.Imports {
		fmt.Printf("L%d: link %s\n", label, i);
		label++;
	} 
}
