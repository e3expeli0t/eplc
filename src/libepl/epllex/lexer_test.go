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
	"testing"
	"strings"
)


func checkType(lexme string, tt TokenType) bool {
	lx := New(strings.NewReader(lexme), "test_lexer.epl")
	
	return lx.Next().Ttype == tt;
}


func TestLexer(t *testing.T) {
	table := []struct {
		inToken string
		outType TokenType
	}{
		{"23", NUM},
		{"2.3", REAL},
		{"HELLO3433", ID},
		{"@MainFunc", CFLAG},
		{"//commant\n/*More commant*/",EOF},
		{"\n\n\"string\"", STRINGLITERAL},
		{"\n'string'", STRINGLITERAL},
		{"if", IF},
		{"else", ELSE},
		{"repeat", REPEAT},
		{"until", UNTIL},
		{"move", MOVE},
		{"in", IN},
		{"static", STATIC},
		{"fixed", FIXED},
		{"decl", DECL},
		{"define", DEFINE},
		{"type", TYPE},
		{"struct", STRUCT},
		{"fnc", FNC},
		{"+", PLUS},
		{"-", MINUS},
		{"/", DEV},
		{"*", MULT},
		{"{", LBRACE},
		{"}", RBRACE},
		{"(", LPAR},
		{")", RPAR},
		{"[", LSUBSCRIPT},
		{"]", RSUBSCRIPT},
		{"+=", PLUSEQUAL},
		{"-=", MINUSEQUAL},
		{"*=", MULTEQUAL},
		{"/=", DEVEQUAL},
		{"=", ASSIGN},
		{">", GT},
		{"<", LT},
		{"==", EQ},
		{">=", GE},
		{"<=", LE},
		{"&&", AND},
		{"||", OR},
		{"!", NOT},
		{"|", UNARYOR},
		{"&", UNARYAND},
		{"<<", LSHIFT},
		{">>", RSHIFT},
		{";", SEMICOLON},
		{":", RETURN_IND},
		{",", COMMA},
		{".", DOT},
		{"int", INT},
		{"int16", INT16},
		{"int32", INT32},
		{"int64", INT64},
		{"float16", FLOAT16},
		{"float64", FLOAT64},
		{"float64", FLOAT64},
		{"float", FLOAT},
		{"long", LONG},
		{"cmx64", CMX64},
		{"cmx", CMX},
		{"uint16", UINT16},
		{"uint32", UINT32},
		{"uint64", UINT64},
		{"uint", UINT},
		{"string", STRING},
	}

	for _, tst := range table{
		if !checkType(tst.inToken, tst.outType) {
			t.Errorf("[!] Couldn't match token '%s' to type %d", tst.inToken, tst.outType)
		}
	}

	t.Log(INT)
}