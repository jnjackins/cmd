%{

package main

import "os/exec"

%}

%union {
	word	string
	words	[]string
	asgn	struct{}
	cmd		*exec.Cmd
	pipe	[]*exec.Cmd
	line	[][]*exec.Cmd
}

%type <word>	word
%type <words>	args
%type <asgn>	asgn
%type <cmd>		cmd
%type <cmd>		expr
%type <pipe>	pipe
%type <line>	line

%token <word> WORD

// APPEND is >>
%left '<' '>' APPEND

%%

top		: '\n'
		| line				{ runLine($1) }

line	: pipe '\n'			{ $$ = [][]*exec.Cmd{$1} }
		| pipe ';' '\n'		{ $$ = [][]*exec.Cmd{$1} }
		| pipe ';' line		{ $$ = append($3, $1) }
		| asgn '\n'			{ updateEnv() }
		| asgn ';' line		{ updateEnv(); $$ = $3 }

pipe	: expr				{ $$ = []*exec.Cmd{$1} }
		| pipe '|' expr		{ connect($1[len($1)-1], $3); $$ = append($1, $3) }

expr	: cmd
		| asgn cmd			{ $$ = $2 }

cmd		: args				{ $$ = &exec.Cmd{Path: $1[0], Args: $1} }
		| cmd '<' word		{ $$.Stdin = open($3, 'r'); defer close($$.Stdin) }
		| cmd '>' word		{ $$.Stdout = open($3, 'w'); defer close($$.Stdout) }
		| cmd APPEND word	{ $$.Stdout = open($3, 'a'); defer close($$.Stdout) }

// TODO: this should not be parsed as asgn if it is a command argument
asgn	: word '=' word		{ env[$1] = $3; $$ = struct{}{} }

args	: word				{ $$ = []string{$1} }
		| args word			{ $$ = append($1, $2) }

word	: WORD
		| word '^' WORD		{ $$ = $1 + $3 }
