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
	reserved["float32"] = FLOAT32
	reserved["float64"] = FLOAT64
	reserved["cmx64"] = CMX64
	reserved["cmx"] = CMX
	reserved["string"] = STRING
	reserved["long"] = LONG
}

func resolveType(buffer bytes.Buffer) Token {
	if buffer.Len() < 2 {
		return Token{Ttype: ID, Lexme: buffer.String()}
	}

	if val, ok := reserved[buffer.String()]; ok {
		return Token{Ttype: val, Lexme: buffer.String()}
	}

	return Token{Ttype: ID, Lexme: buffer.String()}
}
