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

type TableCode int

const EndOfTable TableCode = -1

type TableMap struct {
	Map  map[TableCode]*ScopeSymbolTable
	Last TableCode

	currentTableCode  TableCode
	previousTableCode TableCode
}

func NewTableMap() *TableMap {
	newMap := make(map[TableCode]*ScopeSymbolTable)
	return &TableMap{
		Map:               newMap,
		Last:              0,
		currentTableCode:  0,
		previousTableCode: 0,
	}
}

func (tm *TableMap) Insert(table ScopeSymbolTable) {
	tm.Next()

	tm.Map[tm.currentTableCode] = &table
	tm.Last++
}

//Locate Walks backwards to find the requested symbol
func (tm *TableMap) Locate(name string) SymbolData {
	var symbol SymbolData

	counter := tm.currentTableCode
	for ; counter >= 0; counter-- {
		symbol = tm.Map[counter].Lookup(name)
		if symbol != CantLocate() {
			return symbol
		}
	}

	tm.Restore()

	return CantLocate()
}

//LocateInScope searches for the symbol in the current SymbolMap
func (tm *TableMap) LocateInScope(name string) SymbolData {
	return tm.Map[tm.currentTableCode].Lookup(name)
}

func (tm *TableMap) Reset() {
	tm.previousTableCode = tm.currentTableCode
	tm.currentTableCode = 0
}

func (tm *TableMap) Restore() {
	tmp := tm.previousTableCode

	tm.previousTableCode = tm.currentTableCode
	tm.currentTableCode = tmp
}

func (tm *TableMap) HasNext() bool {
	return tm.Last > tm.currentTableCode
}

func (tm *TableMap) Next() {
	tm.previousTableCode = tm.currentTableCode
	tm.currentTableCode++
}

func (tm *TableMap) GetCurrent() *ScopeSymbolTable {
	return tm.Map[tm.currentTableCode]
}

func (tm *TableMap) SwitchTable(code TableCode) {
	tm.previousTableCode = tm.currentTableCode
	tm.currentTableCode = code
}
