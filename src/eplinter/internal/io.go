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

package internal

import (
	"eplc/src/libepl/Output/color"
	"fmt"
)

func (cl *CommandLine) Print(input...interface{}) {
	fmt.Print(cl.prompt)
	fmt.Print(input...)
}

func (cl *CommandLine) Println(input...interface{}) {
	fmt.Print(cl.prompt)
	fmt.Println(input...)
}
func (cl *CommandLine) PrintError(msg...interface{}) {
	fmt.Print(color.Blink(color.GLightRed("!!")))
	fmt.Print(fmt.Sprintf(color.GLightRed("Error:{}:{} "), cl.lex.Filename, cl.lex.Line))
	fmt.Println(msg...)
}