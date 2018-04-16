# Grammer for aplc 
# not finished yet

type -> int | short | cmx | byte | char | string | int | long

declareation -> [decl_stat] decl id type | [fnc_stat] decl fnc id()[->type][:id] | var_explicit_declaration | fnc_explicit_declaration

decl_stat -> fixed | static

fnc_stat  -> virtual | static

var_explicit_declaration -> id type = val | id = val

fnc_explicit_declaration -> [fnc_stat] fnc id([param_list])[->type][:id]
