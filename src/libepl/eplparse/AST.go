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

package eplparse

//import "eplc/src/libepl/epllex"

/*
	AST (or Abstract Syntax Tree) is whats the parser producing
*/

type Node interface {
	Start()
	End()
}


type VarStat string 

const (
	FixedVar VarStat = "**FIXED**"
)

type ( 
	

	Block struct {

<<<<<<< HEAD
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
=======
func (an *AST) IsData(n *Node) bool {
	return len(n._value) == 0
}

type IF struct {
	ST *SymbolTable
	condition BOOL_EXPR
}
>>>>>>> 9dd32b163eca093275fa40b9d7f69814ffc07d47
