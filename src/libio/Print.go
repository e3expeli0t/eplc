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

package libio

import (
	"eplc/src/libio/color"
	"fmt"
	"os"
)

const EplcVersion = "0.1.3"

var outputStatusOn = false

func SwitchOutputStatus() {
	outputStatusOn = !outputStatusOn
}

func PrintVersion() {
	if !outputStatusOn {
		return
	}

	fmt.Print(color.BBlue("****\t"), color.BCyan("eplc version: "))
	fmt.Print(color.BLightPurple(EplcVersion), color.BBlue(" ---- Development version\t****\n"))
}

func PrintLog(log ...interface{}) {
	if !outputStatusOn {
		return
	}

	fmt.Print(color.BLightGreen("[*] "))
	fmt.Println(log...)
}

func PrintFatalErr(msg ...interface{}) {
	fmt.Print(color.BLightRed("Fatal: "))
	fmt.Println(msg...)
	os.Exit(-1)
}
func PrintErr(msg ...interface{}) {
	fmt.Print(color.BLightGreen("Error: "))
	fmt.Print(msg...)
	fmt.Println()
}
