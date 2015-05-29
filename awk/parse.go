//go:generate go tool yacc -p "awk" action.y

package main

import (
	"log"
	"unicode"
	"unicode/utf8"
)

const eof = 0

var parsedFns []func(record)

func parseAction(action string) []func(record) {
	awkParse(&awkLex{line: action})
	return parsedFns
}

type awkLex struct {
	line string
	peek rune
}

var currentToken string

func (x *awkLex) Lex(lval *awkSymType) int {
	if x.line == "" {
		return eof
	}
	r, n := utf8.DecodeRuneInString(x.line)
	for unicode.IsSpace(r) {
		x.line = x.line[n:]
		r, n = utf8.DecodeRuneInString(x.line)
	}
	switch {
	case r == '"':
		lval.str = x.getString()
		dprintf("lexed string: %#v", lval.str)
		return STR
	case r == '$':
		x.line = x.line[1:]
		lval.str = "$" + x.getSymbol()
		dprintf("lexed field identifier: %s", lval.str)
		return FLD
	case unicode.IsNumber(r):
		lval.num = x.getNum()
		dprintf("lexed number: %#v", lval.num)
		return NUM
	case unicode.IsLetter(r):
		lval.str = x.getSymbol()
		dprintf("lexed symbol: %#v", lval.str)
		if lval.str == "print" {
			return PRINT
		}
		return VAR
	}
	x.line = x.line[n:]
	dprintf("lexed symbol: %c", r)
	currentToken = string(r)
	return int(r)
}

func (x *awkLex) getString() string {
	x.line = x.line[1:]
	var r rune
	var length, n int
	for r != '"' {
		r, n = utf8.DecodeRuneInString(x.line[length:])
		length += n
	}
	s := x.line[:length-1] // don't need the last "
	currentToken = s
	x.line = x.line[length:]
	return s
}

func (x *awkLex) getSymbol() string {
	var n int
	r, length := utf8.DecodeRuneInString(x.line)
	for unicode.IsLetter(r) || unicode.IsNumber(r) || r == '_' {
		length += n
		r, n = utf8.DecodeRuneInString(x.line[length:])
	}

	s := x.line[:length]
	currentToken = s
	x.line = x.line[length:]
	return s
}

func (x *awkLex) getNum() float64 {
	currentToken = "0"
	return 0
}

func (x *awkLex) Error(s string) {
	log.Fatalf("error near token '%s': %s", currentToken, s)
}
