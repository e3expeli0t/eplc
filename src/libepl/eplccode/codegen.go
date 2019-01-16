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
	"eplc/src/libepl/Output"
	"eplc/src/libepl/epllex"
	"eplc/src/libepl/eplparse"
	"eplc/src/libepl/eplparse/ast"
	"fmt"

	"io"
	"reflect"
)


//todo: moe elegant way other then global var
var filename string

/*
	GenerateAIR generates AIR (AVM IR) for optimization and machine
	code generation by AVM
*/
func GenerateAIR(source io.Reader, fname string) {
	lexer := epllex.New(source, fname)
	filename = fname
	parser := eplparse.New(lexer)
	file := parser.ParseProgram()
	ast.Travel(file, generate)
}

func generate(n ast.Node) bool{

	switch n := n.(type) {
	case *ast.Program:
		genProgram(n)
		return true
	default:
		Output.PrintErr("codegen", "Unknown node type '", reflect.TypeOf(n), "' 	expected type ast.Program")
		return false
	}
}

func genProgram(program *ast.Program) {
	var index uint = 0
	writer := Writer{Fname:filename}
	writer.InitializeWriter()

	writer.UpdateLabels(genImport(program.Imports, &index))
	for _, decl := range program.Decls {
		fmt.Println(decl)
		writer.UpdateLabel(genDecls(decl, &index))
	}
	writer.produceST(program.Symbols)
	writer.WriteToTarget()
}

func genImport(node ast.Import, index *uint) []Label {
	var labels []Label
	for _, i := range node.Imports {
		labels = append(labels, CreateLabel(*index, "link "+i))
		*index++
	}

	return labels
}

func genDecls(node ast.Decl, index *uint) Label {
	switch n := node.(type) {
	case *ast.VarDecl:
		return genVarDecl(n, index)
	default:
		Output.PrintErr("codgen", "Unknown node type '", reflect.TypeOf(n), "'")
	}
	return Label{}
}

func genVarDecl(node *ast.VarDecl, index *uint) Label {
	return CreateLabel(*index, fmt.Sprintf("vardecl %s %x", node.Name, node.VarType.Tkey))
}
