package analysis

import (
	"eplc/src/libepl/Types"
	"eplc/src/libepl/eplparse/ast"
	"fmt"
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

type TypeChecker struct {
	Errors []*TypeError
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
		descriptor = fmt.Sprintf("Type %s cant be used in boolean not expression, type should be bool", t1.TypeName)
	case InvalidUseOfUnary:
		descriptor = fmt.Sprintf("Type %s can't be used as number in unary expression, should be int or uint", t1.TypeName)
	default:
		descriptor = fmt.Sprintf("Type error between %s and %s",
			t1.TypeName,
			t2.TypeName,
		)
	}

	tc.Errors = append(tc.Errors, NewError(Major, descriptor))
}

func (tc *TypeChecker) Check(expr *ast.Expression) {

}

func (tc *TypeChecker) WalkExpression(expr ast.Expression) Types.EplType {
	//todo: Divide the boolean nodes into classes
	//todo: Add explicit type check ( for example in BoolGreatEquals Ls and Rs type
	// should be some kind of uint/int)

	switch n := expr.(type) {
	case ast.BoolGreatEquals:
	case ast.BoolGreaterThen:
	case ast.BoolLowerThen:
	case ast.BoolLowerThenEqual:
	case ast.BoolEquals:
	case ast.BoolNotEquals:
	case ast.BoolAnd:
	case ast.BoolOr:
		ExprTypeLeft := tc.WalkExpression(n.Ls)
		ExprTypeRight := tc.WalkExpression(n.Rs)

		if ExprTypeLeft == ExprTypeRight {
			return ExprTypeLeft
		} else {
			//todo: add recovery
			tc.addError(TypeMismatch, ExprTypeLeft, ExprTypeRight)
		}

	case ast.BoolNot:
		ExprType := tc.WalkExpression(n.Expr)

		if ExprType != Types.TypeBool {
			//todo: error recovery
			tc.addError(InvalidUseOfNot, ExprType, Types.EplType{})
		}
		return ExprType

	case ast.BinarySub:
	case ast.BinaryAdd:
	case ast.BinaryDiv:
	case ast.BinaryMul:
		ExprTypeLeft := tc.WalkExpression(n.Ls)
		ExprTypeRight := tc.WalkExpression(n.Rs)

		if ExprTypeLeft == ExprTypeRight {
			return ExprTypeLeft
		} else {
			//todo: add recovery
			tc.addError(TypeMismatch, ExprTypeLeft, ExprTypeRight)
		}

	case ast.UnaryMinus:
	case ast.UnaryPlus:
		ExprType := tc.WalkExpression(n.Rs)

		if ExprType != Types.TypeBool {
			//todo: error recovery
			tc.addError(InvalidUseOfUnary, ExprType, Types.EplType{})
		}
		return ExprType

	case ast.Ident:
		//do nothing
	case ast.Number:
		return Types.TypeInt
	case ast.String:
		return Types.TypeString
	case ast.Boolean:
		return  Types.TypeBool
	}

	return Types.EplType{}
}
