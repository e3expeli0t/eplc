package Types

import "eplc/src/libepl/epllex"
import "encoding/binary"


type EplType struct {
	Tname string //name as text (ex: Int)
	Tkey uint64 //Id for the Type
}

func (et *EplType) ToBasic() BasicType {
	return BasicType{et.Tname, et.Tkey}
}

type BasicType EplType

func BasicTypes() (ta []BasicType) {

	names := []string {"uint", "uint8", "uint16", "uint32", "uint64", "int", "int8", "int16", "int32", "int64",
		"float", "float8", "float16", "float32", "float64", "cmx", "cmx64", }

	for _,n :=range names {
		ta = append(ta, MakeType(n).ToBasic())
	}

	return
}

func IsValidBasicType(token epllex.Token) bool {
	for _, t := range BasicTypes() {
		if ResolveType(token).ToBasic() == t {
			return true
		}
	}

	return false
}

func MakeType(name string) EplType {
	return EplType{name, genKey(name)}
}

func ResolveType(token epllex.Token) EplType {
	return EplType{token.Lexme, genKey(token.Lexme)}
}


func genKey(n string) uint64 {
	data := binary.LittleEndian.Uint64([]byte(n))
	return data ^ 0x45504C54595045
}