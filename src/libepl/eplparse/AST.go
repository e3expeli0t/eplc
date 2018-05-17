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

import "eplc/src/libepl/epllex"

type AST struct {
	_type epllex.TokenType
	_value string
	parent *AST
	childrens []AST
}

func (an *AST) IsRoot() bool  {
	return an.parent == nil 
}

func (an *AST) IsData() bool {
	return len(an._value) == 0
}
