package apllex

/*
	regex:
		+ is 'One or more times'
		* is 'Zero or more time'
		? is 'Zero or one times'
		| is 'Or'
		None is 'Zero times'
		[0-9] is 'All the numbers'
		[a-z] is 'All lowercase letters'
		[A-Z] is 'All uppercase letters'
*/

type TokenType int

const (
	LBRACE         TokenType = iota // {
	RBRACE         TokenType = iota // }
	LPAR           TokenType = iota // (
	RPAR           TokenType = iota // )
	ID             TokenType = iota // ([a-zA-z])+(_)?([0-9])* for example abCd_233
	NUM            TokenType = iota // 0-9
	COMMA          TokenType = iota // ,
	PLUS           TokenType = iota // +
	MINUS          TokenType = iota // -
	MULT           TokenType = iota // *
	DEV            TokenType = iota // /
	PLUSEQUAL      TokenType = iota // +=
	MINUSEQUAL     TokenType = iota // -=
	MULTEQUAL      TokenType = iota // *=
	DEVEQUAL       TokenType = iota // /=
	GT             TokenType = iota // >
	EQ             TokenType = iota // ==
	GE             TokenType = iota // >=
	LE             TokenType = iota // <=
	LT             TokenType = iota // <
	IF             TokenType = iota // if
	ELSE           TokenType = iota // else
	ASSIGN         TokenType = iota // =
	AND            TokenType = iota // &&
	OR             TokenType = iota // ||
	REAL           TokenType = iota // ([0-9])+(.)?([0-9)*
	UNARYOR        TokenType = iota // |
	UNARYAND       TokenType = iota // &
	MOVE           TokenType = iota // move
	UNTIL          TokenType = iota // until
	REPEAT         TokenType = iota // repeat
	IN             TokenType = iota // in
	FNC            TokenType = iota // fnc
	STATIC         TokenType = iota // static
	FIXED          TokenType = iota // fixed
	DECL           TokenType = iota // decl
	SDECL          TokenType = iota // sdecl
	DEFINE         TokenType = iota // define
	VIRTUAL        TokenType = iota // virtual
	STRUCT         TokenType = iota // struct
	TYPE           TokenType = iota // type
	DECLTYPE       TokenType = iota // decltype
	NOT            TokenType = iota // !
	SEMICOLON      TokenType = iota // ;
	DOT            TokenType = iota // .
	STRINGLITERAL  TokenType = iota // "[ID]+[NUM]+ | [NUM]+[ID]+ | [NUM]+ | [ID]+ | NONE"
	LSUBSCRIPT     TokenType = iota // [
	RSUBSCRIPT     TokenType = iota // ]
	STATIC_BLOCK_S TokenType = iota // :
	//DECREASE TokenType = iota // --
	//INCREASE TokenType = iota // ++
	EOF TokenType = -(iota + 1) // '\0'
)
