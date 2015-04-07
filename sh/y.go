//line parse.y:2
package main

import __yyfmt__ "fmt"

//line parse.y:3
import "os/exec"

//line parse.y:9
type shSymType struct {
	yys   int
	word  string
	words []string
	asgn  struct{}
	cmd   *exec.Cmd
	pipe  []*exec.Cmd
	line  [][]*exec.Cmd
}

const WORD = 57346
const IF = 57347
const FOR = 57348
const SWITCH = 57349
const APPEND = 57350

var shToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"WORD",
	"IF",
	"FOR",
	"SWITCH",
	"'|'",
	"'^'",
	"'<'",
	"'>'",
	"APPEND",
	"'\n'",
	"';'",
	"'='",
}
var shStatenames = [...]string{}

const shEofCode = 1
const shErrCode = 2
const shMaxDepth = 200

//line yacctab:1
var shExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
}

const shNprod = 21
const shPrivate = 57344

var shTokenNames []string
var shStates []string

const shLast = 48

var shAct = [...]int{

	7, 8, 3, 19, 9, 19, 17, 16, 9, 18,
	6, 23, 9, 14, 15, 25, 1, 24, 28, 29,
	13, 31, 32, 33, 26, 11, 12, 9, 17, 16,
	20, 21, 22, 5, 30, 4, 2, 10, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 27,
}
var shPact = [...]int{

	23, -1000, -1000, -1000, 12, 0, -1000, -6, 20, -1000,
	8, -1000, 4, 8, -1000, 8, 20, -4, 8, 30,
	8, 8, 8, -4, -1000, -1000, -1000, 8, -1000, -4,
	-1000, -4, -4, -4,
}
var shPgo = [...]int{

	0, 0, 37, 33, 1, 10, 35, 2, 16,
}
var shR1 = [...]int{

	0, 8, 8, 7, 7, 7, 7, 7, 6, 6,
	5, 5, 4, 4, 4, 4, 3, 2, 2, 1,
	1,
}
var shR2 = [...]int{

	0, 1, 1, 2, 3, 3, 2, 3, 1, 3,
	1, 2, 1, 3, 3, 3, 3, 1, 2, 1,
	3,
}
var shChk = [...]int{

	-1000, -8, 13, -7, -6, -3, -5, -1, -4, 4,
	-2, 13, 14, 8, 13, 14, -4, -1, 15, 9,
	10, 11, 12, -1, 13, -7, -5, -3, -7, -1,
	4, -1, -1, -1,
}
var shDef = [...]int{

	0, -2, 1, 2, 0, 0, 8, 17, 10, 19,
	12, 3, 0, 0, 6, 0, 11, 17, 0, 0,
	0, 0, 0, 18, 4, 5, 9, 0, 7, 16,
	20, 13, 14, 15,
}
var shTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	13, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 14,
	10, 15, 11, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 9, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 8,
}
var shTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 12,
}
var shTok3 = [...]int{
	0,
}

//line yaccpar:1

/*	parser for yacc output	*/

var (
	shDebug        = 0
	shErrorVerbose = false
)

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
	state     func() int
}

func (p *shParserImpl) Lookahead() int {
	return p.lookahead()
}

func shNewParser() shParser {
	p := &shParserImpl{
		lookahead: func() int { return -1 },
		state:     func() int { return -1 },
	}
	return p
}

const shFlag = -1000

func shTokname(c int) string {
	if c >= 1 && c-1 < len(shToknames) {
		if shToknames[c-1] != "" {
			return shToknames[c-1]
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

func shErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !shErrorVerbose {
		return "syntax error"
	}
	res := "syntax error: unexpected " + shTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := shPact[state]
	for tok := TOKSTART; tok-1 < len(shToknames); tok++ {
		if n := base + tok; n >= 0 && n < shLast && shChk[shAct[n]] == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if shDef[state] == -2 {
		i := 0
		for shExca[i] != -1 || shExca[i+1] != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; shExca[i] >= 0; i += 2 {
			tok := shExca[i]
			if tok < TOKSTART || shExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if shExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += shTokname(tok)
	}
	return res
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
	shrcvr.state = func() int { return shstate }
	shrcvr.lookahead = func() int { return shchar }
	defer func() {
		// Make sure we report no lookahead when not parsing.
		shstate = -1
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
			shlex.Error(shErrorMessage(shstate, shtoken))
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
		shDollar = shS[shpt-1 : shpt+1]
		//line parse.y:36
		{
			runLine(shDollar[1].line)
		}
	case 3:
		shDollar = shS[shpt-2 : shpt+1]
		//line parse.y:38
		{
			shVAL.line = [][]*exec.Cmd{shDollar[1].pipe}
		}
	case 4:
		shDollar = shS[shpt-3 : shpt+1]
		//line parse.y:39
		{
			shVAL.line = [][]*exec.Cmd{shDollar[1].pipe}
		}
	case 5:
		shDollar = shS[shpt-3 : shpt+1]
		//line parse.y:40
		{
			shVAL.line = append(shDollar[3].line, shDollar[1].pipe)
		}
	case 6:
		shDollar = shS[shpt-2 : shpt+1]
		//line parse.y:41
		{
			updateEnv()
		}
	case 7:
		shDollar = shS[shpt-3 : shpt+1]
		//line parse.y:42
		{
			updateEnv()
			shVAL.line = shDollar[3].line
		}
	case 8:
		shDollar = shS[shpt-1 : shpt+1]
		//line parse.y:44
		{
			shVAL.pipe = []*exec.Cmd{shDollar[1].cmd}
		}
	case 9:
		shDollar = shS[shpt-3 : shpt+1]
		//line parse.y:45
		{
			pconnect(shDollar[1].pipe[len(shDollar[1].pipe)-1], shDollar[3].cmd)
			shVAL.pipe = append(shDollar[1].pipe, shDollar[3].cmd)
		}
	case 11:
		shDollar = shS[shpt-2 : shpt+1]
		//line parse.y:48
		{
			shVAL.cmd = shDollar[2].cmd
		}
	case 12:
		shDollar = shS[shpt-1 : shpt+1]
		//line parse.y:50
		{
			shVAL.cmd = &exec.Cmd{Path: shDollar[1].words[0], Args: shDollar[1].words}
		}
	case 13:
		shDollar = shS[shpt-3 : shpt+1]
		//line parse.y:51
		{
			shVAL.cmd.Stdin = fopen(shDollar[3].word, 'r')
			defer fclose(shVAL.cmd.Stdin)
		}
	case 14:
		shDollar = shS[shpt-3 : shpt+1]
		//line parse.y:52
		{
			shVAL.cmd.Stdout = fopen(shDollar[3].word, 'w')
			defer fclose(shVAL.cmd.Stdout)
		}
	case 15:
		shDollar = shS[shpt-3 : shpt+1]
		//line parse.y:53
		{
			shVAL.cmd.Stdout = fopen(shDollar[3].word, 'a')
			defer fclose(shVAL.cmd.Stdout)
		}
	case 16:
		shDollar = shS[shpt-3 : shpt+1]
		//line parse.y:55
		{
			env[shDollar[1].word] = shDollar[3].word
			shVAL.asgn = struct{}{}
		}
	case 17:
		shDollar = shS[shpt-1 : shpt+1]
		//line parse.y:57
		{
			shVAL.words = []string{shDollar[1].word}
		}
	case 18:
		shDollar = shS[shpt-2 : shpt+1]
		//line parse.y:58
		{
			shVAL.words = append(shDollar[1].words, shDollar[2].word)
		}
	case 20:
		shDollar = shS[shpt-3 : shpt+1]
		//line parse.y:61
		{
			shVAL.word = shDollar[1].word + shDollar[3].word
		}
	}
	goto shstack /* stack new state and value */
}
