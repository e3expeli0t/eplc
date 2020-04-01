package eplparse

import (
	"eplc/src/libepl/Output"
	"eplc/src/libepl/epllex"
	"eplc/src/libepl/eplparse/ast"
)


//----------------------------------------------------------------------------------------------------------------------
//Expressions



/*
Design note:
	Possible starts for math expressions are:
		ID
		(
		-
		+


	Types of expressions:
	there are two major types of expressions. mathematical and boolean
	Mathematical expressions are in the form:
		math_expr := expr math_op expr | ident |number | functioncall | stmt
		note: stmts are void expressions any use of them in expression will be caught during type checking
	Boolean expressions are in the form:
		bool_expr :=
	the Expressions AST hierarchy is like that:
           Node
			|
        Expression --> BinaryDiv, BinaryMul, BinaryPlus, BinaryMinus, UnaryPlus, UnaryMinus, FunctionCall,String, Number
            |
		BoolExpr -->  BoolOr, BoolAnd, Equal, NotEqual, GreaterThen, LowerThen, GreaterEqual, LowerEqual, BoolVal
*/

// maps for easy detection of expr starts
var UnaryStart = []epllex.TokenType{epllex.SUB, epllex.ADD}
//todo: BinaryMap Should also contain UnaryStart?
var BinaryStart = []epllex.TokenType{epllex.ID, epllex.LPAR, epllex.NUM}
var BinaryEnd = []epllex.TokenType{epllex.RPAR, epllex.ID, epllex.NUM}


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

	for p.currentToken.Precedence() > rbp {
		//Left-Denotation: a simple binary expression
		if p.currentToken.IsBinary() && p.currentToken.IsLeftAssociative() {
			switch p.currentToken.Ttype {
			case epllex.ADD:
				p.readNextToken() // skip the op
				left = ast.BinaryAdd{
					Ls: left,
					Rs: p.expr(p.currentToken.Precedence()),
				}
			case epllex.SUB:
				p.readNextToken() // skip the op
				left = ast.BinarySub{
					Ls: left,
					Rs: p.expr(p.currentToken.Precedence()),
				}

			case epllex.MULT:
				p.readNextToken() // skip the op
				left = ast.BinaryMul{
					Ls: left,
					Rs: p.expr(p.currentToken.Precedence()),
				}
			case epllex.DIV:
				p.readNextToken() // skip the op
				left = ast.BinaryDiv{
					Ls: left,
					Rs: p.expr(p.currentToken.Precedence()),
				}
			default:
				p.expect("'+', '-', '/', '*'", p.currentLexme)
			}
		} else {
			break
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
				if the token is some kind of number, real or integer
				note: complex number support will be added in the future
			*/
			p.readNextToken()
			return ast.Number{Value: t.Lexme}
		} else if t.IsIdent() {
			if p.match_n(epllex.LPAR) || p.match_n(epllex.DOT)  {
				//lookahead ==  '.' || lookahead == '('
				return p.ParseFunctionCall()
			} else if p.match_n(epllex.SEMICOLON) {
				return p.ParseSingularExpr()
			} else if p.lookahead.IsLeftAssociative() {
				return p.ParseIdent()
			} else {
				p.expect("';', ident, singular or boolean expression (i.e '>=', '<=', '==' etc.)", p.currentLexme)
			}
		} else  if t.IsString() {
			/*
				this allows expression like 3 * "asdf"
				the type conflicts will be caught during the type checking phase
			*/
			return p.ParseString()
		}  else {
			//todo: change to bad expr?
			p.report("Invalid expression")
		}
	}

	return left
}

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

	if p.match(epllex.ID) && p.match_n(epllex.DOT){
		importPath = append(importPath, p.ParseIdent())
		p.readNextToken()
	}
	// parse function like out.put("asdf")
	for p.match(epllex.DOT) {
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

	if p.match(epllex.LPAR) {
		p.readNextToken() // skip lpar

		for !p.match(epllex.RPAR) {
			paramType := p.TypeHandler.ResolveValueType(p.currentToken)

			if paramType != nil {
				if p.match(epllex.STRINGLITERAL) {
					params = append(params, p.ParseString())
				} else if p.matchTokens(BinaryStart) {
					params = append(params, p.expr(0))
				}
			} else {
				params = append(params, p.ParseIdent())
			}

			if !p.match(epllex.COMMA) && !p.match(epllex.RPAR){
				p.expect("',' or ')'", p.currentLexme)
			} else if p.match(epllex.COMMA){
				p.readNextToken() // skip the comma
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
design note:
	Boolean expression parsing.
	BoolExpr := bool_val | bool_expr bool_op bool_expr| "("bool_expr")" | "!"bool_expr| bool_expr
	Notes:
	*	we parse the expressions using precedence climbing
	*	credit to https://github.com/richardjennings/prattparser/
	*	credit to https://www.engr.mun.ca/~theo/Misc/exp_parsing.htm
	*	credit to https://github.com/golang/go

*/

//todo: Rewrite boolean operator parser. optimize code and clean it
var BoolStartMap = []epllex.TokenType{epllex.ID, epllex.TRUE, epllex.FALSE, epllex.LPAR, epllex.BOOL_NOT}
var BoolEndMap = []epllex.TokenType{epllex.RPAR, epllex.ID, epllex.TRUE, epllex.FALSE}
var SizeOpMap = []epllex.TokenType{epllex.GT, epllex.LT, epllex.LE, epllex.GE}


/*
	Types of boolean expressions:
		SizeOp:
			every operator that RE or LE is mathematical expression
			The operators are:
				>, <, ==, !=, >=, <=
 */

/*
boolExpr(rbp):
	the main function for parsing bool_expr warps the bool_nud and SizeOp parsing
 */
func (p *Parser) boolExpr(rbp int) ast.BoolExpr {
	var left ast.BoolExpr

	left = p.bool_nud(p.currentToken)

	for p.isBinary(left) && p.currentToken.Precedence() > rbp {
		if p.currentToken.IsLeftAssociative() {
			Output.PrintVersion()
			switch p.currentToken.Ttype {
			case epllex.BOOL_AND:
				p.readNextToken()
				left = ast.BoolAnd {
					Le: left,
					Re: p.boolExpr(p.currentToken.Precedence()),
				}

			case epllex.BOOL_OR:
				p.readNextToken()
				left = ast.BoolOr{
					Le: left,
					Re: p.boolExpr(p.currentToken.Precedence()),
				}

				/* TODO: is == right assoc?
				-handles expr == expr types ?
				-boolExpr is an expression? if it does exprs like 12+9987 == 6554
				or exprs like `equal() == true` or even `eq() = !eq()`
				todo: [urg] redesign AST and expression hierarchy
				todo: [urg] optimize expressions parser {apx time: after parser UNIT tests}
				*/
			case epllex.EQ:
				p.readNextToken()

				var re ast.Expression

				if p.matchTokens(BinaryStart) || p.matchTokens(UnaryStart) {
					re = p.expr(0)
				} else if p.matchTokens(BoolStartMap) {
					re = p.boolExpr(0)
				} else {
					p.expect("expression", p.currentLexme)
				}

				left = ast.BoolEquals{
					Le: left,
					Re: re,
				}

			case epllex.NEQ:
				p.readNextToken()
				left = ast.BoolNotEquals{
					Le: left,
					Re: p.boolExpr(p.currentToken.Precedence()),
				}

			}
		}
	}

	return left
}


//note: inefficient function. todo: optimize
// boolean expression version of nud
//what happens when the right of == is bool value
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
			} else if p.matchNTokens(SizeOpMap) {
				p.ParseSizeOP(0)
			}else {
				return p.ParseIdent()
			}
		} else if t.IsBoolVal() {
			return p.ToBoolVal(t)
		} else if p.matchTokens(BinaryStart) {
			return p.ParseSizeOP(0)
		} else {
			p.expect("'(' , '!' or ident", p.currentLexme)
		}

	}
	return left
}


/*
	ParseSizeOP: Handles conditions that include sizing i.e num1 >= num2
	Error handling: if the function didnt locate the correct operator it will return the left expr
	todo: return pointer/ v2.2+
	todo: fix incomplete and bad error handling
	todo: [U-MB] None deterministic design of
*/
//pllod(2+rdp()) >= 23
func (p *Parser) ParseSizeOP(rbp int) ast.BoolExpr {
	var left ast.Expression

	if p.matchTokens(BinaryStart) || p.matchTokens(UnaryStart) {
		left = p.expr(0)
	} else {
		p.expect("expression", p.currentLexme)
	}

	for p.currentToken.Precedence() > rbp {
		if p.currentToken.IsLeftAssociative() {
			switch p.currentToken.Ttype {
			case epllex.GT:
				left = ast.BoolGreaterThen{
					Le: left,
					Re: p.expr(0),
				}
			case epllex.LT:
				left = ast.BoolLowerThen{
					Le: left,
					Re: p.expr(0),
				}
			case epllex.LE:
				left = ast.BoolLowerThenEqual {
					Le: left,
					Re: p.expr(0),
				}
			case epllex.GE:
				left = ast.BoolGreatEquals{
					Le: left,
					Re: p.expr(0),
				}
			}
		}
	}
	return ast.EmptyExpr{}
}
