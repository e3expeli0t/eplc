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

package Types

import (
	"encoding/binary"
	"eplc/src/libepl/eplparse/ast"
)

type ObjectValue interface {
	Convert(t EplType) bool
}

type EvalValue struct {
	ValueType EplType
	Val []byte
}

type UnknownValue struct {
	ValueType EplType
	exp ast.Expression
}

func (uv UnknownValue) Convert(t EplType) bool {
	panic("Unimplemented")
}

//Prototype
////todo: type conversions
func (v *EvalValue) Convert(t EplType) bool {
	if v.ValueType.Attributes == t.Attributes {
		v.ValueType = t
	} else if t.Equals(TypeString) && v.ValueType.IsMathematical() {
		//Strings are represented in memory as sequence of bytes followed by null terminator
		switch v.ValueType {
		case EplType(TypeUint):

		case EplType(TypeUint8):
			//not supported yet
			break
		case EplType(TypeUint16):
		case EplType(TypeUint32):
		case EplType(TypeUint64):
		case EplType(TypeInt):
		case EplType(TypeInt8):
		case EplType(TypeInt16):
		case EplType(TypeInt32):
		case EplType(TypeInt64):
		case EplType(TypeFloat):
		case EplType(TypeFloat8):
		case EplType(TypeUint16):
		case EplType(TypeFloat32):
		case EplType(TypeFloat64):
		default:
			return false
		}
	}
	return true
}

func MakeString(str string) EvalValue {
	return EvalValue{
		ValueType: EplType(TypeString),
		Val:       []byte(str),
	}
}

func MakeUint(val uint) EvalValue {
	slice := make([]byte, 8)

	binary.BigEndian.PutUint64(slice, uint64(val))

	return EvalValue{
		ValueType: EplType(TypeUint),
		Val:       slice,
	}
}


func MakeUint16(val uint16) EvalValue {
	slice := make([]byte, 2)

	binary.BigEndian.PutUint16(slice, val)

	return EvalValue{
		ValueType: EplType(TypeUint16),
		Val:       slice,
	}
}


func MakeUint32(val uint32) EvalValue {
	slice := make([]byte, 4)

	binary.BigEndian.PutUint32(slice, val)

	return EvalValue{
		ValueType: EplType(TypeUint),
		Val:       slice,
	}
}


func MakeUint64(val uint64) EvalValue {
	slice := make([]byte, 8)

	binary.BigEndian.PutUint64(slice, val)

	return EvalValue{
		ValueType: EplType(TypeUint),
		Val:       slice,
	}
}
