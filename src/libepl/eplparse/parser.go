package eplparse

import  (
	"eplc/src/libepl/epllex"
	"eplc/src/libepl/Output"
	"io"
)

func New (source io.Reader, fname string) Parser {
	return Parser {Lexer: epllex.New(source, fname)}
}

type Parser struct {
	Lexer epllex.Lexer
}

func (p *Parser) Construct() {

	tmp := p.readNextToken()

	for p.readNextToken().Ttype != epllex.EOF {
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