/*
*	eplc
*	Copyright (C) 2018 eplc core team
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

package eplparse

import  (
	"eplc/src/libepl/epllex"
	"eplc/src/libepl/Output"
)



//New carete new parser struct
func New (lx epllex.Lexer) Parser {
	return Parser{lx}
}

/*
	Parser. the parser job is to take tokenized stream from the lexer
	and construct a tree form it, the tree is calld AST (Abstract Syntax Tree) 
	by the set of rules that the language grammar produce
	There are couple of parser kinds, in this version of epl (bootstrap) we are going
	to use a parser that calld predictive parser (The grammar class is LL(k)).
	In the future  i'm planning to implement LR(0) parser
*/
type Parser struct {
	Lexer epllex.Lexer
}

//Construct new AST from the token stream
func (p *Parser) Construct() {

	var tmp = p.readNextToken()

	//For debugging prep
	for tmp.Ttype != epllex.EOF {
		Output.PrintLog(tmp.Lexme)
		tmp = p.readNextToken()
	}

}

func (p *Parser) fnc() {
	if p.readNextToken().Ttype != epllex.FNC {
		//perr
	}
	

}

func (p *Parser) match(t epllex.Token) bool{
	return p.readNextToken() == t 
}

func (p *Parser) readNextToken() epllex.Token {
	return p.Lexer.Next()
}