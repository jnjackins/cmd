package main

import (
	"log"
	"path/filepath"
	"strings"
	"unicode"
	"unicode/utf8"
)

const eof = 0

type shLex struct {
	line string
	peek rune
}

// words resulting from a glob expansion that need to be returned by
// following calls to Lex
var leftover []string

func (x *shLex) Lex(yylval *shSymType) int {
	if leftover != nil && len(leftover) > 0 {
		yylval.word = leftover[0]
		leftover = leftover[1:]
		return WORD
	}
	for {
		r := x.next()
		if wordChar(r) && r != eof {
			word := string(r) + x.getWord()
			expanded := expand(word)
			if len(expanded) > 1 {
				leftover = expanded[1:]
			}
			yylval.word = expanded[0]
			return WORD
		} else if r == '\'' {
			yylval.word = x.getQuoted()
			return WORD
		} else if r != '\n' && unicode.IsSpace(r) {
			continue
		} else {
			return int(r)
		}
	}
}

func expand(word string) []string {
	expanded, err := filepath.Glob(word)
	if err != nil || len(expanded) == 0 {
		expanded = []string{word}
	}
	return expanded
}

func (x *shLex) getWord() string {
	var s string
	var r rune
	for {
		r = x.next()
		if wordChar(r) && r != eof {
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
	return !unicode.IsSpace(r) && !strings.ContainsRune("#;&|^$=`'{}()<>", r)
}

func (x *shLex) getQuoted() string {
	var s string
	var r, rr rune
	for {
		r = x.next()
		if r != '\'' && r != eof {
			s += string(r)
		} else if r == '\'' {
			rr = x.next()
			// double single-quote inside a single-quoted string is a literal single-quote
			if rr == '\'' {
				s += string(rr)
			} else {
				x.peek = rr
				break
			}
		} else {
			break
		}
	}
	if r != '\'' && r != eof {
		x.peek = r
	}
	return s
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
