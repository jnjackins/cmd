package main

import (
	"log"
	"os"
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

// for error reporting
var tok string

// Lex returns tokens to the parser. A token is a quoted string,
// a word, or a symbol. Environment variables and globs are expanded
// before they are received by the parser.
func (x *shLex) Lex(lval *shSymType) int {
	if leftover != nil && len(leftover) > 0 {
		lval.word = leftover[0]
		leftover = leftover[1:]
		return WORD
	}
	for {
		r := x.next()
		if wordChar(r) && r != eof {
			word := string(r) + x.getWord()
			expanded := expand(word)
			if len(expanded) == 0 {
				continue
			} else if len(expanded) > 1 {
				leftover = expanded[1:]
			}
			lval.word = expanded[0]
			tok = lval.word
			return WORD
		} else if r == '\'' {
			lval.word = x.getQuoted()
			tok = lval.word
			return WORD
		} else if r != '\n' && unicode.IsSpace(r) {
			continue
		} else {
			switch r {
			case '#':
				x.line = "\n"
				continue
			case '>':
				rr := x.next()
				if rr == '>' {
					tok = ">>"
					return APPEND
				} else {
					x.peek = rr
					tok = string(r)
					return int(r)
				}
			default:
				tok = string(r)
				return int(r)
			}
		}
	}
}

func expand(word string) []string {
	word = os.ExpandEnv(word)
	expanded, err := filepath.Glob(word)
	if err != nil || len(expanded) == 0 {
		if word != "" {
			// TODO: when we have lists, this should be expanded = []string{word}
			expanded = strings.Fields(word)
		}
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
	return !unicode.IsSpace(r) && !strings.ContainsRune("#;&|^=`'{}()<>", r)
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
	log.Printf("token '%s': %s", tok, s)
}
