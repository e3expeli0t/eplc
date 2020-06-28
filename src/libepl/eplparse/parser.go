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

//New create new eplparse struct
func New(lx *epllex.Lexer) *Parser {
	return &Parser{
		Lexer: lx,
		ST:    symboltable.New(),
	}
}

/*
	Parser. the eplparse job is to take tokenized stream from the epllex
	and construct a tree form it, the tree is called AST (Abstract Syntax Tree)
	by the set of rules that the language grammar produce
	There are couple of eplparse kinds, in this version of epl (bootstrap) we are going
	to use a eplparse that called predictive eplparse (The grammar class is LL(k)).
	In the future  i'm planning to implement LR(0) eplparse
*/
type Parser struct {
	Lexer       *epllex.Lexer
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

//Parse main language constructs
//------------------------------------------------------------------------------------------------------------------------------------------------

func (p *Parser) ParseProgramFile() ast.Node {
	//check if the type handler is initialized
	if !p.Ok {
		Output.PrintFatalErr("eplc: Parser couldn't be initialized")
	}

	p.readNextToken()
	p.NewScope()

	var imports []ast.Import
	var decls []ast.Decl
	var fncs []*ast.Fnc
	var MainFunction *ast.Fnc

	for p.match(epllex.IMPORT) {
		Output.PrintLog("Parsing imports")
		imports = append(imports, p.ParseImport())
	}

	if p.match(epllex.DECL) {
		Output.PrintLog("Parsing global variables")
		for p.match(epllex.DECL) {
			p.ST.SetScopeType(symboltable.GLOBAL)
			decls = append(decls, p.ParseVarDecl(false))
		}
	}

	if p.match(epllex.FNC) {
		Output.PrintLog("Parsing functions")
		for p.match(epllex.FNC) {
			fnc := p.ParseFnc()
			//todo: Should change in eplc v0.2 to support cflags
			if fnc.Name.Name == "Main" {
				MainFunction = fnc
				continue
			}
			fncs = append(fncs, fnc)
		}
	}

	//check if any errors occurred
	p.Errors.IsValidFile()

	return &ast.ProgramFile{
		Imports:      &imports,
		GlobalDecls:  &decls,
		Symbols:      &p.ST,
		Functions:    &fncs,
		FileName:     p.Lexer.Filename,
		MainFunction: MainFunction,
	}
}

func (p *Parser) ParseImport() ast.Import {

	var importPath []string
	start := p.currentToken.StartLine

	for !p.match(epllex.SEMICOLON) {
		if p.match(epllex.ID) {
			importPath = append(importPath, p.currentLexme)
		}
		p.readNextToken() // skip the dot
	}

	p.readNextToken() //skip the semicolon
	return ast.Import{StartLoc: start, Imports: &importPath}
}

func (p *Parser) ParseBlock(keepScope bool) *ast.Block {
	if p.match(epllex.LBRACE) {
		p.readNextToken()
	}

	if !keepScope {
		p.NewScope()
		p.ST.SetScopeType(symboltable.BLOCK)
	}

	var contents []ast.Expression

	for !p.match(epllex.RBRACE) && !p.match(epllex.EOF) {
		switch p.currentToken.Ttype {
		case epllex.DECL:
			contents = append(contents, p.ParseVarDecl(false))
		case epllex.IF:
			contents = append(contents, p.ParseIf())
		case epllex.FOR:
			contents = append(contents, p.ParseForLoop())
		case epllex.UNTIL:
			contents = append(contents, p.ParseUntil())
		case epllex.REPEAT:
			contents = append(contents, p.ParseRepeatFam())
		case epllex.LBRACE:
			Output.PrintLog("In loop current token:", p.currentLexme, "Next: ", p.lookahead.Lexme)
			contents = append(contents, p.ParseBlock(true))
		default:
			switch p.currentToken.Ttype {
			case epllex.RETURN:
				p.readNextToken()
				if !p.match(epllex.SEMICOLON) {
					contents = append(contents,
						&ast.Return{Value: p.ParseExpression()},
					)
					continue
				} else {
					p.readNextToken()
					contents = append(contents,
						ast.Return{},
					)
					continue
				}
			case epllex.BREAK:
				p.readNextToken()
				if p.match(epllex.SEMICOLON) {
					p.readNextToken()
				}

				contents = append(contents, &ast.Break{})
				continue
			}
			if p.match(epllex.ID) && p.match_n(epllex.ASSIGN) {
				contents = append(contents, p.ParseAssignStmt())
			} else if p.matchTokens(ExpressionsMAP) {
				contents = append(contents, *p.ParseExpression())
			} else {
				p.report(fmt.Sprintf("Unable to identify: %s", p.currentLexme))
			}
		}
	}

	if p.match(epllex.EOF) {
		p.expect("'}'")
	} else {
		p.readNextToken() // skip }
	}

	return &ast.Block{
		Symbols:  &p.ST,
		ExprList: &contents,
	}
}

// fnc name([param_list])[type] block
func (p *Parser) ParseFnc() *ast.Fnc {
	if p.match(epllex.FNC) {
		p.readNextToken() //Skipping the fnc keyword
	}

	//create new scope for function
	p.NewScope()
	p.ST.SetScopeType(symboltable.FUNCTION)

	var name *ast.Ident
	var params *[]ast.Decl
	var returnType *Types.EplType
	var code *ast.Block

	if p.match(epllex.ID) {
		name = p.ParseIdent()
	} else {
		p.expect("Ident")
	}

	if p.match(epllex.LPAR) {
		params = p.ParseParamList()
	} else {
		p.expect("(")
	}

	if p.match(epllex.RETURN_TYPE_IND) {
		p.readNextToken()
		returnType = p.ParseType()
	} else {
		returnType = p.TypeHandler.MakeType("None")
	}

	if p.match(epllex.LBRACE) {
		code = p.ParseBlock(true)
	} else if p.match(epllex.SEMICOLON) {
		p.readNextToken()
	} else {
		p.expect("'{' or ';'")
	}

	return &ast.Fnc{
		Name:       name,
		ReturnType: returnType,
		Params:     params,
		Body:       code,
	}
}

// param_list := "(" decl [, param_decl] ")"
func (p *Parser) ParseParamList() *[]ast.Decl {
	if p.match(epllex.LPAR) {
		p.readNextToken() // skip LPAR
	}

	var params []ast.Decl

	for p.match(epllex.COMMA) || !p.match(epllex.RPAR) {
		params = append(params, p.decl())

		if p.match(epllex.RPAR) {
			break
		}

		if p.match(epllex.COMMA) {
			p.readNextToken()
		} else {
			p.expect("','")
		}
	}
	p.readNextToken() // skip the rpar token

	for _, decl := range params {
		p.AddToST(decl)
	}
	return &params
}

// decl := "decl" [status] ident type
// var_decl := decl ";" | decl value_assign ";"
// scoped_var_decl := decl [value_assign]
// value_assign := "=" Expression
//status := "fixed" | "mutable"
//todo: support auto value detection
func (p *Parser) ParseVarDecl(scoped bool) ast.Decl {
	var value *ast.Expression

	varDec := p.decl()
	p.AddToST(varDec)

	if p.match(epllex.ASSIGN) {
		value = p.ParseValueAssign(!scoped)

		return &ast.VarExplicitDecl{
			VarDecl: *varDec,
			Value:   value,
		}
	}

	if scoped {
		return varDec
	}

	if p.match(epllex.SEMICOLON) {
		p.readNextToken()
		return varDec
	} else {
		p.expect("';'")
	}

	return varDec
}

func (p *Parser) decl() *ast.VarDecl {
	if p.match(epllex.DECL) {
		p.readNextToken()
	}

	var name *ast.Ident
	var status ast.VarStat
	var varType *Types.EplType

	if p.match(epllex.FIXED) {
		status = ast.Fixed
		p.readNextToken()
	} else {
		status = ast.Mutable
	}

	if p.match(epllex.ID) {
		name = p.ParseIdent()
	} else {
		p.expect("variable name")
	}

	varType = p.ParseType()

	return &ast.VarDecl{
		Name:    name,
		VarType: varType,
		Stat:    status,
	}
}

//if the semicolon flag is set this function uses ParseExpression()
func (p *Parser) ParseValueAssign(semicolon bool) *ast.Expression {
	var value ast.Expression

	if !p.match(epllex.ASSIGN) {
		p.expect("'='")
	}
	p.readNextToken()

	if semicolon {
		value = *p.ParseExpression()
	} else {
		value = p.ParseExpr(0)
	}

	return &value
}

//----------------------------------------------------------------------------------------------------------------------
//Statements

//Design note:
//If is an expression that return None
// IfStmt := "if" Expression block
func (p *Parser) ParseIf() ast.Statement {
	if p.match(epllex.IF) {
		p.readNextToken()
	}

	condition := p.ParseExpr(0)

	if !p.match(epllex.LBRACE) {
		p.expect("{ or ,")
	}

	// Since we don't allow variable declaration in the if header
	// there is no point to keep the current symboltable
	code := p.ParseBlock(false)

	if p.match(epllex.ELSE) {
		elseCode := p.ParseElse()
		return &ast.IfStmt{
			Code:      code,
			Else:      &elseCode,
			Condition: &condition,
		}
	}

	return &ast.IfStmt{
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

	if !p.match(epllex.LBRACE) {
		p.expect("'{'")
	}
	code := p.ParseBlock(false)

	return &ast.ElseStmt{Code: code}
}

//Design note:
// the language does not allow expressions like : expr =  expr therefore the assign is a statement
func (p *Parser) ParseAssignStmt() *ast.AssignStmt {
	var Owner *ast.Ident
	var Value *ast.Expression

	if p.match(epllex.ID) {
		Owner = p.ParseIdent()
	} else {
		p.expect("Ident")
	}

	Value = p.ParseValueAssign(true)
	return &ast.AssignStmt{
		Owner: Owner,
		Value: Value,
	}
}

/*
	There is three kinds of loops
	For: regular loop
	Until: like while
	Repeat: infinite loop
	Repeat-Until: repeats until condition is met
*/

//until := "until" bool_expr "{" expression "}"
func (p *Parser) ParseUntil() *ast.Until {
	if p.match(epllex.UNTIL) {
		p.readNextToken()
	}
	p.NewScope()

	var cond ast.Expression
	var code *ast.Block

	cond = p.ParseExpr(0)

	if !p.match(epllex.LBRACE) {
		p.expect("'{'")
	}

	code = p.ParseBlock(false)

	return &ast.Until{
		Condition: &cond,
		Code:      code,
	}
}

//for := "for" for_header block
func (p *Parser) ParseForLoop() *ast.ForLoop {
	if p.match(epllex.FOR) {
		p.readNextToken() // skip for token
	}

	//New scope for loop
	p.NewScope()

	decl, cond, expr := p.forHeader()

	if !p.match(epllex.LBRACE) {
		p.expect("'{'")
	}
	code := p.ParseBlock(true)

	return &ast.ForLoop{
		VarDef:    &decl,
		Condition: &cond,
		Expr:      &expr,
		Code:      code,
	}
}

//Parses the if header
//for_header :=  [var_decl] ";" [Expression] ";" [Expression] ";"
func (p *Parser) forHeader() (ast.Decl, ast.Expression, ast.Expression) {
	var exp ast.Expression
	var cond ast.Expression
	var decl ast.Decl

	semis := 0 // we need to keep track of the semicolons

	//parse variable declaration
	if !p.match(epllex.SEMICOLON) {
		//note: this call "eats" the semicolon
		//note: we allow the variable declaration to be without decl
		decl = p.ParseVarDecl(false)
	} else {
		semis++
		p.readNextToken() // skip the semicolon
	}

	//parse condition
	if !p.match(epllex.SEMICOLON) {
		cond = *p.ParseExpression()
	} else {
		semis++
		p.readNextToken() //skip the semi
	}

	if !p.match(epllex.RBRACE) {
		exp = p.ParseExpr(0)
	}

	return decl, cond, exp
}

//parse repeat and repeat like loops
func (p *Parser) ParseRepeatFam() ast.Statement {
	if p.match(epllex.REPEAT) {
		p.readNextToken()
	}

	//create new scope here since we have scoped variable declaration
	p.NewScope()

	var varDecl ast.Decl
	var code *ast.Block
	var cond *ast.Expression

	if p.match(epllex.LPAR) {
		p.readNextToken()
		varDecl = p.ParseVarDecl(true)

		if !p.match(epllex.RPAR) {
			p.expect("') in repeat loop variable declaration'")
		}
		p.readNextToken()
	}

	if !p.match(epllex.LBRACE) {
		p.expect("'{'")
	}

	code = p.ParseBlock(true)

	if p.match(epllex.UNTIL) {
		p.readNextToken()

		cond = p.ParseExpression()

		if cond == nil {
			p.expect("condition")
		}

		return &ast.RepeatUntil{
			VarDef:    &varDecl,
			Condition: cond,
			Code:      code,
		}

	}

	return &ast.Repeat{
		VarDef: &varDecl,
		Code:   code,
	}

}

//----------------------------------------------------------------------------------------------------------------------
//State Machine control and support functions

//todo: find other place. Placed here to solve cyclic import...
//convert raw Token value to boolean node
func (p *Parser) ToBoolVal(t epllex.Token) ast.Boolean {
	if t.Ttype == epllex.FALSE {
		return ast.Boolean{Val: ast.BOOL_FALSE}
	}

	return ast.Boolean{Val: ast.BOOL_TRUE}
}

func (p *Parser) AddToST(decl ast.Decl) {
	switch t := decl.(type) {
	case *ast.VarDecl:
		p.ST.AddWOScope(symboltable.NewTypedSymbol(t.Name.Name, *t.VarType))
	}
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

//--------------------------------------------------------------------------------------
//Helper functions
func (p *Parser) ParseType() (Type *Types.EplType) {
	defer p.readNextToken() // skip type
	Type = p.processType()
	return
}

func (p *Parser) processType() (Type *Types.EplType) {

	if p.TypeHandler.IsValidBasicType(p.currentToken) {
		Type = p.TypeHandler.GetType(p.currentLexme)
	} else if p.match(epllex.ID) {
		Type = p.TypeHandler.MakeType(p.currentLexme)
	} else {
		p.expect("type")
	}
	return
}

func (p *Parser) InitializeTypeHandler() {
	p.TypeHandler.Initialize(p.Lexer)
	p.Ok = true
}

//Error handling functions
func (p *Parser) expect(ex string) {
	p.report(fmt.Sprintf("expected %s found %s ", ex, p.currentLexme))
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
