package main

import (
	"bytes"
	"log"
	"strings"
	"unicode"
	"unicode/utf8"
)

const eof = 0

type shLex struct {
	line []byte
	peek rune
}

func (x *shLex) Lex(yylval *shSymType) int {
	for {
		c := x.next()
		if c == eof {
			return eof
		}
		if unicode.IsSpace(c) {
			continue
		}
		if strings.ContainsAny(string(c), "()<>|&;") {
			return x.getSymbol(c, yylval)
		}
		if unicode.IsPrint(c) {
			return x.getWord(c, yylval)
		}
		log.Printf("unrecognized character %q", c)
	}
}

func (x *shLex) getWord(c rune, yylval *shSymType) int {
	add := func(b *bytes.Buffer, c rune) {
		if _, err := b.WriteRune(c); err != nil {
			log.Fatalf("WriteRune: %s", err)
		}
	}
	var b bytes.Buffer
	add(&b, c)

	for {
		c = x.next()
		if strings.ContainsAny(string(c), "()<>|&;") {
			break
		}
		if unicode.IsPrint(c) && !unicode.IsSpace(c) {
			add(&b, c)
		} else {
			break
		}
	}
	if c != eof {
		x.peek = c
	}
	yylval.name = b.String()
	return WORD
}

func (x *shLex) getSymbol(c rune, yylval *shSymType) int {
	switch c {
	case '<':
		yylval.fd = 0
		return REDIR
	case '>':
		yylval.fd = 1
		return REDIR
	case '&':
		d := x.next()
		if d == '&' {
			return AND
		}
		x.peek = d
	case '|':
		d := x.next()
		if d == '|' {
			return OR
		}
		x.peek = d
	}
	return int(c)
}

func (x *shLex) next() rune {
	if x.peek != eof {
		r := x.peek
		x.peek = eof
		return r
	}
	if len(x.line) == 0 {
		return eof
	}
	c, size := utf8.DecodeRune(x.line)
	x.line = x.line[size:]
	if c == utf8.RuneError && size == 1 {
		log.Print("invalid utf8")
		return x.next()
	}
	return c
}

func (x *shLex) Error(s string) {
	log.Printf("parse error: %s", s)
}
