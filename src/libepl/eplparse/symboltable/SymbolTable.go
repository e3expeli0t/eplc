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

package symboltable

import (
	"eplc/src/libepl/eplparse/Types"
)
/*
	Basic  linked list that holds all the symbols inside scopes
*/

type ScopeType string

const (
	//BLOCK = evrey staemeant (if, repeat, move, etc)
	BLOCK ScopeType = "@@@@BLOCK@@@@"

	/*
		FUNCTION Scope is the scope thats starts right after function 
		declaration (includes the arguments)
	*/
	FUNCTION ScopeType = "@@@@FUNCTION@@@@"
	
	//GLOBAL Scope is every symbol that wase declare out side of the function (include the functions itself)
	GLOBAL ScopeType = "@@@@GLOBAL@@@@"
)

/*
	SymbolTable stores the information about the symbols.
	In this case the symboltable data structure is (some kind of) linked list
*/
type SymbolTable struct {
	Table map[string]*SymbolData
	Prev *SymbolTable
	Next *SymbolTable
} 

//SymbolData stores the information about symbols
type SymbolData struct {
	SType Types.EplType
	symbol string
	scope ScopeType
}

//New creates new empty SymbolTable
func New() SymbolTable {
	return SymbolTable{Table: map[string]*SymbolData{}, Prev: nil, Next: nil}
}

func NewBasicSymbol(s string) *SymbolData {
	return &SymbolData{symbol:s}
}

func NewSymbol(s string, t Types.EplType, scope ScopeType) *SymbolData {
	return &SymbolData{t, s, scope}
}


//Add new symbol
func (st *SymbolTable) Add(s *SymbolData) {
	st.Table[s.symbol] = s
}

func (st *SymbolTable) AddType(symbol string, t Types.EplType) {
	st.Table[symbol].SType = t
}

func (st *SymbolTable) AddScope(symbol string, scope ScopeType) {
	st.Table[symbol].scope = scope
}

func (st *SymbolTable) GetType(symbol string) Types.EplType {
	return st.Table[symbol].SType
}

func (st *SymbolTable) GetScope(symbol string) ScopeType {
	return st.Table[symbol].scope
}

//Get symbol
func (st *SymbolTable) Get(symbol string) SymbolData {
	if st.Table[symbol] != nil {
		return *st.Table[symbol]
	}

	return SymbolData{}
} 