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
	"eplc/src/libepl"
	"eplc/src/libepl/Types"
)

type SymbolData struct {
	Line uint
	Offset uint

	Name string
	Type Types.EplType
	Kind SymbolKind
}

func NewTypedSymbol(name string, eplType Types.EplType, info libepl.LocationInfo, kind SymbolKind) *SymbolData {
	return &SymbolData{
		Line:   info.Line,
		Offset: info.Offset,
		Name:   name,
		Type:  	eplType,
		Kind:   kind,
	}
}
