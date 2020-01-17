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
		Else      *Block
	}

	ElseStmt struct {
		Code *Block
	}
)

func (b Block) ExprNode() {
	panic("implement me")
}

func (AssignStmt) Start() uint {
	panic("implement me")
}
func (s AssignStmt) ExprNode() {
	panic("implement me")
}

func (Fnc) ExprNode() {
	panic("implement me")
}

func (ElseStmt) Start() uint {
	panic("implement me")
}

func (ElseStmt) ExprNode() {
	panic("implement me")
}

func (ElseStmt) StmtNode() {
	panic("implement me")
}

func (VarDecl) ExprNode() {
	panic("implement me")
}

func (Fnc) DeclNode() {}

func (Fnc) Start() uint {
	panic("implement me")
}

func (Fnc) StmtNode() {}

func (VarDecl) StmtNode() {}

func (VarDecl) Start() uint {
	return 0
}

func (ProgramFile) Start() uint { return 0 }
func (Block) Start() uint       { return 0xFAC } //Todo: Fix
func (IfStmt) Start() uint      { return 0 }

func (i *Import) Start() uint { return i.StartLoc }
func (VarDecl) DeclNode()     {}

func (VarExplicitDecl) DeclNode() {}
func (IfStmt) StmtNode()          {}
func (IfStmt) ExprNode() {
	panic("implement me")
}

//----------------------------------------------------------------------------------------------------------------------
//Expressions

type (
	Ident struct {
		Name string
	}

	EmptyExpr struct{}

	BoolNot struct {
		Expr *BoolExpr
	}

	BoolEquals struct {
		Le *BoolExpr
		Re *BoolExpr
	}

	BoolGreaterThen struct {
		Le *BoolExpr
		Re *BoolExpr
	}

	BoolGreatEquals struct {
		Le *BoolExpr
		Re *BoolExpr
	}

	BoolLowerThen struct {
		Le *BoolExpr
		Re *BoolExpr
	}

	BoolLowerThenEqual struct {
		Le *BoolExpr
		Re *BoolExpr
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
	}

	Singular struct {
		Symbol Ident
	}

	//Represent number such as 5562 or 2
	Number struct {
		Value string
	}
)

func (n Number) Start() uint {
	panic("implement me")
}

func (n Number) ExprNode() {
	panic("implement me")
}

func (BoolGreaterThen) Start() uint {
	panic("implement me")
}

func (BoolGreaterThen) ExprNode() {
	panic("implement me")
}

func (BoolGreaterThen) BoolExprNode() {
	panic("implement me")
}

func (BoolEquals) Start() uint {
	panic("implement me")
}

func (BoolEquals) ExprNode() {
	panic("implement me")
}

func (BoolEquals) BoolExprNode() {
	panic("implement me")
}

func (BoolGreatEquals) Start() uint {
	panic("implement me")
}

func (BoolGreatEquals) ExprNode() {
	panic("implement me")
}

func (BoolGreatEquals) BoolExprNode() {
	panic("implement me")
}

func (BoolNot) Start() uint {
	panic("implement me")
}

func (BoolNot) ExprNode() {
	panic("implement me")
}

func (BoolNot) BoolExprNode() {
	panic("implement me")
}

func (BoolLowerThen) Start() uint {
	panic("implement me")
}

func (BoolLowerThen) ExprNode() {
	panic("implement me")
}

func (BoolLowerThen) BoolExprNode() {
	panic("implement me")
}

func (BoolLowerThenEqual) Start() uint {
	panic("implement me")
}

func (BoolLowerThenEqual) ExprNode() {
	panic("implement me")
}

func (BoolLowerThenEqual) BoolExprNode() {
	panic("implement me")
}

func (Singular) Start() uint {
	panic("implement me")
}

func (Singular) StmtNode() {
	panic("implement me")
}

func (Singular) ExprNode() {
	panic("implement me")
}

func (FunctionCall) Start() uint {
	panic("implement me")
}

func (FunctionCall) StmtNode() {
	panic("implement me")
}

func (FunctionCall) ExprNode() {
	panic("implement me")
}

func (BinarySub) Start() uint {
	panic("implement me")
}
func (BinarySub) StmtNode() {}
func (BinarySub) ExprNode() {}

func (BinaryAdd) Start() uint {
	panic("implement me")
}
func (BinaryAdd) StmtNode() {}
func (BinaryAdd) ExprNode() {}

func (BinaryDiv) Start() uint {
	panic("implement me")
}
func (BinaryDiv) StmtNode() {}
func (BinaryDiv) ExprNode() {}

func (BinaryMul) Start() uint {
	panic("implement me")
}
func (BinaryMul) StmtNode() {}
func (BinaryMul) ExprNode() {}

func (UnaryPlus) Start() uint {
	panic("implement me")
}
func (UnaryPlus) ExprNode() {}
func (UnaryPlus) StmtNode() {}

func (UnaryMinus) ExprNode() {}
func (UnaryMinus) Start() uint {
	panic("Good")
}
func (UnaryMinus) StmtNode() {}

func (EmptyExpr) Start() uint {
	panic("implement me")
}
func (EmptyExpr) StmtNode() {}
func (EmptyExpr) ExprNode() {}

func (Ident) Start() uint {
	panic("implement me")
}
func (Ident) StmtNode() {}
func (Ident) ExprNode() {}
