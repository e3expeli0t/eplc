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

package Errors

import (
	"fmt"
	"eplc/src/libepl/Output"
)

//TokenError prints error msg with precise info about the token that cause the error
func TokenError(line uint, lineOffset uint, token rune, filename string) {
	Output.PrintErr("Lexical", filename+":"+fmt.Sprint(line)+":"+fmt.Sprint(lineOffset)+": Could't resolve Token '"+string(token)+"'")
}

//Lexical prints lexical error msg
func Lexical(msg ...interface{}) {
	Output.PrintErr("Lexical", msg...)
}