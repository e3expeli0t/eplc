#epl grammar

|symbol | meaning|
|--|--|
|(production) | means one or more occurrences of the production "production"|
|[production] | means zero or one occurrences of the production "production"|
|E|epsilon. empty production |

type ::= "int8" | "int16" | "int32" | "int64"| "int" | "uint8" | "uint16" | "uint32" | "uint64" | "uint" | "long" | "float8" | "float16" | "float32" | "float64" |"float"| "string" | "cmx64" | "cmx" | "bool"

number ::= REAL | NUM | CMX
id ::= ID
str ::= STRINGLITERAL

bool_val ::= TRUE | FALSE
bool_expr ::= bool_val | bool_op | (bool_expr) | ["!"]bool_expr| bool_expr
bool_ops ::= "&&" | "==" | "||" | ">" | "<" | ">=" | "<=" 
bool_op ::= val bool_ops val

else` ::= "else" "{" expression "}"
op_if ::= "if" bool_expr
if_stmt ::= "if" bool_expr "{"expression"}" "else" [op_if] "{"expression"}" [else`]

unary_ops  ::= `"+" | "-" | "!"`   
ops ::= `"+" | "-" | "*" | "/" | ">>" | "<<" | "|" | "&" | "^" | "%"`
assign ::= "+=" |"-=" | "/=" | "*=" | "%=" | "&=" | "|=" | "^="| "<<=" | ">>="

val ::= [unary_ops] id | [unary_ops] number | str
val' ::= [unary_ops]id | [unary_ops] number


repeat_loop_var ::= `"(" scoped_var_decl ")"`
repeat ::= `"repeat" [repeat_loop_var]`

repeat ::= repeat ["("var_explicit_decl")"]` "{" expression "}"
repeat_until ::= repeat` "{" expression "}" "until" bool_expr   
until ::= "until" bool_expr "{" expression "}" 
move ::= "move" id "in" func_call' "{" expression "}"

param_list ::= id param_list | ,param_list | E
func_call ::= id"(" param_list ")";
func_call' ::= id"(" param_list ")"

math_op ::= val' ops val';
var_assign ::= id assign val;

var_stat ::= "fixed"
var_decl ::= "declare" [var_stat] id type; | [var_stat] id type = val; 
fnc_decl ::= fnc id "(" param_list ")" [:type] [->id];
fnc_impl ::= fnc id "(" param_list ")" [:type] [->id] "{" expression "}"

import_list ::= id import_list | . import_list | E
import ::= "import" import_list;

stmt ::= if_stmt | repeat | until | move

block ::= "{"expression"}" | "{"block"}"
expression ::= stmt | func_call | math_op | var_decl | var_assign | block | expression |

program ::= [(import)] [(fnc_decl)] (var_decl) fnc_impl