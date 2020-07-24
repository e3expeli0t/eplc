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
	"eplc/src/libepl/Output"
	"eplc/src/libeplc"
	"flag"
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
)

func main() {
	prof := flag.Bool("profile", false, "profile the compilation process")
	flag.Parse()

	if len(flag.Args()) < 1 {
		fmt.Println("Usage: eplc [options] <file>")
		os.Exit(-1)
	}
	file := flag.Args()[0]

	reader, err := os.Open(file)

	if err != nil {
		Output.PrintFatalErr("eplc", fmt.Sprintf("file: '%s' don't exists", file))
	}

	Output.PrintVersion()
	//For tests
	if *prof {
		f, err := os.Create("Profile/dump.cpu")
		if err != nil {
			Output.PrintFatalErr("could not create CPU profile: ", err.Error())
		}
		defer f.Close()

		if err := pprof.StartCPUProfile(f); err != nil {
			Output.PrintFatalErr("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}
	libeplc.Compile(reader, file)

	f, err := os.Create("Profile/dump.mem")
	if err != nil {
		Output.PrintFatalErr("could not create memory profile: ", err)
	}
	defer f.Close()
	runtime.GC() // get up-to-date statistics
	if err := pprof.WriteHeapProfile(f); err != nil {
		Output.PrintFatalErr("could not write memory profile: ", err)
	}
}
