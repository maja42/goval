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

func Test_MustEvaluator(t *testing.T) {
	variables := map[string]interface{}{
		"var": 21,
	}
	functions := map[string]ExpressionFunction{
		"func": func(args ...interface{}) (interface{}, error) {
			return args[0], nil
		},
	}

	evaluator := NewEvaluator()
	result, err := evaluator.MustEvaluate("func( var ) + 21", variables, functions)
	assert.NoError(t, err)
	assert.Equal(t, 42, result)
	result, err = evaluator.MustEvaluate("1/0", variables, functions)
	assert.Nil(t, result)
	assert.NotNil(t, err)
}
