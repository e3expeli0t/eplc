package epllex

import (
	"testing"
	"strings"
)


func checkType(lexme string, tt TokenType) bool {
	lx := New(strings.NewReader(lexme), "test.epl")
	
	return lx.Next().Ttype == tt;
}


func TestLexer(t *testing.T) {
	table := []struct {
		inToken string
		outType TokenType
	}{
		{"if", IF},
		{"else", ELSE},
		{"repeat", REPEAT},
		{"until", UNTIL},
		{"move", MOVE},
		{"in", IN},
		{"static", STATIC},
		{"fixed", FIXED},
		{"decl", DECL},
		{"define", DEFINE},
		{"type", TYPE},
		{"struct", STRUCT},
		{"fnc", FNC},
		{"+", PLUS},
		{"-", MINUS},
		{"/", DEV},
		{"*", MULT},
		{"{", LBRACE},
		{"}", RBRACE},
		{"(", LPAR},
		{")", RPAR},
		{"[", LSUBSCRIPT},
		{"]", RSUBSCRIPT},
		{"+=", PLUSEQUAL},
		{"-=", MINUSEQUAL},
		{"*=", MULTEQUAL},
		{"/=", DEVEQUAL},
		{"=", ASSIGN},
		{">", GT},
		{"<", LT},
		{"==", EQ},
		{">=", GE},
		{"<=", LE},
		{"&&", AND},
		{"||", OR},
		{"!", NOT},
		{"|", UNARYOR},
		{"&", UNARYAND},
		{"<<", LSHIFT},
		{">>", RSHIFT},
		{";", SEMICOLON},
		{":", RETURN_IND},
		{",", COMMA},
		{".", DOT},
		{"int", INT},
		{"int16", INT16},
		{"int32", INT32},
		{"int64", INT64},
		{"float16", FLOAT16},
		{"float64", FLOAT64},
		{"flot64", FLOAT64},
		{"float", FLOAT},
		{"long", LONG},
		{"cmx64", CMX64},
		{"cmx", CMX},
		{"uint16", UINT16},
		{"uint32", UINT32},
		{"uint64", UINT64},
		{"uint", UINT},
		{"string", STRING},
	}

	for _, tst := range table{
		if !checkType(tst.inToken, tst.outType){
			t.Errorf("Couldn't match token '%s' to type %d", tst.inToken, tst.outType)
		}
	}

	t.Log(INT)
}