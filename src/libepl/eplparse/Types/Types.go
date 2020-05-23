package Types

import (
	"eplc/src/libepl/epllex"
	"strconv"
)

/*
The type system:
	the type system supplies the parser with handful
	functions and info about types

	The type system should:
	* Provide type resolving
	* Disguise between defined types and user defined types
	* For every type create unique id
	* Define type hierarchy
 */
type TypeSystem struct {
	lexer epllex.Lexer
	Fname string
	BasicTypes []BasicType
	TypeMap map[string]*EplType
}

type EplType struct {
	TypeName string //name as text (ex: Int)
	TypeKey  uint64 //Id for the Type
}
type BasicType EplType

/*
design note:
	basic type is a pre defined fundamental data type such as: uint , int etc...
 */
func (et *EplType) ToBasic() BasicType {
	return BasicType{et.TypeName, et.TypeKey}
}

//Todo: change this to more efficient way
func (ts *TypeSystem) Initialize(lex epllex.Lexer) {
	ts.TypeMap = make(map[string]*EplType)

	names := []string{"uint", "uint8", "uint16", "uint32", "uint64", "int", "int8", "int16", "int32", "int64",
		"float", "float8", "float16", "float32", "float64", "cmx", "cmx64","string"}

	for _, n := range names {
		ts.BasicTypes = append(ts.BasicTypes, (ts.MakeType(n)).ToBasic())
	}

	//for informative errors. Really inefficient
	ts.Fname = lex.Filename
	ts.lexer = lex
}

func (ts *TypeSystem) GetType(name string) *EplType {
	return ts.TypeMap[name]
}

func (ts *TypeSystem) IsValidBasicType(token epllex.Token) bool {
	tp := ts.ToType(token).ToBasic()
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
		TypeName: name,
		TypeKey:  ts.genKey(name),
	}

	if !ts.typeDefined(name) {
		ts.TypeMap[name] = _type
	}
	return _type
}

func (ts *TypeSystem) ToType(token epllex.Token) *EplType {
	return &EplType{
		TypeName: token.Lexme,
		TypeKey:  ts.genKey(token.Lexme),
	}
}

// Generate unique type key
func (ts *TypeSystem) genKey(n string) uint64 {
	data, _ := strconv.ParseUint(n, 10, 0)
	return data ^ 0x45504C54595045
}

func (ts *TypeSystem) ResolveValueType(token epllex.Token) *EplType {
	panic("This function is not supported")
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
	}

	return nil
}

func (ts *TypeSystem) typeDefined(name string) bool {
	_, ok := ts.TypeMap[name]
	return ok
}
