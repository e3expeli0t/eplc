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
type AST struct {
	nodes []Node
}

type Node struct {
	_type NType
	_value string
	left *Node
	right *Node
	parent *Node
}

type NType int 

const (
	FNCDECL NType 		= iota //fnc id()[:type] [-> id]
	FNCIMPL NType	 	= iota 
	VARDECL NType 		= iota
	ASSIGN NType 		= iota
	MATHOP 				= iota
	IF		 			= iota
	ADD NType 			= iota
	BLOCK NType 		= iota
	BOOL_EXPR NType 	= iota
	PROGRAM 			= iota
)


//Check if the 
func (an *AST) IsRoot(n *Node) bool  {
	return n.parent == nil 
}

func (an *AST) IsData(n *Node) bool {
	return len(n._value) == 0
}
