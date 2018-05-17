package eplparse

import "eplc/src/libepl/epllex"

type AST struct {
	_type epllex.TokenType
	_value string
	parent *AST
	childrens []AST
}

func (an *AST) IsRoot() bool  {
	return an.parent == nil 
}

func (an *AST) IsData() bool {
	return len(an._value) == 0
}
