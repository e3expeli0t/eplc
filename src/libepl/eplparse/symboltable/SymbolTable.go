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

package symboltable

import (
	"eplc/src/libepl/Types"
)

//Note: Currently this is not in use
//Attribute for more convenient SymbolMap searches
type ScopeType uint
type SymbolKind uint

const (
	//BLOCK = every statement (if, repeat, move, etc)
	BLOCK ScopeType = iota

	/*
		FUNCTION Scope is the scope that's starts right after function
		declaration (including the arguments)
	*/
	FUNCTION

	/*
		GLOBAL Scope is the whole file scope meaning that the symbols that belongs
		to the global scope can be used in the WHOLE file
	*/
	GLOBAL
)

const (
	_ SymbolKind = iota
	Function
	Variable
	Unknown

	//Locates default error value
	EmptySymbol
)

func  CantLocate() SymbolData {
	return SymbolData{Kind: EmptySymbol}
}

func NewScopedSymbolTable(Scope ScopeType) *ScopeSymbolTable {
	return &ScopeSymbolTable {
		Scope: Scope,
		Table: nil,
	}
}

//Every scope has a ScopeSymbolTable
type ScopeSymbolTable struct {
	Scope ScopeType
	Table map[string]*SymbolData
}

func (st *ScopeSymbolTable) Clear() {
	st.Table= make(map[string]*SymbolData)
}

func (st* ScopeSymbolTable) Exists(name string) bool {
	_, found := st.Table[name]
	return found
}

func (st* ScopeSymbolTable) Lookup(name string) SymbolData {
	if !st.Exists(name) {
		return CantLocate()
	}

	return *st.Table[name]
}

func (st* ScopeSymbolTable) GetSymbolType(name string) Types.EplType {
	return st.Table[name].Type
}

func (st* ScopeSymbolTable) SetSymbolType(name string, t Types.EplType) {
	st.Table[name].Type = t
}

func (st* ScopeSymbolTable) Add(data *SymbolData) {
	st.Table[data.Name] = data
}

func (st* ScopeSymbolTable) SetScopeType(scope ScopeType) {
	st.Scope = scope
}
