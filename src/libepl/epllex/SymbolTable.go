package epllex

type SymbolTable struct {
	Table map[string]*SymbolData
} 

type SymbolData struct {
	_type TokenType
	symbol string
	scope string
}


func (st *SymbolTable) Add(s *SymbolData) {
	st.Table[s.symbol] = s
}

func (st *SymbolTable) AddType(symbol string, t TokenType) {
	st.Table[symbol]._type = t
}

func (st *SymbolTable) AddScope(symbol string, scope string) {
	st.Table[symbol].scope = scope
}

func (st *SymbolTable) GetType(symbol string) TokenType {
	return st.Table[symbol]._type
}


func (st *SymbolTable) GetScope(symbol string) string {
	return st.Table[symbol].scope
}

func (st *SymbolTable) Get(symbol string) SymbolData {
	if st.Table[symbol] != nil {
		return *st.Table[symbol]
	}

	return SymbolData{}
} 