%{

package main

%}

%union {
	word	string
	words	[]string
	cmd		cmd
}

%type <words>	args
%type <cmd>		cmd

%token <word> WORD

%%

line:	  '\n'
		| cmd '\n'		{ run($1) }

cmd:	  args			{ $$.args = $1 }
		| cmd '>' WORD	{ $$.stdout = $3 }
		| cmd '<' WORD	{ $$.stdin = $3 }

args:	  WORD			{ $$ = []string{$1} }
		| args WORD		{ $$ = append($1, $2) }

%%
