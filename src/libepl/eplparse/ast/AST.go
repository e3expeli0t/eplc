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

type Node interface {
	Start() uint
}

type VarStat string

type (

	Program struct {
		Symbols symboltable.SymbolTable
		Imports *Import
		Decls []DeclStmt
		Functions []Fnc
		MainFunction *Fnc
	}

	DeclStmt interface {
		Node
	}

	Fnc struct {
		name  string
		returnType 	Types.EplType
		body *Block
	}

	Block struct {
		symbols symboltable.SymbolTable
		nodes []Node
	}
	
	VarDecl struct {
		DeclStmt
		name string
		varType Types.EplType
		stat VarStat 
	}

 	VarExplicitDecl struct {
 		DeclStmt
		VarDecl
		value string
	}

	Import struct {
		StartLoc uint
		Imports []string
	}

	BoolExpr struct {}

	IfElseStmt struct {
		Condition BoolExpr
		Code *Block
		Else *Block	
	}
)

func (Program) Start() uint {
	panic("implement me")
}

func (b *Block) Start() {
}
func (v *VarDecl) Start() {
}
func (v *VarExplicitDecl) Start() {
}
func (i *Import) Start() uint {
	return i.StartLoc
}
func (b *BoolExpr) Start() {
}

