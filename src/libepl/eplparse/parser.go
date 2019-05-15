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
	errCount 	 uint = 0
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

//Todo: change expect so that the function checks if the expected token is coming (if not print err)
//--------------------------------------------------------------------------------------
//Helper functions

func (p *Parser) expect(ex string, fnd string) {
	p.report(fmt.Sprintf("expected %s found %s ", ex, fnd))
}

func (p *Parser) report(msg string) {
	errors.ParsingError(p.Lexer.Filename, p.Lexer.Line, p.Lexer.LineOffset, msg, p.Lexer.GetLine(), currentToken)
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
func (p *Parser) Construct() {
	p.readNextToken()
	var tmp = currentToken

	//For debugging prep
	for tmp.Ttype != epllex.EOF {
		Output.PrintLog(tmp.Ttype, tmp.Lexme)
		p.readNextToken()
		tmp = currentToken
	}

}

func (p *Parser) ParseProgramFile() ast.Node {
	p.readNextToken()
	p.NewScope()

	var imports []ast.Import
	var decls []ast.Decl
	var fncs []ast.Fnc
	var MainFunction ast.Fnc

	Output.PrintLog("Parsing imports")
	for p.match(epllex.IMPORT) {
		imports = append(imports, p.ParseImport())
	}

	Output.PrintLog("Parsing global variables")
	if p.match(epllex.DECL) {
		for p.match(epllex.DECL) {
			decls = append(decls, p.ParseVarDecl(symboltable.GLOBAL))
		}
	}

	Output.PrintLog("Parsing functions")

	if p.match(epllex.FNC) {
		Output.PrintLog("[*] Parsing functions")
		for p.match(epllex.FNC) {
			fnc := p.ParseFnc()

			//todo: Should change in eplc v0.2 to support cflags
			if fnc.Name == "Main" {
				MainFunction = fnc
				continue
			}
			fncs = append(fncs, fnc)
		}
	}

	Output.PrintLog("Completed")
	return &ast.ProgramFile{Imports: &imports, GlobalDecls: &decls, Symbols: &p.ST, Functions: &fncs, FileName: p.Lexer.Filename, MainFunction: &MainFunction}
}
//----------------------------------------------------------------------------------------------------------------------
//Statements

func (p *Parser) ParseImport() ast.Import {

	var importPath []string
	start := currentToken.StartLine

	for !p.match(epllex.SEMICOLON) {
		if p.match(epllex.ID) {
			importPath = append(importPath, currentToken.Lexme)
		}
		p.readNextToken() //SKIP THE DOT
	}

	p.readNextToken() // SKIP THE SEMICOLON
	return ast.Import{StartLoc: start, Imports: importPath}
}


func (p *Parser) ParseFnc() ast.Fnc {
	p.NewScope()
	if p.match(epllex.FNC) {
		p.readNextToken() //Skipping the fnc keyword
	}

	var name string
	var params *[]ast.Decl
	var returnType *Types.EplType


	if p.match(epllex.ID) {
		name = currentToken.Lexme
		fmt.Println("name:", name, "next:",lookahead)
		p.readNextToken()
	} else {
		p.expect("Ident", currentToken.Lexme)
	}

	if p.match(epllex.LPAR) {
		params = p.ParseParamList()
	} else {
		p.expect("(", currentToken.Lexme)
	}

	if p.match(epllex.RETURN_IND) {
		p.readNextToken()

		if Types.IsValidBasicType(currentToken) {
			returnType = Types.ResolveType(currentToken)
		} else {
			returnType = Types.MakeType(currentToken.Lexme)
		}
	}
	p.readNextToken()

	if p.match(epllex.LBRACE) {
		//parse block
	} else if p.match(epllex.SEMICOLON) {
		p.readNextToken()
	} else {
		p.expect("'{' or ';'", currentToken.Lexme)
	}

	//p.ST.Add(symboltable.)

	return ast.Fnc{Name:name, ReturnType: returnType, Params:params}
}


func (p *Parser) ParseParamList() *[]ast.Decl {
	//todo: add symbol table support , fix function
	p.readNextToken() // Skip the LPAR token
	var params []ast.Decl
	var stat string
	var id string
	var Type Types.EplType

	for p.match(epllex.COMMA) || !p.match(epllex.RPAR) {
		if p.match(epllex.FIXED) {
			stat = "fixed"
			p.readNextToken()
		} else if p.match(epllex.ID) {
			stat = "unfixed"
		} else {
			p.expect("Ident or 'fixed' tokens", currentToken.Lexme)
		}

		if p.match(epllex.ID) {
			id = currentToken.Lexme
			p.readNextToken()
		} else {
			p.expect("identifier token", currentToken.Lexme)
		}

		if Types.IsValidBasicType(currentToken) {
			Type = *Types.ResolveType(currentToken)
		} else {
			Type = *Types.MakeType(stat + ":"+currentToken.Lexme)
		}

		p.readNextToken() //skip the type

		params = append(params, &ast.VarDecl{Name: id, VarType: &Type, Stat: ast.VarStat(stat)})
		if p.match(epllex.RPAR) {
			break
		}
		p.readNextToken()//skip the comma
	}
	p.readNextToken() // skip the rpar token

	return &params
}

func (p *Parser) ParseBlock(function bool) {
	if !function {
		p.NewScope()
	}

	if p.match(epllex.LBRACE) {
		p.readNextToken()
	}
}

func (p *Parser) ParseVarDecl(scope symboltable.ScopeType) ast.Decl {
	if p.match(epllex.DECL) { //skip the decl keyword if found
		p.readNextToken()
	}

	var stat string
	var id string
	var Type Types.EplType

	if p.match(epllex.FIXED) {
		stat = "fixed"
		p.readNextToken()
	} else if p.match(epllex.ID) {
		stat = "unfixed"
	} else {
		p.expect("Ident or 'fixed' tokens", currentToken.Lexme)
	}

	if p.match(epllex.ID) {
		id = currentToken.Lexme
		p.readNextToken()
	} else {
		p.expect("identifier token", currentToken.Lexme)
	}

	if Types.IsValidBasicType(currentToken) {
		Type = *Types.ResolveType(currentToken)
	} else {
		Type = *Types.MakeType(stat + ":"+currentToken.Lexme)
	}

	p.readNextToken() //skip the type

	if !p.match(epllex.SEMICOLON) {
		if p.match(epllex.ASSIGN) {
			p.ParseExpression()
		} else {
			p.expect("';'", currentToken.Lexme)
		}
	} else {
		p.readNextToken() //todo: skip the semicolon
	}
	p.ST.Add(symboltable.NewSymbol(id, Type, scope))

	return &ast.VarDecl{Name: id, VarType: &Type, Stat: ast.VarStat(stat)}
}

//----------------------------------------------------------------------------------------------------------------------
//Expressions

/*
	Possible starts are:
		ID
		(
		-
		+

 */

var UnaryStart = []epllex.TokenType{epllex.MINUS, epllex.PLUS}
var BinaryStart = []epllex.TokenType{epllex.ID, epllex.LPAR}

func (p *Parser) ParseExpression() ast.Expression {
	if p.match(epllex.SEMICOLON) {
		p.readNextToken()
		return ast.EmptyExpr{}
	} else if p.match(epllex.ID) && p.match_n(epllex.SEMICOLON) {
		ident := p.ParseIdent()
		p.readNextToken() // skip semicolon
		return ident
	} else if p.matchTokens(UnaryStart) {
		return p.ParseUnaryOp()
	} else if p.matchTokens(BinaryStart) {
		return p.ParseBinaryOp()
	} else if p.match(epllex.ID) {
		if p.match_n(epllex.LPAR) || p.match_n(epllex.COMMA) {
			return p.ParseFunctionCall()
		} else if p.match_n(epllex.SEMICOLON){
			return p.ParseSingularExpr()
		} else {
			p.expect("function call or empty expression", currentToken.Lexme)
		}
	} else {
		p.expect("Ident semicolon unary or binary expression ", currentToken.Lexme)
	}
	return nil
}
func (p *Parser) ParseSingularExpr() ast.Singular {
	ident := p.ParseIdent()
	p.readNextToken() //Skip the semicolon
	return ast.Singular{Symbol: ident}
}

func (p *Parser) ParseIdent() ast.Ident {
	if p.match(epllex.ID) {
		p.readNextToken()
		return ast.Ident{Name: currentToken.Lexme}
	} else {
		p.expect("Ident", currentToken.Lexme)
	}

	return ast.Ident{}
}

func (p *Parser) ParseFunctionCall() ast.FunctionCall {
	var params []ast.Ident
	var importPath []ast.Ident

	for p.match(epllex.COMMA) && !p.match(epllex.RPAR) {
		importPath = append(importPath, p.ParseIdent())
		p.readNextToken()
	}

	if p.match(epllex.LPAR) {
		p.readNextToken()
		for p.match(epllex.COMMA) && !p.match(epllex.RPAR) {
			params = append(params, p.ParseIdent())
		}
	} else {
		p.expect("(", currentToken.Lexme)
	}

	return ast.FunctionCall{PackagePath: importPath, Arguments: params}
}

func (p *Parser) ParseBinaryOp() ast.Expression {
	switch currentToken.Ttype {
	case epllex.LPAR:
		p.readNextToken()
		return p.ParseBinaryOp()
	case epllex.RPAR:
		p.readNextToken()
		return p.ParseBinaryOp()
	}

	if p.matchTokens(UnaryStart) {
		leftExp := p.ParseExpression()

		p.readNextToken() //read the operation token
		op := currentToken.Ttype

		rightExp := p.ParseExpression()

		switch op {
		case epllex.MULT:
			return ast.BinaryMul{Ls: leftExp, Rs: rightExp}
		case epllex.MINUS:
			return ast.BinarySub{Ls: leftExp, Rs: rightExp}
		case epllex.DEV:
			//todo: rename token and ast names v0.1.1
			return ast.BinaryDiv{Ls: leftExp, Rs: rightExp}
		case epllex.PLUS:
			return ast.BinaryAdd{Ls: leftExp, Rs: rightExp}
		}
	}
}

func (p *Parser) ParseUnaryOp() ast.Expression {
	switch currentToken.Ttype {
	case epllex.PLUS:
		p.readNextToken()
		if p.match(epllex.ID) && p.match_n(epllex.SEMICOLON) {
			ident := p.ParseIdent()
			p.readNextToken()
			return ast.UnaryPlus{Rs: ident}
		}
	case epllex.MINUS:
		p.readNextToken()
		if p.match(epllex.ID) && p.match_n(epllex.SEMICOLON) {
			ident := p.ParseIdent()
			p.readNextToken()
			return ast.UnaryMinus{Rs: ident}
		}
	}

	p.expect("'+' or '-' tokens", currentToken.Lexme)
	return ast.EmptyExpr{}
}

//----------------------------------------------------------------------------------------------------------------------
//Statements

func (p *Parser) ParseIf() ast.IfStmt {
	if p.match(epllex.IF) {
		p.readNextToken()
	}
	p.NewScope()

}
//----------------------------------------------------------------------------------------------------------------------
//State Machine control and support functions

func (p *Parser) matchTokens(tokens []epllex.TokenType) bool {
	for _, token := range tokens {
		if p.match(token) {
			return true
		}
	}

	return false
}

func (p *Parser) match(t epllex.TokenType) bool {
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
