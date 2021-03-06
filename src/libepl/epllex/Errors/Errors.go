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

package Errors

import (
	"eplc/src/libio"
	"fmt"
)

//TokenError prints error msg with precise info about the token that cause the error
func TokenError(line uint, lineOffset uint, token rune, filename string) {
	libio.PrintErr(fmt.Sprintf(" in %s:%d:%d: Could't resolve Token %#U", filename, line, lineOffset, token))
}

func ExpError(line uint, lineOffset uint, fname string, cline string, char rune) {
	libio.LexicalPrint(fname, line, lineOffset, cline, "Couldn't resolve token", char)
}

//Prints lexical error msg and quits after that
func FatalLexical(msg ...interface{}) {
	libio.PrintFatalErr(msg...)
}

func EncodingError(line uint, lineOffset uint, filename string, char rune) {
	libio.PrintErr(fmt.Sprintf(" in %s:%d:%d: Encoding error %#U", filename, line, lineOffset, char))
}

//Lexical prints lexical error msg without quiting
func Lexical(msg ...interface{}) {
	libio.PrintErr(msg...)
}
