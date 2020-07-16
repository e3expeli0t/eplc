package eplparse

import (
	"eplc/src/libepl/Types"
	"eplc/src/libepl/epllex"
	"eplc/src/libepl/eplparse/ast"
	"fmt"
)

//Helper maps. Used To efficiently identify expression during
//block parsing
var ExpressionsMAP = []epllex.TokenType{epllex.LPAR, epllex.ID}

//This method enforces the expression to end with semicolon
func (p *Parser) ParseExpression() *ast.Expression {
	expr := p.ParseExpr(0)

	if expr == nil {
		p.expect("expression")
	}

	if !p.match(epllex.SEMICOLON) {
		p.expect("';'")
	}
	p.readNextToken() // skip the semicolon

	return &expr
}

//Expression := unary| Expression binary_op Expression
//binary_op := "+" | "-" | "*" | "/" | "<=" | ">=" | "==" | "<" | ">"
func (p *Parser) ParseExpr(prec int) ast.Expression {
	var left ast.Expression
	var isNil bool
	left = p.ParseUnary()

	for (left != nil) && (p.currentToken.IsBinary() && p.currentToken.Precedence() > prec) {
		left, isNil = p.resolveBinaryOperator(&left)
		if isNil {
			break
		}
	}

	return left
}

//Called on the first element of the expression
//expression parsing stops when the value of unary is nil
//unary := op value | value
//op := "+" | "-"
//value := ident | bool_val | string_literal | digit | float | function_call
func (p *Parser) ParseUnary() (unary ast.Expression){

	if p.match(epllex.LPAR) {
		p.readNextToken()

		unary = p.ParseExpr(0)

		if !p.match(epllex.RPAR) {
			p.expect("(")
		}
		p.readNextToken()
		return
	}

	if p.currentToken.IsUnary() {
		unary = p.resolveOperator()
		return unary
	}

	if p.match(epllex.ID) {
		if p.match_n(epllex.DOT) || p.match_n(epllex.LPAR) {
			return p.ParseFunctionCall()
		}
		return p.ParseIdent()
	}

	return p.parseValue()
}


//Can return nil: this is a sign that the symbol is not expression symbol
func (p *Parser) resolveOperator() (exp ast.Expression) {
	switch p.currentToken.Ttype {
	case epllex.ADD:
		exp = &ast.UnaryPlus{Rs: p.ParseExpr(epllex.HighPrec)}
	case epllex.SUB:
		exp = &ast.UnaryMinus{Rs: p.ParseExpr(epllex.HighPrec)}
	default:
		return nil
	}

	return exp
}


//todo: return pointer
//if the function returns nill this is a sign that the expression has ended
func (p *Parser) resolveBinaryOperator(left* ast.Expression) (exp ast.Expression, nil bool) {
	nil = true

	switch p.currentToken.Ttype {
	case epllex.MULT:
		p.readNextToken() // Skip the operator
		exp = &ast.BinaryMul{
			Ls: *left,
			Rs: p.ParseExpr(p.currentToken.Precedence()),
		}
		nil = false
		break

	case epllex.DIV:
		p.readNextToken() // Skip the operator
		exp = &ast.BinaryDiv{
			Ls: *left,
			Rs: p.ParseExpr(p.currentToken.Precedence()),
		}
		nil = false
		break

	case epllex.ADD:
		p.readNextToken() // Skip the operator
		exp = &ast.BinaryAdd{
			Ls: *left,
			Rs: p.ParseExpr(p.currentToken.Precedence()),
		}
		nil = false
		break

	case epllex.SUB:
		p.readNextToken() // Skip the operator
		exp = &ast.BinarySub{
			Ls: *left,
			Rs: p.ParseExpr(p.currentToken.Precedence()),
		}
		nil = false
		break

	case epllex.EQ:
		p.readNextToken() // Skip the operator
		exp = &ast.BoolEquals {
			Ls: *left,
			Rs: p.ParseExpr(p.currentToken.Precedence()),
		}
		nil = false
		break

	case epllex.NEQ:
		p.readNextToken() // Skip the operator
		exp = ast.BoolNotEquals {
			Ls: *left,
			Rs: p.ParseExpr(p.currentToken.Precedence()),
		}
		nil = false
		break
	case epllex.LT:
		p.readNextToken() // Skip the operator
		exp = &ast.BoolLowerThen {
			Ls: *left,
			Rs: p.ParseExpr(p.currentToken.Precedence()),
		}
		nil = false
		break

	case epllex.GT:
		p.readNextToken() // Skip the operator
		exp = &ast.BoolGreatEquals {
			Ls: *left,
			Rs: p.ParseExpr(p.currentToken.Precedence()),
		}
		nil = false
		break

	case epllex.LE:
		p.readNextToken() // Skip the operator
		exp = &ast.BoolLowerThenEqual {
			Ls: *left,
			Rs: p.ParseExpr(p.currentToken.Precedence()),
		}
		nil = false
		break

	case epllex.GE:
		exp = &ast.BoolGreatEquals {
			Ls: *left,
			Rs: p.ParseExpr(p.currentToken.Precedence()),
		}
		nil = false
		break

	default:
		// Should we set exp to BadExpr?
		break
	}

	return
}

//todo: create function that translate the token to Value ast node
func (p *Parser) parseValue() ast.Expression {
	defer p.readNextToken()

	switch p.currentToken.Ttype {
	case epllex.TRUE, epllex.FALSE:
		return p.ToBoolVal(p.currentToken)
	case epllex.STRINGLITERAL:
		return ast.String{Value: p.currentLexme}
	case epllex.REAL:
		return ast.Number{
			Value: p.currentLexme,
			Real:  true,
		}
	case epllex.NUM:
		return ast.Number{
			Value: p.currentLexme,
			Real:  false,
		}
	case epllex.CMX:
	default:
		p.report(fmt.Sprintf("Can't resolve type of value:%s.", p.currentLexme))
	}

	//Not reached
	panic("Shouldn't be reached")
}

//function_call := package_path "(" args ")"
//package_path := ident "." path_list | ident
func (p *Parser) ParseFunctionCall() *ast.FunctionCall {
	var packagePath []*ast.Ident
	var name ast.Ident
	var args []ast.Expression

	var current *ast.Ident

	//parsing the function path
	for !p.match(epllex.LPAR) && p.match(epllex.ID) {
		if p.match(epllex.ID) {
			current = p.ParseIdent()
			packagePath = append(packagePath, current)
		} else {
			p.expect("ident")
		}
		if p.match(epllex.DOT) {
			p.readNextToken() //assume that the current token is DOT
		} else {
			continue // probably found (
		}
	}

	if packagePath == nil || current == nil {
		p.expect("function name")
	}

	// current != nil
	name = *current

	//Todo: Handle group variable access(i.e std.eof) v0.2+
	if !p.match(epllex.LPAR) {
		p.expect("'('")
	}
	p.readNextToken()

	if !p.match(epllex.RPAR) {
		args = append(args, p.ParseExpr(0))
	}

	for !p.match(epllex.RPAR) && p.match(epllex.COMMA) {
		p.readNextToken()
		args = append(args, p.ParseExpr(0))
	}

	if !p.match(epllex.RPAR) {
		p.expect("')'")
	}
	p.readNextToken() // skip ')'

	return &ast.FunctionCall{
		PackagePath:  packagePath,
		Arguments:    args,
		ReturnType:   Types.EplType{}, //gets filled by the type analyser
		FunctionName: &name,
	}
}

//Ident := ID (Basic string Token)
func (p *Parser) ParseIdent() *ast.Ident {
	if p.match(epllex.ID) {
		defer p.readNextToken()
		return &ast.Ident{Name: p.currentLexme}
	} else {
		p.expect("Ident")
	}
	//Not reached
	panic("Shouldn't be reached")
}