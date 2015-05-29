//line action.y:2
package main

import __yyfmt__ "fmt"

//line action.y:3
//line action.y:7
type awkSymType struct {
	yys  int
	num  float64
	str  string
	fnc  func(record)
	fncs []func(record)
	args []*symbol
}

const NUM = 57346
const STR = 57347
const VAR = 57348
const FLD = 57349
const PRINT = 57350

var awkToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"NUM",
	"STR",
	"VAR",
	"FLD",
	"PRINT",
	"';'",
}
var awkStatenames = [...]string{}

const awkEofCode = 1
const awkErrCode = 2
const awkMaxDepth = 200

//line yacctab:1
var awkExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
}

const awkNprod = 10
const awkPrivate = 57344

var awkTokenNames []string
var awkStates []string

const awkLast = 12

var awkAct = [...]int{

	9, 10, 6, 11, 3, 5, 1, 7, 4, 2,
	0, 8,
}
var awkPact = [...]int{

	-3, -1000, -7, -1000, -1000, -1000, -3, -4, -1000, -1000,
	-1000, -1000,
}
var awkPgo = [...]int{

	0, 9, 4, 8, 7, 6,
}
var awkR1 = [...]int{

	0, 5, 1, 1, 2, 3, 4, 4, 4, 4,
}
var awkR2 = [...]int{

	0, 1, 1, 3, 1, 2, 0, 2, 2, 2,
}
var awkChk = [...]int{

	-1000, -5, -1, -2, -3, 8, 9, -4, -2, 4,
	5, 7,
}
var awkDef = [...]int{

	0, -2, 1, 2, 4, 6, 0, 5, 3, 7,
	8, 9,
}
var awkTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 9,
}
var awkTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8,
}
var awkTok3 = [...]int{
	0,
}

var awkErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}

//line yaccpar:1

/*	parser for yacc output	*/

var (
	awkDebug        = 0
	awkErrorVerbose = false
)

type awkLexer interface {
	Lex(lval *awkSymType) int
	Error(s string)
}

type awkParser interface {
	Parse(awkLexer) int
	Lookahead() int
}

type awkParserImpl struct {
	lookahead func() int
}

func (p *awkParserImpl) Lookahead() int {
	return p.lookahead()
}

func awkNewParser() awkParser {
	p := &awkParserImpl{
		lookahead: func() int { return -1 },
	}
	return p
}

const awkFlag = -1000

func awkTokname(c int) string {
	if c >= 1 && c-1 < len(awkToknames) {
		if awkToknames[c-1] != "" {
			return awkToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func awkStatname(s int) string {
	if s >= 0 && s < len(awkStatenames) {
		if awkStatenames[s] != "" {
			return awkStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func awkErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !awkErrorVerbose {
		return "syntax error"
	}

	for _, e := range awkErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + awkTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := awkPact[state]
	for tok := TOKSTART; tok-1 < len(awkToknames); tok++ {
		if n := base + tok; n >= 0 && n < awkLast && awkChk[awkAct[n]] == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if awkDef[state] == -2 {
		i := 0
		for awkExca[i] != -1 || awkExca[i+1] != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; awkExca[i] >= 0; i += 2 {
			tok := awkExca[i]
			if tok < TOKSTART || awkExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if awkExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += awkTokname(tok)
	}
	return res
}

func awklex1(lex awkLexer, lval *awkSymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = awkTok1[0]
		goto out
	}
	if char < len(awkTok1) {
		token = awkTok1[char]
		goto out
	}
	if char >= awkPrivate {
		if char < awkPrivate+len(awkTok2) {
			token = awkTok2[char-awkPrivate]
			goto out
		}
	}
	for i := 0; i < len(awkTok3); i += 2 {
		token = awkTok3[i+0]
		if token == char {
			token = awkTok3[i+1]
			goto out
		}
	}

out:
	if token == 0 {
		token = awkTok2[1] /* unknown char */
	}
	if awkDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", awkTokname(token), uint(char))
	}
	return char, token
}

func awkParse(awklex awkLexer) int {
	return awkNewParser().Parse(awklex)
}

func (awkrcvr *awkParserImpl) Parse(awklex awkLexer) int {
	var awkn int
	var awklval awkSymType
	var awkVAL awkSymType
	var awkDollar []awkSymType
	awkS := make([]awkSymType, awkMaxDepth)

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	awkstate := 0
	awkchar := -1
	awktoken := -1 // awkchar translated into internal numbering
	awkrcvr.lookahead = func() int { return awkchar }
	defer func() {
		// Make sure we report no lookahead when not parsing.
		awkstate = -1
		awkchar = -1
		awktoken = -1
	}()
	awkp := -1
	goto awkstack

ret0:
	return 0

ret1:
	return 1

awkstack:
	/* put a state and value onto the stack */
	if awkDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", awkTokname(awktoken), awkStatname(awkstate))
	}

	awkp++
	if awkp >= len(awkS) {
		nyys := make([]awkSymType, len(awkS)*2)
		copy(nyys, awkS)
		awkS = nyys
	}
	awkS[awkp] = awkVAL
	awkS[awkp].yys = awkstate

awknewstate:
	awkn = awkPact[awkstate]
	if awkn <= awkFlag {
		goto awkdefault /* simple state */
	}
	if awkchar < 0 {
		awkchar, awktoken = awklex1(awklex, &awklval)
	}
	awkn += awktoken
	if awkn < 0 || awkn >= awkLast {
		goto awkdefault
	}
	awkn = awkAct[awkn]
	if awkChk[awkn] == awktoken { /* valid shift */
		awkchar = -1
		awktoken = -1
		awkVAL = awklval
		awkstate = awkn
		if Errflag > 0 {
			Errflag--
		}
		goto awkstack
	}

awkdefault:
	/* default state action */
	awkn = awkDef[awkstate]
	if awkn == -2 {
		if awkchar < 0 {
			awkchar, awktoken = awklex1(awklex, &awklval)
		}

		/* look through exception table */
		xi := 0
		for {
			if awkExca[xi+0] == -1 && awkExca[xi+1] == awkstate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			awkn = awkExca[xi+0]
			if awkn < 0 || awkn == awktoken {
				break
			}
		}
		awkn = awkExca[xi+1]
		if awkn < 0 {
			goto ret0
		}
	}
	if awkn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			awklex.Error(awkErrorMessage(awkstate, awktoken))
			Nerrs++
			if awkDebug >= 1 {
				__yyfmt__.Printf("%s", awkStatname(awkstate))
				__yyfmt__.Printf(" saw %s\n", awkTokname(awktoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for awkp >= 0 {
				awkn = awkPact[awkS[awkp].yys] + awkErrCode
				if awkn >= 0 && awkn < awkLast {
					awkstate = awkAct[awkn] /* simulate a shift of "error" */
					if awkChk[awkstate] == awkErrCode {
						goto awkstack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if awkDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", awkS[awkp].yys)
				}
				awkp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if awkDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", awkTokname(awktoken))
			}
			if awktoken == awkEofCode {
				goto ret1
			}
			awkchar = -1
			awktoken = -1
			goto awknewstate /* try again in the same state */
		}
	}

	/* reduction by production awkn */
	if awkDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", awkn, awkStatname(awkstate))
	}

	awknt := awkn
	awkpt := awkp
	_ = awkpt // guard against "declared and not used"

	awkp -= awkR2[awkn]
	// awkp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if awkp+1 >= len(awkS) {
		nyys := make([]awkSymType, len(awkS)*2)
		copy(nyys, awkS)
		awkS = nyys
	}
	awkVAL = awkS[awkp+1]

	/* consult goto table to find next state */
	awkn = awkR1[awkn]
	awkg := awkPgo[awkn]
	awkj := awkg + awkS[awkp].yys + 1

	if awkj >= awkLast {
		awkstate = awkAct[awkg]
	} else {
		awkstate = awkAct[awkj]
		if awkChk[awkstate] != -awkn {
			awkstate = awkAct[awkg]
		}
	}
	// dummy call; replaced with literal code
	switch awknt {

	case 1:
		awkDollar = awkS[awkpt-1 : awkpt+1]
		//line action.y:26
		{
			parsedFns = awkDollar[1].fncs
		}
	case 2:
		awkDollar = awkS[awkpt-1 : awkpt+1]
		//line action.y:28
		{
			awkVAL.fncs = []func(record){awkDollar[1].fnc}
		}
	case 3:
		awkDollar = awkS[awkpt-3 : awkpt+1]
		//line action.y:29
		{
			awkVAL.fncs = append(awkDollar[1].fncs, awkDollar[3].fnc)
		}
	case 5:
		awkDollar = awkS[awkpt-2 : awkpt+1]
		//line action.y:33
		{
			awkVAL.fnc = printFn(awkDollar[2].args)
		}
	case 6:
		awkDollar = awkS[awkpt-0 : awkpt+1]
		//line action.y:35
		{
			awkVAL.args = make([]*symbol, 0, 1)
		}
	case 7:
		awkDollar = awkS[awkpt-2 : awkpt+1]
		//line action.y:36
		{
			awkVAL.args = append(awkDollar[1].args, &symbol{fval: awkDollar[2].num, tval: numVal})
		}
	case 8:
		awkDollar = awkS[awkpt-2 : awkpt+1]
		//line action.y:37
		{
			awkVAL.args = append(awkDollar[1].args, &symbol{sval: awkDollar[2].str, tval: strVal})
		}
	case 9:
		awkDollar = awkS[awkpt-2 : awkpt+1]
		//line action.y:38
		{
			awkVAL.args = append(awkDollar[1].args, getSymbol(awkDollar[2].str))
		}
	}
	goto awkstack /* stack new state and value */
}
