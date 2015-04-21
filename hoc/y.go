//line hoc.y:2
package main

import __yyfmt__ "fmt"

//line hoc.y:3
import (
	"fmt"
	"math"
)

//line hoc.y:12
type hocSymType struct {
	yys int
	val float64
	sym string
}

const NUMBER = 57346
const VAR = 57347
const BUILTIN = 57348
const UNDEF = 57349
const UNARYMINUS = 57350

var hocToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"NUMBER",
	"VAR",
	"BUILTIN",
	"UNDEF",
	"'='",
	"'+'",
	"'-'",
	"'*'",
	"'/'",
	"UNARYMINUS",
	"'^'",
	"'('",
	"')'",
}
var hocStatenames = [...]string{}

const hocEofCode = 1
const hocErrCode = 2
const hocMaxDepth = 200

//line yacctab:1
var hocExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
}

const hocNprod = 14
const hocPrivate = 57344

var hocTokenNames []string
var hocStates []string

const hocLast = 49

var hocAct = [...]int{

	2, 9, 10, 11, 12, 15, 13, 13, 16, 18,
	19, 20, 21, 22, 23, 24, 25, 9, 10, 11,
	12, 14, 13, 1, 27, 9, 10, 11, 12, 3,
	13, 0, 26, 4, 17, 6, 4, 5, 6, 8,
	0, 0, 8, 0, 7, 11, 12, 7, 13,
}
var hocPact = [...]int{

	32, -1000, -8, -1000, -1000, 13, -10, 29, 29, 29,
	29, 29, 29, 29, 29, 29, 16, -1000, -7, 34,
	34, -7, -7, -7, -8, 8, -1000, -1000,
}
var hocPgo = [...]int{

	0, 0, 29, 23,
}
var hocR1 = [...]int{

	0, 3, 3, 2, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1,
}
var hocR2 = [...]int{

	0, 1, 1, 3, 1, 1, 4, 3, 3, 3,
	3, 3, 3, 2,
}
var hocChk = [...]int{

	-1000, -3, -1, -2, 4, 5, 6, 15, 10, 9,
	10, 11, 12, 14, 8, 15, -1, 5, -1, -1,
	-1, -1, -1, -1, -1, -1, 16, 16,
}
var hocDef = [...]int{

	0, -2, 1, 2, 4, 5, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 5, 13, 7,
	8, 9, 10, 11, 3, 0, 12, 6,
}
var hocTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	15, 16, 11, 9, 3, 10, 3, 12, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 8, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 14,
}
var hocTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 13,
}
var hocTok3 = [...]int{
	0,
}

//line yaccpar:1

/*	parser for yacc output	*/

var (
	hocDebug        = 0
	hocErrorVerbose = false
)

type hocLexer interface {
	Lex(lval *hocSymType) int
	Error(s string)
}

type hocParser interface {
	Parse(hocLexer) int
	Lookahead() int
}

type hocParserImpl struct {
	lookahead func() int
	state     func() int
}

func (p *hocParserImpl) Lookahead() int {
	return p.lookahead()
}

func hocNewParser() hocParser {
	p := &hocParserImpl{
		lookahead: func() int { return -1 },
		state:     func() int { return -1 },
	}
	return p
}

const hocFlag = -1000

func hocTokname(c int) string {
	if c >= 1 && c-1 < len(hocToknames) {
		if hocToknames[c-1] != "" {
			return hocToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func hocStatname(s int) string {
	if s >= 0 && s < len(hocStatenames) {
		if hocStatenames[s] != "" {
			return hocStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func hocErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !hocErrorVerbose {
		return "syntax error"
	}
	res := "syntax error: unexpected " + hocTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := hocPact[state]
	for tok := TOKSTART; tok-1 < len(hocToknames); tok++ {
		if n := base + tok; n >= 0 && n < hocLast && hocChk[hocAct[n]] == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if hocDef[state] == -2 {
		i := 0
		for hocExca[i] != -1 || hocExca[i+1] != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; hocExca[i] >= 0; i += 2 {
			tok := hocExca[i]
			if tok < TOKSTART || hocExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if hocExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += hocTokname(tok)
	}
	return res
}

func hoclex1(lex hocLexer, lval *hocSymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = hocTok1[0]
		goto out
	}
	if char < len(hocTok1) {
		token = hocTok1[char]
		goto out
	}
	if char >= hocPrivate {
		if char < hocPrivate+len(hocTok2) {
			token = hocTok2[char-hocPrivate]
			goto out
		}
	}
	for i := 0; i < len(hocTok3); i += 2 {
		token = hocTok3[i+0]
		if token == char {
			token = hocTok3[i+1]
			goto out
		}
	}

out:
	if token == 0 {
		token = hocTok2[1] /* unknown char */
	}
	if hocDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", hocTokname(token), uint(char))
	}
	return char, token
}

func hocParse(hoclex hocLexer) int {
	return hocNewParser().Parse(hoclex)
}

func (hocrcvr *hocParserImpl) Parse(hoclex hocLexer) int {
	var hocn int
	var hoclval hocSymType
	var hocVAL hocSymType
	var hocDollar []hocSymType
	hocS := make([]hocSymType, hocMaxDepth)

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	hocstate := 0
	hocchar := -1
	hoctoken := -1 // hocchar translated into internal numbering
	hocrcvr.state = func() int { return hocstate }
	hocrcvr.lookahead = func() int { return hocchar }
	defer func() {
		// Make sure we report no lookahead when not parsing.
		hocstate = -1
		hocchar = -1
		hoctoken = -1
	}()
	hocp := -1
	goto hocstack

ret0:
	return 0

ret1:
	return 1

hocstack:
	/* put a state and value onto the stack */
	if hocDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", hocTokname(hoctoken), hocStatname(hocstate))
	}

	hocp++
	if hocp >= len(hocS) {
		nyys := make([]hocSymType, len(hocS)*2)
		copy(nyys, hocS)
		hocS = nyys
	}
	hocS[hocp] = hocVAL
	hocS[hocp].yys = hocstate

hocnewstate:
	hocn = hocPact[hocstate]
	if hocn <= hocFlag {
		goto hocdefault /* simple state */
	}
	if hocchar < 0 {
		hocchar, hoctoken = hoclex1(hoclex, &hoclval)
	}
	hocn += hoctoken
	if hocn < 0 || hocn >= hocLast {
		goto hocdefault
	}
	hocn = hocAct[hocn]
	if hocChk[hocn] == hoctoken { /* valid shift */
		hocchar = -1
		hoctoken = -1
		hocVAL = hoclval
		hocstate = hocn
		if Errflag > 0 {
			Errflag--
		}
		goto hocstack
	}

hocdefault:
	/* default state action */
	hocn = hocDef[hocstate]
	if hocn == -2 {
		if hocchar < 0 {
			hocchar, hoctoken = hoclex1(hoclex, &hoclval)
		}

		/* look through exception table */
		xi := 0
		for {
			if hocExca[xi+0] == -1 && hocExca[xi+1] == hocstate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			hocn = hocExca[xi+0]
			if hocn < 0 || hocn == hoctoken {
				break
			}
		}
		hocn = hocExca[xi+1]
		if hocn < 0 {
			goto ret0
		}
	}
	if hocn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			hoclex.Error(hocErrorMessage(hocstate, hoctoken))
			Nerrs++
			if hocDebug >= 1 {
				__yyfmt__.Printf("%s", hocStatname(hocstate))
				__yyfmt__.Printf(" saw %s\n", hocTokname(hoctoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for hocp >= 0 {
				hocn = hocPact[hocS[hocp].yys] + hocErrCode
				if hocn >= 0 && hocn < hocLast {
					hocstate = hocAct[hocn] /* simulate a shift of "error" */
					if hocChk[hocstate] == hocErrCode {
						goto hocstack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if hocDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", hocS[hocp].yys)
				}
				hocp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if hocDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", hocTokname(hoctoken))
			}
			if hoctoken == hocEofCode {
				goto ret1
			}
			hocchar = -1
			hoctoken = -1
			goto hocnewstate /* try again in the same state */
		}
	}

	/* reduction by production hocn */
	if hocDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", hocn, hocStatname(hocstate))
	}

	hocnt := hocn
	hocpt := hocp
	_ = hocpt // guard against "declared and not used"

	hocp -= hocR2[hocn]
	// hocp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if hocp+1 >= len(hocS) {
		nyys := make([]hocSymType, len(hocS)*2)
		copy(nyys, hocS)
		hocS = nyys
	}
	hocVAL = hocS[hocp+1]

	/* consult goto table to find next state */
	hocn = hocR1[hocn]
	hocg := hocPgo[hocn]
	hocj := hocg + hocS[hocp].yys + 1

	if hocj >= hocLast {
		hocstate = hocAct[hocg]
	} else {
		hocstate = hocAct[hocj]
		if hocChk[hocstate] != -hocn {
			hocstate = hocAct[hocg]
		}
	}
	// dummy call; replaced with literal code
	switch hocnt {

	case 1:
		hocDollar = hocS[hocpt-1 : hocpt+1]
		//line hoc.y:28
		{
			fmt.Printf("\t%.8v\n", hocDollar[1].val)
		}
	case 3:
		hocDollar = hocS[hocpt-3 : hocpt+1]
		//line hoc.y:31
		{
			hocVAL.val = hocDollar[3].val
			symbols[hocDollar[1].sym] = newVar(hocDollar[3].val)
		}
	case 4:
		hocDollar = hocS[hocpt-1 : hocpt+1]
		//line hoc.y:33
		{
			hocVAL.val = hocDollar[1].val
		}
	case 5:
		hocDollar = hocS[hocpt-1 : hocpt+1]
		//line hoc.y:34
		{
			hocVAL.val = symbols[hocDollar[1].sym].val
		}
	case 6:
		hocDollar = hocS[hocpt-4 : hocpt+1]
		//line hoc.y:35
		{
			hocVAL.val = symbols[hocDollar[1].sym].fn(hocDollar[3].val)
		}
	case 7:
		hocDollar = hocS[hocpt-3 : hocpt+1]
		//line hoc.y:36
		{
			hocVAL.val = hocDollar[1].val + hocDollar[3].val
		}
	case 8:
		hocDollar = hocS[hocpt-3 : hocpt+1]
		//line hoc.y:37
		{
			hocVAL.val = hocDollar[1].val - hocDollar[3].val
		}
	case 9:
		hocDollar = hocS[hocpt-3 : hocpt+1]
		//line hoc.y:38
		{
			hocVAL.val = hocDollar[1].val * hocDollar[3].val
		}
	case 10:
		hocDollar = hocS[hocpt-3 : hocpt+1]
		//line hoc.y:39
		{
			hocVAL.val = hocDollar[1].val / hocDollar[3].val
		}
	case 11:
		hocDollar = hocS[hocpt-3 : hocpt+1]
		//line hoc.y:40
		{
			hocVAL.val = math.Pow(hocDollar[1].val, hocDollar[3].val)
		}
	case 12:
		hocDollar = hocS[hocpt-3 : hocpt+1]
		//line hoc.y:41
		{
			hocVAL.val = hocDollar[2].val
		}
	case 13:
		hocDollar = hocS[hocpt-2 : hocpt+1]
		//line hoc.y:42
		{
			hocVAL.val = -1 * hocDollar[2].val
		}
	}
	goto hocstack /* stack new state and value */
}
