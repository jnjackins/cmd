//line sh.y:2
package main

import __yyfmt__ "fmt"

//line sh.y:3
//line sh.y:7
type shSymType struct {
	yys   int
	word  string
	words []string
	cmd   cmd
}

const WORD = 57346

var shToknames = []string{
	"WORD",
}
var shStatenames = []string{}

const shEofCode = 1
const shErrCode = 2
const shMaxDepth = 200

//line sh.y:30

//line yacctab:1
var shExca = []int{
	-1, 1,
	1, -1,
	-2, 0,
}

const shNprod = 8
const shPrivate = 57344

var shTokenNames []string
var shStates []string

const shLast = 11

var shAct = []int{

	6, 7, 8, 5, 2, 11, 10, 9, 1, 3,
	4,
}
var shPact = []int{

	-1, -1000, -1000, -5, 3, -1000, -1000, 2, 1, -1000,
	-1000, -1000,
}
var shPgo = []int{

	0, 10, 9, 8,
}
var shR1 = []int{

	0, 3, 3, 2, 2, 2, 1, 1,
}
var shR2 = []int{

	0, 1, 2, 1, 3, 3, 1, 2,
}
var shChk = []int{

	-1000, -3, 5, -2, -1, 4, 5, 6, 7, 4,
	4, 4,
}
var shDef = []int{

	0, -2, 1, 0, 3, 6, 2, 0, 0, 7,
	4, 5,
}
var shTok1 = []int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	5, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	7, 3, 6,
}
var shTok2 = []int{

	2, 3, 4,
}
var shTok3 = []int{
	0,
}

//line yaccpar:1

/*	parser for yacc output	*/

var shDebug = 0

type shLexer interface {
	Lex(lval *shSymType) int
	Error(s string)
}

type shParser interface {
	Parse(shLexer) int
	Lookahead() int
}

type shParserImpl struct {
	lookahead func() int
}

func (p *shParserImpl) Lookahead() int {
	return p.lookahead()
}

func shNewParser() shParser {
	p := &shParserImpl{
		lookahead: func() int { return -1 },
	}
	return p
}

const shFlag = -1000

func shTokname(c int) string {
	// 4 is TOKSTART above
	if c >= 4 && c-4 < len(shToknames) {
		if shToknames[c-4] != "" {
			return shToknames[c-4]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func shStatname(s int) string {
	if s >= 0 && s < len(shStatenames) {
		if shStatenames[s] != "" {
			return shStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func shlex1(lex shLexer, lval *shSymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = shTok1[0]
		goto out
	}
	if char < len(shTok1) {
		token = shTok1[char]
		goto out
	}
	if char >= shPrivate {
		if char < shPrivate+len(shTok2) {
			token = shTok2[char-shPrivate]
			goto out
		}
	}
	for i := 0; i < len(shTok3); i += 2 {
		token = shTok3[i+0]
		if token == char {
			token = shTok3[i+1]
			goto out
		}
	}

out:
	if token == 0 {
		token = shTok2[1] /* unknown char */
	}
	if shDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", shTokname(token), uint(char))
	}
	return char, token
}

func shParse(shlex shLexer) int {
	return shNewParser().Parse(shlex)
}

func (shrcvr *shParserImpl) Parse(shlex shLexer) int {
	var shn int
	var shlval shSymType
	var shVAL shSymType
	var shDollar []shSymType
	shS := make([]shSymType, shMaxDepth)

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	shstate := 0
	shchar := -1
	shtoken := -1 // shchar translated into internal numbering
	shrcvr.lookahead = func() int { return shchar }
	defer func() {
		// Make sure we report no lookahead when not parsing.
		shchar = -1
		shtoken = -1
	}()
	shp := -1
	goto shstack

ret0:
	return 0

ret1:
	return 1

shstack:
	/* put a state and value onto the stack */
	if shDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", shTokname(shtoken), shStatname(shstate))
	}

	shp++
	if shp >= len(shS) {
		nyys := make([]shSymType, len(shS)*2)
		copy(nyys, shS)
		shS = nyys
	}
	shS[shp] = shVAL
	shS[shp].yys = shstate

shnewstate:
	shn = shPact[shstate]
	if shn <= shFlag {
		goto shdefault /* simple state */
	}
	if shchar < 0 {
		shchar, shtoken = shlex1(shlex, &shlval)
	}
	shn += shtoken
	if shn < 0 || shn >= shLast {
		goto shdefault
	}
	shn = shAct[shn]
	if shChk[shn] == shtoken { /* valid shift */
		shchar = -1
		shtoken = -1
		shVAL = shlval
		shstate = shn
		if Errflag > 0 {
			Errflag--
		}
		goto shstack
	}

shdefault:
	/* default state action */
	shn = shDef[shstate]
	if shn == -2 {
		if shchar < 0 {
			shchar, shtoken = shlex1(shlex, &shlval)
		}

		/* look through exception table */
		xi := 0
		for {
			if shExca[xi+0] == -1 && shExca[xi+1] == shstate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			shn = shExca[xi+0]
			if shn < 0 || shn == shtoken {
				break
			}
		}
		shn = shExca[xi+1]
		if shn < 0 {
			goto ret0
		}
	}
	if shn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			shlex.Error("syntax error")
			Nerrs++
			if shDebug >= 1 {
				__yyfmt__.Printf("%s", shStatname(shstate))
				__yyfmt__.Printf(" saw %s\n", shTokname(shtoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for shp >= 0 {
				shn = shPact[shS[shp].yys] + shErrCode
				if shn >= 0 && shn < shLast {
					shstate = shAct[shn] /* simulate a shift of "error" */
					if shChk[shstate] == shErrCode {
						goto shstack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if shDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", shS[shp].yys)
				}
				shp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if shDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", shTokname(shtoken))
			}
			if shtoken == shEofCode {
				goto ret1
			}
			shchar = -1
			shtoken = -1
			goto shnewstate /* try again in the same state */
		}
	}

	/* reduction by production shn */
	if shDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", shn, shStatname(shstate))
	}

	shnt := shn
	shpt := shp
	_ = shpt // guard against "declared and not used"

	shp -= shR2[shn]
	// shp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if shp+1 >= len(shS) {
		nyys := make([]shSymType, len(shS)*2)
		copy(nyys, shS)
		shS = nyys
	}
	shVAL = shS[shp+1]

	/* consult goto table to find next state */
	shn = shR1[shn]
	shg := shPgo[shn]
	shj := shg + shS[shp].yys + 1

	if shj >= shLast {
		shstate = shAct[shg]
	} else {
		shstate = shAct[shj]
		if shChk[shstate] != -shn {
			shstate = shAct[shg]
		}
	}
	// dummy call; replaced with literal code
	switch shnt {

	case 2:
		shDollar = shS[shpt-2 : shpt+1]
		//line sh.y:21
		{
			run(shDollar[1].cmd)
		}
	case 3:
		shDollar = shS[shpt-1 : shpt+1]
		//line sh.y:23
		{
			shVAL.cmd.args = shDollar[1].words
		}
	case 4:
		shDollar = shS[shpt-3 : shpt+1]
		//line sh.y:24
		{
			shVAL.cmd.stdout = shDollar[3].word
		}
	case 5:
		shDollar = shS[shpt-3 : shpt+1]
		//line sh.y:25
		{
			shVAL.cmd.stdin = shDollar[3].word
		}
	case 6:
		shDollar = shS[shpt-1 : shpt+1]
		//line sh.y:27
		{
			shVAL.words = []string{shDollar[1].word}
		}
	case 7:
		shDollar = shS[shpt-2 : shpt+1]
		//line sh.y:28
		{
			shVAL.words = append(shDollar[1].words, shDollar[2].word)
		}
	}
	goto shstack /* stack new state and value */
}
