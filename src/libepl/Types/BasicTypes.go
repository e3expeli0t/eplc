package Types

var (
	TypeUint    = EplType{"uint", 0x45504c21303e31}
	TypeUint8   = EplType{"uint8", 0x4550393d37247d}
	TypeUint16  = EplType{"uint16", 0x4525253a2d6173}
	TypeUint32  = EplType{"uint32", 0x4525253a2d6377}
	TypeUint64  = EplType{"uint64", 0x4525253a2d6671}
	TypeInt     = EplType{"int", 0x45504c54303e31}
	TypeInt8    = EplType{"int8", 0x45504c3d37247d}
	TypeInt16   = EplType{"int16", 0x4550253a2d6173}
	TypeInt32   = EplType{"int32", 0x4550253a2d6377}
	TypeInt64   = EplType{"int64", 0x4550253a2d6671}
	TypeFloat   = EplType{"float", 0x45502a38363131}
	TypeFloat8  = EplType{"float8", 0x4536203b38247d}
	TypeFloat16 = EplType{"float16", 0x233c23352d6173}
	TypeFloat32 = EplType{"float32", 0x233c23352d6377}
	TypeFloat64 = EplType{"float64", 0x233c23352d6671}
	TypeCmx     = EplType{"cmx", 0x45504c543a3d3d}
	TypeCmx64   = EplType{"cmx64", 0x45502f39216671}
	TypeString  = EplType{"string", 0x45233826303e22}
	TypeBool    = EplType{"bool", 0x45504c36363f29}
)
