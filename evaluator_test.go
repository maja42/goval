package goval

import (
	"github.com/stretchr/testify/assert"
	"testing"
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
