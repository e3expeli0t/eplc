package eplinter

import (
	"eplc/src/libepl/epllex"
	"eplc/src/libepl/eplparse"
	"eplc/src/libepl/eplparse/ast"
)

type ExecType uint

const (
	COMPILE     ExecType = iota
	PARSE       ExecType = iota
	SYNTAXCHECK ExecType = iota
	EXECUTE     ExecType = iota
)

type ExecUnit struct {
	CurrentExCommand Command
	ExType           ExecType
	Parser           *eplparse.Parser
	Lexer            *epllex.Lexer
	Ast              *ast.Node
}

//todo: Add full independent language parsing (which will lead to python like language)
func (e *ExecUnit) parse() {
	e.Ast = e.Parser.ParseProgramFile()
}

func (e *ExecUnit) run() {
	switch e.ExType {
	case PARSE:

	}
}
