package epllex

import (
	"testing"
	"strings"
	"reflect"
)

func Create(symbols...string) SymbolTable{
	table := map[string]*SymbolData{}

	for _, s := range symbols {
		table[s] = &SymbolData{symbol: s}
	}

	return SymbolTable{Table: table}
}

func TestSymbolTable(t *testing.T) {
	
	tables := []struct {
		inID string 
		SymbolTableStat SymbolTable
	}{
		{"hello", Create("hello")},
		{"world", Create("hello", "world")},
		{"this", Create("hello", "world", "this")},
		{"works", Create("hello", "world", "this", "works")},
	}
	lx := New(strings.NewReader("hello world this works"), "test_symbols.epl")

	for _, table := range tables {
		if lx.Next(); !reflect.DeepEqual(lx.ST, table.SymbolTableStat) {
			t.Errorf(" Symbol table for '%s' symbol is not equal", table.inID)
		}
	}
}