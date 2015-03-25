%{

package main

const (
	normal = iota
	subshell
)

type expr struct{
	kind int
	args []string
}

%}

%union {
	word    string
	words   []string
	expr    expr
	exprs   []expr
}

%type <word>	word
%type <words>	cmd
%type <expr>	expr
%type <exprs>	exprs

%token <word> WORD

%%

line:	  exprs '\n'		{ doline($1) }

exprs:	  expr				{ $$ = []expr{$1} }
		| exprs ';' expr	{ $$ = append($1, $3) }

expr:	  cmd				{ $$.kind = normal; $$.args = $1 }
		| '(' cmd ')'		{ $$.kind = subshell; $$.args = $2 }
	
cmd:	  word				{ $$ = []string{$1} }
		| cmd word			{ $$ = append($1, $2) }

word :	  WORD
		| word '^' WORD			{ $$ = $1 + $3 }

%%

func doline(exprs []expr) {
	for _, expr := range exprs {
		switch expr.kind {
		case normal:
			run(expr.args)
		case subshell:
			// TODO
			run(expr.args)
		}
	}
}
