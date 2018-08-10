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

package epllex

import (
	"fmt"
	"testing"
	"strings"
)


func checkType(lexme string, tt TokenType, start uint, end uint) bool {
	lx := New(strings.NewReader(lexme), "test_lexer.epl")
	tk := lx.Next()
	fmt.Printf("type: %d, lexme: %s, start: %d, end:%d \r", tk.Ttype, tk.Lexme, tk.Start, tk.End)
	fmt.Println()
	return tk.End == end && tk.Start == start && tk.Ttype == tt
}


func TestLexer(t *testing.T) {
	table := []struct {
		inToken string
		outType TokenType
		start uint
		end uint
	}{
		{"23", NUM, 0, 1},
		{"2.3", REAL, 1, 3},
		{"HELLO3433", ID, 0, 9},
		{"h", ID, 0, 1},
		{"@MainFunc", CFLAG, 0, 9},
		{"//commant\n/*More commant*/",EOF, 0, 0},
		{"\n\n\"string\"", STRINGLITERAL, 3, 7},
		{"\n'string'", STRINGLITERAL, 2, 7},
		{"if", IF, 0, 2},
		{"else", ELSE, 0, 4},
		{"repeat", REPEAT, 0, 6},
		{"until", UNTIL, 0, 5},
		{"move", MOVE, 0, 4},
		{"in", IN, 0, 2},
		{"static", STATIC, 0, 6},
		{"fixed", FIXED, 0, 5},
		{"decl", DECL, 0, 4},
		{"define", DEFINE, 0, 6},
		{"type", TYPE, 0, 4},
		{"struct", STRUCT, 0, 6},
		{"fnc", FNC, 0, 3},
		{"+", PLUS, 0, 1},
		{"-", MINUS, 0, 1},
		{"/", DEV, 0, 1},
		{"*", MULT, 0, 1},
		{"{", LBRACE, 0, 1},
		{"}", RBRACE, 0, 1},
		{"(", LPAR, 0, 1},
		{")", RPAR, 0, 1},
		{"[", LSUBSCRIPT, 0, 1},
		{"]", RSUBSCRIPT, 0, 1},
		{"+=", PLUSEQUAL, 0, 2},
		{"-=", MINUSEQUAL, 0, 2},
		{"*=", MULTEQUAL, 0, 2},
		{"/=", DEVEQUAL, 0, 2},
		{"=", ASSIGN, 0, 1},
		{">", GT, 0, 1},
		{"<", LT, 0, 1},
		{"==", EQ, 0, 2},
		{">=", GE, 0, 2},
		{"<=", LE, 0, 2},
		{"&&", AND, 0, 2},
		{"||", OR, 0, 2},
		{"!", NOT, 0, 1},
		{"|", UNARYOR, 0, 1},
		{"&", UNARYAND, 0, 1},
		{"<<", LSHIFT, 0, 2},
		{">>", RSHIFT, 0, 2},
		{";", SEMICOLON, 0, 1},
		{":", RETURN_IND, 0, 1},
		{",", COMMA, 0, 1},
		{".", DOT, 0, 1},
		{"int", INT, 0, 3},
		{"int16", INT16, 0, 5},
		{"int32", INT32, 0, 5},
		{"int64", INT64, 0, 5},
		{"float16", FLOAT16, 0, 7},
		{"float64", FLOAT64, 0, 7},
		{"float64", FLOAT64, 0, 7},
		{"float", FLOAT, 0, 5},
		{"long", LONG, 0, 4},
		{"cmx64", CMX64, 0, 5},
		{"cmx", CMX, 0, 3},
		{"uint16", UINT16, 0, 6},
		{"uint32", UINT32, 0, 6},
		{"uint64", UINT64, 0, 6},
		{"uint", UINT, 0, 4},
		{"bool", BOOL, 0, 4},
		{"true", TRUE, 0, 4},
		{"false", FALSE, 0, 4},
		{"import", IMPORT, 0, 6},
		{"string", STRING, 0, 6},
	}

	for _, tst := range table{
		if !checkType(tst.inToken, tst.outType, tst.start, tst.end) {
			t.Errorf("[!] Couldn't match token '%s' to type %d", tst.inToken, tst.outType)
			break
		}
	}
}