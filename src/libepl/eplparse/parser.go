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
	"eplc/src/libepl/epllex"
	"eplc/src/libepl/eplparse/Types"
	"eplc/src/libepl/eplparse/ast"
	"eplc/src/libepl/eplparse/errors"
	"eplc/src/libepl/eplparse/symboltable"
	"fmt"
)



//New create new parser struct
func New(lx epllex.Lexer) Parser {
	return Parser{
		Lexer:        lx,
		ST:           symboltable.New(),
	}
}

/*
	Parser. the parser job is to take tokenized stream from the lexer
	and construct a tree form it, the tree is called AST (Abstract Syntax Tree)
	by the set of rules that the language grammar produce
	There are couple of parser kinds, in this version of epl (bootstrap) we are going
	to use a parser that called predictive parser (The grammar class is LL(k)).
	In the future  i'm planning to implement LR(0) parser
*/
type Parser struct {
	Lexer       epllex.Lexer
	ST          symboltable.SymbolTable
	Errors      errors.InternalParserError
	TypeHandler Types.TypeSystem
	
	//private fields
	currentToken epllex.Token
	lookahead    epllex.Token
	currentLexme string

	// set to true if type handler is Initialized
	Ok bool
}

//Todo: change expect so that the function checks if the expected token is coming (if not print err)
//--------------------------------------------------------------------------------------
//Helper functions

func(p *Parser) InitializeTypeHandler() {
	p.TypeHandler.Initialize(p.Lexer)
	p.Ok = true
}

//Error handling functions
func (p *Parser) expect(ex string, fnd string) {

	p.report(fmt.Sprintf("expected %s found %s ", ex, fnd))
}

func (p *Parser) report(msg string) {
	p.Errors.ParsingError(p.Lexer.Filename, p.Lexer.Line, p.Lexer.LineOffset, msg, p.Lexer.GetLine(), p.currentToken)
}

//Scope handling functions
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


//Parse main language constructs
//------------------------------------------------------------------------------------------------------------------------------------------------


func (p *Parser) ParseProgramFile() ast.Node {
	//check if the type handler is initialized
	if !p.Ok {
		Output.PrintFatalErr("Error in parser state. Exiting")
	}

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
			decls = append(decls, p.ParseVarDecl(symboltable.GLOBAL, false))
		}
	}

	if p.match(epllex.FNC) {
		Output.PrintLog("Parsing functions")
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

	return &ast.ProgramFile{
		Imports: &imports,
		GlobalDecls: &decls,
		Symbols: &p.ST,
		Functions: &fncs,
		FileName: p.Lexer.Filename,
		MainFunction: &MainFunction,
	}
}


//----------------------------------------------------------------------------------------------------------------------
//Statements

func (p *Parser) ParseImport() ast.Import {

	var importPath []string
	start := p.currentToken.StartLine

	for !p.match(epllex.SEMICOLON) {
		if p.match(epllex.ID) {
			importPath = append(importPath, p.currentLexme)
		}
		p.readNextToken() //SKIP THE DOT
	}

	p.readNextToken() // SKIP THE SEMICOLON
	return ast.Import{StartLoc: start, Imports: importPath}
}

// fnc name([param_list])[type] block
func (p *Parser) ParseFnc() ast.Fnc {
	if p.match(epllex.FNC) {
		p.readNextToken() //Skipping the fnc keyword
	}

	//create new scope for function
	p.NewScope()

	var name string
	var params *[]ast.Decl
	var returnType *Types.EplType
	var code *ast.Block

	if p.match(epllex.ID) {
		name = p.currentLexme
		fmt.Println("name:", name, "next:", p.lookahead)
		p.readNextToken()
	} else {
		p.expect("Ident", p.currentLexme)
	}

	if p.match(epllex.LPAR) {
		params = p.ParseParamList()
	} else {
		p.expect("(", p.currentLexme)
	}


	if p.match(epllex.RETURN_IND) {
		p.readNextToken()

		if p.TypeHandler.IsValidBasicType(p.currentToken) {
			returnType = p.TypeHandler.ResolveType(p.currentToken)
		} else {
			returnType = p.TypeHandler.MakeType(p.currentLexme)
		}
		p.readNextToken() // skip type
	} else {
		returnType = p.TypeHandler.MakeType("None")
	}

	if p.match(epllex.LBRACE) {
		code = p.ParseBlock(true)
	} else if p.match(epllex.SEMICOLON) {
		p.readNextToken()
	} else {
		p.expect("'{' or ';'", p.currentLexme)
	}

	//p.ST.Add(symboltable.)

	return ast.Fnc {
		Name: name,
		ReturnType: returnType,
		Params: params,
		Body: code,
	}
}

//todo: add symbol table support , fix function
func (p *Parser) ParseParamList() *[]ast.Decl {
	if p.match(epllex.LPAR) {
		p.readNextToken() // skip LPAR
	}

	var params []ast.Decl

	for p.match(epllex.COMMA) || !p.match(epllex.RPAR) {
		params = append(params, p.parseParam())

		if p.match(epllex.RPAR) {
			break
		}

		if p.match(epllex.COMMA) {
			p.readNextToken()
		} else {
			p.expect("','", p.currentLexme)
		}
	}
	p.readNextToken() // skip the rpar token

	return &params
}


//this is dup of ParseVarDecl used to simplify the process of parsing paramList
//note: this is a TEMPORARY solution.
func (p *Parser) parseParam() ast.Decl {

	var stat string
	var id string
	var Type Types.EplType

	if p.match(epllex.FIXED) {
		stat = "fixed"
		p.readNextToken()
	} else if p.match(epllex.ID) {
		stat = "unfixed"
	} else {
		p.expect("Ident or 'fixed' tokens", p.currentLexme)
	}

	if p.match(epllex.ID) {
		id = p.currentLexme
		p.readNextToken()
	} else {
		p.expect("identifier token", p.currentLexme)
	}


	if p.TypeHandler.IsValidBasicType(p.currentToken) {
		Type = *p.TypeHandler.ResolveType(p.currentToken)
	} else {
		Type = *p.TypeHandler.MakeType(stat + ":" + p.currentLexme)
	}

	p.readNextToken() //skip the type

	p.ST.Add(symboltable.NewSymbol(id, Type, symboltable.FUNCTION))
	return &ast.VarDecl{Name: id, VarType: &Type, Stat: ast.VarStat(stat)}

}


func (p *Parser) ParseBlock(keepST bool) *ast.Block {
	var scope symboltable.ScopeType

	if !keepST {
		p.NewScope()
		scope = symboltable.BLOCK
	} else {
		scope = symboltable.FUNCTION
	}

	// if found {
	if p.match(epllex.LBRACE) {
		p.readNextToken()
	}

	var contents []ast.Expression

	// as long as there is not } or EOF
	for !p.match(epllex.RBRACE) && !p.match(epllex.EOF){
		switch p.currentToken.Ttype {
		case epllex.LPAR:
			//todo: fix all related ast problems
			contents = append(contents, *p.ParseBlock(false))
		case epllex.IF:
			contents = append(contents, p.ParseIf())
		case epllex.DECL:
			contents = append(contents, p.ParseVarDecl(scope, false))
		case epllex.REPEAT:
			contents = append(contents, p.ParseRepeatFam())
		case epllex.FOR:
			contents = append(contents, p.ParseForLoop())
		case epllex.UNTIL:
			contents = append(contents, p.ParseUntil())
		default:
			if p.match(epllex.ID){
				if p.match_n(epllex.ASSIGN) {
					contents = append(contents, p.ParseAssignStmt())
				} else if p.match_n(epllex.SEMICOLON) {
					contents = append(contents, p.ParseSingularExpr())
				}else {
					contents = append(contents, p.ParseExpression())
				}
			} else {
				break // invalid block structure
			}
		}
	}

	if p.match(epllex.RBRACE) {
		p.readNextToken()
	} else {
		p.expect("}", p.currentLexme)
	}

	st := p.ST.Strip() // for cleaner tree

	return &ast.Block{
		Symbols: &st,
		ExprList:&contents,
	}
}

//todo: code cleanup
//if the header flag is set the function will not search for semicolon after the def
func (p *Parser) ParseVarDecl(scope symboltable.ScopeType, lheader bool) ast.Decl {
	if p.match(epllex.DECL) {
		p.readNextToken() //skip the decl keyword if found
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
		p.expect("Ident or 'fixed' tokens", p.currentLexme)
	}

	if p.match(epllex.ID) {
		id = p.currentLexme
		p.readNextToken()
	} else {
		p.expect("identifier token", p.currentLexme)
	}


	if p.TypeHandler.IsValidBasicType(p.currentToken) {
		Type = *p.TypeHandler.ResolveType(p.currentToken)
	} else {
		Type = *p.TypeHandler.MakeType(stat + ":" + p.currentLexme)
	}

	p.readNextToken() //skip the type

	if !p.match(epllex.SEMICOLON) {
		if p.match(epllex.ASSIGN) {
			p.readNextToken()

			var exp ast.Expression

			if lheader {
				exp = p.expr(0) //ignore semicolon
			} else {
				exp = p.ParseExpression()
			}

			p.ST.Add(symboltable.NewSymbol(id, Type, scope))

			return ast.VarExplicitDecl{
				VarDecl: ast.VarDecl{
					Name:    id,
					VarType:  &Type,
					Stat:    ast.VarStat(stat),
				},

				Value:   &exp,
			}

		} else if !lheader {
			p.expect("';'", p.currentLexme)
		}
	} else {
			p.readNextToken() //skip the semicolon
	}
	p.ST.Add(symboltable.NewSymbol(id, Type, scope))

	return &ast.VarDecl{Name: id, VarType: &Type, Stat: ast.VarStat(stat)}
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

	condition := p.boolExpr(0)

	if p.match(epllex.SEMICOLON) {
		return ast.IfStmt{Condition: &condition}
	}

	if !p.match(epllex.LBRACE) {
		p.expect("{ or ,", p.currentLexme)
	}

	code := p.ParseBlock(false)

	if p.match(epllex.ELSE) {
		elseCode := p.ParseElse()
		return ast.IfStmt{Code: code, Else: &elseCode, Condition: &condition}
	}

	return ast.IfStmt{
		Condition: &condition,
		Code:      code,
		Else:      nil,
	}
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
		p.expect("Ident", p.currentLexme)
	}

	if p.match(epllex.ASSIGN) {
		p.readNextToken()
	} else {
		p.expect("=", p.currentLexme)
	}

	if p.matchTokens(BinaryStart) || p.matchTokens(UnaryStart)  {
		Value = p.ParseExpression()
	} else {
		p.expect("'(', '+', '-' or Ident", p.currentLexme)
	}

	return ast.AssignStmt{
		Owner: Owner,
		Value: &Value,
	}
}

/*
	There is tree kinds of loops
	For: regular loop
	Until: like while
	Repeat: infinite loop
	Repeat-Until: repeat until
 */

func (p *Parser) ParseUntil() ast.Until {
	if p.match(epllex.UNTIL) {
		p.readNextToken()
	}

	p.NewScope()

	var cond ast.BoolExpr
	var code *ast.Block

	if !p.matchTokens(BoolStartMap) {
		p.expect("'true', 'false', '!', '(' or ident", p.currentLexme)
	}

	cond = p.boolExpr(0)

	if !p.match(epllex.LBRACE) {
		p.expect("'{'", p.currentLexme)
	}


	code = p.ParseBlock(true)

	return ast.Until{
		Condition: &cond,
		Code:      code,
	}
}


func (p *Parser) ParseForLoop() ast.ForLoop {
	if p.match(epllex.FOR) {
		p.readNextToken() // skip for token
	}

	//New scope for loop
	p.NewScope()

	/*
	the for header looks like this:
	for [varDef]; [bool_expr]; [expr] "{"

	parsing method:
		if the next token is not SEMICOLON parse the corresponding part
	 */

	var varDecl ast.Decl
	var cond ast.BoolExpr
	var expr ast.Expression
	var code *ast.Block

	//parse vardecl
	if !p.match(epllex.SEMICOLON) {
		if p.match(epllex.DECL) {
			varDecl = p.ParseVarDecl(symboltable.FUNCTION, true)
		} else {
			p.expect("decl", p.currentLexme)
		}
	} else {
		varDecl = nil
		p.readNextToken() // skip the semicolon
	}

	//parse condition
	if !p.match(epllex.SEMICOLON) {
		if p.matchTokens(BoolStartMap) {
			//check: conditions end
			cond = p.boolExpr(0)
		} else {
			p.expect("'true', 'false', '!', '(' or ident", p.currentLexme)
		}
	} else {
		cond = nil
		p.readNextToken() // skip the semicolon
	}

	if !p.match(epllex.LPAR) {
		if p.matchTokens(BinaryStart) {
			if p.match_n(epllex.ASSIGN) {
				expr = p.ParseAssignStmt()
			} else {
				expr = p.expr(0)
			}
		} else {
			p.expect("'(', '+', '-' or Ident", p.currentLexme)
		}
	}

	code = p.ParseBlock(true)

	return ast.ForLoop{
		VarDef:    varDecl,
		Condition: cond,
		Expr:      expr,
		Code:      code,
	}
}

//parse repeat and repeat like loops
func (p *Parser) ParseRepeatFam() ast.Statement {
	if p.match(epllex.REPEAT) {
		p.readNextToken()
	}

	p.NewScope()

	var varDecl ast.Decl
	var code *ast.Block
	var cond ast.BoolExpr

	if p.match(epllex.LPAR) {
		varDecl = p.ParseVarDecl(symboltable.FUNCTION, true)
		if !p.match(epllex.RPAR) {
			p.expect("')'", p.currentLexme)
		}
	}

	if !p.match(epllex.LBRACE) {
		p.expect("'{'", p.currentLexme)
	}

	code = p.ParseBlock(true)

	if p.match(epllex.UNTIL) {
		if !p.matchTokens(BoolStartMap) {
			p.expect("'true', 'false', '!', '(' or ident", p.currentLexme)
		}

		cond = p.boolExpr(0)

		return ast.RepeatUntil{
			VarDef:    varDecl,
			Condition: &cond,
			Code:      code,
		}

	}
	return ast.Repeat{
		VarDef: varDecl,
		Code:   code,
	}
}


//----------------------------------------------------------------------------------------------------------------------
//State Machine control and support functions

//todo: find other place. Placed here to solve cyclic import...
//convert raw Token value to boolean node
func (p *Parser) ToBoolVal(t epllex.Token) ast.Boolean {
	if t.Ttype == epllex.FALSE {
		return  ast.Boolean{Val: ast.BOOL_FALSE}
	}

	return ast.Boolean{Val:ast.BOOL_TRUE}
}

func (p *Parser) isBinary(exp ast.Expression) bool {
	switch exp.(type) {
	case ast.BinaryAdd, ast.BinaryDiv, ast.BinaryMul, ast.BinarySub:
		return true
	case ast.BoolGreatEquals, ast.BoolGreaterThen,
	ast.BoolEquals, ast.BoolLowerThen, ast.BoolLowerThenEqual,
	ast.BoolNotEquals:
		return true

	case ast.Ident, ast.UnaryMinus,
	ast.UnaryPlus, ast.String, ast.FunctionCall, ast.Number:
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
// check if the p.lookahead equal to one of the list elem's
func (p *Parser) matchNTokens(tokens []epllex.TokenType) bool {
	for _, token := range tokens {
		if p.match_n(token) {
			return true
		}
	}

	return false
}

func (p *Parser) match(t epllex.TokenType) bool {
	return p.currentToken.Ttype == t
}

func (p *Parser) match_n(t epllex.TokenType) bool {
	return p.lookahead.Ttype == t
}

func (p *Parser) readNextToken() {
	if (epllex.Token{}) == p.currentToken {
		p.currentToken = p.Lexer.Next()
	} else {
		p.currentToken = p.lookahead
	}
	p.lookahead = p.Lexer.Next()
	p.currentLexme = p.currentToken.Lexme
}
