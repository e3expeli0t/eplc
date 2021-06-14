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

package analysis

import (
	"eplc/src/libepl/Types"
	"eplc/src/libepl/eplparse/ast"
	"eplc/src/libepl/eplparse/symboltable"
	"eplc/src/libio"

	"fmt"
	"reflect"
)

/*
	Epl's type system:
		Every EplType extends the type ObjectValue.
		Meaning that int and string can be both represented as ObjectValue

		The type hierarchy looks like that:
				 ObjectValue
					  |
                ______|_______
				|			  |
			BasicType		EplType

		BasicType is form of EplType that represents basic language types
		like int and float
*/

type TypeErrorCase uint

const (
	TypeMismatch    TypeErrorCase = iota
	InvalidUseOfNot TypeErrorCase = iota
	InvalidUseOfUnary TypeErrorCase = iota
)

func NewTypeChecker(table *symboltable.TableMap) *TypeChecker {
	if table.Last == 0 {
		panic("Couldn't load symbol table")
	}
	return &TypeChecker{
		SymbolMap: table,
	}
}

type TypeChecker struct {
	SymbolMap *symboltable.TableMap
	Current  symboltable.ScopeSymbolTable
	Errors    []*TypeError
}

func (tc *TypeChecker) addError(c TypeErrorCase, t1 Types.EplType, t2 Types.EplType) {
	var descriptor string

	switch c {
	case TypeMismatch:
		descriptor = fmt.Sprintf("Type mismatch between types %s and %s",
			t1.TypeName,
			t2.TypeName,
		)
	case InvalidUseOfNot:
		descriptor = fmt.Sprintf(
			"Type %s cant be used in boolean not expression, type should be bool",
			t1.TypeName,
			)
	case InvalidUseOfUnary:
		descriptor = fmt.Sprintf(
			"Type %s can't be used as number in unary expression, should be int or uint",
			t1.TypeName,
			)
	default:
		descriptor = fmt.Sprintf("Type error between %s and %s",
			t1.TypeName,
			t2.TypeName,
		)
	}

	tc.Errors = append(tc.Errors, NewError(Major, descriptor))
}

func(tc *TypeChecker) HasErrors() bool {
	return len(tc.Errors) > 0
}

func (tc *TypeChecker) WalkExpression(expr ast.Expression) Types.EplType {
	//note: since go doesn't allow fallthrough in type switches
	// we have this amazing chunk of redundant code
	switch n := expr.(type) {
	case *ast.BoolGreatEquals:
		status, t := tc.HandleBinaryGeneral(&n.Rs, &n.Ls)

		if status {
			return t
		}
		break
	case *ast.BoolGreaterThen:
		status, t := tc.HandleBinaryGeneral(&n.Rs, &n.Ls)

		if status {
			return t
		}
		break
	case *ast.BoolLowerThen:
		status, t := tc.HandleBinaryGeneral(&n.Rs, &n.Ls)

		if status {
			return t
		}
		break
	case *ast.BoolLowerThenEqual:
		status, t := tc.HandleBinaryGeneral(&n.Rs, &n.Ls)

		if status {
			return t
		}
		break
	case *ast.BoolEquals:
		status, t := tc.HandleBinaryGeneral(&n.Rs, &n.Ls)

		if status {
			return t
		}
		break
	case *ast.BoolNotEquals:
		status, t := tc.HandleBinaryGeneral(&n.Rs, &n.Ls)

		if status {
			return t
		}
		break
	case *ast.BoolAnd:
		status, t := tc.HandleBinaryGeneral(&n.Rs, &n.Ls)

		if status {
			return t
		}
		break
	case *ast.BoolOr:
		status, t := tc.HandleBinaryGeneral(&n.Rs, &n.Ls)

		if status {
			return t
		}
		break
	case *ast.BoolNot:
		ExprType := tc.WalkExpression(n.Expr)

		if ExprType.Equals(Types.TypeBool) {
			//todo: error recovery
			tc.addError(InvalidUseOfNot, ExprType, Types.EplType{})
		}
		return ExprType

	case *ast.BinarySub:
		status, t := tc.HandleBinaryMathematical(&n.Rs, &n.Ls)

		if status {
			return t
		}
		break
	case *ast.BinaryAdd:
		status, t := tc.HandleBinaryMathematical(&n.Rs, &n.Ls)

		if status {
			return t
		}
		break
	case *ast.BinaryDiv:
		status, t := tc.HandleBinaryMathematical(&n.Rs, &n.Ls)

		if status {
			return t
		}
		break
	case *ast.BinaryMul:
		status, t := tc.HandleBinaryMathematical(&n.Rs, &n.Ls)

		if status {
			return t
		}
		break

	case *ast.UnaryMinus:
		status, t := tc.HandleUnaryMathematical(&n.Rs)

		if status {
			return t
		}
		break
	case *ast.UnaryPlus:
		status, t := tc.HandleUnaryMathematical(&n.Rs)

		if status {
			return t
		}
		break
	case *ast.FunctionCall:

		libio.PrintLog("Found call to: "+n.FunctionName.Name)

		if len(n.PackagePath) == 0 {
			return tc.SymbolMap.Locate(n.FunctionName.Name).Type
		}

		return tc.SymbolMap.Locate(n.ConstructFullPath()).Type

	case *ast.Ident:

	case ast.Number:
		//Add system cpu bits resolver
		return Types.TypeInt.AsEplType()
	case ast.String:
		return Types.TypeString.AsEplType()
	case ast.Boolean:
		return Types.TypeBool.AsEplType()
	default:
		panic(fmt.Sprintf("Unexpected node type: %s", reflect.TypeOf(n)))
	}

	return Types.EplType{}
}

func (tc* TypeChecker) HandleBinaryGeneral(left *ast.Expression, right *ast.Expression) (bool, Types.EplType) {
	ExprTypeLeft := tc.WalkExpression(*left)
	ExprTypeRight := tc.WalkExpression(*right)

	if ExprTypeRight == ExprTypeLeft {
		return true, ExprTypeLeft
	}
	tc.addError(TypeMismatch, ExprTypeLeft, ExprTypeRight)
	return false, ExprTypeRight
}


func (tc* TypeChecker) HandleBinaryMathematical(left *ast.Expression, right *ast.Expression) (bool, Types.EplType) {
	ExprTypeLeft := tc.WalkExpression(*left)
	ExprTypeRight := tc.WalkExpression(*right)


	if ExprTypeLeft.IsMathematical() && ExprTypeRight.IsMathematical() {
		return true, ExprTypeLeft
	}
	tc.addError(TypeMismatch, ExprTypeLeft, ExprTypeRight)
	return false, ExprTypeRight
}

func (tc* TypeChecker) HandleUnaryMathematical(left *ast.Expression) (bool, Types.EplType) {
	ExprType := tc.WalkExpression(*left)

	if ExprType.IsMathematical(){
		return true, ExprType
	}
	tc.addError(InvalidUseOfUnary, ExprType, Types.EplType{})
	return false, ExprType
}
