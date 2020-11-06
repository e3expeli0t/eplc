package symboltable

type TableCode uint

type TableMap struct {
	Map map[TableCode]*ScopeSymbolTable
	Size uint64

	currentTableCode  TableCode
	previousTableCode TableCode
}

func (tm *TableMap) Insert(table ScopeSymbolTable) {
	tm.previousTableCode = tm.currentTableCode
	tm.currentTableCode++
	tm.Map[tm.currentTableCode] = &table
}

//Locate Walks backwards to find the requested symbol
func (tm *TableMap) Locate(name string) SymbolData {
	var symbol SymbolData

	counter := tm.currentTableCode
	for ; counter >= 0; counter-- {
		symbol = tm.Map[counter].Lookup(name)
		if symbol != CantLocate() {
			return symbol
		}
	}

	tm.Restore()

	return CantLocate()
}

//LocateInScope searches for the symbol in the current SymbolTable
func (tm *TableMap) LocateInScope(name string) SymbolData {
	return tm.Map[tm.currentTableCode].Lookup(name)
}

func (tm *TableMap) Reset() {
	tm.previousTableCode = tm.currentTableCode
	tm.currentTableCode = 0
}

func (tm *TableMap) Restore() {
	tmp := tm.previousTableCode

	tm.previousTableCode = tm.currentTableCode
	tm.currentTableCode = tmp
}