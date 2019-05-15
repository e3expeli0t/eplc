package eplparse

import (
	"eplc/src/libepl/epllex"
	"strings"
	"testing"
)

func TestParser_ParseBlock(t *testing.T) {

}

func TestParser_ParseExpression(t *testing.T) {

}

func TestParser_ParseFnc(t *testing.T) {
	lx := epllex.New(strings.NewReader("fnc exec(command uint, command string, args string): Proc {\n}"), "Test")
	p := New(lx)

}

func TestParser_ParseIdent(t *testing.T) {

}

func TestParser_ParseImport(t *testing.T) {

}

func TestParser_ParseParamList(t *testing.T) {

}

func TestParser_ParseUnaryOp(t *testing.T) {

}

func TestParser_ParseVarDecl(t *testing.T) {

}

func TestParser_ParseProgramFile(t *testing.T) {

}
