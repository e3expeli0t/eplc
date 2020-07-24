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

package errors

import (
	"eplc/src/libepl/Output"
	"eplc/src/libepl/epllex"
)

type InternalParserError struct {
	ErrCount uint
}

//todo: support position marker (marks the token in specific position)
func (ipe *InternalParserError) ParsingError(filename string, line uint, lineOffset uint, errorMsg string, currentline string, token epllex.Token) {
	descriptor := Output.ErrorDescriptor{
		Fname: filename,
		Line: line,
		LineOffset: lineOffset,
		CurrentLine: currentline,
		ErrorMSG: errorMsg,
		Token: token.Lexme,
	}
	Output.ParserIntelligentError(descriptor)
	ipe.ErrCount++
}

func (ipe *InternalParserError) IsValidFile() {
	if ipe.ErrCount != 0 {
		Output.PrintFatalErr("Parser: Aborting due to previous errors")
	}
}
