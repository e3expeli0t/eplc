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
	"strings"
	"testing"
)

func checkType(lexme string, tt TokenType, startLine uint, startOffset uint) bool {
	lx := New(strings.NewReader(lexme), "test_lexer.epl")
	tk := lx.Next()
	//fmt.Printf("type: %d, lexme: %s, start Line: %d, start offset:%d \n", tk.Ttype, tk.Lexme, tk.StartLine, tk.StartOffset)
	return tk.StartLine == startLine && tk.StartOffset == startOffset && tk.Ttype == tt
}

func TestLexer(t *testing.T) {
	table := []struct {
		inToken     string
		outType     TokenType
		startLine   uint
		startOffset uint
	}{
		{"23", NUM, 0, 0},
		{"2.3", REAL, 0, 0},
		{"HELLO3433", ID, 0, 0},
		{"h", ID, 0, 0},
		{"@MainFunc", CFLAG, 0, 0},
		{"//comment\n/*More comment*/", EOF, 0, 0},
		{"\n\n\"string\"", STRINGLITERAL, 2, 0},
		{"\n'string'", STRINGLITERAL, 1, 0},
		{"if", IF, 0, 0},
		{"else", ELSE, 0, 0},
		{"repeat", REPEAT, 0, 0},
		{"until", UNTIL, 0, 0},
		{"move", MOVE, 0, 0},
		{"in", IN, 0, 0},
		{"static", STATIC, 0, 0},
		{"fixed", FIXED, 0, 0},
		{"decl", DECL, 0, 0},
		{"define", DEFINE, 0, 0},
		{"type", TYPE, 0, 0},
		{"struct", STRUCT, 0, 0},
		{"fnc", FNC, 0, 0},
		{"+", ADD, 0, 0},
		{"-", SUB, 0, 0},
		{"/", DIV, 0, 0},
		{"*", MULT, 0, 0},
		{"{", LBRACE, 0, 0},
		{"}", RBRACE, 0, 0},
		{"(", LPAR, 0, 0},
		{")", RPAR, 0, 0},
		{"[", LSUBSCRIPT, 0, 0},
		{"]", RSUBSCRIPT, 0, 0},
		{"+=", PLUSEQUAL, 0, 0},
		{"-=", MINUSEQUAL, 0, 0},
		{"*=", MULTEQUAL, 0, 0},
		{"/=", DEVEQUAL, 0, 0},
		{"=", ASSIGN, 0, 0},
		{">", GT, 0, 0},
		{"<", LT, 0, 0},
		{"==", EQ, 0, 0},
		{">=", GE, 0, 0},
		{"<=", LE, 0, 0},
		{"&&", BOOL_AND, 0, 0},
		{"||", BOOL_OR, 0, 0},
		{"!", BOOL_NOT, 0, 0},
		{"|", UNARYOR, 0, 0},
		{"&", UNARYAND, 0, 0},
		{"<<", LSHIFT, 0, 0},
		{">>", RSHIFT, 0, 0},
		{";", SEMICOLON, 0, 0},
		{":", RETURN_IND, 0, 0},
		{",", COMMA, 0, 0},
		{".", DOT, 0, 0},
		{"int", INT, 0, 0},
		{"int16", INT16, 0, 0},
		{"int32", INT32, 0, 0},
		{"int64", INT64, 0, 0},
		{"float16", FLOAT16, 0, 0},
		{"float64", FLOAT64, 0, 0},
		{"float64", FLOAT64, 0, 0},
		{"float", FLOAT, 0, 0},
		{"long", LONG, 0, 0},
		{"cmx64", CMX64, 0, 0},
		{"cmx", CMX, 0, 0},
		{"uint16", UINT16, 0, 0},
		{"uint32", UINT32, 0, 0},
		{"uint64", UINT64, 0, 0},
		{"uint", UINT, 0, 0},
		{"bool", BOOL, 0, 0},
		{"true", TRUE, 0, 0},
		{"false", FALSE, 0, 0},
		{"import", IMPORT, 0, 0},
		{"string", STRING, 0, 0},
	}

	for _, tst := range table {
		if !checkType(tst.inToken, tst.outType, tst.startLine, tst.startOffset) {
			t.Errorf("[!] Couldn't match token '%s' to type %d", tst.inToken, tst.outType)
			break
		}
	}
}
