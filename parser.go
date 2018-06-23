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

const LITERAL_NIL = 57346
const LITERAL_BOOL = 57347
const LITERAL_NUMBER = 57348
const LITERAL_STRING = 57349
const IDENT = 57350
const AND = 57351
const OR = 57352
const EQL = 57353
const NEQ = 57354
const LSS = 57355
const GTR = 57356
const LEQ = 57357
const GEQ = 57358
const SHL = 57359
const SHR = 57360
const BIT_NOT = 57361

var yyToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"LITERAL_NIL",
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
	"SHL",
	"SHR",
	"BIT_NOT",
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

//line parser.go.y:136

//line yacctab:1
var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyPrivate = 57344

const yyLast = 419

var yyAct = [...]int{

	43, 2, 79, 75, 76, 74, 73, 68, 40, 39,
	45, 74, 37, 38, 7, 6, 46, 47, 48, 49,
	50, 51, 52, 53, 54, 55, 56, 57, 58, 59,
	60, 61, 62, 63, 64, 65, 66, 67, 42, 69,
	30, 31, 24, 25, 26, 27, 28, 29, 35, 36,
	5, 32, 34, 33, 19, 20, 21, 22, 23, 4,
	37, 38, 3, 19, 20, 21, 22, 23, 83, 37,
	38, 1, 0, 0, 0, 80, 0, 81, 82, 72,
	0, 0, 0, 0, 84, 30, 31, 24, 25, 26,
	27, 28, 29, 35, 36, 0, 32, 34, 33, 19,
	20, 21, 22, 23, 0, 37, 38, 21, 22, 23,
	0, 37, 38, 77, 30, 31, 24, 25, 26, 27,
	28, 29, 35, 36, 0, 32, 34, 33, 19, 20,
	21, 22, 23, 0, 37, 38, 0, 0, 70, 30,
	31, 24, 25, 26, 27, 28, 29, 35, 36, 0,
	32, 34, 33, 19, 20, 21, 22, 23, 0, 37,
	38, 78, 30, 31, 24, 25, 26, 27, 28, 29,
	35, 36, 0, 32, 34, 33, 19, 20, 21, 22,
	23, 0, 37, 38, 30, 0, 24, 25, 26, 27,
	28, 29, 35, 36, 0, 32, 34, 33, 19, 20,
	21, 22, 23, 0, 37, 38, 24, 25, 26, 27,
	28, 29, 35, 36, 0, 32, 34, 33, 19, 20,
	21, 22, 23, 0, 37, 38, 24, 25, 26, 27,
	28, 29, 35, 36, 0, 0, 34, 33, 19, 20,
	21, 22, 23, 0, 37, 38, 10, 11, 12, 13,
	9, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 18, 0, 0, 0, 0, 16, 0, 0, 0,
	17, 0, 14, 0, 8, 0, 15, 44, 24, 25,
	26, 27, 28, 29, 35, 36, 0, 0, 0, 33,
	19, 20, 21, 22, 23, 0, 37, 38, 10, 11,
	12, 13, 9, 0, 0, 10, 11, 12, 13, 9,
	0, 0, 0, 18, 0, 0, 0, 0, 16, 0,
	18, 0, 17, 0, 14, 16, 8, 71, 15, 17,
	0, 14, 41, 8, 0, 15, 10, 11, 12, 13,
	9, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 18, 0, 0, 0, 0, 16, 0, 0, 0,
	17, 0, 14, 0, 8, 0, 15, 24, 25, 26,
	27, 28, 29, 35, 36, 0, 0, 0, 0, 19,
	20, 21, 22, 23, 0, 37, 38, 26, 27, 28,
	29, 35, 36, 0, 0, 0, 0, 19, 20, 21,
	22, 23, 0, 37, 38, 35, 36, 0, 0, 0,
	0, 19, 20, 21, 22, 23, 0, 37, 38,
}
var yyPact = [...]int{

	332, -1000, 153, -1000, -1000, -1000, -1000, -1000, 332, -24,
	-1000, -1000, -1000, -1000, 301, 242, 332, 332, 332, 332,
	332, 332, 332, 332, 332, 332, 332, 332, 332, 332,
	332, 332, 332, 332, 332, 332, 332, -1, 332, 105,
	294, -1000, -25, 153, -1000, -32, 76, -17, -17, -17,
	82, 82, -17, -17, -17, 374, 374, 388, 388, 388,
	388, 195, 175, 215, 356, 267, 40, 40, -1000, 130,
	-1000, -1000, -31, -1000, 332, -1000, 332, 332, -1000, -1000,
	153, 31, 153, 332, 153,
}
var yyPgo = [...]int{

	0, 71, 0, 62, 59, 50, 15, 14, 38, 10,
}
var yyR1 = [...]int{

	0, 1, 2, 2, 2, 2, 2, 2, 2, 2,
	3, 3, 3, 3, 3, 3, 3, 3, 4, 4,
	4, 4, 4, 4, 5, 5, 5, 5, 5, 5,
	5, 5, 5, 6, 6, 6, 6, 6, 6, 7,
	7, 7, 8, 8, 9, 9,
}
var yyR2 = [...]int{

	0, 1, 1, 1, 1, 1, 1, 3, 3, 4,
	1, 1, 1, 1, 2, 3, 2, 3, 2, 3,
	3, 3, 3, 3, 2, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 2, 1,
	3, 4, 1, 3, 3, 5,
}
var yyChk = [...]int{

	-1000, -1, -2, -3, -4, -5, -6, -7, 32, 8,
	4, 5, 6, 7, 30, 34, 24, 28, 19, 23,
	24, 25, 26, 27, 11, 12, 13, 14, 15, 16,
	9, 10, 20, 22, 21, 17, 18, 29, 30, -2,
	32, 31, -8, -2, 35, -9, -2, -2, -2, -2,
	-2, -2, -2, -2, -2, -2, -2, -2, -2, -2,
	-2, -2, -2, -2, -2, -2, -2, -2, 8, -2,
	33, 33, -8, 31, 36, 35, 36, 37, 31, 33,
	-2, -2, -2, 37, -2,
}
var yyDef = [...]int{

	0, -2, 1, 2, 3, 4, 5, 6, 0, 39,
	10, 11, 12, 13, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 14, 0, 42, 16, 0, 0, 18, 24, 38,
	19, 20, 21, 22, 23, 25, 26, 27, 28, 29,
	30, 31, 32, 33, 34, 35, 36, 37, 40, 0,
	7, 8, 0, 15, 0, 17, 0, 0, 41, 9,
	43, 0, 44, 0, 45,
}
var yyTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 28, 3, 3, 3, 27, 22, 3,
	32, 33, 25, 23, 36, 24, 29, 26, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 37, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 30, 3, 31, 21, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 34, 20, 35,
}
var yyTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19,
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
		//line parser.go.y:62
		{
			yyVAL.expr = yyDollar[1].expr
			yylex.(*Lexer).result = yyVAL.expr
		}
	case 7:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:74
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 8:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:75
		{
			yyVAL.expr = callFunction(yylex.(*Lexer).functions, yyDollar[1].token.literal, []interface{}{})
		}
	case 9:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.go.y:76
		{
			yyVAL.expr = callFunction(yylex.(*Lexer).functions, yyDollar[1].token.literal, yyDollar[3].exprList)
		}
	case 10:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:80
		{
			yyVAL.expr = nil
		}
	case 11:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:81
		{
			yyVAL.expr = yyDollar[1].token.value
		}
	case 12:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:82
		{
			yyVAL.expr = yyDollar[1].token.value
		}
	case 13:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:83
		{
			yyVAL.expr = yyDollar[1].token.value
		}
	case 14:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:84
		{
			yyVAL.expr = []interface{}{}
		}
	case 15:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:85
		{
			yyVAL.expr = yyDollar[2].exprList
		}
	case 16:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:86
		{
			yyVAL.expr = map[string]interface{}{}
		}
	case 17:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:87
		{
			yyVAL.expr = yyDollar[2].exprMap
		}
	case 18:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:91
		{
			yyVAL.expr = unaryMinus(yyDollar[2].expr)
		}
	case 19:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:92
		{
			yyVAL.expr = add(yyDollar[1].expr, yyDollar[3].expr)
		}
	case 20:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:93
		{
			yyVAL.expr = sub(yyDollar[1].expr, yyDollar[3].expr)
		}
	case 21:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:94
		{
			yyVAL.expr = mul(yyDollar[1].expr, yyDollar[3].expr)
		}
	case 22:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:95
		{
			yyVAL.expr = div(yyDollar[1].expr, yyDollar[3].expr)
		}
	case 23:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:96
		{
			yyVAL.expr = mod(yyDollar[1].expr, yyDollar[3].expr)
		}
	case 24:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:100
		{
			yyVAL.expr = !asBool(yyDollar[2].expr)
		}
	case 25:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:101
		{
			yyVAL.expr = deepEqual(yyDollar[1].expr, yyDollar[3].expr)
		}
	case 26:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:102
		{
			yyVAL.expr = !deepEqual(yyDollar[1].expr, yyDollar[3].expr)
		}
	case 27:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:103
		{
			yyVAL.expr = compare(yyDollar[1].expr, yyDollar[3].expr, "<")
		}
	case 28:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:104
		{
			yyVAL.expr = compare(yyDollar[1].expr, yyDollar[3].expr, ">")
		}
	case 29:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:105
		{
			yyVAL.expr = compare(yyDollar[1].expr, yyDollar[3].expr, "<=")
		}
	case 30:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:106
		{
			yyVAL.expr = compare(yyDollar[1].expr, yyDollar[3].expr, ">=")
		}
	case 31:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:107
		{
			left := asBool(yyDollar[1].expr)
			right := asBool(yyDollar[3].expr)
			yyVAL.expr = left && right
		}
	case 32:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:108
		{
			left := asBool(yyDollar[1].expr)
			right := asBool(yyDollar[3].expr)
			yyVAL.expr = left || right
		}
	case 33:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:112
		{
			yyVAL.expr = asInteger(yyDollar[1].expr) | asInteger(yyDollar[3].expr)
		}
	case 34:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:113
		{
			yyVAL.expr = asInteger(yyDollar[1].expr) & asInteger(yyDollar[3].expr)
		}
	case 35:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:114
		{
			yyVAL.expr = asInteger(yyDollar[1].expr) ^ asInteger(yyDollar[3].expr)
		}
	case 36:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:115
		{
			l := asInteger(yyDollar[1].expr)
			r := asInteger(yyDollar[3].expr)
			if r >= 0 {
				yyVAL.expr = l << uint(r)
			} else {
				yyVAL.expr = l >> uint(-r)
			}
		}
	case 37:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:116
		{
			l := asInteger(yyDollar[1].expr)
			r := asInteger(yyDollar[3].expr)
			if r >= 0 {
				yyVAL.expr = l >> uint(r)
			} else {
				yyVAL.expr = l << uint(-r)
			}
		}
	case 38:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:117
		{
			yyVAL.expr = ^asInteger(yyDollar[2].expr)
		}
	case 39:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:121
		{
			yyVAL.expr = accessVar(yylex.(*Lexer).variables, yyDollar[1].token.literal)
		}
	case 40:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:122
		{
			yyVAL.expr = accessField(yyDollar[1].expr, yyDollar[3].token.literal)
		}
	case 41:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.go.y:123
		{
			yyVAL.expr = accessField(yyDollar[1].expr, yyDollar[3].expr)
		}
	case 42:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:127
		{
			yyVAL.exprList = []interface{}{yyDollar[1].expr}
		}
	case 43:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:128
		{
			yyVAL.exprList = append(yyDollar[1].exprList, yyDollar[3].expr)
		}
	case 44:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:132
		{
			yyVAL.exprMap = make(map[string]interface{})
			yyVAL.exprMap[asObjectKey(yyDollar[1].expr)] = yyDollar[3].expr
		}
	case 45:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:133
		{
			yyVAL.exprMap = addObjectMember(yyDollar[1].exprMap, yyDollar[3].expr, yyDollar[5].expr)
		}
	}
	goto yystack /* stack new state and value */
}
