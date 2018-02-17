package aplparse

import "../apllex"

type Parser struct {
	Lexer apllex.Lexer
}

func (p *Parser) readNextToken() apllex.Token {
	return p.Lexer.Next()
}

/*
func (p *Parser) ConstructAST() AST {
}
*/