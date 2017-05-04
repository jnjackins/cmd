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
	c := x.next(false)
	for {
		if unicode.IsSpace(c) {
			c = x.next(false)
		} else {
			break
		}
	}
	if c == eof {
		return eof
	}

	switch c {
	case ';', '(', ')':
		return int(c)

	case '|', '&':
		d := x.next(false)
		if d == c {
			if c == '&' {
				return AND
			}
			if c == '|' {
				return OR
			}
		}
		x.peek = d
		return int(c)

	case '<':
		yylval.tree = mkLeaf(REDIR, 0, "")
		return REDIR

	case '>':
		if peek := x.next(false); peek == '>' {
			yylval.tree = mkLeaf(APPEND, 1, "")
			return APPEND
		} else {
			yylval.tree = mkLeaf(REDIR, 1, "")
			x.peek = peek
			return REDIR
		}

	case '\'':
		var b bytes.Buffer
		for {
			c = x.next(true)
			if c == eof {
				panic("TODO: multiline quoted strings")
			}
			if c == '\'' {
				if d := x.next(false); d != '\'' {
					x.peek = d
					break
				}
			}
			b.WriteRune(c)
		}

		yylval.tree = mkLeaf(QUOTE, 0, b.String())
		return QUOTE
	}

	// if we made it this far, it's a regular word.
	return x.readWord(c, yylval)
}

func (x *shLex) readWord(c rune, yylval *shSymType) int {
	var b bytes.Buffer
	b.WriteRune(c)

	for {
		c = x.next(false)
		if strings.ContainsAny(string(c), "()<>|&;'") {
			break
		}
		if unicode.IsPrint(c) && !unicode.IsSpace(c) {
			b.WriteRune(c)
		} else {
			break
		}
	}
	if c != eof {
		x.peek = c
	}

	s := b.String()

	switch s {
	case "if":
		return IF
	case "then":
		return THEN
	case "fi":
		return FI
	case "for":
		return FOR
	case "in":
		return IN
	case "do":
		return DO
	case "done":
		return DONE
	}

	yylval.tree = mkLeaf(WORD, 0, s)
	return WORD
}

func (x *shLex) next(inQuote bool) rune {
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
	if c == '#' && !inQuote {
		// nothing more to do for this line
		return eof
	}
	if c == utf8.RuneError && size == 1 {
		log.Print("invalid utf8")
		return x.next(inQuote)
	}
	return c
}

func (x *shLex) Error(s string) {
	log.Printf("parse error: %s", s)
}
