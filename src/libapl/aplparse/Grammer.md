# Grammer for aplc 
# not finished yet

type -> int | short | cmx | byte | char | string | int | long

declareation -> [decl_stat] decl id type | [fnc_stat] decl fnc id()[->type][:id] | var_explicit_declaration | fnc_explicit_declaration

decl_stat -> fixed | static

fnc_stat  -> virtual | static

var_explicit_declaration -> id type = val | id = val

fnc_explicit_declaration -> [fnc_stat] fnc id([decl_param_list])[->type][:id] {expression}

stmt -> if_stmt | repeat_stmt | mov_stmt | until_stmt | stmt

if_stmt  -> if bool_expr { expression }; | if bool_expr { expression } else {expression}; |  if bool_expr { expression } else if bool_expr {expression}; | if bool_expr { expression } else if bool_expr {expression} else {expression};

repeat_stmt -> repeat {expression};

until_stmt  -> repeat {expression} until bool_expr;

mov_stmt -> mov id in id {expression}; | mov id in fnc_call {expression};

fnc_call -> id([param_list]); | id([param_list])

param_list -> id , param_list | id

decl_param_list ->   id type, decl_param_list| id type

change_value -> id type = val; |id = val; | id += val; id = val + val;  | id -= val; | id /= val; | id *= val; | id ^= val; | id &=val | id %= val; | id >>= val; | id <<=val;

bool_expr -> id == val | fnc_call == val | (bool_expr) | bool_expr || bool_expr | bool_expr && bool_expr | id | bool_value

bool_value -> true | false

expression ->  stmt | change_value| expressio