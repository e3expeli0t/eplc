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

package analysis

import (
	"eplc/src/libepl"
	"eplc/src/libepl/Types"
	"eplc/src/libepl/eplparse/ast"
	"eplc/src/libepl/eplparse/symboltable"
	"eplc/src/libio"
	"fmt"
)

func NewAnalyzer(syntaxTree ast.Node) *Analysis {
	return &Analysis{
		AST:         syntaxTree,
		phase:   libepl.Analysis,
	}
}

type Analysis struct {
	AST       ast.Node
	SymbolMap *symboltable.TableMap
	TChecker  *TypeChecker
	phase     libepl.PhaseIndicator
	FileInfo  libepl.InfoStruct
}

func (ana* Analysis) Init() {
	switch n := ana.AST.(type) {
	case *ast.ProgramFile:
		ana.SymbolMap = n.SymbolTableMap
		//Output.PrintFatalErr(ana.SymbolMap.ToString())
		ana.FileInfo = libepl.NewInfoStruct(n.FileName)
	default:
		panic(fmt.Sprintf("Couldn't initialize Analyser. Node %s is not recognized", n))
	}

 	ana.TChecker = NewTypeChecker(ana.SymbolMap)
}

func (ana* Analysis) Run() {
	ana.phase = libepl.TypeChecker
	libio.PrintLog("Starting type analysis")

	ast.Travel(ana.AST, ana.TravelAST)
}

//Visitor function. Preforms basic type analysis and other checks
func (ana* Analysis) TravelAST(node ast.Node) bool {
	switch n := node.(type) {
	case *ast.ProgramFile:
		//todo: handle st frame change
		ana.TravelAST(n.MainFunction)
		for _,function := range *n.Functions {
			ana.TravelAST(function)
		}

	case *ast.IfStmt:
		//todo: handle st frame change
		t := ana.TChecker.WalkExpression(*n.Condition)

		if t != Types.TypeBool.AsEplType() {
			if ana.TChecker.HasErrors() {
				ana.PrintError(ana.DumpErrors(ana.TChecker.Errors))
			} else {
				ana.PrintError("Expected boolean expression got expression of type:" + t.TypeName)
			}

			return false
		}

	case *ast.ForLoop:
		//todo: handle st frame change
		t := ana.TChecker.WalkExpression(*n.Condition)

		if t != Types.TypeBool.AsEplType() {
			if ana.TChecker.HasErrors() {
				ana.PrintError(ana.DumpErrors(ana.TChecker.Errors))
			} else {
				ana.PrintError("Expected boolean expression got expression of type:" + t.TypeName)
			}

			return false
		}
	case *ast.RepeatUntil:
		//todo: handle st frame change
		t := ana.TChecker.WalkExpression(*n.Condition)

		if t != Types.TypeBool.AsEplType() {
			if ana.TChecker.HasErrors() {
				ana.PrintError(ana.DumpErrors(ana.TChecker.Errors))
			} else {
				ana.PrintError("Expected boolean expression got expression of type:" + t.TypeName)
			}

			return false
		}
	}

	return true
}


func (ana Analysis) DumpErrors(errors []*TypeError) string {
	var errorString string

	for i, e := range errors {
		errorString += fmt.Sprintf("\t%d: %s\n", i, e.Descriptor)
	}

	return errorString
}

//Very basic  error function
//todo: implement IntelligentError
func (ana* Analysis) PrintError(err string) {
	errorString := "While performing: "
	switch ana.phase {
	case libepl.TypeChecker:
		errorString = "Type analysis. "
	default:
		errorString += "Analysis. "
	}
	errorString += "\nGot:\n"+err
	libio.PrintErr(errorString)
}