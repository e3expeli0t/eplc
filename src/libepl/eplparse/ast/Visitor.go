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

type Visitor interface {
	Visit(node Node) (v Visitor)
}

type Traveler func(Node) bool

func (f Traveler) Visit(node Node) Visitor {
	if f(node) {
		return f
	}
	return nil
}

func Travel(node Node, f func(Node) bool) {
	Walk(Traveler(f), node)
}
