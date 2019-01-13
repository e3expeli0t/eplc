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

import (
	"eplc/src/libepl/Output"
	"eplc/src/libepl/eplparse/Types"
	"eplc/src/libepl/eplparse/errors"
	"fmt"

	//"eplc/src/libepl/eplparse/errors"
	"eplc/src/libepl/epllex"
	"eplc/src/libepl/eplparse/ast"
	"eplc/src/libepl/eplparse/symboltable"
)

var (
	currentToken epllex.Token
	lookahead    epllex.Token
)

//New carete new parser struct
func New(lx epllex.Lexer) Parser {
	return Parser{Lexer: lx, ST: symboltable.New()}
}

/*
	Parser. the parser job is to take tokensized stream from the lexer
	and construct a tree form it, the tree is called AST (Abstract Syntax Tree)
	by the set of rules that the language grammar produce
	There are couple of parser kinds, in this version of epl (bootstrap) we are going
	to use a parser that called predictive parser (The grammar class is LL(k)).
	In the future  i'm planning to implement LR(0) parser
*/
type Parser struct {
	Lexer epllex.Lexer
	ST    symboltable.SymbolTable
}

func (p *Parser) report(msg string)  {
	errors.ParsingError(p.Lexer.Filename, p.Lexer.Line, p.Lexer.LineOffset, msg)
}

func (p *Parser) NewScope() {
	currentST := p.ST

	if p.ST.Next == nil {
		p.ST = symboltable.New()
		p.ST.Prev = &currentST
	} else {
		for currentST.Next != nil {
			currentST = *currentST.Next
		}

		p.ST = symboltable.New()
		p.ST.Prev = &currentST
	}
}

func (p *Parser) PreviousScope() {
	p.ST = *p.ST.Prev
}

func (p *Parser) NextScope() {
	if p.ST.Next != nil {
		p.ST = *p.ST.Next
	}
	//TODO: Make the method go to the first scope in case the next scope is nil
}


//Construct new AST from the token stream
func (p *Parser) Construct(){
	p.readNextToken()
	var tmp = currentToken

	//For debugging prep
	for tmp.Ttype != epllex.EOF {
		Output.PrintLog(tmp.Ttype, tmp.Lexme)
		p.readNextToken()
		tmp = currentToken
	}

}

func (p *Parser) ParseProgram() ast.Node {
	p.readNextToken()
	p.NewScope()

	var imports ast.Import
	var decls []ast.Decl

	if p.match(epllex.IMPORT) {
		imports = p.ParseImport()
	}

	if p.match(epllex.DECL) {
		for p.match(epllex.DECL) {
			decls = append(decls, p.ParseVarDecl(symboltable.GLOBAL))
		}
	} else {
		p.report(fmt.Sprintf("Found %s", currentToken.Lexme))
	}

	return ast.Program{Imports: &imports, Decls:decls}
}

func (p *Parser)ParseImport() ast.Import {

	var importList []string
	start := currentToken.StartLine
	
	for !p.match(epllex.SEMICOLON) {
		if p.match(epllex.ID) {
			importList = append(importList, currentToken.Lexme)
		}
		p.readNextToken()
	}

	p.readNextToken()
	return ast.Import{start,importList}
}


func (p *Parser) ParseVarDecl(scope symboltable.ScopeType) ast.Decl {
	var stat string
	var id string
	var Type Types.EplType

	if p.match_n(epllex.FIXED) {
		stat = "fixed"
		p.readNextToken()
	} else if p.match_n(epllex.ID){
		stat = "unfixed"
		p.readNextToken()
	} else {
		p.report(fmt.Sprintf("Found %s expected ID or fixed tokens", currentToken.Lexme))
	}

	if p.match(epllex.ID) {
		id = currentToken.Lexme
	} else {
		p.report(fmt.Sprintf("found %s expected ID", currentToken.Lexme))
	}
	p.readNextToken()

	if Types.IsValidBasicType(currentToken) {
		Type = *Types.ResolveType(currentToken)
	} else {
		Type = *Types.MakeType(stat+currentToken.Lexme)
	}

	p.readNextToken()

	if !p.match(epllex.SEMICOLON) {
		p.readNextToken()
	}

	p.readNextToken()
	p.ST.Add(symboltable.NewSymbol(id, Type, scope))

	return ast.VarDecl{id, Type, ast.VarStat(stat)}
}

func (p *Parser) match(t epllex.TokenType) bool{
	return currentToken.Ttype == t 
}

func (p *Parser) match_n(t epllex.TokenType) bool {
	return lookahead.Ttype == t
}

func (p *Parser) readNextToken() {

	if (epllex.Token{}) == currentToken {
		currentToken = p.Lexer.Next()
	} else {
		currentToken = lookahead
	}
	lookahead = p.Lexer.Next()
}
