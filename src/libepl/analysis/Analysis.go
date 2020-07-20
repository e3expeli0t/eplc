package analysis

import (
	"eplc/src/libepl"
	"eplc/src/libepl/Output"
	"eplc/src/libepl/Types"
	"eplc/src/libepl/eplparse/ast"
	"eplc/src/libepl/eplparse/symboltable"
	"fmt"
)

func NewAnalyzer(syntaxTree ast.Node) *Analysis {
	return &Analysis{
		AST:         syntaxTree,
		phase:   libepl.Analysis,
	}
}

type Analysis struct {
	AST ast.Node
	SymbolTable *symboltable.SymbolTable
	TChecker *TypeChecker
	phase libepl.PhaseIndicator
	FileInfo libepl.InfoStruct
}

func (ana* Analysis) Init() {
	switch n := ana.AST.(type) {
	case *ast.ProgramFile:
		ana.SymbolTable= n.Symbols
		ana.FileInfo = libepl.NewInfoStruct(n.FileName)
	default:
		Output.PrintFatalErr(fmt.Sprintf("Couldn't initialize Analyse. Node %s is not recognized", n))
	}

 	ana.TChecker = NewTypeChecker(ana.SymbolTable)
}

func (ana* Analysis) Run() {
	ana.phase = libepl.TypeChecker
	Output.PrintLog("Starting type analysis")

	ast.Travel(ana.AST, ana.TravelAST)
}

//Visitor function. Preforms basic type analysis and other checks
func (ana* Analysis) TravelAST(node ast.Node) bool {
	switch n := node.(type) {
	case *ast.IfStmt:
		t := ana.TChecker.WalkExpression(*n.Condition)
		if t != Types.TypeBool.AsEplType() {
			if ana.TChecker.HasErrors() {
				ana.PrintError(ana.DumpErrors(ana.TChecker.Errors))
			} else {
				ana.PrintError("Expected boolean expression got expression of type:" + t.TypeName)
			}
			return false
		}
	case ast.ForLoop:
		t := ana.TChecker.WalkExpression(*n.Condition)
		if t != Types.TypeBool.AsEplType() {
			if ana.TChecker.HasErrors() {
				ana.PrintError(ana.DumpErrors(ana.TChecker.Errors))
			} else {
				ana.PrintError("Expected boolean expression got expression of type:" + t.TypeName)
			}
			return false
		}
	case ast.RepeatUntil:
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
	Output.PrintErr(errorString)
}