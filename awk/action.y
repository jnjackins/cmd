%{

package main

%}

%union {
	num		float64
	str		string
	fnc		func(record)
	fncs	[]func(record)
	args	[]*symbol
}

%type <fncs>	stmts
%type <fnc>		stmt
%type <fnc>		print
%type <args>	args

%token <num> NUM
%token <str> STR VAR FLD
%token PRINT

%%

awk		: stmts				{ parsedFns = $1 }

stmts	: stmt				{ $$ = []func(record){$1} }
		| stmts ';' stmt	{ $$ = append($1, $3)}

stmt	: print

print	: PRINT args		{ $$ = printFn($2) }

args	: 					{ $$ = make([]*symbol, 0, 1) }
		| args NUM			{ $$ = append($1, &symbol{fval: $2, tval: numVal}) }
		| args STR			{ $$ = append($1, &symbol{sval: $2, tval: strVal}) }
		| args FLD			{ $$ = append($1, getSymbol($2)) }
