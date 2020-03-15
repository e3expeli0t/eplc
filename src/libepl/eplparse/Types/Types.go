package Types

import (
	"eplc/src/libepl/epllex"
	"eplc/src/libepl/eplparse/Types/errors"
	"fmt"
	"strconv"
)


type TypeSystem struct {
	lexer epllex.Lexer
	Fname string
	BasicTypes []BasicType
	TypeMap map[string]EplType
}

type EplType struct {
	Tname string //name as text (ex: Int)
	Tkey  uint64 //Id for the Type
}

/*
design note:
	basic type is a pre defined fundamental data type such as: uint , int etc...
 */
func (et *EplType) ToBasic() BasicType {
	return BasicType{et.Tname, et.Tkey}
}

type BasicType EplType

//Todo: change this to more efficient way
func (ts *TypeSystem) Initialize(lex epllex.Lexer) {
	ts.TypeMap = make(map[string]EplType)

	names := []string{"uint", "uint8", "uint16", "uint32", "uint64", "int", "int8", "int16", "int32", "int64",
		"float", "float8", "float16", "float32", "float64", "cmx", "cmx64",}

	for _, n := range names {
		ts.BasicTypes = append(ts.BasicTypes, (ts.MakeType(n)).ToBasic())
	}
	ts.Fname = lex.Filename
	ts.lexer = lex
}

func (ts *TypeSystem) IsValidBasicType(token epllex.Token) bool {
	tp := ts.ResolveType(token).ToBasic()
	for _, t := range ts.BasicTypes {
		if tp == t {
			return true
		}
	}

	return false
}

//create new type with name and type key
func (ts *TypeSystem) MakeType(name string) *EplType {
	_type := &EplType{
		Tname: name,
		Tkey:  ts.genKey(name),
	}

	if !ts.typeDefined(name) {
		ts.TypeMap[name] = *_type
	}
	return _type
}

func (ts *TypeSystem)ResolveType(token epllex.Token) *EplType {
	return &EplType{
		Tname: token.Lexme,
		Tkey:  ts.genKey(token.Lexme),
	}
}

// Generate unique type key
func (ts *TypeSystem) genKey(n string) uint64 {
	data, _ := strconv.ParseUint(n, 10, 0)
	return data ^ 0x45504C54595045
}


//todo: rewrite type system on version 0.2
func (ts *TypeSystem) ResolveValueType(token epllex.Token) *EplType {
	switch token.Ttype {
	case epllex.NUM:
		//the default number value is int ( the size is defined by the target system)
		return ts.MakeType("int") // the call does not matter
	case epllex.REAL:
		//the default  real number value is float ( the size is defined by the target system)
		return ts.MakeType("float") // the call does not matter
	case epllex.STRINGLITERAL:
		//todo: define string handling system
		//string is special runtime defined type
		return ts.MakeType("string")
	case epllex.ID: // never should hit
	break
	default:
		errors.UnresolvedTypeError(
			fmt.Sprintf("Couldn't resolve type of '%s'", token.Lexme),
			token,
			ts.Fname,
			ts.lexer.GetLine(),
			)

	}

	return nil
}


func (ts *TypeSystem) typeDefined(name string) bool {
	_, ok := ts.TypeMap[name]
	return ok
}
