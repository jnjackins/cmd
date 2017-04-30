//line syntax.y:1

// generated by goyacc

package main

import __yyfmt__ "fmt"

//line syntax.y:4
//line syntax.y:11
type shSymType struct {
	yys  int
	tree *treeNode
}

const IF = 57346
const THEN = 57347
const FI = 57348
const FOR = 57349
const IN = 57350
const DO = 57351
const DONE = 57352
const WORD = 57353
const QUOTE = 57354
const REDIR = 57355
const SIMPLE = 57356
const WORDS = 57357
const PAREN = 57358
const AND = 57359
const OR = 57360

var shToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"IF",
	"THEN",
	"FI",
	"FOR",
	"IN",
	"DO",
	"DONE",
	"WORD",
	"QUOTE",
	"REDIR",
	"SIMPLE",
	"WORDS",
	"PAREN",
	"AND",
	"OR",
	"';'",
	"'&'",
	"'|'",
	"'('",
	"')'",
}
var shStatenames = [...]string{}

const shEofCode = 1
const shErrCode = 2
const shInitialStackSize = 16

//line yacctab:1
var shExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
	-1, 34,
	5, 3,
	6, 3,
	10, 3,
	-2, 8,
}

const shPrivate = 57344

const shLast = 58

var shAct = [...]int{

	24, 12, 8, 32, 10, 20, 21, 11, 6, 31,
	22, 13, 14, 17, 18, 34, 16, 5, 17, 18,
	15, 16, 9, 13, 14, 13, 14, 35, 27, 30,
	43, 40, 41, 36, 37, 28, 29, 39, 33, 38,
	22, 1, 42, 25, 3, 26, 4, 2, 3, 7,
	4, 0, 19, 3, 0, 4, 0, 23,
}
var shPact = [...]int{

	0, -1000, -1000, 1, 0, -16, -7, -1000, 14, 0,
	0, 17, -1000, -1000, -1000, -1000, -1000, 0, 0, -1000,
	0, -2, -1000, -20, 33, -4, 0, 25, -16, -16,
	-7, -1000, -1000, 0, -1000, -1000, 14, 31, 12, -1000,
	23, 0, 20, -1000,
}
var shPgo = [...]int{

	0, 0, 47, 45, 43, 17, 8, 49, 2, 1,
	41,
}
var shR1 = [...]int{

	0, 10, 10, 1, 1, 2, 2, 2, 3, 3,
	4, 4, 4, 5, 5, 6, 6, 7, 7, 7,
	7, 8, 8, 9, 9,
}
var shR2 = [...]int{

	0, 0, 1, 2, 2, 1, 1, 2, 2, 2,
	1, 3, 3, 1, 3, 1, 3, 1, 3, 5,
	8, 1, 2, 1, 1,
}
var shChk = [...]int{

	-1000, -10, -2, -4, -3, -5, -6, -7, -8, 22,
	4, 7, -9, 11, 12, 19, 20, 17, 18, -2,
	21, 13, -9, -2, -1, -4, -3, 11, -5, -5,
	-6, 11, 23, 5, 19, -1, 8, -1, -8, 6,
	19, 9, -1, 10,
}
var shDef = [...]int{

	1, -2, 2, 5, 6, 10, 13, 15, 17, 0,
	0, 0, 21, 23, 24, 8, 9, 0, 0, 7,
	0, 0, 22, 0, 0, 0, 0, 0, 11, 12,
	14, 16, 18, 0, -2, 4, 0, 0, 0, 19,
	0, 0, 0, 20,
}
var shTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 20, 3,
	22, 23, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 19,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 21,
}
var shTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18,
}
var shTok3 = [...]int{
	0,
}

var shErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}

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
	lval  shSymType
	stack [shInitialStackSize]shSymType
	char  int
}

func (p *shParserImpl) Lookahead() int {
	return p.char
}

func shNewParser() shParser {
	return &shParserImpl{}
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

	for _, e := range shErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
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
	var shVAL shSymType
	var shDollar []shSymType
	_ = shDollar // silence set and not used
	shS := shrcvr.stack[:]

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	shstate := 0
	shrcvr.char = -1
	shtoken := -1 // shrcvr.char translated into internal numbering
	defer func() {
		// Make sure we report no lookahead when not parsing.
		shstate = -1
		shrcvr.char = -1
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
	if shrcvr.char < 0 {
		shrcvr.char, shtoken = shlex1(shlex, &shrcvr.lval)
	}
	shn += shtoken
	if shn < 0 || shn >= shLast {
		goto shdefault
	}
	shn = shAct[shn]
	if shChk[shn] == shtoken { /* valid shift */
		shrcvr.char = -1
		shtoken = -1
		shVAL = shrcvr.lval
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
		if shrcvr.char < 0 {
			shrcvr.char, shtoken = shlex1(shlex, &shrcvr.lval)
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
			shrcvr.char = -1
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
		//line syntax.y:21
		{
			execute(shDollar[1].tree)
		}
	case 4:
		shDollar = shS[shpt-2 : shpt+1]
		//line syntax.y:26
		{
			shVAL.tree = mkTree(';', shDollar[1].tree, shDollar[2].tree)
		}
	case 7:
		shDollar = shS[shpt-2 : shpt+1]
		//line syntax.y:30
		{
			shVAL.tree = mkTree(';', shDollar[1].tree, shDollar[2].tree)
		}
	case 9:
		shDollar = shS[shpt-2 : shpt+1]
		//line syntax.y:33
		{
			shVAL.tree = mkTree('&', shDollar[1].tree)
		}
	case 11:
		shDollar = shS[shpt-3 : shpt+1]
		//line syntax.y:36
		{
			shVAL.tree = mkTree(AND, shDollar[1].tree, shDollar[3].tree)
		}
	case 12:
		shDollar = shS[shpt-3 : shpt+1]
		//line syntax.y:37
		{
			shVAL.tree = mkTree(OR, shDollar[1].tree, shDollar[3].tree)
		}
	case 14:
		shDollar = shS[shpt-3 : shpt+1]
		//line syntax.y:40
		{
			shVAL.tree = mkTree('|', shDollar[1].tree, shDollar[3].tree)
		}
	case 16:
		shDollar = shS[shpt-3 : shpt+1]
		//line syntax.y:43
		{
			shVAL.tree = shDollar[1].tree
			shDollar[1].tree.redirect(shDollar[2].tree.int, shDollar[3].tree.string)
		}
	case 17:
		shDollar = shS[shpt-1 : shpt+1]
		//line syntax.y:45
		{
			shVAL.tree = mkTree(SIMPLE, shDollar[1].tree)
		}
	case 18:
		shDollar = shS[shpt-3 : shpt+1]
		//line syntax.y:46
		{
			shVAL.tree = mkTree(PAREN, shDollar[2].tree)
		}
	case 19:
		shDollar = shS[shpt-5 : shpt+1]
		//line syntax.y:47
		{
			shVAL.tree = mkTree(IF, shDollar[2].tree, shDollar[4].tree)
		}
	case 20:
		shDollar = shS[shpt-8 : shpt+1]
		//line syntax.y:48
		{
			shVAL.tree = mkTree(FOR, shDollar[2].tree, shDollar[4].tree, shDollar[7].tree)
		}
	case 21:
		shDollar = shS[shpt-1 : shpt+1]
		//line syntax.y:50
		{
			shVAL.tree = mkTree(WORDS, shDollar[1].tree)
		}
	case 22:
		shDollar = shS[shpt-2 : shpt+1]
		//line syntax.y:51
		{
			shVAL.tree = shDollar[1].tree
			shDollar[1].tree.children = append(shDollar[1].tree.children, shDollar[2].tree)
		}
	}
	goto shstack /* stack new state and value */
}
