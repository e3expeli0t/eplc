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

package libepl

func NewInfoStruct (filename string) InfoStruct {
	return InfoStruct{
		Filename:  filename,
		SystemCPU: ResolveBits(),
	}
}

//todo: temporary solution
func ResolveBits() int {
	const is64Bit = uint64(^uintptr(0)) == ^uint64(0)

	if is64Bit {
		return 64
	} else {
		return 32
	}
}

type InfoStruct struct {
	Filename string
	SystemCPU int
}
