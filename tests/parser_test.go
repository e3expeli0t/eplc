/*
*	Copyright (C) 2018-2020 Elia Ariaz
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

package tests

import (
	"eplc/src/libepl/epllex"
	"eplc/src/libepl/eplparse"
	"eplc/src/libepl/eplparse/ast"
	"strings"
	"testing"
)

var (
	currentInput string
	lastError string
)


func matchRes(output ast.Node, input ast.Node) {
	switch t := output.(type) {
	case ast.ProgramFile:
		if !(t == input) {

		}
	}
	
}

func TestParser_ParseBlock(t *testing.T) {

}

func TestParser_ParseExpression(t *testing.T) {

}

func TestParser_ParseFnc(t *testing.T) {
	lx := epllex.New(strings.NewReader("fnc exec(command uint, command string, args string): Proc {\n}"), "Test")
	p := eplparse.New(lx)

}

func TestParser_ParseIdent(t *testing.T) {

}

func TestParser_ParseImport(t *testing.T) {

}

func TestParser_ParseParamList(t *testing.T) {

}

func TestParser_ParseUnaryOp(t *testing.T) {

}

func TestParser_ParseVarDecl(t *testing.T) {

}

func TestParser_ParseProgramFile(t *testing.T) {

}
