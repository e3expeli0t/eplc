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

import "fmt"

type ErrorLevel uint

const (
	Major ErrorLevel = iota
	Minor ErrorLevel = iota
	Fatal ErrorLevel = iota
)

func (l ErrorLevel) ToString() string {
	switch l {
	case Fatal:
		return "Fatal"
	case Major:
		return "Major"
	case Minor:
		return "Minor"
	default:
		return ""
	}
}

type Error interface {
	IsFatal() bool
	ToString() string
}

func NewError(level ErrorLevel, msg string) *TypeError {
	return &TypeError{
		Descriptor: msg,
		Level:      level,
	}
}
type TypeError struct {
	Line       int
	Offset     int
	Descriptor string
	Level ErrorLevel
}

func (t TypeError) IsFatal() bool {
	return t.Level != Fatal
}
//as of now the compiler don't keep track of line numbers
func (t TypeError) ToString() string {
	return fmt.Sprintf("%s Error: %s", t.Level.ToString(), t.Descriptor)
}
