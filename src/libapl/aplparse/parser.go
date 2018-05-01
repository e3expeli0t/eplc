package aplparse

import "aplc/src/libapl/apllex"

type Parser struct {
	Lexer apllex.Lexer
}

func (p *Parser) readNextToken() apllex.Token {
	return p.Lexer.Next()
}