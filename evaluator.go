package goval

import (
	"errors"

	"github.com/maja42/goval/internal"
)

// EvalConfig specifies authorized type for evaluation result
type EvalConfig struct {
	StringResult  bool
	ArrayResult   bool
	IntegerResult bool
	FloatResult   bool
	ObjectResult  bool
	BoolResult    bool
}

// NewEvaluator creates a new evaluator.
func NewEvaluator() *Evaluator {
	defaultConf := EvalConfig{StringResult: true, ArrayResult: true, IntegerResult: true, FloatResult: true, ObjectResult: true, BoolResult: true}

	return &Evaluator{config: defaultConf}
}

// Evaluator is used to evaluate expression strings.
type Evaluator struct {
	config EvalConfig
}

// ExpressionFunction can be called from within expressions.
//
// The returned object needs to have one of the following types: `nil`, `bool`, `int`, `float64`, `[]interface{}` or `map[string]interface{}`.
type ExpressionFunction = func(args ...interface{}) (interface{}, error)

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
	result, err = internal.Evaluate(str, variables, functions)
	if err != nil {
		return nil, err
	}
	if err = e.checkResultType(result); err != nil {
		return nil, err
	}

	return result, nil
}

// checkResultType verifies that a result value type is authorized in config
func (e *Evaluator) checkResultType(result interface{}) error {
	if _, ok := result.(int); ok && e.config.IntegerResult {
		return nil
	}
	if _, ok := result.(float64); ok && e.config.FloatResult {
		return nil
	}
	if _, ok := result.(string); ok && e.config.StringResult {
		return nil
	}
	if _, ok := result.(bool); ok && e.config.BoolResult {
		return nil
	}
	if _, ok := result.([]interface{}); ok && e.config.ArrayResult {
		return nil
	}
	if _, ok := result.(map[string]interface{}); ok && e.config.ObjectResult {
		return nil
	}
	return errors.New("Type not accepted")
}

// Configure stores current configuration for evaluator
func (e *Evaluator) Configure(c EvalConfig) {
	e.config = c
}
