package goval

import (
	"runtime"
)

func init() {
	// yyDebug = 4
	yyErrorVerbose = true
}

// NewEvaluator creates a new evaluator.
func NewEvaluator() *Evaluator {
	return &Evaluator{
		parser: yyNewParser(),
	}
}

// Evaluator is used to evaluate expression strings.
type Evaluator struct {
	parser yyParser
}

// ExpressionFunction can be called from within expressions.
//
// The returned object needs to have one of the following types: `nil`, `bool`, `int`, `float64`, `[]interface{}` or `map[string]interface{}`.
type ExpressionFunction func(args ...interface{}) (interface{}, error)

// Evaluate the given expression string.
//
// Optionally accepts a list of variables (accessible but not modifiable from within expressions).
//
// Optionally accepts a list of expression functions (can be called from within expressions).
//
// Returns the resulting object or an error.
//
// Stateless. Can be called concurrently. If expression functions modify variables, concurrent execution requires additional synchronization.
func (e *Evaluator) Evaluate(str string, variables map[string]interface{}, functions map[string]ExpressionFunction) (result interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			err = r.(error)
		}
	}()

	lexer := newLexer(str, variables, functions)

	e.parser.Parse(lexer)
	return lexer.result, nil
}
