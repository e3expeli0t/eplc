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
	"eplc/src/libepl/Types"
	"eplc/src/libepl/eplparse/symboltable"
)

/*
	AST (or Abstract Syntax Tree) is the output of  eplparse
*/

type VarStat string
const (
	Fixed   VarStat = "fixed"
	Mutable VarStat = "mutable"
)

type (
	Node interface {
		Start() uint
	}

	Expression interface {
		Node
		ExprNode()
	}

	//Design note: experimental: every statement is an expression that return void
	Statement interface {
		Expression
		StmtNode()
	}

	Decl interface {
		Statement
		DeclNode()
	}
)


//Todo: replace all string names with symbol table references
type (
	ProgramFile struct {
		FileName     string
		Symbols      *symboltable.SymbolTable
		Imports      *[]Import
		GlobalDecls  *[]Decl
		Functions    *[]*Fnc
		MainFunction *Fnc
	}

	Fnc struct {
		Name       *Ident
		ReturnType *Types.EplType
		Params     *[]Decl
		Body       *Block
	}

	Block struct {
		Symbols  *symboltable.SymbolTable
		ExprList *[]Expression
	}

	VarDecl struct {
		Name    *Ident
		VarType *Types.EplType
		Stat    VarStat
	}

	//Declaration and assignment
	VarExplicitDecl struct {
		VarDecl
		Value *Expression
	}

	Import struct {
		StartLoc uint
		Imports  *[]string
	}

	AssignStmt struct {
		Owner *Ident
		Value *Expression
	}

	IfStmt struct {
		Condition *Expression
		Code      *Block
		Else      *Statement // to support if-else-if
	}

	ElseStmt struct {
		Code *Block
	}

	//todo: IMPL in v0.2++
	MoveLoop struct {}

	ForLoop struct {
		VarDef    *Decl
		Condition *Expression
		Expr  	  *Expression // Assign or expr
		Code      *Block
	}

/*
	Example:
		repeat (i int = 0) {
		//code block
		}
 */

	Repeat struct {
		VarDef  *Decl
		Code   *Block
	}

/*
	Example:
		until (1 == 2) {
			//empty
		}
 */

	Until struct {
		Condition *Expression
		Code      *Block
	}


/*
	Example:
		repeat (i int = 0) {
			out.put(i);
			i += 1;
		} until (i == 3);
 */

	RepeatUntil struct {
		VarDef    *Decl
		Condition *Expression
		Code      *Block
	}
	
	Return struct {
		Value *Expression
	}
	Break struct {}
)

func (Break) Start() uint {
	panic("Invalid call")
}
func (Break) ExprNode() {}
func (Break) StmtNode() {}

func (Return) Start() uint {
	panic("Invalid call")
}
func (Return) ExprNode() {}
func (Return) StmtNode() {}

func (Until) Start() uint {
	panic("Invalid call")
}
func (Until) ExprNode() {}

func (RepeatUntil) Start() uint {
	panic("Invalid Call")
}
func (RepeatUntil) ExprNode() {}
func (RepeatUntil) StmtNode() {}


func (Repeat) Start() uint {
	panic("Invalid call")
}
func (Repeat) ExprNode() {}
func (Repeat) StmtNode() {}


func (ForLoop) Start() uint {
	panic("Invalid call")
}
func (ForLoop) ExprNode() {}
func (ForLoop) StmtNode() {}

func (Block) ExprNode() {}

func (AssignStmt) Start() uint {
	panic("Invalid call")
}
func (AssignStmt) ExprNode() {}

func (ElseStmt) Start() uint {
	panic("Invalid call")
}
func (ElseStmt) ExprNode() {}
func (ElseStmt) StmtNode() {}

func (Fnc) Start() uint {
	panic("Invalid call")
}
func (Fnc) DeclNode() {}
func (Fnc) ExprNode() {}
func (Fnc) StmtNode() {}

func (VarDecl) Start() uint {
	panic("Invalid call")
}
func (VarDecl) StmtNode() {}
func (VarDecl) ExprNode() {}
func (VarDecl) DeclNode(){}
func (VarExplicitDecl) DeclNode() {}

func (ProgramFile) Start() uint {
	panic("Invalid call")
}

func (Block) Start() uint {
	panic("Invalid call")
}

func (IfStmt) Start() uint {
	panic("Invalid call")
}
func (IfStmt) StmtNode(){}
func (IfStmt) ExprNode() {}

func (Import) Start() uint {
	panic("Invalid call")
}
//----------------------------------------------------------------------------------------------------------------------
//Expressions

//Boolean values
type BoolValue uint

const (
	BOOL_FALSE BoolValue = 0
	BOOL_TRUE BoolValue  = 1
)

//Nodes definition
type (
	Ident struct {
		Name string
	}

	EmptyExpr struct{}

	//Note: Expressions must be passed by reference
	BoolNot struct {
		Expr Expression
	}

	BoolAnd struct {
		Ls Expression
		Rs Expression
	}

	BoolOr struct {
		Ls Expression
		Rs Expression
	}

	BoolEquals struct {
		Ls Expression
		Rs Expression
	}

	BoolNotEquals struct {
		Ls Expression
		Rs Expression
	}
	
	BoolGreaterThen struct {
		Ls Expression
		Rs Expression
	}

	BoolGreatEquals struct {
		Ls Expression
		Rs Expression
	}

	BoolLowerThen struct {
		Ls Expression
		Rs Expression
	}

	BoolLowerThenEqual struct {
		Ls Expression
		Rs Expression
	}

	Boolean struct {
		Val BoolValue
	}

	BinaryMul struct {
		Ls Expression
		Rs Expression
	}

	BinaryDiv struct {
		Ls Expression
		Rs Expression
	}

	BinaryAdd struct {
		Ls Expression
		Rs Expression
	}

	BinarySub struct {
		Ls Expression
		Rs Expression
	}

	UnaryPlus struct {
		Rs Expression
	}

	UnaryMinus struct {
		Rs Expression
	}

	FunctionCall struct {
		PackagePath []*Ident

		/*
		argument can be any thing that returns value.
		None value is not allowed. And will be caught during type checking
		*/
		Arguments    []Expression
		ReturnType   Types.EplType //the return type is set during type analysis
		FunctionName *Ident
	}

	Singular struct {
		Symbol Ident
	}

	//Represent number such as 5562 or 233.9973
	Number struct {
		Value string
		Real bool // If the number is real (i.e 3.2552) this is true
	}

	String struct {
		Value string
	}
)

func (Boolean) Start() uint {
	panic("implement me")
}

func (Boolean) ExprNode() {}

func (BoolOr) Start() uint {
	panic("Invalid call")
}
func (BoolOr) ExprNode() {}

func (BoolAnd) Start() uint {
	panic("implement me")
}
func (BoolAnd) ExprNode() {}

func (BoolNotEquals) Start() uint {
	panic("Invalid call")
}
func (BoolNotEquals) ExprNode() {}

func (String) Start() uint {
	panic("Invalid call")
}
func (String) ExprNode() {}

func (Number) Start() uint {
	panic("Invalid call")
}
func (Number) ExprNode() {}

func (BoolGreaterThen) Start() uint {
	panic("Invalid call")
}
func (BoolGreaterThen) ExprNode() {}

func (BoolEquals) Start() uint {
	panic("Invalid call")
}
func (BoolEquals) ExprNode() {}

func (BoolGreatEquals) Start() uint {
	panic("Invalid call")
}
func (BoolGreatEquals) ExprNode() {}

func (BoolNot) Start() uint {
	panic("Invalid call")
}
func (BoolNot) ExprNode() {}

func (BoolLowerThen) Start() uint {
	panic("Invalid call")
}
func (BoolLowerThen) ExprNode() {}

func (BoolLowerThenEqual) Start() uint {
	panic("Invalid call")
}
func (BoolLowerThenEqual) ExprNode() {}

func (Singular) Start() uint {
	panic("Invalid call")
}
func (Singular) StmtNode() {}
func (Singular) ExprNode() {}

func (FunctionCall) Start() uint {
	panic("Invalid call")
}
func (FunctionCall) StmtNode() {}
func (FunctionCall) ExprNode() {}

func (BinarySub) Start() uint {
	panic("Invalid call")
}
func (BinarySub) StmtNode() {}
func (BinarySub) ExprNode() {}

func (BinaryAdd) Start() uint {
	panic("Invalid call")
}
func (BinaryAdd) StmtNode() {}
func (BinaryAdd) ExprNode() {}

func (BinaryDiv) Start() uint {
	panic("Invalid call")
}
func (BinaryDiv) StmtNode() {}
func (BinaryDiv) ExprNode() {}

func (BinaryMul) Start() uint {
	panic("Invalid call")
}
func (BinaryMul) StmtNode() {}
func (BinaryMul) ExprNode() {}

func (UnaryPlus) Start() uint {
	panic("Invalid call")
}
func (UnaryPlus) ExprNode() {}
func (UnaryPlus) StmtNode() {}

func (UnaryMinus) Start() uint {
	panic("Invalid call")
}
func (UnaryMinus) ExprNode() {}
func (UnaryMinus) StmtNode() {}

func (EmptyExpr) Start() uint {
	panic("Invalid call")
}
func (EmptyExpr) StmtNode() {}
func (EmptyExpr) ExprNode() {}

func (Ident) Start() uint {
	panic("Invalid call")
}
func (Ident) StmtNode() {}
func (Ident) ExprNode() {}

