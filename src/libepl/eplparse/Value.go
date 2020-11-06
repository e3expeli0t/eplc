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

package eplparse

import (
	"encoding/binary"
	"eplc/src/libepl/Types"
)

type ObjectValue interface {
	Convert(t Types.EplType) bool
}

type EvalValue struct {
	ValueType Types.EplType
	Val       []byte
}

//Prototype
////todo: type conversions
func (v *EvalValue) Convert(t Types.EplType) bool {
	if v.ValueType.Attributes == t.Attributes {
		v.ValueType = t
	} else if t.Equals(Types.TypeString) && v.ValueType.IsMathematical() {
		//Strings are represented in memory as sequence of bytes followed by null terminator
		switch v.ValueType {
		case Types.EplType(Types.TypeUint):

		case Types.EplType(Types.TypeUint8):
			//not supported yet
			break
		case Types.EplType(Types.TypeUint16):
		case Types.EplType(Types.TypeUint32):
		case Types.EplType(Types.TypeUint64):
		case Types.EplType(Types.TypeInt):
		case Types.EplType(Types.TypeInt8):
		case Types.EplType(Types.TypeInt16):
		case Types.EplType(Types.TypeInt32):
		case Types.EplType(Types.TypeInt64):
		case Types.EplType(Types.TypeFloat):
		case Types.EplType(Types.TypeFloat8):
		case Types.EplType(Types.TypeUint16):
		case Types.EplType(Types.TypeFloat32):
		case Types.EplType(Types.TypeFloat64):
		default:
			return false
		}
	}
	return true
}

func MakeString(str string) EvalValue {
	return EvalValue{
		ValueType: Types.EplType(Types.TypeString),
		Val:       []byte(str),
	}
}

func MakeUint(val uint) EvalValue {
	slice := make([]byte, 8)

	binary.BigEndian.PutUint64(slice, uint64(val))

	return EvalValue{
		ValueType: Types.EplType(Types.TypeUint),
		Val:       slice,
	}
}


func MakeUint16(val uint16) EvalValue {
	slice := make([]byte, 2)

	binary.BigEndian.PutUint16(slice, val)

	return EvalValue{
		ValueType: Types.EplType(Types.TypeUint16),
		Val:       slice,
	}
}


func MakeUint32(val uint32) EvalValue {
	slice := make([]byte, 4)

	binary.BigEndian.PutUint32(slice, val)

	return EvalValue{
		ValueType: Types.EplType(Types.TypeUint),
		Val:       slice,
	}
}


func MakeUint64(val uint64) EvalValue {
	slice := make([]byte, 8)

	binary.BigEndian.PutUint64(slice, val)

	return EvalValue{
		ValueType: Types.EplType(Types.TypeUint),
		Val:       slice,
	}
}
