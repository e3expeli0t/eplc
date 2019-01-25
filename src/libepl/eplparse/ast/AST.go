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

	Statement interface {
		Node
		StmtNode()
	}

	Decl interface {
		Statement
		DeclNode()
	}

	Expression interface {
		Statement
		ExprNode()
	}


)

//Todo: replace all string names in symbol table references
type (

	Program struct {
		Symbols *symboltable.SymbolTable
		Imports *Import
		Decls *[]Decl
		Functions *[]Fnc
		MainFunction *Fnc
	}

	Fnc struct {
		Name  string
		ReturnType 	*Types.EplType
		Body *Block
	}

	Block struct {
		Symbols *symboltable.SymbolTable
		Nodes *[]Node
	}
	
	VarDecl struct {
		Name string
		VarType *Types.EplType
		Stat *VarStat
	}

 	VarExplicitDecl struct {
		VarDecl
		Value *Expression
	}


	Import struct {
		StartLoc uint
		Imports []string
	}

	AssignStmt struct {
		Owner string
		Value *Expression
	}


	IfStmt struct {
		Condition *BoolExpr
		Code *Block
		Else *Block
	}

	ElseStmt struct {
		Code *Block
	}

	BoolExpr struct {}
)

func (VarDecl) StmtNode() {
	panic("implement me")
}

func (VarDecl) Start() uint {
	panic("implement me")
}

func (*Program) Start() uint {return 0}
func (*Block) Start() uint{return 0xFAC} //Todo: Fix
func (*IfStmt) Start() uint {return 0}

func (i *Import) Start() uint {return i.StartLoc}
func (*VarDecl) DeclNode() {}

func (*VarExplicitDecl) DeclNode() {}
func (*IfStmt) StmtNode() {}
