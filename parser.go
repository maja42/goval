//line parser.go.y:2
package main

import __yyfmt__ "fmt"

//line parser.go.y:2
//line parser.go.y:6
type yySymType struct {
	yys      int
	token    Token
	expr     interface{}
	exprList []interface{}
	exprMap  map[string]interface{}
}

const LITERAL_BOOL = 57346
const LITERAL_NUMBER = 57347
const LITERAL_STRING = 57348
const IDENT = 57349
const AND = 57350
const OR = 57351
const EQL = 57352
const NEQ = 57353
const LSS = 57354
const GTR = 57355
const LEQ = 57356
const GEQ = 57357

var yyToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"LITERAL_BOOL",
	"LITERAL_NUMBER",
	"LITERAL_STRING",
	"IDENT",
	"AND",
	"OR",
	"EQL",
	"NEQ",
	"LSS",
	"GTR",
	"LEQ",
	"GEQ",
	"'|'",
	"'^'",
	"'&'",
	"'+'",
	"'-'",
	"'*'",
	"'/'",
	"'%'",
	"'!'",
	"'.'",
	"'['",
	"']'",
	"'('",
	"')'",
	"'{'",
	"'}'",
	"','",
	"':'",
}
var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16

//line parser.go.y:127

//line yacctab:1
var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyPrivate = 57344

const yyLast = 327

var yyAct = [...]int{

	39, 2, 72, 68, 69, 67, 66, 61, 36, 35,
	41, 67, 33, 34, 7, 42, 43, 44, 45, 46,
	47, 48, 49, 50, 51, 52, 53, 54, 55, 56,
	57, 58, 59, 60, 6, 62, 28, 29, 22, 23,
	24, 25, 26, 27, 30, 32, 31, 17, 18, 19,
	20, 21, 5, 33, 34, 4, 3, 38, 19, 20,
	21, 76, 33, 34, 1, 0, 0, 0, 73, 0,
	74, 75, 0, 0, 0, 0, 0, 77, 28, 29,
	22, 23, 24, 25, 26, 27, 30, 32, 31, 17,
	18, 19, 20, 21, 65, 33, 34, 0, 0, 0,
	0, 0, 0, 70, 28, 29, 22, 23, 24, 25,
	26, 27, 30, 32, 31, 17, 18, 19, 20, 21,
	0, 33, 34, 0, 0, 63, 28, 29, 22, 23,
	24, 25, 26, 27, 30, 32, 31, 17, 18, 19,
	20, 21, 0, 33, 34, 71, 28, 29, 22, 23,
	24, 25, 26, 27, 30, 32, 31, 17, 18, 19,
	20, 21, 0, 33, 34, 28, 0, 22, 23, 24,
	25, 26, 27, 30, 32, 31, 17, 18, 19, 20,
	21, 0, 33, 34, 22, 23, 24, 25, 26, 27,
	30, 32, 31, 17, 18, 19, 20, 21, 0, 33,
	34, 22, 23, 24, 25, 26, 27, 0, 32, 31,
	17, 18, 19, 20, 21, 0, 33, 34, 10, 11,
	12, 9, 24, 25, 26, 27, 0, 0, 0, 17,
	18, 19, 20, 21, 15, 33, 34, 0, 16, 0,
	13, 0, 8, 0, 14, 40, 22, 23, 24, 25,
	26, 27, 0, 0, 31, 17, 18, 19, 20, 21,
	0, 33, 34, 10, 11, 12, 9, 0, 0, 0,
	10, 11, 12, 9, 0, 10, 11, 12, 9, 15,
	0, 0, 0, 16, 0, 13, 15, 8, 64, 14,
	16, 15, 13, 37, 8, 16, 14, 13, 0, 8,
	0, 14, 22, 23, 24, 25, 26, 27, 0, 0,
	0, 17, 18, 19, 20, 21, 0, 33, 34, 17,
	18, 19, 20, 21, 0, 33, 34,
}
var yyPact = [...]int{

	271, -1000, 138, -1000, -1000, -1000, -1000, -1000, 271, -20,
	-1000, -1000, -1000, 266, 214, 271, 271, 271, 271, 271,
	271, 271, 271, 271, 271, 271, 271, 271, 271, 271,
	271, 271, 271, 0, 271, 96, 259, -1000, -21, 138,
	-1000, -28, 70, -13, -13, 37, 37, -13, -13, -13,
	210, 210, 300, 300, 300, 300, 174, 157, 191, 292,
	236, -1000, 118, -1000, -1000, -27, -1000, 271, -1000, 271,
	271, -1000, -1000, 138, 28, 138, 271, 138,
}
var yyPgo = [...]int{

	0, 64, 0, 56, 55, 52, 34, 14, 57, 10,
}
var yyR1 = [...]int{

	0, 1, 2, 2, 2, 2, 2, 2, 2, 2,
	3, 3, 3, 3, 3, 3, 3, 4, 4, 4,
	4, 4, 4, 5, 5, 5, 5, 5, 5, 5,
	5, 5, 6, 6, 6, 7, 7, 7, 8, 8,
	9, 9,
}
var yyR2 = [...]int{

	0, 1, 1, 1, 1, 1, 1, 3, 3, 4,
	1, 1, 1, 2, 3, 2, 3, 2, 3, 3,
	3, 3, 3, 2, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 1, 3, 4, 1, 3,
	3, 5,
}
var yyChk = [...]int{

	-1000, -1, -2, -3, -4, -5, -6, -7, 28, 7,
	4, 5, 6, 26, 30, 20, 24, 19, 20, 21,
	22, 23, 10, 11, 12, 13, 14, 15, 8, 9,
	16, 18, 17, 25, 26, -2, 28, 27, -8, -2,
	31, -9, -2, -2, -2, -2, -2, -2, -2, -2,
	-2, -2, -2, -2, -2, -2, -2, -2, -2, -2,
	-2, 7, -2, 29, 29, -8, 27, 32, 31, 32,
	33, 27, 29, -2, -2, -2, 33, -2,
}
var yyDef = [...]int{

	0, -2, 1, 2, 3, 4, 5, 6, 0, 35,
	10, 11, 12, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 13, 0, 38,
	15, 0, 0, 17, 23, 18, 19, 20, 21, 22,
	24, 25, 26, 27, 28, 29, 30, 31, 32, 33,
	34, 36, 0, 7, 8, 0, 14, 0, 16, 0,
	0, 37, 9, 39, 0, 40, 0, 41,
}
var yyTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 24, 3, 3, 3, 23, 18, 3,
	28, 29, 21, 19, 32, 20, 25, 22, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 33, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 26, 3, 27, 17, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 30, 16, 31,
}
var yyTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15,
}
var yyTok3 = [...]int{
	0,
}

var yyErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}

//line yaccpar:1

/*	parser for yacc output	*/

var (
	yyDebug        = 0
	yyErrorVerbose = false
)

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

type yyParser interface {
	Parse(yyLexer) int
	Lookahead() int
}

type yyParserImpl struct {
	lval  yySymType
	stack [yyInitialStackSize]yySymType
	char  int
}

func (p *yyParserImpl) Lookahead() int {
	return p.char
}

func yyNewParser() yyParser {
	return &yyParserImpl{}
}

const yyFlag = -1000

func yyTokname(c int) string {
	if c >= 1 && c-1 < len(yyToknames) {
		if yyToknames[c-1] != "" {
			return yyToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func yyStatname(s int) string {
	if s >= 0 && s < len(yyStatenames) {
		if yyStatenames[s] != "" {
			return yyStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func yyErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !yyErrorVerbose {
		return "syntax error"
	}

	for _, e := range yyErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + yyTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := yyPact[state]
	for tok := TOKSTART; tok-1 < len(yyToknames); tok++ {
		if n := base + tok; n >= 0 && n < yyLast && yyChk[yyAct[n]] == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if yyDef[state] == -2 {
		i := 0
		for yyExca[i] != -1 || yyExca[i+1] != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; yyExca[i] >= 0; i += 2 {
			tok := yyExca[i]
			if tok < TOKSTART || yyExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if yyExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += yyTokname(tok)
	}
	return res
}

func yylex1(lex yyLexer, lval *yySymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = yyTok1[0]
		goto out
	}
	if char < len(yyTok1) {
		token = yyTok1[char]
		goto out
	}
	if char >= yyPrivate {
		if char < yyPrivate+len(yyTok2) {
			token = yyTok2[char-yyPrivate]
			goto out
		}
	}
	for i := 0; i < len(yyTok3); i += 2 {
		token = yyTok3[i+0]
		if token == char {
			token = yyTok3[i+1]
			goto out
		}
	}

out:
	if token == 0 {
		token = yyTok2[1] /* unknown char */
	}
	if yyDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", yyTokname(token), uint(char))
	}
	return char, token
}

func yyParse(yylex yyLexer) int {
	return yyNewParser().Parse(yylex)
}

func (yyrcvr *yyParserImpl) Parse(yylex yyLexer) int {
	var yyn int
	var yyVAL yySymType
	var yyDollar []yySymType
	_ = yyDollar // silence set and not used
	yyS := yyrcvr.stack[:]

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yyrcvr.char = -1
	yytoken := -1 // yyrcvr.char translated into internal numbering
	defer func() {
		// Make sure we report no lookahead when not parsing.
		yystate = -1
		yyrcvr.char = -1
		yytoken = -1
	}()
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	if yyDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", yyTokname(yytoken), yyStatname(yystate))
	}

	yyp++
	if yyp >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyS[yyp] = yyVAL
	yyS[yyp].yys = yystate

yynewstate:
	yyn = yyPact[yystate]
	if yyn <= yyFlag {
		goto yydefault /* simple state */
	}
	if yyrcvr.char < 0 {
		yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
	}
	yyn += yytoken
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = yyAct[yyn]
	if yyChk[yyn] == yytoken { /* valid shift */
		yyrcvr.char = -1
		yytoken = -1
		yyVAL = yyrcvr.lval
		yystate = yyn
		if Errflag > 0 {
			Errflag--
		}
		goto yystack
	}

yydefault:
	/* default state action */
	yyn = yyDef[yystate]
	if yyn == -2 {
		if yyrcvr.char < 0 {
			yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
		}

		/* look through exception table */
		xi := 0
		for {
			if yyExca[xi+0] == -1 && yyExca[xi+1] == yystate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			yyn = yyExca[xi+0]
			if yyn < 0 || yyn == yytoken {
				break
			}
		}
		yyn = yyExca[xi+1]
		if yyn < 0 {
			goto ret0
		}
	}
	if yyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			yylex.Error(yyErrorMessage(yystate, yytoken))
			Nerrs++
			if yyDebug >= 1 {
				__yyfmt__.Printf("%s", yyStatname(yystate))
				__yyfmt__.Printf(" saw %s\n", yyTokname(yytoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for yyp >= 0 {
				yyn = yyPact[yyS[yyp].yys] + yyErrCode
				if yyn >= 0 && yyn < yyLast {
					yystate = yyAct[yyn] /* simulate a shift of "error" */
					if yyChk[yystate] == yyErrCode {
						goto yystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if yyDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", yyS[yyp].yys)
				}
				yyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if yyDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", yyTokname(yytoken))
			}
			if yytoken == yyEofCode {
				goto ret1
			}
			yyrcvr.char = -1
			yytoken = -1
			goto yynewstate /* try again in the same state */
		}
	}

	/* reduction by production yyn */
	if yyDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", yyn, yyStatname(yystate))
	}

	yynt := yyn
	yypt := yyp
	_ = yypt // guard against "declared and not used"

	yyp -= yyR2[yyn]
	// yyp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if yyp+1 >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyVAL = yyS[yyp+1]

	/* consult goto table to find next state */
	yyn = yyR1[yyn]
	yyg := yyPgo[yyn]
	yyj := yyg + yyS[yyp].yys + 1

	if yyj >= yyLast {
		yystate = yyAct[yyg]
	} else {
		yystate = yyAct[yyj]
		if yyChk[yystate] != -yyn {
			yystate = yyAct[yyg]
		}
	}
	// dummy call; replaced with literal code
	switch yynt {

	case 1:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:57
		{
			yyVAL.expr = yyDollar[1].expr
			yylex.(*Lexer).result = yyVAL.expr
		}
	case 7:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:69
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 8:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:70
		{
			yyVAL.expr = callFunction(yylex.(*Lexer).functions, yyDollar[1].token.literal, []interface{}{})
		}
	case 9:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.go.y:71
		{
			yyVAL.expr = callFunction(yylex.(*Lexer).functions, yyDollar[1].token.literal, yyDollar[3].exprList)
		}
	case 10:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:75
		{
			yyVAL.expr = yyDollar[1].token.value
		}
	case 11:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:76
		{
			yyVAL.expr = yyDollar[1].token.value
		}
	case 12:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:77
		{
			yyVAL.expr = yyDollar[1].token.value
		}
	case 13:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:78
		{
			yyVAL.expr = []interface{}{}
		}
	case 14:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:79
		{
			yyVAL.expr = yyDollar[2].exprList
		}
	case 15:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:80
		{
			yyVAL.expr = map[string]interface{}{}
		}
	case 16:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:81
		{
			yyVAL.expr = yyDollar[2].exprMap
		}
	case 17:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:85
		{
			yyVAL.expr = unaryMinus(yyDollar[2].expr)
		}
	case 18:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:86
		{
			yyVAL.expr = add(yyDollar[1].expr, yyDollar[3].expr)
		}
	case 19:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:87
		{
			yyVAL.expr = sub(yyDollar[1].expr, yyDollar[3].expr)
		}
	case 20:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:88
		{
			yyVAL.expr = mul(yyDollar[1].expr, yyDollar[3].expr)
		}
	case 21:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:89
		{
			yyVAL.expr = div(yyDollar[1].expr, yyDollar[3].expr)
		}
	case 22:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:90
		{
			yyVAL.expr = mod(yyDollar[1].expr, yyDollar[3].expr)
		}
	case 23:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:94
		{
			yyVAL.expr = !asBool(yyDollar[2].expr)
		}
	case 24:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:95
		{
			yyVAL.expr = deepEqual(yyDollar[1].expr, yyDollar[3].expr)
		}
	case 25:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:96
		{
			yyVAL.expr = !deepEqual(yyDollar[1].expr, yyDollar[3].expr)
		}
	case 26:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:97
		{
			yyVAL.expr = compare(yyDollar[1].expr, yyDollar[3].expr, "<")
		}
	case 27:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:98
		{
			yyVAL.expr = compare(yyDollar[1].expr, yyDollar[3].expr, ">")
		}
	case 28:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:99
		{
			yyVAL.expr = compare(yyDollar[1].expr, yyDollar[3].expr, "<=")
		}
	case 29:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:100
		{
			yyVAL.expr = compare(yyDollar[1].expr, yyDollar[3].expr, ">=")
		}
	case 30:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:101
		{
			left := asBool(yyDollar[1].expr)
			right := asBool(yyDollar[3].expr)
			yyVAL.expr = left && right
		}
	case 31:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:102
		{
			left := asBool(yyDollar[1].expr)
			right := asBool(yyDollar[3].expr)
			yyVAL.expr = left || right
		}
	case 32:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:106
		{
			yyVAL.expr = asInteger(yyDollar[1].expr) | asInteger(yyDollar[3].expr)
		}
	case 33:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:107
		{
			yyVAL.expr = asInteger(yyDollar[1].expr) & asInteger(yyDollar[3].expr)
		}
	case 34:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:108
		{
			yyVAL.expr = asInteger(yyDollar[1].expr) ^ asInteger(yyDollar[3].expr)
		}
	case 35:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:112
		{
			yyVAL.expr = accessVar(yylex.(*Lexer).variables, yyDollar[1].token.literal)
		}
	case 36:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:113
		{
			yyVAL.expr = accessField(yyDollar[1].expr, yyDollar[3].token.literal)
		}
	case 37:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.go.y:114
		{
			yyVAL.expr = accessField(yyDollar[1].expr, yyDollar[3].expr)
		}
	case 38:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:118
		{
			yyVAL.exprList = []interface{}{yyDollar[1].expr}
		}
	case 39:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:119
		{
			yyVAL.exprList = append(yyDollar[1].exprList, yyDollar[3].expr)
		}
	case 40:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:123
		{
			yyVAL.exprMap = make(map[string]interface{})
			yyVAL.exprMap[asObjectKey(yyDollar[1].expr)] = yyDollar[3].expr
		}
	case 41:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:124
		{
			yyVAL.exprMap = addObjectMember(yyDollar[1].exprMap, yyDollar[3].expr, yyDollar[5].expr)
		}
	}
	goto yystack /* stack new state and value */
}
