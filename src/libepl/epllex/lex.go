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

	"eplc/src/libepl/epllex/Errors"
	"bufio"
	"bytes"
	"io"
	"unicode/utf8"
	"fmt"
)

var (
	line       uint
	lineOffset uint
	errCount  uint
)

/*
	Lexer. the job of the lexer is to break the input stream into
	meaningful parts that later will be used by the parser and the SDT
*/
type Lexer struct {
	Buffer *bufio.Reader
	filename string
	ST SymbolTable
}

//New Lexer
func New(file io.Reader, name string) Lexer {
	return Lexer{Buffer: bufio.NewReader(file), ST: SymbolTable{map[string]*SymbolData{}}, filename: name}
}

//Checks if the charecter is valid utf
func (l *Lexer) checkEncodiung(ch rune) bool {
	return utf8.ValidRune(ch)
}

func (l *Lexer) scanNumbers() Token {
	var buf bytes.Buffer
	ch := l.read()
	irn := false

	for isNum(ch) || ch == '.' {
		if ch == '.' {
			irn = true
		}
		buf.WriteRune(ch)
		ch = l.read()
	}
	l.unread()

	if irn {
		return Token{Ttype: REAL, Lexme: buf.String()}
	}

	return Token{Ttype: NUM, Lexme: buf.String()}
}

//Next reads the input stream and resolve the type of the readed character
func (l *Lexer) Next() Token {
	l.skipWhiteSpaces()
	ch := l.read()

	if ch == -1 {
		return Token{Ttype: EOF, Lexme: ""}
	} else if ch == '\n' {
		line++
	} else if isLetter(ch) {
		l.unread()
		return l.scanID()
	} else if isNum(ch) {
		l.unread()
		return l.scanNumbers()
	} else {

		switch ch {
		case '!':
			return Token{Ttype: NOT, Lexme: "!"}
		case ';':
			return Token{Ttype: SEMICOLON, Lexme: ";"}
		case '.':
			return Token{Ttype: DOT, Lexme: "."}
		case ',':
			return Token{Ttype: COMMA, Lexme: ","}
		case ':':
			return Token{Ttype: RETURN_IND, Lexme: ":"}
		case '[':
			return Token{Ttype: LSUBSCRIPT, Lexme: "["}
		case ']':
			return Token{Ttype: RSUBSCRIPT, Lexme: "]"}
		case '{':
			return Token{Ttype: LBRACE, Lexme: "{"}
		case '}':
			return Token{Ttype: RBRACE, Lexme: "}"}
		case '(':
			return Token{Ttype: LPAR, Lexme: "("}
		case ')':
			return Token{Ttype: RPAR, Lexme: ")"}
		case '|':
			if ch = l.read(); ch == '|' {
				return Token{Ttype: OR, Lexme: "||"}
			}
			l.unread()

			return Token{Ttype: UNARYOR, Lexme: "|"}
		case '&':
			if ch = l.read(); ch == '&' {
				return Token{Ttype: AND, Lexme: "&&"}
			}
			l.unread()

			return Token{Ttype: UNARYAND, Lexme: "&"}
		case '+':
			if ch = l.read(); ch == '=' {
				return Token{Ttype: PLUSEQUAL, Lexme: "+="}
			}
			l.unread()

			return Token{Ttype: PLUS, Lexme: "+"}

		case '-':
			if ch = l.read(); ch == '=' {
				return Token{Ttype: MINUSEQUAL, Lexme: "-="}
			}
			l.unread()

			return Token{Ttype: MINUS, Lexme: "-"}
		case '*':
			if ch = l.read(); ch == '=' {
				return Token{Ttype: MULTEQUAL, Lexme: "*="}
			}
			l.unread()

			return Token{Ttype: MULT, Lexme: "*"}
		case '/':
			if ch = l.read(); ch == '=' {
				return Token{Ttype: DEVEQUAL, Lexme: "/="}
			} else if ch == '/' {
				l.unread()
				l.readLine()
				return l.Next()
			} else if ch == '*' {
					l.skipMltLinesComment()
					return l.Next()
			}
			l.unread()

			return Token{Ttype: DEV, Lexme: "/"}
		case '>':
			if ch = l.read(); ch == '=' {
				return Token{Ttype: GE, Lexme: ">="}
			}
			l.unread();  
			if ch = l.read(); ch == '>' {
				return Token{Ttype: RSHIFT, Lexme: ">>"}
			}
			l.unread()

			return Token{Ttype: GT, Lexme: ">"}
		case '<':
			if ch = l.read(); ch == '=' {
				return Token{Ttype: LE, Lexme: "<="}
			}
			l.unread();  
			
			if ch = l.read(); ch == '<' {
				return Token{Ttype: LSHIFT, Lexme: "<<"}
			}
			l.unread()

			return Token{Ttype: LT, Lexme: "<"}
		case '=':
			if ch = l.read(); ch == '=' {
				return Token{Ttype: EQ, Lexme: "=="}
			}
			l.unread()

			return Token{Ttype: ASSIGN, Lexme: "="}
		case '\'':
			return l.matchBy('\'')
		case '"':
			return l.matchBy('"')
		}
		Errors.TokenError(line, lineOffset, ch, l.filename)
		errCount++;
		
		if errCount > 5 {
			Errors.FatalLexical("To many errors")
		}
		l.Next()
	}
	
	Errors.FatalLexical("To many errors")
	return Token{"", EOF}
}


func (l *Lexer) matchBy(s rune) Token {
	var buffer bytes.Buffer
	ch := l.read()
	buffer.WriteRune(s)

	for ch != s {
		buffer.WriteRune(ch)
		ch = l.read()
	}
	buffer.WriteRune(ch)

	return Token{Ttype: STRINGLITERAL, Lexme: buffer.String()}
}

//Skips the whitespaces
func (l *Lexer) skipWhiteSpaces() {
	ch := l.read()

	for ch == '\n' || ch == '\t' || ch == '\r' || ch == ' ' {
		if ch == '\n'{
			lineOffset = 0;
			line++;
		}
		ch = l.read()
	}
	l.unread()
}

func (l *Lexer) readLine() {
	ch := l.read()

	for ch != '\n' {
		ch = l.read()
	}

	l.unread()
}

func (l *Lexer) read() rune {
	char, _, err := l.Buffer.ReadRune()

	if err != nil {
		return -1
	}

	if !l.checkEncodiung(char) {
		Errors.Lexical( l.filename+":"+fmt.Sprint(line)+":"+fmt.Sprint(lineOffset)+": Invalid character: U", uint(char))
	}

	lineOffset++
	return char
}

func (l *Lexer) scanID() Token {
	var buf bytes.Buffer
	ch := l.read()

	for isLetter(ch) || isNum(ch) {

		buf.WriteRune(ch)
		ch = l.read()
	}
	l.unread()

	if l.ST.Get(buf.String()) != (SymbolData{}) {
		return Token{buf.String(), ID}
	}

	tmp := resolveType(buf)
	
	if tmp.Ttype == ID {
		l.ST.Add(&SymbolData{symbol: tmp.Lexme})
	}

	return tmp
}

func (l *Lexer) skipMltLinesComment() {
	ch := l.read()

	for {
		if ch == '*' && l.read() == '/' {
			break;
		} else if ch == -1 {
			break;
		} else {
			ch = l.read()
		}
	}

	if ch == -1 {
		Errors.Lexical(":"+l.filename+": error found EOF expected '*/' ")
	}
}

func (l *Lexer) unread() {
	lineOffset--
	l.Buffer.UnreadRune()
}
