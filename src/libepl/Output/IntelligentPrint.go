/*
*	Copyright (C) 2018-2020 Elia Ariaz
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

package Output

import (
	"eplc/src/libepl"
	"eplc/src/libepl/Output/color"
	"fmt"
	"strings"
)

type LexicalErrorDescriptor struct {
	fname       string
	line        uint
	lineOffset  uint
	currentLine string
	errorMSG    string
	ch          rune
}


//Todo: Make this more elegant
func LexicalPrint(fname string, line uint, lineOffset uint, currentLine string, errorMSG string, ch rune) {
	ed := LexicalErrorDescriptor{fname, line, lineOffset,
		currentLine, errorMSG, ch}
	LexerIntelligentError(ed)
}

type ErrorDescriptor struct {
	Fname       string
	Line        uint
	LineOffset  uint
	CurrentLine string
	ErrorMSG    string
	Token       string
}

func (e *ErrorDescriptor) basicInfoPrinter() string {
	return color.BGreen(fmt.Sprintf(":%s:%d:%d: ", e.Fname, e.Line, e.LineOffset))
}

func (e *LexicalErrorDescriptor) basicInfoPrinter() string {
	return color.BGreen(fmt.Sprintf(":%s:%d:%d: ", e.fname, e.line, e.lineOffset))
}

func prints(s string, n int) string {
	return strings.Repeat(s, n)
}

//Todo: make it support Cflags (there token don't have the @ char)
func (e *ErrorDescriptor) TokenMarker() string {
	str := e.CurrentLine + "\n"

	for _, token := range strings.Split(str, " ") {
		if token == "" {
			str += prints(color.BLightPurple("-"), 1)
		} else if token[0] == e.Token[0] {
			index := 0
			for (index <= len(e.Token)-1 && index <= len(token)-1) && token[index] == e.Token[index] {
				index++
			}

			if index == len(e.Token) {
				str += prints(color.BLightPurple("^"), len(e.Token))
				if len(e.Token) < len(token) {
					str += prints(color.BLightPurple("~"), len(token)-len(e.Token))
				}
			} else {
				str += prints(color.BLightPurple("~"), len(token))
			}
		} else {
			str += prints(color.BLightPurple("~"), len(token)+1) // It works for some reason
		}
	}

	//pad
	if len(str) < len(e.CurrentLine) {
		str += prints(color.BLightPurple("~"), len(e.CurrentLine)-len(str))
	}
	return str
}

func (e *LexicalErrorDescriptor) LexicalMarker() string {
	str := e.currentLine + "\n"

	for _, token := range e.currentLine {
		if token != e.ch {
			str += prints(color.BLightPurple("~"), 1)
		} else {
			str += prints(color.BLightPurple("^"), 1)
		}
	}

	if len(str) < len(e.currentLine) {
		str += prints(color.BLightPurple("~"), len(e.currentLine)-len(str))
	}

	return str
}

//todo: support multiple error printing ?
func LexerIntelligentError(err LexicalErrorDescriptor) {
	//Basic information about the error
	fmt.Print(color.BGreen("Error:") + err.basicInfoPrinter())
	fmt.Print(color.BGreen(err.errorMSG) + "\n")

	//Mark the problematic char
	fmt.Println(err.LexicalMarker())

	//Print tips?
	fmt.Println(specialError(string(err.ch), libepl.Lexer))
}

//todo: support multiple error printing ?
func ParserIntelligentError(err ErrorDescriptor) {
	//Basic information about the error
	fmt.Print(color.BGreen("Error:") + err.basicInfoPrinter())
	fmt.Println(color.BGreen(err.ErrorMSG) + "\n")

	//Mark the problematic token
	fmt.Println(err.TokenMarker())

	//Print tips?
	fmt.Println(specialError(err.Token, libepl.Parser))
	PrintFatalErr("SyntaxError", "To many errors")
}

func TypeIntelligentError(err ErrorDescriptor) {
	//Basic information about the error
	fmt.Print(color.BGreen("Error:") + err.basicInfoPrinter())
	fmt.Println(color.BGreen(err.ErrorMSG) + "\n")

	//Mark the problematic token
	fmt.Println(err.TokenMarker())

	//Print tips?
	fmt.Println(specialError(err.Token, libepl.TypeChecker))
	PrintFatalErr("TypeError", "To many errors")
}

func specialError(toke string, d libepl.PhaseIndicator) (err string) {


	switch toke {
	case "EOF": //Note: in the future this  will be changed to much smarter version
		err = "This error is possibly because you mismatched block indicators (i.e '{' and '}'"
	default:
		switch d {
		case libepl.Parser:
			err = "Sorry about that mate. You probably have syntax problem.\n\tTry fix it"
		case libepl.Lexer:
			err = "Sorry about that mate. You probably have weird letters.\n\tTry check your file encoding"
		case libepl.TypeChecker:
			err = "Sorry about that mate. You probably used the wrong type.\n\tTry checking your types again"
		}
	}

	return color.Blink("[?] ")+color.BWhite(err)
}