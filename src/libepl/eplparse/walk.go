package eplparse

import "eplc/src/libepl/eplparse/ast"

type Visitor interface {
	Visit(node ast.Node) (v Visitor)
}

func walkImports(v Visitor, list *[]ast.Import) {
	for _, i := range *list {
		Walk(v, &i)
	}
}

func walkExprList(v Visitor, exprs *[]ast.Expression) {
	for _, e := range *exprs {
		Walk(v, e)
	}
}
func walkDecls(v Visitor, decls *[]ast.Decl) {
	for _, i := range *decls {
		Walk(v, i)
	}
}

func Walk(v Visitor, node ast.Node) {
	if v = v.Visit(node); v == nil {
		return
	}

	switch n := node.(type) {
	//Statements
	case *ast.ProgramFile:
		walkImports(v, n.Imports)
		walkDecls(v, n.GlobalDecls)

		for _, i := range *n.Functions {
			Walk(v, &i)
		}
		Walk(v, n.MainFunction)

	case *ast.Fnc:
		Walk(v, n.Body)
	case *ast.Block:
		walkExprList(v, n.ExprList)
	case *ast.IfStmt:
		Walk(v, *n.Condition)
		Walk(v, n.Code)
		Walk(v, *n.Else)
	case *ast.AssignStmt:
		Walk(v, n.Owner)
		Walk(v, *n.Value)
	//Expressions
	case *ast.BinaryAdd:
		Walk(v, n.Ls)
		Walk(v, n.Rs)
	case *ast.BinarySub:
		Walk(v, n.Ls)
		Walk(v, n.Rs)
	case *ast.BinaryMul:
		Walk(v, n.Ls)
		Walk(v, n.Rs)
	case *ast.BinaryDiv:
		Walk(v, n.Ls)
		Walk(v, n.Rs)
	case *ast.UnaryMinus:
		Walk(v, n.Rs)
	case *ast.UnaryPlus:
		Walk(v, n.Rs)
	case *ast.Ident:
		//do nothing
	}
}
