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

var (
	reserved = map[string]TokenType{}
)

///Token struct
type Token struct {
	Lexme string
	Ttype TokenType
}

func isLetter(ch rune) bool {
	return ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z'
}

func isNum(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

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

func resolveType(buffer bytes.Buffer) Token {
	reserve()

	if buffer.Len() < 2 {
		return Token{Ttype: ID, Lexme: buffer.String()}
	}

	if val, ok := reserved[buffer.String()]; ok {
		return Token{Ttype: val, Lexme: buffer.String()}
	}

	return Token{Ttype: ID, Lexme: buffer.String()}
}
