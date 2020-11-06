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

package libeplc

import (
	"eplc/src/libepl/analysis"
	"eplc/src/libepl/eplccode"
	"eplc/src/libepl/epllex"
	"eplc/src/libepl/eplparse"
	"eplc/src/libio"

	"fmt"
	"io"
	"os"
	"runtime"
	"runtime/debug"
	"runtime/pprof"
)

func eplRecover(r interface{}, stack bool) {

	errorString := fmt.Sprintf("got an '%v' error during execution.\n"+
		"This is an internal compiler error, please contact one of the project developers.\n"+
		"You can rerun the compiler with the flag --verbose to get stack trace", r)
	libio.PrintErr(errorString)

	if stack {
		fmt.Println("----------------------Stack----------------------")
		fmt.Println(string(debug.Stack()))
		fmt.Println("-------------------------------------------------")
	}

	libio.PrintFatalErr("Quiting due to previous error")
}

func Compile(source io.Reader, fileName string, options Flag) {
	defer func() {
		if r := recover(); r != nil {
			eplRecover(r, options.Contains(Verbose))
		}
	}()

	if options.Contains(Profile) {
		f, err := os.Create("Profile/dump.cpu")
		if err != nil {
			libio.PrintFatalErr("could not create CPU profile: ", err.Error())
		}
		defer f.Close()

		if err := pprof.StartCPUProfile(f); err != nil {
			libio.PrintFatalErr("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	lexer := epllex.New(source, fileName)
	parser := eplparse.New(lexer)

	parser.InitializeTypeHandler()
	file := parser.ParseProgramFile()

	Analyzer := analysis.NewAnalyzer(file)
	Analyzer.Init()
	Analyzer.Run()

	eplccode.GenerateAIR(file)

	if options.Contains(Profile) {
		f, err := os.Create("Profile/dump.mem")
		if err != nil {
			libio.PrintFatalErr("could not create memory profile: ", err)
		}
		defer f.Close()

		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			libio.PrintFatalErr("could not write memory profile: ", err)
		}
	}
}
