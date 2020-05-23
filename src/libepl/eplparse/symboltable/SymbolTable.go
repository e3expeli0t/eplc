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

type ScopeType string

const (
	//BLOCK = every statement (if, repeat, move, etc)
	BLOCK ScopeType = "@@@@BLOCK@@@@"

	/*
		FUNCTION Scope is the scope that's starts right after function
		declaration (including the arguments)
	*/
	FUNCTION ScopeType = "@@@@FUNCTION@@@@"

	/*
		GLOBAL Scope is the whole file scope meaning that the symbols that belongs
		to the global scope can be used in the WHOLE file
	*/
	GLOBAL ScopeType = "@@@@GLOBAL@@@@"
)

/*
	SymbolTable stores the information about the symbols.
	In this case the symboltable data structure is (some kind of) linked list.
	This simple linked lis holds all the information about symbols that belongs to
	the current file
*/
type SymbolTable struct {
	Table map[string]*SymbolData
	Prev  *SymbolTable
	Next  *SymbolTable
	CurrentScope ScopeType
}

//SymbolData stores the information about symbols
type SymbolData struct {
	symbol string
	scope  ScopeType
	SType  Types.EplType
}

//New creates new empty SymbolTable
func New() SymbolTable {
	return SymbolTable{Table: map[string]*SymbolData{}, Prev: nil, Next: nil}
}

func NewBasicSymbol(s string) *SymbolData {
	return &SymbolData{symbol: s}
}

func NewSymbol(s string, t Types.EplType, scope ScopeType) *SymbolData {
	return &SymbolData{
		symbol: s,
		scope:  scope,
		SType:  t,
	}
}

func NewTypedSymbol(s string, t Types.EplType) *SymbolData {
	return &SymbolData{
		SType:  t,
		symbol: s,
	}
}

//returns the current Symbol Table without the prev
func (st *SymbolTable) Strip() SymbolTable {
	tmp := *st
	tmp.Prev = nil
	tmp.Next = nil

	return tmp
}

func (st *SymbolTable) First() SymbolTable {
	tmp := *st
	for tmp.Prev != nil {
		tmp = *tmp.Prev
	}

	return tmp
}

func (st *SymbolTable) Last() SymbolTable {
	tmp := *st
	for tmp.Next != nil {
		tmp = *tmp.Next
	}

	return tmp
}

//Add new symbol
func (st *SymbolTable) Add(s *SymbolData) {
	st.Table[s.symbol] = s
}

//Sets the SymbolData scope to the current scope
func (st *SymbolTable) AddWOScope(s *SymbolData) {
	s.scope = st.CurrentScope
	st.Table[s.symbol] = s
}

func (st *SymbolTable) AddType(symbol string, t Types.EplType) {
	st.Table[symbol].SType = t
}

func (st *SymbolTable) SetSymbolScope(symbol string, scope ScopeType) {
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

func (st *SymbolTable) SetScopeType(scope ScopeType) {
	st.CurrentScope = scope
}
