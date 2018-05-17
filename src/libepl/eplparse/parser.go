package eplparse

import  (
	"eplc/src/libepl/epllex"
	"eplc/src/libepl/Output"
)

func New (lx epllex.Lexer) Parser {
	return Parser{lx}
}

type Parser struct {
	Lexer epllex.Lexer
}

func (p *Parser) Construct() {

	var tmp = p.readNextToken()

	for tmp.Ttype != epllex.EOF {
		Output.PrintLog(tmp.Lexme)
		tmp = p.readNextToken()
	}

}

func (p *Parser) expression() {

}

func (p *Parser) match(t epllex.Token) bool{
	return p.readNextToken() == t 
}

func (p *Parser) readNextToken() epllex.Token {
	return p.Lexer.Next()
}