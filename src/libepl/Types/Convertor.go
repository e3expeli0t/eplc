package Types

import (
	"encoding/binary"
)

//todo: type conversions
func FromUint(val []byte, target []byte){
	// For now we assume that we are on 64bit cpu
	binary.BigEndian.Uint64(val)

}