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
	Lexer  epllex.Lexer
	ST     symboltable.SymbolTable
	Errors errors.InternalParserError
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
	p.TypeHandler.Initialize(p.Lexer.Filename)
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

	return &ast.ProgramFile{Imports: &imports, GlobalDecls: &decls, Symbols: &p.ST, Functions: &fncs, FileName: p.Lexer.Filename, MainFunction: &MainFunction}
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
	p.NewScope()
	if p.match(epllex.FNC) {
		p.readNextToken() //Skipping the fnc keyword
	}

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

	// set the types error handler line descriptor
	p.TypeHandler.SetLine(p.Lexer.GetLine())

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
			p.expect("Ident or 'fixed' tokens", p.currentLexme)
		}

		if p.match(epllex.ID) {
			id = p.currentLexme
			p.readNextToken()
		} else {
			p.expect("identifier token", p.currentLexme)
		}

		// set the types error handler line descriptor
		p.TypeHandler.SetLine(p.Lexer.GetLine())

		if p.TypeHandler.IsValidBasicType(p.currentToken) {
			Type = *p.TypeHandler.ResolveType(p.currentToken)
		} else {
			Type = *p.TypeHandler.MakeType(stat + ":" + p.currentLexme)
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
		p.expect("Ident or 'fixed' tokens", p.currentLexme)
	}

	if p.match(epllex.ID) {
		id = p.currentLexme
		p.readNextToken()
	} else {
		p.expect("identifier token", p.currentLexme)
	}

	// set the types error handler line descriptor
	p.TypeHandler.SetLine(p.Lexer.GetLine())

	if p.TypeHandler.IsValidBasicType(p.currentToken) {
		Type = *p.TypeHandler.ResolveType(p.currentToken)
	} else {
		Type = *p.TypeHandler.MakeType(stat + ":" + p.currentLexme)
	}

	p.readNextToken() //skip the type

	if !p.match(epllex.SEMICOLON) {
		if p.match(epllex.ASSIGN) {
			p.ParseExpression()
		} else {
			p.expect("';'", p.currentLexme)
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
	if p.match(epllex.SEMICOLON) {
		p.readNextToken() //Skip the semicolon
	} else {
		p.expect("';'", p.currentLexme)
	}

	return ast.Singular{Symbol: ident}
}

//Ident := ID (Basic string Token)
func (p *Parser) ParseIdent() ast.Ident {
	if p.match(epllex.ID) {
		tmp := ast.Ident{Name: p.currentLexme}
		p.readNextToken()
		return tmp
	} else {
		p.expect("Ident", p.currentLexme)
	}

	return ast.Ident{}
}

func (p *Parser) ParseString() ast.String {
	if p.match(epllex.STRINGLITERAL) {
		tmp := ast.String{Value: p.currentLexme}
		p.readNextToken()
		return tmp
	} else {
		p.expect("String", p.currentLexme)
	}

	return ast.String{}
}


//design note: the current token needs to be the first ident of the import path
func (p *Parser) ParseFunctionCall() ast.FunctionCall {

	var params []ast.Expression
	var importPath []ast.Ident
	var name ast.Ident

	// parse function like out.put("asdf")
	for p.match(epllex.COMMA) {
		importPath = append(importPath, p.ParseIdent())
		p.readNextToken()
	}

	/*
	todo: revisit module path system.
		should the function called globally as [path]funcname() or locally as funcname
	*/
	if p.match(epllex.ID) {
		name = p.ParseIdent()
	} else {
		//todo: function name or ident
		p.expect("function name ", p.currentLexme)
	}

	// set the types error handler line descriptor
	p.TypeHandler.SetLine(p.Lexer.GetLine())


	if p.match(epllex.LPAR) {
		p.readNextToken()
		for p.match(epllex.COMMA) && !p.match(epllex.RPAR) {
			paramType := p.TypeHandler.ResolveValueType(p.currentToken)

			if paramType != nil {
				if p.match(epllex.STRINGLITERAL) {
					params = append(params, p.ParseString())
				} else if p.match(epllex.NUM) || p.match(epllex.REAL) {
					params = append(params, p.expr(0))
				}
			} else {
				params = append(params, p.ParseIdent())
			}
		}

		if p.match(epllex.RPAR) {
			p.readNextToken()
		} else {
			p.expect(")", p.currentLexme)
		}
	} else {
		p.expect("(", p.currentLexme)
	}

	return ast.FunctionCall{PackagePath: importPath, Arguments: params, FunctionName:name}
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

	// all expressions must end with semicolon
	if !p.match(epllex.SEMICOLON) {
		p.expect("';'", p.currentLexme)
	}
	p.readNextToken() // skip the semi

	return exprs
}

func (p *Parser) expr(rbp int) ast.Expression {
	var left ast.Expression

	left = p.nud(p.currentToken)

	for p.isBinary(left) && p.currentToken.Precedence() > rbp {
		//Left-Denotation: a simple binary expression
		if p.currentToken.IsLeftAssociative() {
			switch p.currentToken.Ttype {
			case epllex.ADD:
				left = ast.BinaryAdd{
					Ls: left,
					Rs: p.expr(p.currentToken.Precedence()),
				}
			case epllex.SUB:
				left = ast.BinarySub{
					Ls: left,
					Rs: p.expr(p.currentToken.Precedence()),
				}

			case epllex.MULT:
				left = ast.BinaryMul{
					Ls: left,
					Rs: p.expr(p.currentToken.Precedence()),
				}
			case epllex.DIV:
				left = ast.BinaryDiv{
					Ls: left,
					Rs: p.expr(p.currentToken.Precedence()),
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
			p.expect(")", p.currentLexme)
		}

		p.readNextToken() //skip RPAR
	default:
		if t.IsUnary() {
			
			//if the token is an unary operator
			if t.Ttype == epllex.ADD {
				p.readNextToken()
				left = ast.UnaryPlus{Rs: p.expr(epllex.HighPrec)}
			} else if t.Ttype == epllex.SUB {
				p.readNextToken()
				left = ast.UnaryMinus{Rs: p.expr(epllex.HighPrec)}
			}
		} else if t.IsScalar() {
			
			/*
			if the token is some ind of number, real or integer 
			note: complex number support will be added in the future
			*/
			p.readNextToken()
			return ast.Number{Value: t.Lexme}
		} else if t.IsIdent() {
			if p.match_n(epllex.LPAR) || p.match_n(epllex.DOT) {
				
				//lookahead ==  '/' || lookahead == '('
				return p.ParseFunctionCall()
			} else if p.match_n(epllex.SEMICOLON) {
				return ast.Singular{Symbol: p.ParseIdent()}
			} else {
				p.expect("';' or ident", p.currentLexme)
			}
		} else  if t.IsString() {
			/*
			this allows expression like 3 * "asdf"
			the type conflicts will be caught during the type checking phase
			 */
			return p.ParseString()
		}
	}

	return left
}



func (p *Parser) boolExpr(rbp int) ast.BoolExpr {
	var left ast.BoolExpr

	left = p.bool_nud(p.currentToken)

	for p.isBinary(left) && p.currentToken.Precedence() > rbp {
		if p.currentToken.IsLeftAssociative() {
			switch p.currentToken.Ttype {
			case epllex.BOOL_AND:
				left = ast.BoolAnd{
					Le: left,
					Re: p.boolExpr(p.currentToken.Precedence()),
				}

			case epllex.BOOL_OR:
				left = ast.BoolOr{
					Le: left,
					Re: p.boolExpr(p.currentToken.Precedence()),
				}

			//TODO: is == right assoc?
			case epllex.EQ:
				left = ast.BoolEquals{
					Le: left,
					Re: p.boolExpr(p.currentToken.Precedence()),
				}
			case epllex.NEQ:
				left = ast.BoolNotEquals{
					Le: left,
					Re: p.boolExpr(p.currentToken.Precedence()),
				}
			case epllex.GT:
				left = ast.BoolGreaterThen{
					Le: left,
					Re: p.boolExpr(p.currentToken.Precedence()),
				}

			case epllex.LT:
				left = ast.BoolLowerThen{
					Le: left,
					Re: p.boolExpr(p.currentToken.Precedence()),
				}
			case epllex.LE:
				left = ast.BoolLowerThenEqual{
					Le: left,
					Re: p.boolExpr(p.currentToken.Precedence()),
				}
			}
		}
		p.readNextToken()
	}

	return left
}

// boolean expression version of nud
func (p* Parser) bool_nud(t epllex.Token) ast.BoolExpr {
	var left ast.BoolExpr

	switch t.Ttype {
	case epllex.LPAR:
		p.readNextToken()

		//call bool expr to parse the inner expr
		left = p.boolExpr(0)

		// t != )
		if !p.match(epllex.RPAR) {
			p.expect("')'", p.currentLexme)
		}
		p.readNextToken()
	case epllex.BOOL_NOT:
		left = ast.BoolNot{Expr:p.boolExpr(epllex.HighPrec)}
	default:
		if  t.IsIdent() {
			if p.match_n(epllex.LPAR) || p.match_n(epllex.DOT) {
				return p.ParseFunctionCall()
			} else {
				p.expect("'(' or '.'", p.lookahead.Lexme)
			}
		} else {
			p.expect("'(', '!' or ident", p.currentLexme)
		}

	}
	return left
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

//----------------------------------------------------------------------------------------------------------------------
//State Machine control and support functions

func (p *Parser) isBinary(exp ast.Expression) bool {
	switch exp.(type) {
	case ast.BinaryAdd, ast.BinaryDiv, ast.BinaryMul, ast.BinarySub:
		return true
	case ast.BoolGreatEquals, ast.BoolGreaterThen,
	ast.BoolEquals, ast.BoolLowerThen, ast.BoolLowerThenEqual, ast.BoolNotEquals:
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
