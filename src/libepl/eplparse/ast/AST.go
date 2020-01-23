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

package ast

import (
	"eplc/src/libepl/eplparse/Types"
	"eplc/src/libepl/eplparse/symboltable"
)

/*
	AST (or Abstract Syntax Tree) is whats the parser producing
*/

type VarStat string

type (
	Node interface {
		Start() uint
	}

	Expression interface {
		Node
		ExprNode()
	}

	//Design note: experimental: every statement is an expression that return void
	Statement interface {
		Expression
		StmtNode()
	}

	Decl interface {
		Statement
		DeclNode()
	}
	BoolExpr interface {
		Expression
		BoolExprNode()
	}
)

//Todo: replace all string names in symbol table references
//Todo: loops
type (
	ProgramFile struct {
		FileName     string
		Symbols      *symboltable.SymbolTable
		Imports      *[]Import
		GlobalDecls  *[]Decl
		Functions    *[]Fnc
		MainFunction *Fnc
	}

	Fnc struct {
		Name       string
		ReturnType *Types.EplType
		Params     *[]Decl
		Body       *Block
	}

	Block struct {
		Symbols  *symboltable.SymbolTable
		ExprList *[]Expression
	}

	VarDecl struct {
		Name    string
		VarType *Types.EplType
		Stat    VarStat
	}

	//todo: IMPL in parser
	VarExplicitDecl struct {
		VarDecl
		Value *Expression
	}

	Import struct {
		StartLoc uint
		Imports  []string
	}

	AssignStmt struct {
		Owner Ident
		Value *Expression
	}

	IfStmt struct {
		Condition *BoolExpr
		Code      *Block
		Else      *Statement // to support if-else-if
	}

	ElseStmt struct {
		Code *Block
	}


	//todo: IMPL in v0.2++
	MoveLoop struct {}
)

func (b Block) ExprNode() {}

func (AssignStmt) Start() uint {
	panic("Invalid call")
}
func (AssignStmt) ExprNode() {}

func (ElseStmt) Start() uint {
	panic("Invalid call")
}

func (ElseStmt) ExprNode() {}
func (ElseStmt) StmtNode() {}

func (Fnc) Start() uint {
	panic("Invalid call")
}
func (Fnc) DeclNode() {}
func (Fnc) ExprNode() {}
func (Fnc) StmtNode() {}

func (vd VarDecl) Start() uint {
	return 0
}
func (VarDecl) StmtNode() {}
func (VarDecl) ExprNode() {}

func (ProgramFile) Start() uint { return 0 }
func (Block) Start() uint       { return 0xFAC } //Todo: Fix
func (IfStmt) Start() uint      { return 0 }

func (i *Import) Start() uint { return i.StartLoc }
func (VarDecl) DeclNode()     {}

func (VarExplicitDecl) DeclNode() {}
func (IfStmt) StmtNode(){}
func (IfStmt) ExprNode() {}

//----------------------------------------------------------------------------------------------------------------------
//Expressions

type (
	//todo: change bool expressions to refs
	Ident struct {
		Name string
	}

	EmptyExpr struct{}

	BoolNot struct {
		Expr BoolExpr
	}

	BoolAnd struct {
		Le BoolExpr
		Re BoolExpr
	}

	BoolOr struct {
		Le BoolExpr
		Re BoolExpr
	}

	BoolEquals struct {
		Le BoolExpr
		Re BoolExpr
	}

	BoolNotEquals struct {
		Le BoolExpr
		Re BoolExpr
	}
	
	BoolGreaterThen struct {
		Le BoolExpr
		Re BoolExpr
	}

	BoolGreatEquals struct {
		Le BoolExpr
		Re BoolExpr
	}

	BoolLowerThen struct {
		Le BoolExpr
		Re BoolExpr
	}

	BoolLowerThenEqual struct {
		Le BoolExpr
		Re BoolExpr
	}

	BinaryMul struct {
		Ls Expression
		Rs Expression
	}

	BinaryDiv struct {
		Ls Expression
		Rs Expression
	}

	BinaryAdd struct {
		Ls Expression
		Rs Expression
	}

	BinarySub struct {
		Ls Expression
		Rs Expression
	}

	UnaryPlus struct {
		Rs Expression
	}

	UnaryMinus struct {
		Rs Expression
	}

	FunctionCall struct {
		PackagePath []Ident
		Arguments   []Ident
		ReturnType  Types.EplType //todo: Version 0.2+
		FunctionName Ident
	}



Singular struct {
		Symbol Ident
	}

	//Represent number such as 5562 or 2
	Number struct {
		Value string
	}

	String struct {
		Value string
	}


)


func (b BoolOr) Start() uint {
	panic("implement me")
}

func (b BoolOr) ExprNode() {}
func (b BoolOr) BoolExprNode() {}


func (b BoolAnd) Start() uint {
	panic("implement me")
}

func (b BoolAnd) ExprNode() {}
func (b BoolAnd) BoolExprNode() {}

func (b BoolNotEquals) Start() uint {
	panic("Invalid call")
}

func (b BoolNotEquals) ExprNode() {}
func (c BoolNotEquals) BoolExprNode() {}

// both bool expr and expr because its can be boolean function
func (c FunctionCall) BoolExprNode() {}

func (s String) Start() uint {
	panic("Invalid call")
}

func (s String) ExprNode() {}
func (n Number) Start() uint {
	panic("Invalid call")
}

func (n Number) ExprNode() {}

func (BoolGreaterThen) Start() uint {
	panic("Invalid call")
}

func (BoolGreaterThen) ExprNode() {}
func (BoolGreaterThen) BoolExprNode() {}

func (BoolEquals) Start() uint {
	panic("Invalid call")
}

func (BoolEquals) ExprNode() {}
func (BoolEquals) BoolExprNode() {}

func (BoolGreatEquals) Start() uint {
	panic("Invalid call")
}

func (BoolGreatEquals) ExprNode() {}

func (BoolGreatEquals) BoolExprNode() {}

func (BoolNot) Start() uint {
	panic("Invalid call")
}

func (BoolNot) ExprNode() {}
func (BoolNot) BoolExprNode() {}

func (BoolLowerThen) Start() uint {
	panic("Invalid call")
}

func (BoolLowerThen) ExprNode() {}
func (BoolLowerThen) BoolExprNode() {}

func (BoolLowerThenEqual) Start() uint {
	panic("Invalid call")
}

func (BoolLowerThenEqual) ExprNode() {}
func (BoolLowerThenEqual) BoolExprNode() {}

func (Singular) Start() uint {
	panic("Invalid call")
}

func (Singular) StmtNode() {}
func (Singular) ExprNode() {}

func (FunctionCall) Start() uint {
	panic("Invalid call")
}

func (FunctionCall) StmtNode() {}
func (FunctionCall) ExprNode() {}

func (BinarySub) Start() uint {
	panic("Invalid call")
}
func (BinarySub) StmtNode() {}
func (BinarySub) ExprNode() {}

func (BinaryAdd) Start() uint {
	panic("Invalid call")
}
func (BinaryAdd) StmtNode() {}
func (BinaryAdd) ExprNode() {}

func (BinaryDiv) Start() uint {
	panic("Invalid call")
}
func (BinaryDiv) StmtNode() {}
func (BinaryDiv) ExprNode() {}

func (BinaryMul) Start() uint {
	panic("Invalid call")
}
func (BinaryMul) StmtNode() {}
func (BinaryMul) ExprNode() {}

func (UnaryPlus) Start() uint {
	panic("Invalid call")
}
func (UnaryPlus) ExprNode() {}
func (UnaryPlus) StmtNode() {}

func (UnaryMinus) Start() uint {
	panic("Invalid call")
}
func (UnaryMinus) ExprNode() {}
func (UnaryMinus) StmtNode() {}

func (EmptyExpr) Start() uint {
	panic("Invalid call")
}
func (EmptyExpr) StmtNode() {}
func (EmptyExpr) ExprNode() {}

func (Ident) Start() uint {
	panic("Invalid call")
}
func (Ident) StmtNode() {}
func (Ident) ExprNode() {}
