package goval

import (
	"go/token"
	"go/scanner"
	"strconv"
	"errors"
	"fmt"
	"strings"
	"unsafe"
)

const bitSizeOfInt = int(unsafe.Sizeof(0)) * 8

type tokenInf struct {
	literal string
	value   interface{}
}

type lexer struct {
	scanner scanner.Scanner
	result  interface{}

	nextTokenType int
	nextTokenInfo tokenInf

	variables map[string]interface{}
	functions map[string]ExpressionFunction
}

func newLexer(src string, variables map[string]interface{}, functions map[string]ExpressionFunction) *lexer {
	if variables == nil {
		variables = map[string]interface{}{}
	}
	if functions == nil {
		functions = map[string]ExpressionFunction{}
	}

	lexer := &lexer{
		variables: variables,
		functions: functions,
	}

	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(src))

	lexer.scanner.Init(file, []byte(src), nil, 0)
	return lexer
}

func (l *lexer) scan() (token.Pos, token.Token, string) {
	for {
		pos, tok, lit := l.scanner.Scan()
		if tok == token.SEMICOLON && lit == "\n" {
			// go/scanner automatically inserted this token --> ignore it
			continue
		}
		if tok.IsKeyword() {
			// go knows about keywords, we don't. So we treat them as simple identifiers
			tok = token.IDENT
		}
		return pos, tok, lit
	}
}
func (l *lexer) Lex(lval *yySymType) int {
	var tokenType int
	var err error

	if l.nextTokenType > 0 {
		// The last scan-operation returned multiple tokens, so we return the remaining one
		tokenType = l.nextTokenType
		l.nextTokenType = 0
		lval.token = l.nextTokenInfo
		return tokenType
	}

	pos, tok, lit := l.scan()

	tokenInfo := tokenInf{
		value:   nil,
		literal: lit,
	}

	switch tok {

	case token.EOF:
		tokenType = 0

		// Literals

	case token.INT:
		tokenType = _LITERAL_NUMBER
		hex := strings.TrimPrefix(lit, "0x")
		if len(hex) < len(lit) {
			var hexVal uint64
			hexVal, err = strconv.ParseUint(hex, 16, bitSizeOfInt)
			tokenInfo.value = int(hexVal)
		} else {
			tokenInfo.value, err = strconv.Atoi(lit)
		}
		if err != nil {
			l.Perrorf(pos, "parse error: cannot parse integer")
		}
	case token.FLOAT:
		tokenType = _LITERAL_NUMBER
		tokenInfo.value, err = strconv.ParseFloat(lit, 64)
		if err != nil {
			l.Perrorf(pos, "parse error: cannot parse float")
		}

	case token.STRING:
		tokenType = _LITERAL_STRING
		tokenInfo.value, err = strconv.Unquote(lit)
		if err != nil {
			l.Perrorf(pos, "parse error: cannot unquote string literal")
		}

		// Arithmetic

	case token.ADD, token.SUB, token.MUL, token.QUO, token.REM:
		tokenType = int(tok.String()[0])

		// Logic

	case token.NOT:
		tokenType = int(tok.String()[0])

	case token.LAND:
		tokenType = _AND
	case token.LOR:
		tokenType = _OR

	case token.EQL:
		tokenType = _EQL
	case token.NEQ:
		tokenType = _NEQ

	case token.LSS:
		tokenType = _LSS
	case token.GTR:
		tokenType = _GTR
	case token.LEQ:
		tokenType = _LEQ
	case token.GEQ:
		tokenType = _GEQ

	case token.ARROW:
		// This token is known by go, but not within our expressions.
		// Instead, we treat it as two tokens (less and unary-minus).
		tokenType = _LSS
		tokenInfo.literal = "<"
		// Remember the minus-operator and omit it the next time:
		l.nextTokenType = int('-')
		l.nextTokenInfo = tokenInf{
			value:   nil,
			literal: "-",
		}

		// Bit manipulations

	case token.AND, token.OR, token.XOR:
		tokenType = int(tok.String()[0])

	case token.SHL:
		tokenType = _SHL
	case token.SHR:
		tokenType = _SHR

	case token.IDENT:
		if lit == "nil" {
			tokenType = _LITERAL_NIL
		} else if lit == "true" {
			tokenType = _LITERAL_BOOL
			tokenInfo.value = true
		} else if lit == "false" {
			tokenType = _LITERAL_BOOL
			tokenInfo.value = false
		} else if lit == "in" || lit == "IN" {
			tokenType = _IN
		} else {
			tokenType = _IDENT
		}

	case token.PERIOD:
		tokenType = int('.')

	case token.COMMA:
		tokenType = int(',')

	case token.COLON:
		tokenType = int(':')

	case token.LBRACK, token.RBRACK,
		token.LBRACE, token.RBRACE,
		token.LPAREN, token.RPAREN:
		tokenType = int(tok.String()[0])

	case token.ILLEGAL:
		if lit == "~" {
			tokenType = _BIT_NOT
			break
		}
		fallthrough

	default:
		l.Perrorf(pos, "unknown token %q (%q)", tok.String(), lit)
	}

	lval.token = tokenInfo
	return tokenType
}

func (l *lexer) Error(e string) {
	panic(errors.New(e))
}

func (l *lexer) Perrorf(pos token.Pos, format string, a ...interface{}) {
	if pos.IsValid() {
		format = format + " at position " + strconv.Itoa(int(pos))
	}
	panic(fmt.Errorf(format, a...))
}
