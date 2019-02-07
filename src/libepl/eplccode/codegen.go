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


//todo: Make this more efficient
/*
	GenerateAIR generates AIR (AVM IR) for optimization and machine
	code generation by AVM
*/
func GenerateAIR(source io.Reader, fname string) {
	lexer := epllex.New(source, fname)
	parser := eplparse.New(lexer)
	file := parser.ParseProgramFile()
	generate(file)
}

func generate(n ast.Node) bool {

	switch n := n.(type) {
	case *ast.ProgramFile:
		genProgram(n)
		return true
	default:
		Output.PrintErr("Unknown node type '", reflect.TypeOf(n), "' 	expected type ast.ProgramFile")
		return false
	}
}

func genProgram(program *ast.ProgramFile) {
	var index uint = 0
	writer := Writer{Fname:program.FileName}
	writer.InitializeWriter()

	writer.UpdateLabels(genImport(program.Imports, &index))
	for _, decl := range *program.GlobalDecls {
		fmt.Println(decl)
		writer.UpdateLabel(genDecls(decl, &index))
	}
	writer.produceST(program.Symbols)
	writer.WriteToTarget()
}

/*
	Generates link command to link external libraries to the program
 */
func genImport(node *ast.Import, index *uint) []Label {
	var labels []Label
	for _, i := range node.Imports {
		labels = append(labels, CreateLabel(*index, "link", i, ""))
		*index++
	}

	return labels
}

/*
	Generate vardecl statements
 */
func genDecls(node ast.Decl, index *uint) Label {
	switch n := node.(type) {
	case *ast.VarDecl:
		return genVarDecl(n, index)
	case *ast.VarExplicitDecl:
		var labels []Label
		labels = append(labels, genVarDecl(&ast.VarDecl{n.Name, n.VarType,n.Stat}, index))
		labels = append(labels, genAssignStmt(n.Name, n.Value, index)...)
	default:
		Output.PrintErr("Unknown node type '", reflect.TypeOf(n), "'")
	}
	return Label{}
}


/*
	Jump to runtimeResolve Location and move resReg to the variable
 */
func genAssignStmt(varname string, expression *ast.Expression, index *uint) []Label {
	var labeles []Label
	labeles = append(labeles, genJump(*index+2, index))
	*index++
	*index++
	labeles = append(labeles, genMove("$[resReg]", varname, index))
	*index++
	labeles = append(labeles, genExpression(expression, index))
	*index++

	return labeles
}

func genVarDecl(node *ast.VarDecl, index *uint) Label {
	return CreateLabel(*index, "vardecl", node.Name, node.VarType.Tname)
}

func genExpression(expression *ast.Expression, index *uint) Label {
	return CreateLabel(*index, "runtimeResolve", asString(*index), "")
}
//----------------------------------------------------------------------------------------------------------------------
//Simple IS


func genJump(loc uint, index *uint) Label {
	return CreateLabel(*index, "jump", asString(loc), "")
}

func genMove(loc string, dest string, index *uint) Label {
	return CreateLabel(*index, "move", loc, dest)
}


//----------------------------------------------------------------------------------------------------------------------
// Support functions

func asString(n uint) string {
	return fmt.Sprintf("%d", n)
}