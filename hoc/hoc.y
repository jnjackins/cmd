%{

package main

import (
	"fmt"
	"math"
)

%}

%union {
	val float64
	sym string
}

%token	<val> NUMBER
%token	<sym> VAR BUILTIN UNDEF
%type	<val> expr asgn
%right	'='
%left	'+' '-'
%left	'*' '/'
%left	UNARYMINUS
%right	'^' // exponentiation

%%

hoc		: expr						{ fmt.Printf("\t%.8v\n", $1) }
		| asgn

asgn	: VAR '=' expr				{ $$ = $3; symbols[$1] = newVar($3) }

expr	: NUMBER					{ $$ = $1 }
		| VAR						{ $$ = symbols[$1].val }
		| BUILTIN '(' expr ')'		{ $$ = symbols[$1].fn($3) }
		| expr '+' expr				{ $$ = $1 + $3 }
		| expr '-' expr				{ $$ = $1 - $3 }
		| expr '*' expr				{ $$ = $1 * $3 }
		| expr '/' expr				{ $$ = $1 / $3 }
		| expr '^' expr				{ $$ = math.Pow($1, $3) }
		| '(' expr ')'				{ $$ = $2 }
		| '-' expr %prec UNARYMINUS	{ $$ = -1 * $2 }
