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

package eplccode

import (
	"eplc/src/libepl/Output"
	"eplc/src/libepl/eplparse/ast"
	"fmt"
	"reflect"
)

//todo: Make this more efficient
//todo: Automated jump locations allocation
/*
	GenerateAIR generates AIR (AVM IR) for optimization and machine
	code generation by AVM
*/

func GenerateAIR(file ast.Node) {
	generate(file)
}

func generate(n ast.Node) bool {
	switch n := n.(type) {
	case *ast.ProgramFile:
		Output.PrintLog("Generating AIR")
		genProgram(n)
		return true
	default:
		Output.PrintErr("Unknown node type '", reflect.TypeOf(n), "' 	expected type ast.ProgramFile")
		return false
	}
}

func genProgram(program *ast.ProgramFile) {
	var index uint = 0
	writer := Writer{Fname: program.FileName}
	writer.InitializeWriter()

	writer.UpdateLabels(genImport(*program.Imports, &index))
	for _, decl := range *program.GlobalDecls {
		writer.UpdateLabels(genDecls(decl, &index))
		Output.PrintLog("decl generated")
	}
	for _, fnc := range *program.Functions {
		writer.UpdateLabels(genDecls(fnc, &index))
		fmt.Println("Added new Function decl")
	}
	if *program.MainFunction != (ast.Fnc{}) {
		writer.UpdateLabels(genDecls(program.MainFunction, &index))
	} else {
		Output.PrintFatalErr("Couldn't find main function")
	}

	writer.ProduceST(program.Symbols)
	writer.ProduceAST(program)
	writer.WriteToTarget()
}

/*
	Generates link command to link external libraries to the program
*/
func genImport(node []ast.Import, index *uint) []Label {
	var labels []Label
	for _, i := range node {
		for _, lib := range *i.Imports {
			labels = append(labels, CreateLabel(*index, "link", lib, ""))
			*index++
		}
	}

	return labels
}

func genBlock(index *uint, block *ast.Block) []Label {
	var labels []Label
	labels = append(labels, CreateLabel(*index, "$!!", "$[#]", "$[#]")) //Empty command
	return labels
}

//----------------------------------------------------------------------------------------------------------------------
//Decls
/*
	Generate vardecl statements
*/
func genDecls(node ast.Decl, index *uint) []Label {
	switch n := node.(type) {
	case *ast.VarDecl:
		return []Label{genVarDecl(n, index)}
	case *ast.VarExplicitDecl:
		var labels []Label
		labels = append(labels, genVarDecl(&ast.VarDecl{n.Name, n.VarType, n.Stat}, index))
		labels = append(labels, genAssignStmt(n.Name.Name, n.Value, index)...)
		return labels
	case ast.Fnc:
		return genFncDecl(&n, index)
	case *ast.Fnc:
		return genFncDecl(n, index)
	default:
		Output.PrintErr("Unknown node type '", reflect.TypeOf(n), "'")
	}
	return []Label{}
}

func genFncDecl(node *ast.Fnc, index *uint) []Label {
	var labels []Label
	labels = append(labels, CreateLabel(*index, "fncdecl", node.Name.Name, node.ReturnType.TypeName))
	*index++
	labels = append(labels, genParamList(index, DeclsDeepConvert(node.Params))...)
	labels = append(labels, genJump(*index+2, index))
	*index++
	labels = append(labels, genBlock(index, node.Body)...)
	*index++
	return labels
}

func genVarDecl(node *ast.VarDecl, index *uint) Label {
	return CreateLabel(*index, "vardecl", node.Name.Name, node.VarType.TypeName)
}

//-----------------------------------------------------------------------------------------------------------------------
//Statements

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
//Helpers

func genParamList(index *uint, list *[]ast.VarDecl) []Label {
	var labels []Label

	for _, param := range *list {
		labels = append(labels, CreateLabel(*index, "DeclParam", param.Name.Name, param.VarType.TypeName))
		*index++
	}

	return labels
}

//----------------------------------------------------------------------------------------------------------------------
// Support functions

func DeclsDeepConvert(decls *[]ast.Decl) *[]ast.VarDecl {
	var converted []ast.VarDecl

	for _, decl := range *decls {
		switch n := decl.(type) {
		case *ast.VarDecl:
			converted = append(converted, *n)
		default:
			break
		}
	}

	return &converted
}

func asString(n uint) string {
	return fmt.Sprintf("%d", n)
}
