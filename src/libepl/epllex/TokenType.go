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

package epllex

/*
	regex:
		+ is 'One or more times'
		* is 'Zero or more time'
		? is 'Zero or one times'
		| is 'Or'
		None is 'Zero times'
		[0-9] is 'All the numbers'
		[a-z] is 'All lowercase letters'
		[A-Z] is 'All uppercase letters'
*/

type TokenType int

const (
	BOOL   			TokenType = iota //bool type
	INT   			TokenType = iota //int type
	INT8   			TokenType = iota //int8 type
	INT16  			TokenType = iota //int16 type
	INT32   		TokenType = iota //int32 type
	INT64   		TokenType = iota //int64 type
	UINT   			TokenType = iota //uint type
	UINT8   		TokenType = iota //uint8 type
	UINT16   		TokenType = iota //uint16 type
	UINT32   		TokenType = iota //uint32 type
	UINT64   		TokenType = iota //uint64 type
	FLOAT16   	   	TokenType = iota //float16 type	
	FLOAT32   		TokenType = iota //float32 type
	FLOAT64   	   	TokenType = iota //float64 type
	FLOAT   	   	TokenType = iota //float type
	CMX64   	   	TokenType = iota //cmx64 type
	CMX   		   	TokenType = iota //cmx type
	LONG 			TokenType = iota // long type
	STRING   	   	TokenType = iota //string type
	
	LBRACE         	TokenType = iota // {
	RBRACE         	TokenType = iota // }
	LPAR           	TokenType = iota // (
	RPAR           	TokenType = iota // )
	LSUBSCRIPT     	TokenType = iota // [
	RSUBSCRIPT     	TokenType = iota // ]
	
	ID             	TokenType = iota // ([a-zA-z])+(_)?([0-9])* for example abCd_233
	STRINGLITERAL  	TokenType = iota // "[ID]+[NUM]+ | [NUM]+[ID]+ | [NUM]+ | [ID]+ | NONE"
	NUM            	TokenType = iota // 0-9
	REAL           	TokenType = iota // ([0-9])+(.)?([0-9)*
	
	PLUS           	TokenType = iota // +
	MINUS          	TokenType = iota // -
	MULT           	TokenType = iota // *
	DEV            	TokenType = iota // /

	PLUSEQUAL      	TokenType = iota // +=
	MINUSEQUAL     	TokenType = iota // -=
	MULTEQUAL      	TokenType = iota // *=
	DEVEQUAL       	TokenType = iota // /=
	ASSIGN         	TokenType = iota // =

	GT             	TokenType = iota // >
	EQ             	TokenType = iota // ==
	GE             	TokenType = iota // >=
	LE             	TokenType = iota // <=
	LT             	TokenType = iota // <
	AND            	TokenType = iota // &&
	OR             	TokenType = iota // ||
	NOT            	TokenType = iota // !
	
	UNARYOR        	TokenType = iota // |
	UNARYAND       	TokenType = iota // &
	LSHIFT		   	TokenType = iota // <<
	RSHIFT		   	TokenType = iota // >>

	IF             	TokenType = iota // if
	ELSE           	TokenType = iota // else
	
	MOVE           	TokenType = iota // move
	UNTIL          	TokenType = iota // until
	REPEAT         	TokenType = iota // repeat
	IN             	TokenType = iota // in
	FNC            	TokenType = iota // fnc
	STATIC         	TokenType = iota // static
	FIXED          	TokenType = iota // fixed
	DECL           	TokenType = iota // decl
	DEFINE         	TokenType = iota // define
	STRUCT         	TokenType = iota // struct
	TYPE           	TokenType = iota // type
	
	SEMICOLON      	TokenType = iota // ;
	DOT            	TokenType = iota // .
	COMMA          	TokenType = iota // ,
	RETURN_IND 		TokenType = iota // :
	CFLAG			TokenType = iota // @id
	
	EOF 			TokenType = -(iota + 1) // '\0'
)
