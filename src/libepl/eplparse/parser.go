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
	errCount     uint = 0
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
	Lexer  epllex.Lexer
	ST     symboltable.SymbolTable
	Errors errors.InternalParserError
}

//Todo: change expect so that the function checks if the expected token is coming (if not print err)
//--------------------------------------------------------------------------------------
//Helper functions

func (p *Parser) expect(ex string, fnd string) {
	p.report(fmt.Sprintf("expected %s found %s ", ex, fnd))
}

func (p *Parser) report(msg string) {
	p.Errors.ParsingError(p.Lexer.Filename, p.Lexer.Line, p.Lexer.LineOffset, msg, p.Lexer.GetLine(), currentToken)
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

func (p *Parser) ParseProgramFile() ast.Node {
	p.readNextToken()
	p.NewScope()

	var imports []ast.Import
	var decls []ast.Decl
	var fncs []ast.Fnc
	var MainFunction ast.Fnc

	for p.match(epllex.IMPORT) {
		Output.PrintLog("Parsing imports")
		imports = append(imports, p.ParseImport())
	}

	if p.match(epllex.DECL) {
		Output.PrintLog("Parsing global variables")
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
			Output.PrintLog("Found: ", fnc.Name)
			if fnc.Name == "Main" {
				Output.PrintLog("Found Main function")
				MainFunction = fnc
				continue
			}
			fncs = append(fncs, fnc)
		}
	}

	//check if any errors occurred
	p.Errors.IsValidFile()

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

// fnc name([param_list])[type] block
func (p *Parser) ParseFnc() ast.Fnc {
	p.NewScope()
	if p.match(epllex.FNC) {
		p.readNextToken() //Skipping the fnc keyword
	}

	var name string
	var params *[]ast.Decl
	var returnType *Types.EplType
	var code *ast.Block

	if p.match(epllex.ID) {
		name = currentToken.Lexme
		fmt.Println("name:", name, "next:", lookahead)
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
		p.readNextToken() // skip type
	} else {
		returnType = Types.MakeType("None")
	}

	if p.match(epllex.LBRACE) {
		code = p.ParseBlock(true)
	} else if p.match(epllex.SEMICOLON) {
		p.readNextToken()
	} else {
		p.expect("'{' or ';'", currentToken.Lexme)
	}

	//p.ST.Add(symboltable.)

	return ast.Fnc{Name: name, ReturnType: returnType, Params: params, Body: code}
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
			Type = *Types.MakeType(stat + ":" + currentToken.Lexme)
		}

		p.readNextToken() //skip the type

		p.ST.Add(symboltable.NewSymbol(id, Type, symboltable.FUNCTION))
		params = append(params, &ast.VarDecl{Name: id, VarType: &Type, Stat: ast.VarStat(stat)})

		if p.match(epllex.RPAR) {
			break
		}
		p.readNextToken() //skip the comma
	}
	p.readNextToken() // skip the rpar token

	return &params
}

/*
	todo: implement:
			loops
			VarExplicitDecl
*/
func (p *Parser) ParseBlock(function bool) *ast.Block {
	var scope symboltable.ScopeType

	if !function {
		p.NewScope()
		scope = symboltable.BLOCK
	} else {
		scope = symboltable.FUNCTION
	}

	//skipping block contents for testing prep

	// if found {
	if p.match(epllex.LBRACE) {
		p.readNextToken()
	}

	var contents []ast.Expression

	// as long as there is not } or EOF
	for !p.match(epllex.RBRACE) && !p.match(epllex.EOF){
		switch currentToken.Ttype {
		case epllex.LPAR:
			//todo: fix all related ast problems
			contents = append(contents, *p.ParseBlock(false))
		case epllex.IF:
			contents = append(contents, p.ParseIf())
		case epllex.DECL:
			contents = append(contents, p.ParseVarDecl(scope))
		default:
			if p.match(epllex.ID){
				if p.match_n(epllex.ASSIGN) {
					contents = append(contents, p.ParseAssignStmt())
				} else if p.match_n(epllex.SEMICOLON) {
					contents = append(contents, p.ParseSingularExpr())
				} else {
					contents = append(contents, p.ParseExpression())
				}
			}

		}
	}

	if p.match(epllex.RBRACE) {
		p.readNextToken()
	} else {
		p.expect("}", currentToken.Lexme)
	}
	return &ast.Block{Symbols: &p.ST}
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
		Type = *Types.MakeType(stat + ":" + currentToken.Lexme)
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
Design note:
	Possible starts are:
		ID
		(
		-
		+

*/

var UnaryStart = []epllex.TokenType{epllex.SUB, epllex.ADD}
var BinaryStart = []epllex.TokenType{epllex.ID, epllex.LPAR, epllex.NUM}
var BinaryEnd = []epllex.TokenType{epllex.RPAR, epllex.ID, epllex.NUM}

//Singular := Ident";"
func (p *Parser) ParseSingularExpr() ast.Singular {
	ident := p.ParseIdent()
	p.readNextToken() //Skip the semicolon
	return ast.Singular{Symbol: ident}
}

//Ident := ID (Basic string Token)
func (p *Parser) ParseIdent() ast.Ident {
	if p.match(epllex.ID) {
		p.readNextToken() //todo: logic dont make scense
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
		if p.match(epllex.RPAR) {
			p.readNextToken()
		} else {
			p.expect(")", currentToken.Lexme)
		}
	} else {
		p.expect("(", currentToken.Lexme)
	}

	return ast.FunctionCall{PackagePath: importPath, Arguments: params}
}
/*
Design note:
Expression := UnaryExpr Expression | UnaryExpr | Expression op Expression | ID | Singular | FunctionCall | e

Notes:
*	we parse the expressions using precedence climbing
*	credit to https://github.com/richardjennings/prattparser/
*	credit to https://www.engr.mun.ca/~theo/Misc/exp_parsing.htm
*	credit to https://github.com/golang/go
*/
func (p *Parser) ParseExpression() ast.Expression {
	exprs := p.expr(0)
	if !p.matchNTokens(BinaryEnd) {
		p.expect("')', ';', ident or number ", currentToken.Lexme)
		p.readNextToken()
	}

	return exprs
}

func (p *Parser) expr(rbp int) ast.Expression {
	var left ast.Expression

	left = p.nud(currentToken)

	for p.isBinary(left) && currentToken.Precedence() > rbp {
		//Left-Denotation: a simple binary expression
		if currentToken.IsLeftAssociative() {
			switch currentToken.Ttype {
			case epllex.ADD:
				left = ast.BinaryAdd{
					Ls: left,
					Rs: p.expr(currentToken.Precedence()),
				}
			case epllex.SUB:
				left = ast.BinarySub{
					Ls: left,
					Rs: p.expr(currentToken.Precedence()),
				}

			case epllex.MULT:
				left = ast.BinaryMul{
					Ls: left,
					Rs: p.expr(currentToken.Precedence()),
				}
			case epllex.DIV:
				left = ast.BinaryDiv{
					Ls: left,
					Rs: p.expr(currentToken.Precedence()),
				}
			}
			p.readNextToken()
		}
	}
	return left
}

//Null-Denotation: a simple unary expression
func (p *Parser) nud(t epllex.Token) ast.Expression {

	var left ast.Expression

	switch t.Ttype {
	case epllex.LPAR:
		p.readNextToken() //skip paren
		left = p.expr(0)

		// if there is not a ) symbol after expression
		if !p.match(epllex.RPAR) {
			p.expect(")", currentToken.Lexme)
		}

		p.readNextToken() //skip RPAR
	default:
		if t.IsUnary() {
			if t.Ttype == epllex.ADD {
				p.readNextToken()
				left = ast.UnaryPlus{Rs: p.expr(epllex.HighPrec)}
			} else if t.Ttype == epllex.SUB {
				p.readNextToken()
				left = ast.UnaryMinus{Rs: p.expr(epllex.HighPrec)}
			}
		} else if t.IsScalar() {
			p.readNextToken()
			return ast.Number{Value: t.Lexme}
		} else if t.IsIdent() {
			return ast.Singular{Symbol: p.ParseIdent()}
		}
	}

	return left
}

func (p *Parser) ParseBoolExpr() ast.BoolExpr {
	return ast.BoolEquals{} //todo: accualy parse
}

//----------------------------------------------------------------------------------------------------------------------
//Statements

//Design note:
//If is an expression that return None
// IfStmt := "if" boolExpr block
func (p *Parser) ParseIf() ast.Statement {
	if p.match(epllex.IF) {
		p.readNextToken()
	}

	condition := p.ParseBoolExpr()

	if p.match(epllex.SEMICOLON) {
		return ast.IfStmt{Condition: &condition}
	}

	if !p.match(epllex.LBRACE) {
		p.expect("{ or ,", currentToken.Lexme)
	}

	code := p.ParseBlock(false)

	if p.match(epllex.ELSE) {
		elseCode := p.ParseBlock(false)
		return ast.IfStmt{Code: code, Else: elseCode, Condition: &condition}
	}

	return ast.EmptyExpr{}
}

func (p *Parser) ParseElse() ast.Statement {
	if p.match(epllex.ELSE) {
		p.readNextToken()
	}

	if p.match(epllex.IF) {
		return p.ParseIf()
	}

	code := p.ParseBlock(false)

	return ast.ElseStmt{Code: code}
}

//Design note:
// the language does not allow expressions like : expr =  expr therefore the assign is a statement
func (p *Parser) ParseAssignStmt() ast.AssignStmt {
	var Owner ast.Ident
	var Value ast.Expression

	if p.match(epllex.ID) {
		Owner = p.ParseIdent()
	} else {
		p.expect("Ident", currentToken.Lexme)
	}

	if p.match(epllex.ASSIGN) {
		p.readNextToken()
	} else {
		p.expect("=", currentToken.Lexme)
	}

	if p.matchTokens(BinaryStart) || p.matchTokens(UnaryStart)  {
		Value = p.ParseExpression()
	} else {
		p.expect("'(', '+', '-' or Ident", currentToken.Lexme)
	}

	return ast.AssignStmt{
		Owner: Owner,
		Value: &Value,
	}
}

//----------------------------------------------------------------------------------------------------------------------
//State Machine control and support functions

func (p *Parser) isBinary(exp ast.Expression) bool {
	switch exp.(type) {
	case ast.BinaryAdd, ast.BinaryDiv, ast.BinaryMul, ast.BinarySub:
		return true
	}

	return false
}

func (p *Parser) matchTokens(tokens []epllex.TokenType) bool {
	for _, token := range tokens {
		if p.match(token) {
			return true
		}
	}

	return false
}

//Design note:
// check if the lookahead equal to one of the list elem's
func (p *Parser) matchNTokens(tokens []epllex.TokenType) bool {
	for _, token := range tokens {
		if p.match_n(token) {
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
