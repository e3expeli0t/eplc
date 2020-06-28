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

package tests

import (
	"eplc/src/libepl/epllex"
	"strings"
	"testing"
)

func checkType(lexme string, tt epllex.TokenType, startLine uint, startOffset uint) bool {
	lx := epllex.New(strings.NewReader(lexme), "test_lexer.epl")
	tk := lx.Next()
	return tk.StartLine == startLine && tk.StartOffset == startOffset && tk.Ttype == tt
}

func TestLexer(t *testing.T) {
	table := []struct {
		inToken     string
		outType     epllex.TokenType
		startLine   uint
		startOffset uint
	}{
		{"23", epllex.NUM, 0, 0},
		{"2.3", epllex.REAL, 0, 0},
		{"HELLO3433", epllex.ID, 0, 0},
		{"h", epllex.ID, 0, 0},
		{"@MainFunc", epllex.CFLAG, 0, 0},
		{"//comment\n/*More comment*/", epllex.EOF, 0, 0},
		{"\n\n\"string\"", epllex.STRINGLITERAL, 2, 0},
		{"\n'string'", epllex.STRINGLITERAL, 1, 0},
		{"if", epllex.IF, 0, 0},
		{"else", epllex.ELSE, 0, 0},
		{"repeat", epllex.REPEAT, 0, 0},
		{"until", epllex.UNTIL, 0, 0},
		{"move", epllex.MOVE, 0, 0},
		{"in", epllex.IN, 0, 0},
		{"static", epllex.STATIC, 0, 0},
		{"fixed", epllex.FIXED, 0, 0},
		{"decl", epllex.DECL, 0, 0},
		{"define", epllex.DEFINE, 0, 0},
		{"type", epllex.TYPE, 0, 0},
		{"struct", epllex.STRUCT, 0, 0},
		{"fnc", epllex.FNC, 0, 0},
		{"+", epllex.ADD, 0, 0},
		{"-", epllex.SUB, 0, 0},
		{"/", epllex.DIV, 0, 0},
		{"*", epllex.MULT, 0, 0},
		{"{", epllex.LBRACE, 0, 0},
		{"}", epllex.RBRACE, 0, 0},
		{"(", epllex.LPAR, 0, 0},
		{")", epllex.RPAR, 0, 0},
		{"[", epllex.LSUBSCRIPT, 0, 0},
		{"]", epllex.RSUBSCRIPT, 0, 0},
		{"+=", epllex.PLUSEQUAL, 0, 0},
		{"-=", epllex.MINUSEQUAL, 0, 0},
		{"*=", epllex.MULTEQUAL, 0, 0},
		{"/=", epllex.DEVEQUAL, 0, 0},
		{"=", epllex.ASSIGN, 0, 0},
		{">", epllex.GT, 0, 0},
		{"<", epllex.LT, 0, 0},
		{"==", epllex.EQ, 0, 0},
		{">=", epllex.GE, 0, 0},
		{"<=", epllex.LE, 0, 0},
		{"&&", epllex.BOOL_AND, 0, 0},
		{"||", epllex.BOOL_OR, 0, 0},
		{"!", epllex.BOOL_NOT, 0, 0},
		{"|", epllex.UNARYOR, 0, 0},
		{"&", epllex.UNARYAND, 0, 0},
		{"<<", epllex.LSHIFT, 0, 0},
		{">>", epllex.RSHIFT, 0, 0},
		{";", epllex.SEMICOLON, 0, 0},
		{":", epllex.RETURN_TYPE_IND, 0, 0},
		{",", epllex.COMMA, 0, 0},
		{".", epllex.DOT, 0, 0},
		{"int", epllex.INT, 0, 0},
		{"int16", epllex.INT16, 0, 0},
		{"int32", epllex.INT32, 0, 0},
		{"int64", epllex.INT64, 0, 0},
		{"float16", epllex.FLOAT16, 0, 0},
		{"float64", epllex.FLOAT64, 0, 0},
		{"float64", epllex.FLOAT64, 0, 0},
		{"float", epllex.FLOAT, 0, 0},
		{"long", epllex.LONG, 0, 0},
		{"cmx64", epllex.CMX64, 0, 0},
		{"cmx", epllex.CMX, 0, 0},
		{"uint16", epllex.UINT16, 0, 0},
		{"uint32", epllex.UINT32, 0, 0},
		{"uint64", epllex.UINT64, 0, 0},
		{"uint", epllex.UINT, 0, 0},
		{"bool", epllex.BOOL, 0, 0},
		{"true", epllex.TRUE, 0, 0},
		{"false", epllex.FALSE, 0, 0},
		{"import", epllex.IMPORT, 0, 0},
		{"string", epllex.STRING, 0, 0},
	}

	for _, tst := range table {
		if !checkType(tst.inToken, tst.outType, tst.startLine, tst.startOffset) {
			t.Errorf("[!] Couldn't match token '%s' to type %d", tst.inToken, tst.outType)
			break
		}
	}
}
