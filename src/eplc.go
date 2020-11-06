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

package main

import (
	"eplc/src/libeplc"
	"eplc/src/libio"
	"flag"
	"fmt"
	"os"
)

func main() {
	var flags = libeplc.Default

	prof := flag.Bool("profile", false, "profile the compilation process")
	verbose := flag.Bool("verbose", false, "print verbose output")
	flag.Parse()

	if len(flag.Args()) < 1 {
		fmt.Println("Usage: eplc [options] <file>")
		os.Exit(-1)
	}
	file := flag.Args()[0]

	reader, err := os.Open(file)
	if err != nil {
		libio.PrintFatalErr("eplc", fmt.Sprintf("file: '%s' don't exists", file))
	}
	libio.PrintVersion()

	if *prof {
		flags = flags.Add(libeplc.Profile)
	}

	if *verbose {
		flags = flags.Add(libeplc.Verbose)
	}

	libeplc.Compile(reader, file, flags)

}