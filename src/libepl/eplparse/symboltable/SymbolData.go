package symboltable

import (
	"eplc/src/libepl"
	"eplc/src/libepl/Types"
)

type SymbolData struct {
	Line uint
	Offset uint

	Name string
	Type Types.EplType
	Kind SymbolKind
}

func NewTypedSymbol(name string, eplType Types.EplType, info libepl.LocationInfo, kind SymbolKind) *SymbolData {
	return &SymbolData{
		Line:   info.Line,
		Offset: info.Offset,
		Name:   name,
		Type:  	eplType,
		Kind:   kind,
	}
}
