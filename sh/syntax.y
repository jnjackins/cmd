%{// generated by goyacc

package main
%}

%union{
    name string
    fd    int
    node  node
    redir redirecter
    pipe  connecter
    arg   *argNode
}

%type <node> list cmd line
%type <pipe> pipe
%type <redir> item redir
%type <arg> words

%token <name> WORD
%token <fd> REDIR

%left OR
%left AND

%%

tree    : /* empty */
        | line             { execute($1) }

line    : list
        | cmd
        | cmd line         { $$ = &listNode{typ: typeListSequence, left: $1, right: $2} }

cmd     : list ';'
        | list '&'         { $$ = &forkNode{tree: $1} }

list    : pipe             { $$ = $1.(node) }
        | list AND pipe    { $$ = &listNode{typ: typeListAnd, left: $1, right: $3} }
        | list OR pipe     { $$ = &listNode{typ: typeListOr, left: $1, right: $3} }

pipe    : redir            { $$ = $1.(connecter) }
        | pipe '|' redir   { $$ = &pipeNode{left: $1.(connecter), right: $3.(connecter)} }

redir   : item
        | item REDIR WORD  { $$ = $1; $1.redirect($2, $3) }

item    : words            { $$ = &simpleNode{args: $1} }
        | '(' cmd ')'      { $$ = &parenNode{tree: $2} }

words   : WORD             { $$ = &argNode{val: $1} }
        | WORD words       { $$ = &argNode{val: $1}; $$.next = $2 }