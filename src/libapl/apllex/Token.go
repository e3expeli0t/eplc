package apllex

import (
	"bytes"
)

var (
	reserved = map[string]TokenType{}
)

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
	reserved["virtual"] = VIRTUAL
	reserved["repeat"] = REPEAT
	reserved["until"] = UNTIL
	reserved["fnc"] = FNC
	reserved["if"] = IF
	reserved["else"] = ELSE
	reserved["in"] = IN
	reserved["define"] = DEFINE
	reserved["sdecl"] = SDECL
	reserved["struct"] = STRUCT
	reserved["static"] = STATIC
	reserved["type"] = TYPE
	reserved["decltype"] = DECLTYPE
	reserved["fixed"] = FIXED
	reserved["decl"] = DECL
	reserved["move"] = MOVE

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
