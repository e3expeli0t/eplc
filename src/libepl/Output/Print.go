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

package Output

import (
	"fmt"
	"eplc/src/libepl/Output/color"
	"os"
)

const startmsg = "Starting aplc v0.001:"

func PrintLog(log ...interface{}) {
	fmt.Print(color.BLightGreen("[*] "))
	fmt.Println(log...)
}

func PrintErr(seg string, msg...interface{}) {
	fmt.Print(color.BLightRed("Fatal: "), "<"+seg+">")
	fmt.Print(msg...)
	fmt.Println()
	fmt.Println(color.BLightBlue("Quiting!!"))
	os.Exit(-1)
}

func PrintStartMSG() {
	fmt.Println(color.LightBlue(startmsg))
}