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
	"eplc/src/libepl/eplccode"
	"eplc/src/libepl/epllex"
	"eplc/src/libepl/eplparse"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

type CodeJunk struct {
	Lines []string
	count uint
}

func (c *CodeJunk) Add(data string) {
	c.Lines = append(c.Lines, data)
	c.count++
}

func (c *CodeJunk) Remove() {
	if c.count > 0 {
		c.Lines = c.Lines[:c.count-1]
	}
}

func (c *CodeJunk) ToString() string {
	var str = "{"

	for _, line := range c.Lines {
		str += line + "\n"
	}

	str += "}"
	return str
}

//todo: implement history
//todo: use CLI v2 ?
type CommandLine struct {
	prompt   string // can be changed by the REPL
	codeJunk *CodeJunk
	lex      *epllex.Lexer
	parse    *eplparse.Parser
	writer   *eplccode.Writer // Used for debugging and faster code
	exit     bool
}

func NewCommandLine(name string) *CommandLine {
	lexer := epllex.New(os.Stdin, name)
	parse := eplparse.New(lexer)
	writer := &eplccode.Writer{
		Fname:      name,
		TargetName: fmt.Sprintf("%s.out", name),
		Labels:     nil,
	}
	return &CommandLine{
		prompt:   "epl->> ",
		codeJunk: &CodeJunk{Lines: nil},
		lex:      lexer,
		parse:    parse,
		writer:   writer,
		exit:     false,
	}
}

func (cl *CommandLine) SetPrompt(s string) {
	cl.prompt = s
}

func (cl *CommandLine) process(data string) {
	cl.lex.Reset(strings.NewReader(data))
}

func (cl *CommandLine) startInterHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cl.exitNow()
		os.Exit(0)

	}()

}

func (cl *CommandLine) ProcessAll() {
	cl.lex.Reset(strings.NewReader(cl.codeJunk.ToString()))
}

//TODO: implement interpreter logic (after the compiler is finished)
func (cl *CommandLine) Execute() bool {
	/*block := cl.parse.ParseBlock(true)
	b, err := json.MarshalIndent(block, "", "\t")

	if err != nil {
		cl.exit = true
		return false
	}
*/
	cl.SetPrompt("++ ")
	var currentToken epllex.Token

	for currentToken.Ttype != epllex.EOF {
		currentToken = cl.lex.Next()
		cl.Println(currentToken.Lexme)
	}
	cl.SetPrompt("epl->> ")
	return true
}

func (cl *CommandLine) Run() {
	var currentLine string
	cl.writer.InitializeWriter()
	cl.startInterHandler()

	for !cl.exit {
		fmt.Print(cl.prompt)
		count, err := fmt.Scan(&currentLine)

		if err != nil {
			fmt.Println(err)
			continue
		}
		cl.SetPrompt("!! ")
		cl.Println(currentLine)
		cl.SetPrompt("epl--> ")

		if count > 0 {
			cl.codeJunk.Add(currentLine)
		}

		cl.process(currentLine)
		cl.Execute()
	}
}

func (cl *CommandLine) exitNow() {
	cl.exit = true

	cl.prompt = "@ "
	cl.Println("Goodbye")
	os.Exit(0)
}
