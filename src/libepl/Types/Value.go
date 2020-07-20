package Types

import "encoding/binary"

type Value struct {
	ValueType EplType
	Val []byte
}

//Prototype
////todo: type conversions
func (v *Value) Convert(t EplType) bool {
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

func MakeString(str string) Value {
	return Value{
		ValueType: EplType(TypeString),
		Val:       []byte(str),
	}
}

func MakeUint(val uint) Value {
	slice := make([]byte, 8)

	binary.BigEndian.PutUint64(slice, uint64(val))

	return Value{
		ValueType: EplType(TypeUint),
		Val:       slice,
	}
}


func MakeUint16(val uint16) Value {
	slice := make([]byte, 2)

	binary.BigEndian.PutUint16(slice, val)

	return Value{
		ValueType: EplType(TypeUint16),
		Val:       slice,
	}
}


func MakeUint32(val uint32) Value {
	slice := make([]byte, 4)

	binary.BigEndian.PutUint32(slice, val)

	return Value{
		ValueType: EplType(TypeUint),
		Val:       slice,
	}
}


func MakeUint64(val uint64) Value {
	slice := make([]byte, 8)

	binary.BigEndian.PutUint64(slice, val)

	return Value{
		ValueType: EplType(TypeUint),
		Val:       slice,
	}
}
