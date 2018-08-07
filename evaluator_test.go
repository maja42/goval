package goval

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Evaluator(t *testing.T) {
	variables := map[string]interface{}{
		"var": 21,
	}
	functions := map[string]ExpressionFunction{
		"func": func(args ...interface{}) (interface{}, error) {
			return args[0], nil
		},
	}

	evaluator := NewEvaluator()
	result, err := evaluator.Evaluate("func( var ) + 21", variables, functions)
	assert.NoError(t, err)
	assert.Equal(t, 42, result)
}

func Test_EvaluatorConfig(t *testing.T) {
	evaluator := NewEvaluator()
	evaluator.Configure(EvalConfig{})

	cases := []string{
		"1+1",
		"1.0+1.0",
		"\"toto\"",
	}

	for _, v := range cases {
		_, err := evaluator.Evaluate(v, nil, nil)
		assert.EqualError(t, err, "Type not accepted")
	}

	evaluator.Configure(EvalConfig{IntegerResult: true, FloatResult: true, StringResult: true})

	for _, v := range cases {
		_, err := evaluator.Evaluate(v, nil, nil)
		assert.NoError(t, err)
	}

}
