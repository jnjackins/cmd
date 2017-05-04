%{// generated by goyacc

package main
%}

%term IF THEN FI FOR IN DO DONE
%term WORD QUOTE REDIR APPEND
%term SIMPLE WORDS PAREN
%left AND OR

%union{
        tree *treeNode
}
%type <tree> term line cmd list pipe redir item words word
%type <tree> IF THEN FI FOR IN DO DONE
%type <tree> WORD QUOTE REDIR APPEND

%%

tree    : /* empty */
        | line                                   { execute($1) }

// A term is like a line, but always terminated by a semicolon.
// Terms are used in compound statements like if and for statements.
term    : list ';'
        | cmd term                               { $$ = mkTree(';', $1, $2) }

line    : list
        | cmd
        | cmd line                               { $$ = mkTree(';', $1, $2) }

cmd     : list ';'
        | list '&'                               { $$ = mkTree('&', $1) }

list    : pipe
        | list AND pipe                          { $$ = mkTree(AND, $1, $3) }
        | list OR pipe                           { $$ = mkTree(OR, $1, $3) }

pipe    : redir
        | pipe '|' redir                         { $$ = mkTree('|', $1, $3) }

redir   : item
        | redir REDIR WORD                       { $$ = $1; $1.redirect($2.int,$3.string, false) }
        | redir APPEND WORD                      { $$ = $1; $1.redirect($2.int,$3.string, true) }

item    : words                                  { $$ = mkTree(SIMPLE, $1) }
        | '(' line ')'                           { $$ = mkTree(PAREN, $2) }
        | IF term THEN term FI                   { $$ = mkTree(IF, $2, $4) }
        | FOR WORD IN words ';' DO term DONE     { $$ = mkTree(FOR, $2, $4, $7) }

words   : word                                   { $$ = mkTree(WORDS, $1) }
        | words word                             { $$ = $1; $1.children = append($1.children, $2) }

word    : WORD
        | QUOTE
