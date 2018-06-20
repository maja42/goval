package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"strconv"
	"fmt"
	"errors"
)

//go:generate goyacc.exe -o parser.go parser.go.y

func Test_Literals(t *testing.T) {
	assertEvaluation(t, nil, true, "true")
	assertEvaluation(t, nil, false, "false")

	assertEvaluation(t, nil, 42, "42")

	assertEvaluation(t, nil, 4.2, "4.2")
	assertEvaluation(t, nil, 42.0, "42.0")
	assertEvaluation(t, nil, 42.0, "4.2e1")
	assertEvaluation(t, nil, 400.0, "4e2")

	assertEvaluation(t, nil, "text", `"text"`)
	assertEvaluation(t, nil, "", `""`)
	assertEvaluation(t, nil, `te"xt`, `"te\"xt"`)
	assertEvaluation(t, nil, `text\`, `"text\\"`)

	assertEvaluation(t, nil, "text", "`text`")
	assertEvaluation(t, nil, "", "``")
	assertEvaluation(t, nil, `text\`, "`text\\`")

	assertEvaluation(t, nil, "Hello, 世界", `"Hello, 世界"`)
	assertEvaluation(t, nil, "\t\t\n\xFF\u0100.+=!", `"\t	\n\xFF\u0100.+=!"`)

	assertEvaluation(t, nil, []interface{}{}, `[]`)
	assertEvaluation(t, nil, []interface{}{1, 2, 3}, `[1, 2, 3]`)
	assertEvaluation(t, nil, []interface{}{true, false, 42, 4.2, "text"}, `[true, false, 42, 4.2, "text"]`)
	assertEvaluation(t, nil, []interface{}{[]interface{}{1, 2}, []interface{}{3}, []interface{}{}}, `[ [1,2], [3], [] ]`)

}

func Test_LiteralsOutOfRange(t *testing.T) {
	assertEvalError(t, nil, "parse error: cannot parse integer at position 1", "9999999999999999999999999999")
	assertEvalError(t, nil, "parse error: cannot parse float at position 1", "9.9e999")
}

func Test_MissingOperator(t *testing.T) {
	assertEvalError(t, nil, "syntax error: unexpected LITERAL_BOOL", "true false")
	assertEvalError(t, nil, "syntax error: unexpected '!'", "true!")
	assertEvalError(t, nil, "syntax error: unexpected LITERAL_NUMBER", "42 42")
	assertEvalError(t, nil, "syntax error: unexpected IDENT", "42 var")
	assertEvalError(t, nil, "syntax error: unexpected IDENT", `42text`)
	assertEvalError(t, nil, "syntax error: unexpected LITERAL_STRING", `"text" "text"`)
}

func Test_InvalidLiterals(t *testing.T) {
	assertEvalError(t, nil, "var error: variable \"bool\" does not exist", "bool")
	assertEvalError(t, nil, "syntax error: unexpected LITERAL_NUMBER", `4.2.0`)

	assertEvalError(t, nil, "unknown token \"CHAR\" (\"'t'\") at position 1", `'t'`)
	assertEvalError(t, nil, "unknown token \"CHAR\" (\"'text'\") at position 1", `'text'`)
	assertEvalError(t, nil, "parse error: cannot unquote string literal at position 1", `"`)
	assertEvalError(t, nil, "parse error: cannot unquote string literal at position 1", `"text`)
	assertEvalError(t, nil, "parse error: cannot unquote string literal at position 5", `text"`)

	assertEvalError(t, nil, "syntax error: unexpected $end", `[`)
	assertEvalError(t, nil, "syntax error: unexpected ']'", `]`)
	assertEvalError(t, nil, "syntax error: unexpected ']'", `[1, ]`)
	assertEvalError(t, nil, "syntax error: unexpected ','", `[, 1]`)
}

func Test_Bool_Not(t *testing.T) {
	vars := getTestVars()
	assertEvaluation(t, vars, false, "!true")
	assertEvaluation(t, vars, true, "!false")

	assertEvaluation(t, vars, true, "!!true")
	assertEvaluation(t, vars, false, "!!false")

	// via variables:
	assertEvaluation(t, vars, false, "!tr")
	assertEvaluation(t, vars, true, "!fl")

	assertEvaluation(t, vars, true, "(!(!(true)))")
	assertEvaluation(t, vars, false, "(!(!(false)))")
}

func Test_Bool_Not_NotApplicable(t *testing.T) {
	assertEvalError(t, nil, "type error: required bool, but was number", "!0")
	assertEvalError(t, nil, "type error: required bool, but was number", "!1")

	assertEvalError(t, nil, "type error: required bool, but was string", `!"text"`)
	assertEvalError(t, nil, "type error: required bool, but was number", "!1.0")
	assertEvalError(t, nil, "type error: required bool, but was array", "![]")
	assertEvalError(t, nil, "type error: required bool, but was array", "![false]")
}

func Test_String_Concat(t *testing.T) {
	assertEvaluation(t, nil, "text", `"te" + "xt"`)
	assertEvaluation(t, nil, "00", `"0" + "0"`)
	assertEvaluation(t, nil, "text", `"t" + "e" + "x" + "t"`)
	assertEvaluation(t, nil, "", `"" + ""`)

	assertEvaluation(t, nil, "text42", `"text" + 42`)
	assertEvaluation(t, nil, "42text", `42 + "text"`)

	assertEvaluation(t, nil, "texttrue", `"text" + true`)
	assertEvaluation(t, nil, "textfalse", `"text" + false`)
	assertEvaluation(t, nil, "truetext", `true + "text"`)
	assertEvaluation(t, nil, "falsetext", `false + "text"`)

	assertEvaluation(t, nil, "truetext42false", `true +  "text" + 42 + false`)
}

func Test_Arithmetic_Add(t *testing.T) {
	// int + int
	assertEvaluation(t, nil, 42, "21 + 21")
	assertEvaluation(t, nil, 4, "0 + 4")
	// float + float
	assertEvaluation(t, nil, 4.2, "2.1 + 2.1")
	assertEvaluation(t, nil, 0.4, "0.0 + 0.4")
	// int + float
	assertEvaluation(t, nil, 23.1, "21 + 2.1")
	assertEvaluation(t, nil, 0.4, "0 + 0.4")
	// float + int
	assertEvaluation(t, nil, 23.1, "2.1 + 21")
	assertEvaluation(t, nil, 0.4, "0.4 + 0")

	assertEvaluation(t, nil, 63, "21 + 21 + 21")
	assertEvaluation(t, nil, 6.4, "2.1 + 2.1 + 2.2")
}

func Test_Add_WithUnaryMinus(t *testing.T) {
	assertEvaluation(t, nil, 21, "42 + -21")
	assertEvaluation(t, nil, 2.1, "4.2 + -2.1")

	assertEvaluation(t, nil, -1, "-4+3")
	assertEvaluation(t, nil, -1, "(-4)+3")
	assertEvaluation(t, nil, -7, "-(4+3)")
}

func Test_Array_Concat(t *testing.T) {
	vars := map[string]interface{}{
		"arr": []interface{}{true, 42},
	}
	assertEvaluation(t, vars, []interface{}{}, `[] + []`)
	assertEvaluation(t, vars, []interface{}{0, 1, 2, 3}, `[0, 1] + [2, 3]`)

	assertEvaluation(t, vars, []interface{}{true, 42, true, 42, true, 42}, `[] + arr + [] + arr + arr`)
	assert.Len(t, vars["arr"], 2)

	assertEvaluation(t, vars, []interface{}{true, 42, 0, 1, true, 42}, `arr + [0, 1] + arr`)
	assert.Len(t, vars["arr"], 2)
}

func Test_Add_IncompatibleTypes(t *testing.T) {
	vars := getTestVars()
	assertEvalError(t, vars, "type error: cannot add or concatenate type bool and bool", `false + false`)
	assertEvalError(t, vars, "type error: cannot add or concatenate type bool and bool", `false + true`)
	assertEvalError(t, vars, "type error: cannot add or concatenate type bool and number", `false + 42`)
	assertEvalError(t, vars, "type error: cannot add or concatenate type bool and array", `false + arr`)
	assertEvalError(t, vars, "type error: cannot add or concatenate type bool and object", `false + obj`)

	assertEvalError(t, vars, "type error: cannot add or concatenate type number and bool", `42 + false`)
	assertEvalError(t, vars, "type error: cannot add or concatenate type bool and bool", `true + false`)
	assertEvalError(t, vars, "type error: cannot add or concatenate type number and bool", `42 + false`)
	assertEvalError(t, vars, "type error: cannot add or concatenate type array and bool", `arr + false`)
	assertEvalError(t, vars, "type error: cannot add or concatenate type object and bool", `obj + false`)

	assertEvalError(t, vars, "type error: cannot add or concatenate type array and object", `arr + obj`)
	assertEvalError(t, vars, "type error: cannot add or concatenate type object and object", `obj + obj`)
	assertEvalError(t, vars, "type error: cannot add or concatenate type object and array", `obj + arr`)
}

func Test_UnaryMinus(t *testing.T) {
	vars := getTestVars()
	assertEvaluation(t, vars, -42, "-42")
	assertEvaluation(t, vars, -4.2, "-4.2")
	assertEvaluation(t, vars, -42.0, "-42.0")
	assertEvaluation(t, vars, -42.0, "-4.2e1")
	assertEvaluation(t, vars, -400.0, "-4e2")

	assertEvaluation(t, vars, -42, "-int")
	assertEvaluation(t, vars, -4.2, "-float")

	assertEvaluation(t, vars, -42, "(-(42))")
	assertEvaluation(t, vars, -4.2, "(-(4.2))")
}

func Test_UnaryMinus_IncompatibleTypes(t *testing.T) {
	vars := getTestVars()
	assertEvalError(t, vars, "type error: unary minus requires number, but was bool", "-true")
	assertEvalError(t, vars, "type error: unary minus requires number, but was bool", "-false")
	assertEvalError(t, vars, "type error: unary minus requires number, but was string", `-"0"`)

	assertEvalError(t, vars, "type error: unary minus requires number, but was array", `-arr`)
	assertEvalError(t, vars, "type error: unary minus requires number, but was object", `-obj`)
}

func Test_Arithmetic_Subtract(t *testing.T) {
	// int - int
	assertEvaluation(t, nil, 21, "42 - 21")
	assertEvaluation(t, nil, -4, "0 - 4")
	// float - float
	assertEvaluation(t, nil, 2.1, "4.2 - 2.1")
	assertEvaluation(t, nil, -0.4, "0.0 - 0.4")
	// int - float
	assertEvaluation(t, nil, 18.9, "21 - 2.1")
	assertEvaluation(t, nil, -0.4, "0 - 0.4")
	// float - int
	assertEvaluation(t, nil, -18.9, "2.1 - 21")
	assertEvaluation(t, nil, 0.4, "0.4 - 0")

	assertEvaluation(t, nil, 22, "42 - 12 - 8")
	assertEvaluation(t, nil, 2.2, "4.2 - 1.2 - 0.8")
}

func Test_Subtract_WithUnaryMinus(t *testing.T) {
	assertEvaluation(t, nil, 42, "21 - -21")
	assertEvaluation(t, nil, 4.2, "2.1 - -2.1")
}

func Test_Arithmetic_Multiply(t *testing.T) {
	// int * int
	assertEvaluation(t, nil, 8, "4 * 2")
	assertEvaluation(t, nil, 0, "0 * 4")
	assertEvaluation(t, nil, -8, "-2 * 4")
	assertEvaluation(t, nil, 8, "-2 * -4")
	// float * float
	assertEvaluation(t, nil, 10.5, "4.2 * 2.5")
	assertEvaluation(t, nil, 0.0, "0.0 * 2.4")
	assertEvaluation(t, nil, -0.8, "-2.0 * 0.4")
	assertEvaluation(t, nil, 0.8, "-2.0 * -0.4")
	// int * float
	assertEvaluation(t, nil, 50.0, "20 * 2.5")
	assertEvaluation(t, nil, -5.0, "10 * -0.5")
	// float * int
	assertEvaluation(t, nil, 50.0, "2.5 * 20")
	assertEvaluation(t, nil, 6.0, "0.5 * 12")

	assertEvaluation(t, nil, 24, "2 * 3 * 4")
	assertEvaluation(t, nil, 9.0, "1.2 * 2.5 * 3")
}

func Test_Arithmetic_Divide(t *testing.T) {
	// int / int
	assertEvaluation(t, nil, 1, "4 / 3")
	assertEvaluation(t, nil, 3, "12 / 4")
	assertEvaluation(t, nil, -2, "-4 / 2")
	assertEvaluation(t, nil, 2, "-4 / -2")
	// float / float
	assertEvaluation(t, nil, 2.75, "5.5 / 2.0")
	assertEvaluation(t, nil, 3.0, "12.0 / 4.0")
	assertEvaluation(t, nil, -2/4.5, "-2.0 / 4.5")
	assertEvaluation(t, nil, 2/4.5, "-2.0 / -4.5")
	// int / float
	assertEvaluation(t, nil, 2/4.5, "2 / 4.5")
	// float / int
	assertEvaluation(t, nil, 2.75, "5.5 / 2")

	assertEvaluation(t, nil, 2, "144 / 12 / 6")
	assertEvaluation(t, nil, 1.2/2.5/3, "1.2 / 2.5 / 3")
}

func Test_Arithmetic_InvalidTypes(t *testing.T) {
	vars := getTestVars()
	allTypes := []string{"true", "false", "42", "4.2", `"text"`, `"0"`, "[0]", "[]", "arr", "obj"}
	typeOfAllTypes := []string{"bool", "bool", "number", "number", "string", "string", "array", "array", "array", "object"}

	for idx1, t1 := range allTypes {
		for idx2, t2 := range allTypes {
			typ1 := typeOfAllTypes[idx1]
			typ2 := typeOfAllTypes[idx2]

			if typ1 == "number" && typ2 == "number" {
				continue
			}

			// + --> tested separately
			// -
			expectedErr := fmt.Sprintf("type error: cannot subtract type %s and %s", typ1, typ2)
			assertEvalError(t, vars, expectedErr, t1+"-"+t2)
			// *
			expectedErr = fmt.Sprintf("type error: cannot multiply type %s and %s", typ1, typ2)
			assertEvalError(t, vars, expectedErr, t1+"*"+t2)
			// /
			expectedErr = fmt.Sprintf("type error: cannot divide type %s and %s", typ1, typ2)
			assertEvalError(t, vars, expectedErr, t1+"/"+t2)
		}

	}
}

func Test_Arithmetic_Order(t *testing.T) {
	assertEvaluation(t, nil, 8, "2 + 2 * 3")
	assertEvaluation(t, nil, 8, "2 * 3 + 2")

	assertEvaluation(t, nil, 6, "4 + 8 / 4")
	assertEvaluation(t, nil, 6, "8 / 4 + 4")
}

func Test_Arithmetic_Parenthesis(t *testing.T) {
	assertEvaluation(t, nil, 8, "2 + (2 * 3)")
	assertEvaluation(t, nil, 12, "(2 + 2) * 3")
	assertEvaluation(t, nil, 8, "(2 * 3) + 2")
	assertEvaluation(t, nil, 10, "2 * (3 + 2)")

	assertEvaluation(t, nil, 6, "4 + (8 / 4)")
	assertEvaluation(t, nil, 3, "(4 + 8) / 4")
	assertEvaluation(t, nil, 6, "(8 / 4) + 4")
	assertEvaluation(t, nil, 1, "8 / (4 + 4)")
}

func Test_Literals_Parenthesis(t *testing.T) {
	assertEvaluation(t, nil, true, "(true)")
	assertEvaluation(t, nil, false, "(false)")

	assertEvaluation(t, nil, 42, "(42)")
	assertEvaluation(t, nil, 4.2, "(4.2)")

	assertEvaluation(t, nil, "text", `("text")`)
}

func Test_And(t *testing.T) {
	assertEvaluation(t, nil, false, "false && false")
	assertEvaluation(t, nil, false, "false && true")
	assertEvaluation(t, nil, false, "true && false")
	assertEvaluation(t, nil, true, "true && true")

	assertEvaluation(t, nil, false, "true && false && true")
}

func Test_Or(t *testing.T) {
	assertEvaluation(t, nil, false, "false || false")
	assertEvaluation(t, nil, true, "false || true")
	assertEvaluation(t, nil, true, "true || false")
	assertEvaluation(t, nil, true, "true || true")

	assertEvaluation(t, nil, true, "true || false || true")
}

func Test_AndOr_Order(t *testing.T) {
	// AND has precedes over OR
	assertEvaluation(t, nil, true, "true || false && false")
	assertEvaluation(t, nil, true, "false && false || true")
}

func Test_AndOr_InvalidTypes(t *testing.T) {
	vars := getTestVars()
	allTypes := []string{"true", "false", "42", "4.2", `"text"`, `"0"`, "[0]", "[]", "arr", "obj"}
	typeOfAllTypes := []string{"bool", "bool", "number", "number", "string", "string", "array", "array", "array", "object"}

	for idx1, t1 := range allTypes {
		for idx2, t2 := range allTypes {
			typ1 := typeOfAllTypes[idx1]
			typ2 := typeOfAllTypes[idx2]

			if typ1 == "bool" && typ2 == "bool" {
				continue
			}

			nonBoolType := typ1
			if typ1 == "bool" {
				nonBoolType = typ2
			}

			// and
			expectedErr := fmt.Sprintf("type error: required bool, but was %s", nonBoolType)
			assertEvalError(t, vars, expectedErr, t1+"&&"+t2)
			// or
			expectedErr = fmt.Sprintf("type error: required bool, but was %s", nonBoolType)
			assertEvalError(t, vars, expectedErr, t1+"||"+t2)

			result, err := evaluate(t1+"||"+t2, vars, nil)
			assert.Errorf(t, err, "%v || %v\n", t1, t2)
			assert.Nil(t, result)
		}

	}
}

func assertEquality(t *testing.T, variables map[string]interface{}, equal bool, v1, v2 string) {
	assertEvaluation(t, variables, equal, v1+"=="+v2)
	assertEvaluation(t, variables, !equal, v1+"!="+v2)
}

func Test_Equality_Simple(t *testing.T) {
	assertEquality(t, nil, true, "false", "false")
	assertEquality(t, nil, true, "true", "true")
	assertEquality(t, nil, false, "false", "true")

	assertEquality(t, nil, true, "42", "42")
	assertEquality(t, nil, false, "42", "41")
	assertEquality(t, nil, false, "1", "-1")

	assertEquality(t, nil, true, "4.2", "4.2")
	assertEquality(t, nil, false, "4.2", "4.1")

	assertEquality(t, nil, true, "42", "42.0")
	assertEquality(t, nil, true, "42.0", "42")

	assertEquality(t, nil, false, "42", "42.1")
	assertEquality(t, nil, false, "42.1", "42")

	assertEquality(t, nil, true, `""`, `""`)
	assertEquality(t, nil, true, `"text"`, ` "text"`)
	assertEquality(t, nil, true, `"text"`, ` "te" + "xt"`)

	assertEquality(t, nil, false, `"text"`, ` "Text"`)
	assertEquality(t, nil, false, `"0"`, ` 0`)
	assertEquality(t, nil, false, `""`, ` 0`)
}

func Test_Equality_Arrays(t *testing.T) {
	vars := map[string]interface{}{
		"null":     nil,
		"emptyArr": []interface{}{},

		"arr1a": []interface{}{false, true, 42, 4.2, "text", []interface{}{34.0}, map[string]interface{}{"A": 45, "B": 1.2}},
		"arr1b": []interface{}{false, true, 42.0, 4.2, "text", []interface{}{34}, map[string]interface{}{"B": 1.2, "A": 45}},

		"arr2": []interface{}{[]interface{}{34.0}, map[string]interface{}{"A": 45, "B": 1.2}, false, true, 42, 4.2, "text"},
		"arr3": []interface{}{false, true, 42, 4.2, "text"},
		"arr4": []interface{}{false, true, 42, 4.2, ""},
	}

	assertEquality(t, vars, true, `emptyArr`, `emptyArr`)
	assertEquality(t, vars, true, `[]`, `emptyArr`)
	assertEquality(t, vars, true, `emptyArr`, `[]`)
	assertEquality(t, vars, true, `arr1a`, `arr1a`)
	assertEquality(t, vars, true, `arr1b`, `arr1b`)
	assertEquality(t, vars, true, `arr2`, `arr2`)
	assertEquality(t, vars, true, `arr3`, `arr3`)
	assertEquality(t, vars, true, `arr3`, `[false, true, 42, 4.2, "text"]`)
	assertEquality(t, vars, true, `arr4`, `arr4`)
	assertEquality(t, vars, true, `arr4`, `[false, true, 42, 4.2, ""]`)

	assertEquality(t, vars, true, `arr1a`, `arr1b`)
	assertEquality(t, vars, true, `arr1b`, `arr1b`)

	assertEquality(t, vars, false, `arr1a`, `arr2`)
	assertEquality(t, vars, false, `arr1a`, `arr3`)
	assertEquality(t, vars, false, `arr2`, `arr3`)
	assertEquality(t, vars, false, `arr3`, `arr4`)

	assertEquality(t, vars, false, `emptyArr`, `null`)
	assertEquality(t, vars, false, `emptyArr`, `0`)
	assertEquality(t, vars, false, `emptyArr`, `arr1a`)
	assertEquality(t, vars, false, `emptyArr`, `""`)
}

func Test_Equal_Objects(t *testing.T) {
	vars := map[string]interface{}{
		"null":     nil,
		"emptyObj": map[string]interface{}{},

		"obj1a": map[string]interface{}{"a": false, "b": true, "c": 42, "d": 4.2, "e": "text", "f": []interface{}{34.0}, "g": map[string]interface{}{"A": 45, "B": 1.2}},
		"obj1b": map[string]interface{}{"b": true, "a": false, "c": 42.0, "d": 4.2, "e": "text", "f": []interface{}{34}, "g": map[string]interface{}{"A": 45, "B": 1.2}},

		"obj2": map[string]interface{}{"a": false, "b": true, "c": 42, "d": 4.2, "e": "text"},
		"obj3": map[string]interface{}{"a": false, "b": true, "c": 42, "d": 4.2, "e": ""},
	}

	assertEquality(t, vars, true, "emptyObj", "emptyObj")
	assertEquality(t, vars, true, "obj1a", "obj1a")
	assertEquality(t, vars, true, "obj1b", "obj1b")
	assertEquality(t, vars, true, "obj2", "obj2")
	assertEquality(t, vars, true, "obj3", "obj3")

	assertEquality(t, vars, true, "obj1a", "obj1b")
	assertEquality(t, vars, true, "obj1b", "obj1b")

	assertEquality(t, vars, false, "obj1a", "obj2")
	assertEquality(t, vars, false, "obj1a", "obj3")
	assertEquality(t, vars, false, "obj2", "obj3")

	assertEquality(t, vars, false, "emptyObj", "null")
	assertEquality(t, vars, false, "emptyObj", "0")
	assertEquality(t, vars, false, "emptyObj", "obj1a")
	assertEquality(t, vars, false, `emptyObj`, `""`)
}

func Test_VariableAccess_Simple(t *testing.T) {
	vars := getTestVars()
	for key, val := range vars {
		assertEvaluation(t, vars, val, key)
		assertEvaluation(t, vars, val, "("+key+")")
	}
}

func Test_VariableAccess_DoesNotExist(t *testing.T) {
	assertEvalError(t, nil, "var error: variable \"var\" does not exist", "var")
	assertEvalError(t, nil, "var error: variable \"varName\" does not exist", "varName")

	assertEvalError(t, nil, "var error: variable \"var\" does not exist", "var.field")
	assertEvalError(t, nil, "var error: variable \"var\" does not exist", "var[0]")
	assertEvalError(t, nil, "var error: variable \"var\" does not exist", "var[fieldName]")
}

func Test_VariableAccess_Arithmetic(t *testing.T) {
	vars := getTestVars()
	assertEvaluation(t, vars, 84, "int + int")
	assertEvaluation(t, vars, 8.4, "float + float")
	assertEvaluation(t, vars, 88.2, "int + float + int")
}

func Test_VariableAccess_DotSyntax(t *testing.T) {
	vars := getTestVars()

	// access object fields
	for key, val := range vars["obj"].(map[string]interface{}) {
		assertEvaluation(t, vars, val, "obj."+key)
	}
}

func Test_VariableAccess_DotSyntax_DoesNotExist(t *testing.T) {
	vars := getTestVars()
	assertEvalError(t, vars, "var error: object has no member \"key\"", "obj.key")
	assertEvalError(t, vars, "var error: object has no member \"key\"", "obj.key.field")
	assertEvalError(t, vars, "var error: object has no member \"key\"", "obj.key[0]")
	assertEvalError(t, vars, "var error: object has no member \"key\"", "obj.key[fieldName]")
}

func Test_VariableAccess_DotSyntax_InvalidType(t *testing.T) {
	vars := getTestVars()
	assertEvalError(t, vars, "syntax error: unexpected LITERAL_NUMBER", "obj.0")
}

func Test_VariableAccess_DotSyntax_InvalidSyntax(t *testing.T) {
	vars := getTestVars()
	assertEvalError(t, vars, "syntax error: unexpected '[', expecting IDENT", "obj.[b]")
}

func Test_VariableAccess_ArraySyntax(t *testing.T) {
	vars := getTestVars()

	// access object fields
	for key, val := range vars["obj"].(map[string]interface{}) {
		assertEvaluation(t, vars, val, `obj["`+key+`"]`)
		assertEvaluation(t, vars, val, `obj[("`+key+`")]`)
	}

	// access array elements
	for idx, val := range vars["arr"].([]interface{}) {
		strIdx := strconv.Itoa(idx)
		// with int:
		assertEvaluation(t, vars, val, `arr[`+strIdx+`]`)
		assertEvaluation(t, vars, val, `arr[(`+strIdx+`)]`)
		// with float:
		assertEvaluation(t, vars, val, `arr[`+strIdx+`.0]`)
		assertEvaluation(t, vars, val, `arr[(`+strIdx+`.0)]`)
	}

	// access array literals
	assertEvaluation(t, vars, false, `[false, 42,  "text"][0]`)
	assertEvaluation(t, vars, 42, `[false, 42,  "text"][1]`)
	assertEvaluation(t, vars, "text", `[false, 42, "text"][2]`)
	assertEvaluation(t, vars, 0.0, `([0.0])[0]`)
}

func Test_VariableAccess_ArraySyntax_DoesNotExist(t *testing.T) {
	vars := getTestVars()
	assertEvalError(t, vars, "var error: object has no member \"key\"", `obj["key"]`)
	assertEvalError(t, vars, "var error: object has no member \"key\"", `obj["key"].field`)
	assertEvalError(t, vars, "var error: object has no member \"key\"", `obj["key"][0]`)
	assertEvalError(t, vars, "var error: object has no member \"key\"", `obj["key"][fieldName]`)

	assertEvalError(t, vars, "var error: array index 5 is out of range [0, 4]", `arr[5]`)
	assertEvalError(t, vars, "var error: array index 6 is out of range [0, 4]", `arr[6]`)
	assertEvalError(t, vars, "var error: array index 0 is out of range [0, 0]", `[][0]`)
	assertEvalError(t, vars, "var error: array index 41 is out of range [0, 1]", `[1][41]`)
}

func Test_VariableAccess_ArraySyntax_InvalidType(t *testing.T) {
	vars := getTestVars()
	assertEvalError(t, vars, "syntax error: object key must be string, but was bool", `obj[true]`)
	assertEvalError(t, vars, "syntax error: object key must be string, but was number", `obj[0]`)
	assertEvalError(t, vars, "syntax error: object key must be string, but was array", `obj[arr]`)
	assertEvalError(t, vars, "syntax error: object key must be string, but was object", `obj[obj]`)

	assertEvalError(t, vars, "syntax error: array index must be number, but was bool", `arr[true]`)
	assertEvalError(t, vars, "syntax error: array index must be number, but was string", `arr["0"]`)
	assertEvalError(t, vars, "syntax error: array index must be number, but was string", `["0"]["0"]`)
	assertEvalError(t, vars, "syntax error: array index must be number, but was array", `arr[arr]`)
	assertEvalError(t, vars, "syntax error: array index must be number, but was object", `arr[obj]`)
}

func Test_VariableAccess_ArraySyntax_FloatHasDecimals(t *testing.T) {
	vars := getTestVars()
	assertEvalError(t, vars, "eval error: array index must be whole number, but was 0.100000", `arr[0.1]`)
	assertEvalError(t, vars, "eval error: array index must be whole number, but was 0.500000", `arr[0.5]`)
	assertEvalError(t, vars, "eval error: array index must be whole number, but was 0.900000", `arr[0.9]`)
	assertEvalError(t, vars, "eval error: array index must be whole number, but was 2.000100", `arr[2.0001]`)
}

func Test_VariableAccess_Nested(t *testing.T) {
	vars := map[string]interface{}{
		"arr": []interface{}{
			10, "a",
			[]interface{}{
				11, "b",
			},
			map[string]interface{}{
				"a": 13,
				"b": "c",
			},
		},
		"obj": map[string]interface{}{
			"a": 20,
			"b": "a",
			"c": []interface{}{
				22, 23,
			},
			"d": map[string]interface{}{
				"a": 24,
				"b": "b",
			},
		},
	}

	// array:
	assertEvaluation(t, vars, 10, `arr[0]`)
	assertEvaluation(t, vars, "a", `arr[1]`)
	assertEvaluation(t, vars, 11, `arr[2][0]`)
	assertEvaluation(t, vars, "b", `arr[2][1]`)
	assertEvaluation(t, vars, 13, `arr[3].a`)
	assertEvaluation(t, vars, 13, `arr[3]["a"]`)
	assertEvaluation(t, vars, "c", `arr[3].b`)
	assertEvaluation(t, vars, "c", `arr[3]["b"]`)
	// object:
	assertEvaluation(t, vars, 20, `obj.a`)
	assertEvaluation(t, vars, 20, `obj["a"]`)
	assertEvaluation(t, vars, "a", `obj.b`)
	assertEvaluation(t, vars, "a", `obj["b"]`)
	assertEvaluation(t, vars, 22, `obj.c[0]`)
	assertEvaluation(t, vars, 23, `obj["c"][1]`)
	assertEvaluation(t, vars, 24, `obj.d.a`)
	assertEvaluation(t, vars, 24, `obj.d["a"]`)
	assertEvaluation(t, vars, "b", `obj["d"].b`)
	assertEvaluation(t, vars, "b", `obj["d"]["b"]`)
}

func Test_VariableAccess_DynamicAccess(t *testing.T) {
	vars := map[string]interface{}{
		"num0": 0,
		"num1": 1,
		"letA": "a",
		"letB": "b",

		"arr": []interface{}{
			0, 4, "a", "abc", 42,
		},

		"obj": map[string]interface{}{
			"a":   0,
			"b":   4,
			"c":   "a",
			"d":   "abc",
			"abc": 43,
		},
	}

	assertEvaluation(t, vars, 0, `arr[num0]`)
	assertEvaluation(t, vars, 4, `arr[num1]`)
	assertEvaluation(t, vars, "a", `arr[num1 + 1]`)
	assertEvaluation(t, vars, "abc", `arr[num1 + 1 + num1]`)

	assertEvaluation(t, vars, 0, `obj[letA]`)
	assertEvaluation(t, vars, 4, `obj[letB]`)
	assertEvaluation(t, vars, 43, `obj[letA + letB + "c"]`)

	assertEvaluation(t, vars, 0, `arr[ obj.a ]`)
	assertEvaluation(t, vars, 42, `arr[ obj["b"] ]`)
	assertEvaluation(t, vars, 42, `arr[ obj[letB] ]`)
	assertEvaluation(t, vars, 0, `arr[ obj[arr[2]] ]`)
	assertEvaluation(t, vars, 0, `arr[ arr[obj.a] ]`)

	assertEvaluation(t, vars, 0, `obj[ arr[2] ]`)
	assertEvaluation(t, vars, 43, `obj[ arr[num1 + num1 + 1] ]`)
	assertEvaluation(t, vars, 43, `obj[ arr[obj.a + 3] ]`)
	assertEvaluation(t, vars, 43, `obj[ arr[obj["a"] + 3] ]`)
}

func Test_FunctionCall_Simple(t *testing.T) {
	var shouldReturn interface{}
	var expectedArg interface{}

	functions := map[string]ExpressionFunction{
		"func1": func(args ...interface{}) (interface{}, error) {
			return shouldReturn, nil
		},
		"func2": func(args ...interface{}) (interface{}, error) {
			assert.Equal(t, expectedArg, args[0])
			return args[0], nil
		},
		"func3": func(args ...interface{}) (interface{}, error) {
			return []interface{}{len(args), args}, nil
		},
		"func4": func(args ...interface{}) (interface{}, error) {
			return nil, errors.New("simulated error")
		},
	}

	tests := map[string]interface{}{`true`: true, `false`: false, `42`: 42, `4.2`: 4.2, `"text"`: "text", `"0"`: "0"}

	for expr, expected := range tests {
		shouldReturn = expected
		assertEvaluationFuncs(t, nil, functions, expected, `func1()`)
		expectedArg = expected
		assertEvaluationFuncs(t, nil, functions, expected, `func2(`+expr+`)`)
	}

	expectedReturn := []interface{}{6, []interface{}{true, false, 42, 4.2, "text", "0"}}
	assertEvaluationFuncs(t, nil, functions, expectedReturn, `func3(true, false, 42, 4.2, "text", "0")`)

	assertEvalErrorFuncs(t, nil, functions, "function error: \"func4\" - simulated error", "func4()")
}

func Test_FunctionCall_Nested(t *testing.T) {
	functions := map[string]ExpressionFunction{
		"func": func(args ...interface{}) (interface{}, error) {
			var allArgs = make([]interface{}, 0)

			for _, arg := range args {
				multi, ok := arg.([]interface{})
				if ok {
					allArgs = append(allArgs, multi...)
				} else {
					allArgs = append(allArgs, arg)
				}
			}
			return allArgs, nil
		},
	}

	assertEvaluationFuncs(t, nil, functions, []interface{}{1, 2, 3, 4}, `func(1, 2, 3, 4)`)
	assertEvaluationFuncs(t, nil, functions, []interface{}{1, 2, 3, 4}, `func([1, 2], [3], 4)`)
	assertEvaluationFuncs(t, nil, functions, []interface{}{1, 2, 3, 4}, `func(func(1, 2, 3, 4))`)
	assertEvaluationFuncs(t, nil, functions, []interface{}{1, 2, 3, 4}, `func(func(1, 2), func(3, 4))`)
	assertEvaluationFuncs(t, nil, functions, []interface{}{1, 2, 3, 4}, `func(func(1, func(2), func()), func(), func(3, 4))`)
}

func Test_FunctionCall_Variables(t *testing.T) {
	vars := getTestVars()

	functions := map[string]ExpressionFunction{
		"func": func(args ...interface{}) (interface{}, error) {
			assert.Len(t, args, 2)
			varName := args[0].(string)
			varValue := args[1]
			assert.Equal(t, vars[varName], varValue)
			return varValue, nil
		},
	}
	for name, val := range vars {
		assertEvaluationFuncs(t, vars, functions, val, `func("`+name+`", `+name+` )`)
	}

	// function with same name as variable:
	vars["func"] = "foo"
	assertEvaluationFuncs(t, vars, functions, "foo", `func("func", func)`)
}

func Test_InvalidFunctionCalls(t *testing.T) {
	vars := map[string]interface{}{"func": nil}
	functions := map[string]ExpressionFunction{
		"func": func(args ...interface{}) (interface{}, error) {
			return nil, nil
		},
	}

	assertEvalErrorFuncs(t, vars, functions, "syntax error: unexpected $end", `func(`)
	assertEvalErrorFuncs(t, vars, functions, "syntax error: unexpected ')'", `func)`)
	assertEvalErrorFuncs(t, vars, functions, "syntax error: unexpected ','", `func((1, 2))`)
}

// func tokenize(src string) {
// 	var scanner scanner.Scanner
// 	fset := token.NewFileSet()
// 	file := fset.AddFile("", fset.Base(), len(src))
// 	scanner.Init(file, []byte(src), nil, 0)
//
// 	for {
// 		pos, tok, lit := scanner.Scan()
// 		fmt.Printf("%3d %20s %q\n", pos, tok.String(), lit)
// 		if tok == token.EOF {
// 			return
// 		}
// 	}
// }

func evaluate(str string, variables map[string]interface{}, functions map[string]ExpressionFunction) (result interface{}, err error) {
	evaluator := NewEvaluator()
	return evaluator.Evaluate(str, variables, functions)
}

func assertEvaluation(t *testing.T, variables map[string]interface{}, expected interface{}, str string) {
	result, err := evaluate(str, variables, nil)
	if assert.NoError(t, err) {
		assert.Equal(t, expected, result)
	}
}

func assertEvaluationFuncs(t *testing.T, variables map[string]interface{}, functions map[string]ExpressionFunction, expected interface{}, str string) {
	result, err := evaluate(str, variables, functions)
	if assert.NoError(t, err) {
		assert.Equal(t, expected, result)
	}
}

func assertEvalError(t *testing.T, variables map[string]interface{}, expectedErr string, str string) {
	assertEvalErrorFuncs(t, variables, nil, expectedErr, str)
}

func assertEvalErrorFuncs(t *testing.T, variables map[string]interface{}, functions map[string]ExpressionFunction, expectedErr string, str string) {
	result, err := evaluate(str, variables, functions)
	if assert.Error(t, err) {
		assert.Equal(t, expectedErr, err.Error())
	}
	assert.Nil(t, result)
}

func getTestVars() map[string]interface{} {
	return map[string]interface{}{
		"tr":    true,
		"fl":    false,
		"int":   42,
		"float": 4.2,
		"str":   "text",
		"arr":   []interface{}{true, 21, 2.1, "txt"},
		"obj": map[string]interface{}{
			"b": false,
			"i": 51,
			"f": 5.1,
			"s": "tx",
		},
	}
}
