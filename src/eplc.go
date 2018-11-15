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

package main

import (
	"eplc/src/libepl"
	"eplc/src/libepl/Output"
	"io"
	"os"
	"fmt"
	"github.com/pkg/profile"
)

func main() {
	defer profile.Start(profile.CPUProfile).Stop()
	//defer profile.Start(profile.MemProfile).Stop()

	if len(os.Args) == 1 {
		fmt.Println("Usage: eplc <file>")
		os.Exit(-1)
	}

	args := os.Args[1:]
	file := args[0]
	Output.PrintStartMSG()

	var reader io.Reader
	reader, err := os.Open(file)

	if err != nil {
		Output.PrintFatalErr("eplc", fmt.Sprintf("file: '%s' don't exists", file))
	}

	libepl.Compile(reader, file)
}
