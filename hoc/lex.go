package main

import (
	"fmt"
	"io"
	"strings"
	"unicode"
)

const eof = 0

var (
	lineNum   int
	lastToken string
)

type hocLex struct {
	r    *strings.Reader
	peek rune
}

func (x *hocLex) Lex(lval *hocSymType) int {
	for {
		c, _, err := x.r.ReadRune()
		if err != nil {
			if err == io.EOF {
				return eof
			}
			elog.Print(err)
		}
		switch {
		case unicode.IsSpace(c):
			continue
		case c == '.' || unicode.IsNumber(c):
			if err := x.r.UnreadRune(); err != nil {
				elog.Print(err)
			}
			n, err := fmt.Fscanf(x.r, "%f", &lval.val)
			if err != nil {
				elog.Print(err)
			}
			if n != 1 {
				elog.Print("error parsing number")
			}
			lastToken = fmt.Sprint(lval.val)
			return NUMBER
		case unicode.IsLetter(c):
			if err := x.r.UnreadRune(); err != nil {
				elog.Print(err)
			}
			word := x.getWord()
			lval.sym = word
			return VAR
		}
		if c == '\n' {
			lineNum++
		}
		lastToken = string(c)
		return int(c)
	}
}

func (x *hocLex) getWord() string {
	var word string
	for {
		c, _, err := x.r.ReadRune() 
		if err == io.EOF {
			break
		}
		if err != nil {
			elog.Print(err)
			break
		}
		if unicode.IsLetter(c) || unicode.IsNumber(c) {
			word += string(c)
		} else {
			x.r.UnreadRune()
			break
		}
	}
	
	return word
}

func (x *hocLex) Error(s string) {
	elog.Printf("error near line %d: token %#v: %s", lineNum, lastToken, s)
}
