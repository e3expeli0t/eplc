package eplinter

import (
	"eplc/src/libepl/epllex"
	"eplc/src/libepl/eplparse"
	"fmt"
	"os"
)

type Command struct {
	commandText string
	commandRes  uint
}

func makeCommand(cnt string) Command {
	return Command{commandText: cnt, commandRes: 0xdeadbeef}
}

type Session struct {
	PrevCommands []Command

	Pt string
}

func NewSession() Session {
	reader, err := os.Open(".eplcache/inter.epl")
	if err != nil {
		panic("[+] Couldn't load base cache")
	}
	lex := epllex.New(reader, "base file")
	parser := eplparse.New(lex)
	return Session{PrevCommands: []Command{}, Parser: &parser, Lexer: &lex, Pt: "epl->"}
}

func (ses *Session) loop() {
	var currentCommand Command
	var userInput string
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()

	for currentCommand.commandText != "exit" {
		fmt.Print(ses.Pt)
		_, _ = fmt.Scanln(&userInput)
		currentCommand = makeCommand(userInput)
		ses.PrevCommands = append(ses.PrevCommands, currentCommand)

	}
}
