%{

package main

import (
	"io"
	"os/exec"
)

%}

%union {
	word	string
	words	[]string
	cmd		*exec.Cmd
	pipe	[]*exec.Cmd
	line	[][]*exec.Cmd
}

%type <words>	args
%type <cmd>		cmd
%type <pipe>	pipe
%type <line>	line

%token <word> WORD

%%

top		: '\n'
		| line				{ runLine($1) }

line	: pipe '\n'			{ $$ = [][]*exec.Cmd{$1} }
		| pipe ';' '\n'		{ $$ = [][]*exec.Cmd{$1} }
		| pipe ';' line		{ $$ = append($3, $1) }

pipe	: cmd				{ $$ = []*exec.Cmd{$1} }
		| pipe '|' cmd		{ connect($1[len($1)-1], $3); $$ = append($1, $3) }

cmd		: args				{ $$ = &exec.Cmd{Path: $1[0], Args: $1} }
		| cmd '>' WORD		{ $$.Stdout = create($3); defer close($$.Stdout.(io.Closer)) }
		| cmd '<' WORD		{ $$.Stdin = open($3); defer close($$.Stdin.(io.Closer)) }

args	: WORD				{ $$ = []string{$1} }
		| args WORD			{ $$ = append($1, $2) }
