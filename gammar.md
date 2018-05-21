type ::= "int8" | "int16" | "int32" | "int64"| "int" | "uint8" | "uint16" | "uint32" | "uint64" | "uint" | "long" | "float8" | "float16" | "float32" | "float64" |"float"| "string" | "cmx64" | "cmx" | "bool"

number ::= REAL | NUM | CMX

id ::= ID

str ::= STRINGLITERAL
bool_val ::= "true" | "false"

bool_expr ::= bool_val | bool_op | (bool_expr) | ["!"]bool_expr| bool_expr

bool_ops ::= "&&" | "==" | "||"

bool_op ::= val bool_ops val

if_stmt ::= "if" bool_expr "{"expression"}" "else" [op_if] "{"expression"}" [else`]

else` ::= "else" "{" expression "}"

op_if ::= "if" bool_expr

repeat_loop_var ::= "("id int = val`")"

repeat` ::= "repeat" [repeat_loop_var]

repeat ::= repeat` "{" expression "}"

repeat_until ::= repeat` "{" expression "}" "until" bool_expr   

until ::= "until" bool_expr "{" expression "}" 

move ::= "move" id "in" func_call' "{" expression "}"

func_call ::= id"(" param_list ")";

func_call' ::= id"(" param_list ")"

param_list ::= id param_list | ,param_list | E

unary_ops ::= "+" | "-" | "!"

ops ::= "+" | "-" | "*" | "/" | ">>" | "<<" | "|" | "&" | "^" | "%"

math_op ::= val' ops val';

assign ::= "+=" |"-=" | "/=" | "*=" | "%=" | "&=" | "|=" | "^="| "<<=" | ">>="

var_decl ::= "declare" id type; | id type = val; 

var_assign ::= id assign val;

val ::= [unary_ops] id | [unary_ops] number | str

val' ::= [unary_ops]id | [unary_ops] number

stmt ::= if_stmt | repeat | until | move

expression ::= stmt | func_call | math_op | var_decl | var_assign | expression

