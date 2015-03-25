%{

package main

%}

%union {
	word	string
	words	[]string
	cmd		cmd
	line	[]cmd
}

%type <words>	args
%type <cmd>		cmd
%type <line>	line

%token <word> WORD

%%

top:	  '\n'
		| line				{ doline($1) }

line:	  cmd '\n'			{ $$ = []cmd{$1} }
		| cmd ';' '\n'		{ $$ = []cmd{$1} }
		| cmd ';' line		{ $$ = append($3, $1) }

cmd:	  args				{ $$.args = $1 }
		| cmd '>' WORD		{ $$.stdout = $3 }
		| cmd '<' WORD		{ $$.stdin = $3 }

args:	  WORD				{ $$ = []string{$1} }
		| args WORD			{ $$ = append($1, $2) }

%%

func doline(line []cmd) {
	// TODO: change parsing rules to avoid iterating backwards?
	for i := len(line)-1; i >= 0; i-- {
		run(line[i])
	}
}