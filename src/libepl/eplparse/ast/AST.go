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

//import "eplc/src/libepl/epllex"

/*
	AST (or Abstract Syntax Tree) is whats the parser producing
*/

type Node interface {
	Start() uint
}

type VarStat string 

type (

	Block struct {
		nodes []Node
	}
	
	VarDecl struct {
		name string
		stat VarStat 
	}

 	VarExplicitDecl struct {
		VarDecl
		value string
	}

	Import struct {
		StartLoc uint
		Imports []string
	}

	BoolExpr struct {}

	IfElseStatment struct {
		Condition BoolExpr
		Code *Block
		Else *Block	
	}

	IfStatement struct {
		Condition BoolExpr
		Code *Block
	}
	IfElseIfStatement struct{
		Stmts []IfElseIfStatement
	}
)

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
func (i *IfElseStatment) Start() {
}
func (i *IfElseIfStatement) Start() {
}
func (i *IfStatement) Start() {
}

