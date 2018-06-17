package main

import (
	"runtime"
	"go/token"
)

func init() {
	// yyDebug = 4
	yyErrorVerbose = true
}

type Evaluator interface {
	Evaluate(str string, variables map[string]interface{}) (interface{}, error)
}

func NewEvaluator() Evaluator {
	return &evaluator{
		parser: yyNewParser(),
	}
}

type evaluator struct {
	parser yyParser
}

func (e *evaluator) Evaluate(str string, variables map[string]interface{}) (result interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			err = r.(error)
		}
	}()

	lexer := NewLexer(str, variables)

	e.parser.Parse(lexer)

	pos, tok, _ := lexer.scan()
	if tok != token.EOF {
		lexer.Perrorf(pos, "syntax error")
	}
	return lexer.result, nil
}
