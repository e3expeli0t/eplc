/*
*	Copyright (C) 2018-2020 Elia Ariaz
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

package ast

import (
	"eplc/src/libepl/Output"
	"reflect"
)

func walkImports(v Visitor, list *[]Import) {
	for _, i := range *list {
		Walk(v, &i)
	}
}

func walkExprList(v Visitor, exprs *[]Expression) {
	for _, e := range *exprs {
		Walk(v, e)
	}
}
func walkDecls(v Visitor, decls *[]Decl) {
	for _, i := range *decls {
		Walk(v, i)
	}
}

func Walk(v Visitor, node Node) {
	if v = v.Visit(node); v == nil {
		return
	}

	switch n := node.(type) {
	//Statements
	case nil:
	case *ProgramFile:
		walkImports(v, n.Imports)
		walkDecls(v, n.GlobalDecls)

		for _, i := range *n.Functions {
			Walk(v, i)
		}
		Walk(v, n.MainFunction)
	case *Import:
		//do nothing
	case *Fnc:
		Walk(v, n.Name)
		walkDecls(v, n.Params)
		Walk(v, n.Body)
	case *Block:
		walkExprList(v, n.ExprList)
	case *ForLoop:
		Walk(v, *n.VarDef)
		Walk(v, *n.Condition)
		Walk(v, *n.Expr)
		Walk(v, n.Code)
	case *VarDecl:
		Walk(v, n.Name)
	case *VarExplicitDecl:
		Walk(v, n.Name)
		Walk(v, *n.Value)
	case *AssignStmt:
		Walk(v, n.Owner)
		Walk(v, *n.Value)
	case *IfStmt:
		Walk(v, *n.Condition)
		Walk(v, n.Code)
		if n.Else != nil {
			Walk(v, *n.Else)
		}
	case *ElseStmt:
		Output.PrintFatalErr("ElseStmt node is deprecated")
	case *Repeat:
		{
			if n.VarDef != nil {
				Walk(v, *n.VarDef)
			}
			Walk(v, n.Code)
		}
	case *RepeatUntil:
		{
			if n.VarDef != nil {
				Walk(v, *n.VarDef)
			}
			Walk(v, n.Code)
			Walk(v, *n.Condition)
		}
	case *Until:
		Walk(v, *n.Condition)
		Walk(v, n.Code)
	case *BoolGreatEquals:
		Walk(v, n.Ls)
		Walk(v, n.Rs)
	case *BoolGreaterThen:
		Walk(v, n.Ls)
		Walk(v, n.Rs)
	case *BoolLowerThen:
		Walk(v, n.Ls)
		Walk(v, n.Rs)
	case *BoolLowerThenEqual:
		Walk(v, n.Ls)
		Walk(v, n.Rs)
	case *BoolEquals:
		Walk(v, n.Ls)
		Walk(v, n.Rs)
	case *BoolNotEquals:
		Walk(v, n.Ls)
		Walk(v, n.Rs)
	case *BoolAnd:
		Walk(v, n.Ls)
		Walk(v, n.Rs)
	case *BoolOr:
		Walk(v, n.Ls)
		Walk(v, n.Rs)
	case *BoolNot:
		Walk(v, n.Expr)
	case *BinarySub:
		Walk(v, n.Ls)
		Walk(v, n.Rs)
	case *BinaryAdd:
		Walk(v, n.Ls)
		Walk(v, n.Rs)
	case *BinaryDiv:
		Walk(v, n.Ls)
		Walk(v, n.Rs)
	case *BinaryMul:
		Walk(v, n.Ls)
		Walk(v, n.Rs)
	case *UnaryMinus:
		Walk(v, n.Rs)
	case *UnaryPlus:
		Walk(v, n.Rs)
	case *FunctionCall:
		for _, i := range n.PackagePath {
			Walk(v, i)
		}
		Walk(v, n.FunctionName)
		walkExprList(v, &n.Arguments)

	case *Return:
		Walk(v, *n.Value)
	case *Ident:
		//do nothing
	case Number:
		//do nothing
	case String:
		//do nothing
	case Boolean:
		//do nothing
	case *Break:
		//do nothing
	case Singular:
		Output.PrintFatalErr("Singular node is deprecated")
	case EmptyExpr:
		//do nothing
	default:
		Output.PrintFatalErr("Unexpected node", reflect.TypeOf(n))
	}
	v.Visit(nil)
}
