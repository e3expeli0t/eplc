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
)

var prevOffset uint = 0
/*
	Lexer. the job of the lexer is to break the input stream into
	meaningful parts that later will be used by the parser and the IR generator
	The lexer is Deterministic finite state machine (Deterministic finite automata)
*/
type Lexer struct {
	Buffer     *bufio.Reader
	Filename   string
	Line       uint
	LineOffset uint
	ErrCount   uint
}

//New Lexer
func New(file io.Reader, name string) Lexer {
	return Lexer{Buffer: bufio.NewReader(file), Filename: name, Line: 0, LineOffset: 0, ErrCount: 0}
}

//Checks if the character is valid utf
func (l *Lexer) checkEncoding(ch rune) bool {
	return utf8.ValidRune(ch)
}

//Next reads the input stream and resolve the type of the read character
func (l *Lexer) Next() Token {
	l.skipWhiteSpaces() //if there are whitespaces skip them

	//Save the Line and the Line offset for the parser error handling
	startOffset := l.LineOffset
	startLine := l.Line

	ch := l.read() //read one char

	if ch == -1 {
		return Token{Ttype: EOF, Lexme: "", StartOffset: 0, StartLine: 0}
	} else if isLetter(ch) {
		l.unread()
		return l.scanID(false)
	} else if isNum(ch) {
		l.unread()
		return l.scanNumbers()
	} else {

		switch ch {
		case '@':
			//Read the string that followed the @ char and return it as compiler flag
			return l.scanID(true)
		case '!':
			return Token{Ttype: NOT, Lexme: "!", StartOffset: startOffset, StartLine: startLine}
		case ';':
			return Token{Ttype: SEMICOLON, Lexme: ";", StartOffset: startOffset, StartLine: startLine}
		case '.':
			return Token{Ttype: DOT, Lexme: ".", StartOffset: startOffset, StartLine: startLine}
		case ',':
			return Token{Ttype: COMMA, Lexme: ",", StartOffset: startOffset, StartLine: startLine}
		case ':':
			return Token{Ttype: RETURN_IND, Lexme: ":", StartOffset: startOffset, StartLine: startLine}
		case '[':
			return Token{Ttype: LSUBSCRIPT, Lexme: "[", StartOffset: startOffset, StartLine: startLine}
		case ']':
			return Token{Ttype: RSUBSCRIPT, Lexme: "]", StartOffset: startOffset, StartLine: startLine}
		case '{':
			return Token{Ttype: LBRACE, Lexme: "{", StartOffset: startOffset, StartLine: startLine}
		case '}':
			return Token{Ttype: RBRACE, Lexme: "}", StartOffset: startOffset, StartLine: startLine}
		case '(':
			return Token{Ttype: LPAR, Lexme: "(", StartOffset: startOffset, StartLine: startLine}
		case ')':
			return Token{Ttype: RPAR, Lexme: ")", StartOffset: startOffset, StartLine: startLine}
		case '|':
			if ch = l.read(); ch == '|' {
				return Token{Ttype: OR, Lexme: "||", StartOffset: startOffset, StartLine: startLine}
			}
			l.unread()

			return Token{Ttype: UNARYOR, Lexme: "|", StartOffset: startOffset, StartLine: startLine}
		case '&':
			if ch = l.read(); ch == '&' {
				return Token{Ttype: AND, Lexme: "&&", StartOffset: startOffset, StartLine: startLine}
			}
			l.unread()

			return Token{Ttype: UNARYAND, Lexme: "&", StartOffset: startOffset, StartLine: startLine}
		case '+':
			if ch = l.read(); ch == '=' {
				return Token{Ttype: PLUSEQUAL, Lexme: "+=", StartOffset: startOffset, StartLine: startLine}
			}
			l.unread()

			return Token{Ttype: PLUS, Lexme: "+", StartOffset: startOffset, StartLine: startLine}

		case '-':
			if ch = l.read(); ch == '=' {
				return Token{Ttype: MINUSEQUAL, Lexme: "-=", StartOffset: startOffset, StartLine: startLine}
			}
			l.unread()

			return Token{Ttype: MINUS, Lexme: "-", StartOffset: startOffset, StartLine: startLine}
		case '*':
			if ch = l.read(); ch == '=' {
				return Token{Ttype: MULTEQUAL, Lexme: "*=", StartOffset: startOffset, StartLine: startLine}
			}
			l.unread()

			return Token{Ttype: MULT, Lexme: "*", StartOffset: startOffset, StartLine: startLine}
		case '/':
			if ch = l.read(); ch == '=' {
				return Token{Ttype: DEVEQUAL, Lexme: "/=", StartOffset: startOffset, StartLine: startLine}
			} else if ch == '/' {
				l.unread()
				l.readLine()
				return l.Next()
			} else if ch == '*' {
				l.skipMltLinesComment()
				return l.Next()
			}
			l.unread()

			return Token{Ttype: DEV, Lexme: "/", StartOffset: startOffset, StartLine: startLine}
		case '>':
			if ch = l.read(); ch == '=' {
				return Token{Ttype: GE, Lexme: ">=", StartOffset: startOffset, StartLine: startLine}
			}
			l.unread()
			if ch = l.read(); ch == '>' {
				return Token{Ttype: RSHIFT, Lexme: ">>", StartOffset: startOffset, StartLine: startLine}
			}
			l.unread()

			return Token{Ttype: GT, Lexme: ">", StartOffset: startOffset, StartLine: startLine}
		case '<':
			if ch = l.read(); ch == '=' {
				return Token{Ttype: LE, Lexme: "<=", StartOffset: startOffset, StartLine: startLine}
			}
			l.unread()

			if ch = l.read(); ch == '<' {
				return Token{Ttype: LSHIFT, Lexme: "<<", StartOffset: startOffset, StartLine: startLine}
			}
			l.unread()

			return Token{Ttype: LT, Lexme: "<", StartOffset: startOffset, StartLine: startLine}
		case '=':
			if ch = l.read(); ch == '=' {
				return Token{Ttype: EQ, Lexme: "==", StartOffset: startOffset, StartLine: startLine}
			}
			l.unread()

			return Token{Ttype: ASSIGN, Lexme: "=", StartOffset: startOffset, StartLine: startLine}
		case '\'':
			return l.matchBy('\'')
		case '"':
			return l.matchBy('"')
		}
		Errors.TokenError(l.Line, l.LineOffset, ch, l.Filename)
		l.ErrCount++

		if l.ErrCount > 5 {
			Errors.FatalLexical("To many errors")
		}
		l.Next()
	}

	Errors.FatalLexical("To many errors")
	return Token{"", EOF, l.Line, l.LineOffset}
}

func (l *Lexer) matchBy(s rune) Token {
	var buffer bytes.Buffer

	startOffset := l.LineOffset
	startLine := l.Line

	ch := l.read()
	buffer.WriteRune(s)

	for ch != s {
		buffer.WriteRune(ch)
		ch = l.read()
	}
	buffer.WriteRune(ch)

	return Token{Ttype: STRINGLITERAL, Lexme: buffer.String(), StartOffset: startOffset - 1, StartLine: startLine}
}

func (l *Lexer) scanNumbers() Token {

	startOffset := l.LineOffset
	startLine := l.Line

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
		return Token{Ttype: REAL, Lexme: buf.String(), StartOffset: startOffset, StartLine: startLine}
	}

	return Token{Ttype: NUM, Lexme: buf.String(), StartOffset: startOffset, StartLine: startLine}
}

func (l *Lexer) scanID(cf bool) Token {

	startOffset := l.LineOffset
	startLine := l.Line

	var buf bytes.Buffer
	ch := l.read()

	for isLetter(ch) || isNum(ch) {

		buf.WriteRune(ch)
		ch = l.read()
	}
	l.unread()

	if cf {
		return Token{Lexme: buf.String(), Ttype: CFLAG, StartOffset: startOffset - 1, StartLine: startLine}
	}

	/*
	if l.ST.Get(buf.String()) != (SymbolData{}) {
		return Token{buf.String(), ID}
	}
*/
	tmp := resolveType(buf, startLine, startOffset)
	/*
		if tmp.Ttype == ID {
			l.ST.Add(&SymbolData{symbol: tmp.Lexme})
		}
	*/
	return tmp
}

//Skips the whitespaces
func (l *Lexer) skipWhiteSpaces() {
	ch := l.read()

	for ch == '\n' || ch == '\t' || ch == '\r' || ch == ' ' {
		if ch == '\n' {
			prevOffset = l.LineOffset
			l.LineOffset = 0
			l.Line++
		}
		ch = l.read()
	}
	l.unread()
}

func (l *Lexer) skipMltLinesComment() {
	ch := l.read()

	for {
		if ch == '*' && l.read() == '/' {
			break
		} else if ch == -1 {
			break
		} else {
			ch = l.read()
		}
	}

	if ch == -1 {
		Errors.FatalLexical(":" + l.Filename + ": error found EOF expected '*/' ")
	}
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

	if !l.checkEncoding(char) {
		Errors.EncodingError(l.Line, l.LineOffset, l.Filename, char)
		l.ErrCount++
	}

	l.LineOffset++
	return char
}

func (l *Lexer) unread() {
	if l.LineOffset == 0 {
		l.Line--
		l.LineOffset = prevOffset
	} else {
		l.LineOffset--
	}
	l.Buffer.UnreadRune()
}
