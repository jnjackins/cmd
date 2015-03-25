package main

import (
	"log"
	"strings"
	"unicode/utf8"
)

const eof = 0

type shLex struct {
	line string
	peek rune
}

func (x *shLex) Lex(yylval *shSymType) int {
	for {
		r := x.next()
		if wordChar(r) {
			s := string(r) + x.getStr()
			yylval.word = s
			return WORD
		} else if r == ' ' {
			continue
		} else {
			return int(r)
		}
	}
}

func (x *shLex) getStr() string {
	var s string
	var r rune
	for {
		r = x.next()
		if wordChar(r) {
			s += string(r)
		} else {
			break
		}
	}
	if r != eof {
		x.peek = r
	}
	return s
}

func wordChar(r rune) bool {
	return !strings.ContainsRune("\n \t#;&|^$=`'{}()<>", r) && r != eof
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
	c, size := utf8.DecodeRuneInString(x.line)
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
