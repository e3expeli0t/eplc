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
	"bytes"
)

//todo: remove to inside the lexer struct
var (
	reserved = map[string]TokenType{}
)

///Token struct
type Token struct {
	Lexme       string
	Ttype       TokenType
	StartLine   uint
	StartOffset uint
}

const (
	ExpLowPrec = 0
	HighPrec   = 5
)

//todo: add equal prec and parsing
func (t *Token) Precedence() int {
	switch t.Ttype {
	case BOOL_OR:
		return 1
	case BOOL_AND:
		return 2
	case EQ, NEQ, GT, GE, LT, LE:
		return 3
	case ADD, SUB:
		return 4
	case MULT, DIV:
		return 5
		
	}

	return ExpLowPrec
}

func (t *Token) IsScalar() bool {
	return t.Ttype == NUM || t.Ttype == REAL
}

func (t *Token) IsIdent() bool {
	return t.Ttype == ID
}

func ( t *Token) IsString() bool {
	return t.Ttype == STRINGLITERAL
}

func (t *Token) IsUnary() bool {
	return t.Ttype == ADD || t.Ttype == SUB
}

func (t *Token) IsLeftAssociative() bool {
	switch t.Ttype {
	case MULT, ADD, SUB, DIV:
		return true
	}
	return false
}

func isLetter(ch rune) bool {
	return ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z'
}

func isNum(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

//these are all the reserved names
func reserve() {
	reserved["repeat"] = REPEAT
	reserved["until"] = UNTIL
	reserved["define"] = DEFINE
	reserved["fnc"] = FNC
	reserved["if"] = IF
	reserved["else"] = ELSE
	reserved["in"] = IN
	reserved["struct"] = STRUCT
	reserved["static"] = STATIC
	reserved["type"] = TYPE
	reserved["fixed"] = FIXED
	reserved["decl"] = DECL
	reserved["move"] = MOVE
	reserved["bool"] = BOOL
	reserved["int"] = INT
	reserved["int8"] = INT8
	reserved["int16"] = INT16
	reserved["int32"] = INT32
	reserved["int64"] = INT64
	reserved["uint"] = UINT
	reserved["uint8"] = UINT8
	reserved["uint16"] = UINT16
	reserved["uint32"] = UINT32
	reserved["uint64"] = UINT64
	reserved["float"] = FLOAT
	reserved["float16"] = FLOAT16
	reserved["float32"] = FLOAT32
	reserved["float64"] = FLOAT64
	reserved["cmx64"] = CMX64
	reserved["cmx"] = CMX
	reserved["string"] = STRING
	reserved["long"] = LONG
	reserved["bool"] = BOOL
	reserved["import"] = IMPORT
	reserved["true"] = TRUE
	reserved["false" ] = FALSE
}

func resolveType(buffer bytes.Buffer, startLine uint, startOffset uint) Token {
	reserve()

	if buffer.Len() < 2 {
		return Token{Ttype: ID, Lexme: buffer.String()}
	}

	if val, ok := reserved[buffer.String()]; ok {
		return Token{Ttype: val, Lexme: buffer.String()}
	}

	return Token{Ttype: ID, Lexme: buffer.String(), StartLine: startLine, StartOffset: startOffset}
}
