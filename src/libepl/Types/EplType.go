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

import "strconv"

type TypeAttribute uint

const (
	TypeAttrMathematical TypeAttribute = 1 << iota
	TypeAttrObject
	TypeAttrBasic // basic language type like uint
	TypeAttrInherits // It is an extension of other type
	TypeAttrExtendable // Other Types can extend it
	TypeAttrDefault
)

func NewType(name string, attributes TypeAttribute) EplType {
	data, _ := strconv.ParseUint(name, 10, 0)

	return EplType{
		TypeName: name,
		TypeKey: data ^ 0x45504C54595045,
		Attributes: TypeAttrObject | TypeAttrExtendable | TypeAttrInherits | attributes,
	}
}

//Add basic type attribute?
func NewBasicType(name string) BasicType {
	return BasicType(NewType(name,TypeAttrDefault))
}

func NewMathematicalBasicType(name string) BasicType {
	return BasicType(NewType(name, TypeAttrDefault | TypeAttrMathematical))
}

type EplType struct {
	TypeName string //name as text (ex: Int)
	TypeKey  uint64 //Id for the Type
	Attributes TypeAttribute // Type attributes
}
type BasicType EplType

/*
design note:
	basic type is a pre defined fundamental data type such as: uint , int etc...
*/

func (bt *BasicType) Equals(et EplType) bool {
	return bt.TypeKey == et.TypeKey
}

func (bt *BasicType) AsEplType() EplType {
	return EplType(*bt)
}

func (et *EplType) Equals(t BasicType) bool {
	return et.TypeKey == t.TypeKey
}

func (et *EplType) ToBasic() BasicType {
	return BasicType{TypeName: et.TypeName, TypeKey: et.TypeKey, Attributes: et.Attributes}
}

func (et *EplType) IsMathematical() bool {
	return et.HasAttr(TypeAttrMathematical)
}

func (et *EplType) SetAttribute(attr TypeAttribute) {
	et.Attributes = et.Attributes | attr
}

func (et *EplType) HasAttr(attr TypeAttribute) bool {
	return et.Attributes & attr != 0
}

func (et *EplType) ResetAttr() {
	et.Attributes = TypeAttrObject | TypeAttrExtendable | TypeAttrInherits
}