package internal

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func Test_Literals_Simple(t *testing.T) {
	assertEvaluation(t, nil, nil, "nil")

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
}

func Test_Literals_Hex(t *testing.T) {
	assertEvaluation(t, nil, 0, "0x0")
	assertEvaluation(t, nil, 1, "0x01")
	assertEvaluation(t, nil, 10, "0x0A")
	assertEvaluation(t, nil, 255, "0xFF")
	assertEvaluation(t, nil, 42330, "0xA55A")
	assertEvaluation(t, nil, 23205, "0x5AA5")
	assertEvaluation(t, nil, 65535, "0xFFFF") // 16bit

	result, err := Evaluate("0x7FFFFFFF", nil, nil) // 32bit, leading zero
	if assert.NoError(t, err) {
		assert.Equal(t, int64(2147483647), int64(result.(int)))
	}

	if BitSizeOfInt == 32 {
		result, err = Evaluate("0x80000000", nil, nil) // 32bit, leading one (highest negative)
		if assert.NoError(t, err) {
			assert.Equal(t, int32(-2147483648), int32(result.(int)))
		}

		result, err = Evaluate("0xFFFFFFFF", nil, nil) // 32bit, leading one (lowest negative)
		if assert.NoError(t, err) {
			assert.Equal(t, int32(-1), int32(result.(int)))
		}
	}

	if BitSizeOfInt >= 64 {
		result, err = Evaluate("0xFFFFFFFF", nil, nil) // 32bit
		if assert.NoError(t, err) {
			assert.Equal(t, int64(4294967295), int64(result.(int)))
		}

		result, err = Evaluate("0x7FFFFFFFFFFFFFFF", nil, nil) // 64bit, leading zero (highest positive)
		if assert.NoError(t, err) {
			assert.Equal(t, int64(9223372036854775807), int64(result.(int)))
		}

		result, err = Evaluate("0x8000000000000000", nil, nil) // 64bit, leading one (highest negative)
		if assert.NoError(t, err) {
			assert.Equal(t, int64(-9223372036854775808), int64(result.(int)))
		}

		result, err = Evaluate("0xFFFFFFFFFFFFFFFF", nil, nil) // 64bit, leading one (lowest negative)
		if assert.NoError(t, err) {
			assert.Equal(t, int64(-1), int64(result.(int)))
		}
	}
}

func Test_Literals_Arrays(t *testing.T) {
	assertEvaluation(t, nil, []interface{}{}, `[]`)
	assertEvaluation(t, nil, []interface{}{1, 2, 3}, `[1, 2, 3]`)
	assertEvaluation(t, nil, []interface{}{true, false, 42, 4.2, "text"}, `[true, false, 42, 4.2, "text"]`)
	assertEvaluation(t, nil, []interface{}{[]interface{}{1, 2}, []interface{}{3}, []interface{}{}}, `[ [1,2], [3], [] ]`)

	assertEvaluation(t, nil, []interface{}{map[string]interface{}{}}, `[{}]`)
	assertEvaluation(t, nil, []interface{}{map[string]interface{}{"a": 1}}, `[{"a": 1}]`)

}

func Test_Literals_Objects(t *testing.T) {
	assertEvaluation(t, nil, map[string]interface{}{}, `{}`)
	assertEvaluation(t, nil, map[string]interface{}{"a": true}, `{"a": true}`)

	expected := map[string]interface{}{
		"a": false, "b": true, "c": 42, "d": 4.2, "e": "text",
	}
	assertEvaluation(t, nil, expected, `{"a": false, "b": true, "c": 42, "d": 4.2, "e": "text"}`)

	expected = map[string]interface{}{
		"a": []interface{}{34.0},
		"b": map[string]interface{}{"A": 45, "B": 1.2},
	}
	assertEvaluation(t, nil, expected, `{"a": [34.0], "b": {"A": 45, "B": 1.2}}`)
}

func Test_Literals_Objects_DynamicKeys(t *testing.T) {
	vars := map[string]interface{}{
		"str": "text",
	}
	assertEvaluation(t, vars, map[string]interface{}{"ab": true}, `{"a" + "b": true}`)
	assertEvaluation(t, vars, map[string]interface{}{"key42": true}, `{"key" + 42: true}`)
	assertEvaluation(t, vars, map[string]interface{}{"text": true}, `{str: true}`)
}

func Test_LiteralsOutOfRange(t *testing.T) {
	if BitSizeOfInt == 32 {
		assertEvalError(t, nil, "parse error: cannot parse integer at position 1", "0x100000000") // 33bit
	} else {
		assertEvalError(t, nil, "parse error: cannot parse integer at position 1", "0x10000000000000000") // 65bit
	}

	assertEvalError(t, nil, "parse error: cannot parse integer at position 1", "9999999999999999999999999999")
	assertEvalError(t, nil, "parse error: cannot parse float at position 1", "9.9e999")
}

func Test_Literals_Objects_DuplicateKey(t *testing.T) {
	assertEvalError(t, nil, "syntax error: duplicate object key \"a\"", `{"a": 0, "a": 0}`)
}

func Test_Literals_Objects_InvalidKeyType(t *testing.T) {
	assertEvalError(t, nil, "syntax error: object key must be string, but was nil", `{nil: 0}`)
	assertEvalError(t, nil, "syntax error: object key must be string, but was number", `{0: 0}`)
	assertEvalError(t, nil, "syntax error: object key must be string, but was number", `{"a": 0, 1: 0}`)
}

func Test_MissingOperator(t *testing.T) {
	assertEvalError(t, nil, "syntax error: unexpected LITERAL_BOOL", "true false")
	assertEvalError(t, nil, "syntax error: unexpected '!'", "true!")
	assertEvalError(t, nil, "syntax error: unexpected LITERAL_NUMBER", "42 42")
	assertEvalError(t, nil, "syntax error: unexpected LITERAL_NIL", "nil nil")
	assertEvalError(t, nil, "syntax error: unexpected IDENT", "42 var")
	assertEvalError(t, nil, "syntax error: unexpected IDENT", `42text`)
	assertEvalError(t, nil, "syntax error: unexpected LITERAL_STRING", `"text" "text"`)
}

func Test_UnsupportedTokens(t *testing.T) {
	assertEvalError(t, nil, "unknown token \"ILLEGAL\" (\"§\") at position 3", "0 § 0")
	assertEvalError(t, nil, "unknown token \"...\" (\"\") at position 3", "0 ... 0")
	assertEvalError(t, nil, "unknown token \"+=\" (\"\") at position 3", "0 += 0")
}

func Test_InvalidLiterals(t *testing.T) {
	assertEvalError(t, nil, "var error: variable \"bool\" does not exist", "bool")
	assertEvalError(t, nil, "var error: variable \"null\" does not exist", "null")
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

	assertEvalError(t, nil, "syntax error: unexpected $end", `{`)
	assertEvalError(t, nil, "syntax error: unexpected '}'", `}`)
	assertEvalError(t, nil, "syntax error: unexpected '}'", `{"a":}`)
	assertEvalError(t, nil, "syntax error: unexpected '}'", `{"a"}`)
	assertEvalError(t, nil, "syntax error: unexpected ':'", `{:1}`)
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
	// string + string
	assertEvaluation(t, nil, "text", `"te" + "xt"`)
	assertEvaluation(t, nil, "00", `"0" + "0"`)
	assertEvaluation(t, nil, "text", `"t" + "e" + "x" + "t"`)
	assertEvaluation(t, nil, "", `"" + ""`)

	// string + number
	assertEvaluation(t, nil, "text42", `"text" + 42`)
	assertEvaluation(t, nil, "text4.2", `"text" + 4.2`)
	assertEvaluation(t, nil, "42text", `42 + "text"`)
	assertEvaluation(t, nil, "4.2text", `4.2 + "text"`)

	// string + bool
	assertEvaluation(t, nil, "texttrue", `"text" + true`)
	assertEvaluation(t, nil, "textfalse", `"text" + false`)
	assertEvaluation(t, nil, "truetext", `true + "text"`)
	assertEvaluation(t, nil, "falsetext", `false + "text"`)

	// string + nil
	assertEvaluation(t, nil, "textnil", `"text" + nil`)
	assertEvaluation(t, nil, "niltext", `nil + "text"`)

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

func Test_Object_Concat(t *testing.T) {
	vars := map[string]interface{}{
		"obj1": map[string]interface{}{"a": 1, "b": 2},
		"obj2": map[string]interface{}{"b": 3, "c": 4},
	}

	assertEvaluation(t, vars, map[string]interface{}{}, `{} + {}`)
	assertEvaluation(t, vars, map[string]interface{}{"a": 1, "b": 3, "c": 4}, `{"a": 1, "b": 2} + {"b": 3, "c": 4}`)

	assertEvaluation(t, vars, map[string]interface{}{"a": 1, "b": 3, "c": 4}, `obj1 + obj2`)
	assert.Equal(t, map[string]interface{}{"a": 1, "b": 2}, vars["obj1"])

	assertEvaluation(t, vars, map[string]interface{}{"a": 1, "b": 2, "c": 4}, `obj2 + obj1`)
	assertEvaluation(t, vars, map[string]interface{}{"a": 1, "b": 3, "c": 4}, `{} + obj1 + {} + obj2 + obj2`)
	assertEvaluation(t, vars, map[string]interface{}{"a": 1, "b": 3, "c": 4, "d": 42}, `obj1 + {"d": 42} + obj2`)

	assert.Equal(t, map[string]interface{}{"a": 1, "b": 2}, vars["obj1"])
	assert.Equal(t, map[string]interface{}{"b": 3, "c": 4}, vars["obj2"])
}

func Test_Add_IncompatibleTypes(t *testing.T) {
	vars := getTestVars()
	assertEvalError(t, vars, "type error: cannot add or concatenate type number and nil", `0 + nil`)
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
	assertEvalError(t, vars, "type error: cannot add or concatenate type object and array", `obj + arr`)
	assertEvalError(t, vars, "type error: cannot add or concatenate type nil and array", `nil + arr`)
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
	assertEvalError(t, vars, "type error: unary minus requires number, but was nil", "-nil")
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

func Test_Arithmetic_Modulo(t *testing.T) {
	// int % int
	assertEvaluation(t, nil, 1, "4 % 3")
	assertEvaluation(t, nil, 0, "12 % -4")
	assertEvaluation(t, nil, -55, "-140 % 85")
	assertEvaluation(t, nil, -1, "-7 % -2")
	assertEvaluation(t, nil, 8, "8 % 13")
	// float % float
	assertEvaluation(t, nil, 1.5, "5.5 % 2.0")
	assertEvaluation(t, nil, 0.0, "12.0 % 4.0")
	assertEvaluation(t, nil, -2.0, "-2.0 % 4.5")
	assertEvaluation(t, nil, -4.0, "-12.5 % -4.25")
	// int % float
	assertEvaluation(t, nil, 1.0, "10 % 4.5")
	// float % int
	assertEvaluation(t, nil, 1.5, "5.5 % 2")

	assertEvaluation(t, nil, 4, "154 % 12 % 6")
	assertEvaluation(t, nil, 0.5, "1.5 % 2.5 % 1")
}

func Test_Arithmetic_InvalidTypes(t *testing.T) {
	vars := getTestVars()
	allTypes := []string{"nil", "true", "false", "42", "4.2", `"text"`, `"0"`, "[0]", "[]", "arr", `{"a":0}`, "{}", "obj"}
	typeOfAllTypes := []string{"nil", "bool", "bool", "number", "number", "string", "string", "array", "array", "array", "object", "object", "object"}

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
			// %
			expectedErr = fmt.Sprintf("type error: cannot perform modulo on type %s and %s", typ1, typ2)
			assertEvalError(t, vars, expectedErr, t1+"%"+t2)
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
	allTypes := []string{"nil", "true", "false", "42", "4.2", `"text"`, `"0"`, "[0]", "[]", "arr", `{"a":0}`, "{}", "obj"}
	typeOfAllTypes := []string{"nil", "bool", "bool", "number", "number", "string", "string", "array", "array", "array", "object", "object", "object"}

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

			result, err := Evaluate(t1+"||"+t2, vars, nil)
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
	assertEquality(t, nil, true, "nil", "nil")
	assertEquality(t, nil, false, "nil", "false")
	assertEquality(t, nil, false, "false", "nil")

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

		"arr1a": []interface{}{nil, false, true, 42, 4.2, "text", []interface{}{34.0}, map[string]interface{}{"A": 45, "B": 1.2}},
		"arr1b": []interface{}{nil, false, true, 42.0, 4.2, "text", []interface{}{34}, map[string]interface{}{"B": 1.2, "A": 45}},

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

		"obj1a": map[string]interface{}{"n": nil, "a": false, "b": true, "c": 42, "d": 4.2, "e": "text", "f": []interface{}{34.0}, "g": map[string]interface{}{"A": 45, "B": 1.2}},
		"obj1b": map[string]interface{}{"n": nil, "b": true, "a": false, "c": 42.0, "d": 4.2, "e": "text", "f": []interface{}{34}, "g": map[string]interface{}{"A": 45, "B": 1.2}},

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

func assertComparison(t *testing.T, variables map[string]interface{}, v1, v2 interface{}) {
	int1, ok := v1.(int)
	if ok {
		int2, ok := v2.(int)
		if ok {
			assertEvaluation(t, variables, int1 < int2, fmt.Sprintf("%d<%d", int1, int2))
			assertEvaluation(t, variables, int1 <= int2, fmt.Sprintf("%d<=%d", int1, int2))
			assertEvaluation(t, variables, int1 > int2, fmt.Sprintf("%d>%d", int1, int2))
			assertEvaluation(t, variables, int1 >= int2, fmt.Sprintf("%d>=%d", int1, int2))
		} else {
			float2 := v2.(float64)
			assertEvaluation(t, variables, float64(int1) < float2, fmt.Sprintf("%d<%f", int1, float2))
			assertEvaluation(t, variables, float64(int1) <= float2, fmt.Sprintf("%d<=%f", int1, float2))
			assertEvaluation(t, variables, float64(int1) > float2, fmt.Sprintf("%d>%f", int1, float2))
			assertEvaluation(t, variables, float64(int1) >= float2, fmt.Sprintf("%d>=%f", int1, float2))
		}
		return
	}

	float1 := v1.(float64)
	int2, ok := v2.(int)
	if ok {
		assertEvaluation(t, variables, float1 < float64(int2), fmt.Sprintf("%f<%d", float1, int2))
		assertEvaluation(t, variables, float1 <= float64(int2), fmt.Sprintf("%f<=%d", float1, int2))
		assertEvaluation(t, variables, float1 > float64(int2), fmt.Sprintf("%f>%d", float1, int2))
		assertEvaluation(t, variables, float1 >= float64(int2), fmt.Sprintf("%f>=%d", float1, int2))
	} else {
		float2 := v2.(float64)
		assertEvaluation(t, variables, float1 < float2, fmt.Sprintf("%f<%f", float1, float2))
		assertEvaluation(t, variables, float1 <= float2, fmt.Sprintf("%f<=%f", float1, float2))
		assertEvaluation(t, variables, float1 > float2, fmt.Sprintf("%f>%f", float1, float2))
		assertEvaluation(t, variables, float1 >= float2, fmt.Sprintf("%f>=%f", float1, float2))
	}
	return
}

func Test_Compare(t *testing.T) {
	// int, int
	assertComparison(t, nil, 3, 4)
	assertComparison(t, nil, -4, 2)
	assertComparison(t, nil, 4, 3)
	assertComparison(t, nil, 2, -4)
	assertComparison(t, nil, 2, 2)

	// float, float
	assertComparison(t, nil, 3.5, 3.51)
	assertComparison(t, nil, -4.9, 2.0)
	assertComparison(t, nil, 3.51, 3.5)
	assertComparison(t, nil, 2.1, -4.0)
	assertComparison(t, nil, 2.0, 2.0)

	// int, float
	assertComparison(t, nil, 3, 3.1)
	assertComparison(t, nil, -4, 2.0)
	assertComparison(t, nil, 4, 3.5)
	assertComparison(t, nil, 2, -4.0)
	assertComparison(t, nil, 2, 2.0)

	// float, int
	assertComparison(t, nil, 3.5, 4)
	assertComparison(t, nil, -4.9, 2)
	assertComparison(t, nil, 3.51, 3)
	assertComparison(t, nil, 2.1, -4)
	assertComparison(t, nil, 2.0, 2)
}

func Test_CompareHugeIntegers(t *testing.T) {
	// these integers can't be represented accurately as floats:
	i := 999999999999999998
	j := 999999999999999999
	assert.True(t, i < j)
	assert.False(t, float64(i) < float64(j))

	// ... we should be able to deal with them:
	assertEvaluation(t, nil, true, fmt.Sprintf("%d < %d", i, j))
	assertEvaluation(t, nil, false, fmt.Sprintf("%d.0 < %d.0", i, j))

	assertEvaluation(t, nil, true, fmt.Sprintf("%d <= %d", i, j))
	assertEvaluation(t, nil, true, fmt.Sprintf("%d.0 <= %d.0", i, j))

	assertEvaluation(t, nil, true, fmt.Sprintf("%d > %d", j, i))
	assertEvaluation(t, nil, false, fmt.Sprintf("%d.0 > %d.0", j, i))

	assertEvaluation(t, nil, true, fmt.Sprintf("%d >= %d", j, i))
	assertEvaluation(t, nil, true, fmt.Sprintf("%d.0 >= %d.0", j, i))
}

func Test_Compare_InvalidTypes(t *testing.T) {
	vars := getTestVars()
	allTypes := []string{"nil", "true", "false", "42", "4.2", `"text"`, `"0"`, "[0]", "[]", "arr", `{"a":0}`, "{}", "obj"}
	typeOfAllTypes := []string{"nil", "bool", "bool", "number", "number", "string", "string", "array", "array", "array", "object", "object", "object"}

	for idx1, t1 := range allTypes {
		for idx2, t2 := range allTypes {
			typ1 := typeOfAllTypes[idx1]
			typ2 := typeOfAllTypes[idx2]

			if typ1 == "number" && typ2 == "number" {
				continue
			}

			// <
			expectedErr := fmt.Sprintf("type error: cannot compare type %s and %s", typ1, typ2)
			assertEvalError(t, vars, expectedErr, t1+"<"+t2)
			// <=
			expectedErr = fmt.Sprintf("type error: cannot compare type %s and %s", typ1, typ2)
			assertEvalError(t, vars, expectedErr, t1+"<="+t2)
			// >
			expectedErr = fmt.Sprintf("type error: cannot compare type %s and %s", typ1, typ2)
			assertEvalError(t, vars, expectedErr, t1+">"+t2)
			// >=
			expectedErr = fmt.Sprintf("type error: cannot compare type %s and %s", typ1, typ2)
			assertEvalError(t, vars, expectedErr, t1+">="+t2)
		}

	}
}

func Test_BitManipulation_Or(t *testing.T) {
	assertEvaluation(t, nil, 0, "0|0")
	assertEvaluation(t, nil, 10, "8|2")
	assertEvaluation(t, nil, 11, "8|2|1")
	assertEvaluation(t, nil, 15, "8|4|2|1")
	assertEvaluation(t, nil, 13, "9|5")

	assertEvaluation(t, nil, 10, "8|2.0")
	assertEvaluation(t, nil, 10, "8.0|2")
	assertEvaluation(t, nil, 10, "8.0|2.0")

	assertEvalError(t, nil, "type error: cannot cast floating point number to integer without losing precision", "8|2.1")
	assertEvalError(t, nil, "type error: cannot cast floating point number to integer without losing precision", "8.1|2")
	assertEvalError(t, nil, "type error: cannot cast floating point number to integer without losing precision", "8.1|2.1")
}

func Test_BitManipulation_And(t *testing.T) {
	assertEvaluation(t, nil, 0, "8&2")
	assertEvaluation(t, nil, 8, "13&10")
	assertEvaluation(t, nil, 2, "10&15&2")

	assertEvaluation(t, nil, 2, "15&2.0")
	assertEvaluation(t, nil, 2, "15.0&2")
	assertEvaluation(t, nil, 2, "15.0&2.0")

	assertEvalError(t, nil, "type error: cannot cast floating point number to integer without losing precision", "15&2.1")
	assertEvalError(t, nil, "type error: cannot cast floating point number to integer without losing precision", "15.1&2")
	assertEvalError(t, nil, "type error: cannot cast floating point number to integer without losing precision", "15.1&2.1")
}

func Test_BitManipulation_XOr(t *testing.T) {
	assertEvaluation(t, nil, 10, "8^2")
	assertEvaluation(t, nil, 7, "13^10")
	assertEvaluation(t, nil, 0, "15^15")
	assertEvaluation(t, nil, 4, "10^15^1")

	assertEvaluation(t, nil, 7, "13^10.0")
	assertEvaluation(t, nil, 7, "13.0^10")
	assertEvaluation(t, nil, 7, "13.0^10.0")

	assertEvalError(t, nil, "type error: cannot cast floating point number to integer without losing precision", "13^10.1")
	assertEvalError(t, nil, "type error: cannot cast floating point number to integer without losing precision", "13.1^10")
	assertEvalError(t, nil, "type error: cannot cast floating point number to integer without losing precision", "13.1^10.1")
}

func Test_BitManipulation_Not(t *testing.T) {
	assertEvaluation(t, nil, 0, "~-1")
	assertEvaluation(t, nil, 0, "~-1.0")
	assertEvaluation(t, nil, 0x5AA5, "(~0xA55A) & 0xFFFF")
	assertEvaluation(t, nil, 0xA55A, "(~0x5AA5) & 0xFFFF")

	if BitSizeOfInt == 32 {
		assertEvaluation(t, nil, -1, "~0")
		assertEvaluation(t, nil, 0, "~0xFFFFFFFF")
	} else if BitSizeOfInt == 64 {
		assertEvaluation(t, nil, -1, "~0")
		assertEvaluation(t, nil, 0, "~0xFFFFFFFFFFFFFFFF")
	}
}

func Test_BitManipulation_Not_InvalidTypes(t *testing.T) {
	assertEvalError(t, nil, "type error: required number of type integer, but was nil", "~nil")
	assertEvalError(t, nil, "type error: required number of type integer, but was bool", "~true")
	assertEvalError(t, nil, "type error: required number of type integer, but was bool", "~false")
	assertEvalError(t, nil, "type error: cannot cast floating point number to integer without losing precision", "~4.2")
	assertEvalError(t, nil, "type error: required number of type integer, but was string", `~"text"`)
	assertEvalError(t, nil, "type error: required number of type integer, but was array", "~[]")
	assertEvalError(t, nil, "type error: required number of type integer, but was object", "~{}")
}

func Test_BitManipulation_Shift(t *testing.T) {
	assertEvaluation(t, nil, 1, "0x01 << 0")
	assertEvaluation(t, nil, 2, "0x01 << 1")
	assertEvaluation(t, nil, 4, "0x01 << 2")
	assertEvaluation(t, nil, 24, "0x03 << 3")

	if BitSizeOfInt == 32 {
		assertEvaluation(t, nil, -2147483648, "0x01 << 31") // 32bit, leading one (highest negative)
		assertEvaluation(t, nil, 0, "0x01 << 32")           // 32bit, truncated
	} else if BitSizeOfInt == 64 {
		assertEvaluation(t, nil, -9223372036854775808, "0x01 << 63") // 64bit, leading one (highest negative)
		assertEvaluation(t, nil, 0, "0x01 << 64")                    // 64bit, truncated
	}

	if BitSizeOfInt == 32 {
		assertEvaluation(t, nil, 1, "0x40000000 >> 30")
		assertEvaluation(t, nil, 2, "0x40000000 >> 28")
		assertEvaluation(t, nil, 4, "0x40000000 >> 28")
		assertEvaluation(t, nil, 12, "0x60000000 >> 27")

		assertEvaluation(t, nil, 0, "0x40000000 >> 31")  // underflow
		assertEvaluation(t, nil, -1, "0x80000000 >> 31") // sign extension
	} else if BitSizeOfInt == 64 {
		assertEvaluation(t, nil, 1, "0x4000000000000000 >> 62")
		assertEvaluation(t, nil, 2, "0x4000000000000000 >> 61")
		assertEvaluation(t, nil, 4, "0x4000000000000000 >> 60")
		assertEvaluation(t, nil, 12, "0x6000000000000000 >> 59")

		assertEvaluation(t, nil, 0, "0x4000000000000000 >> 63")  // underflow
		assertEvaluation(t, nil, -1, "0x8000000000000000 >> 63") // sign extension
	}
}

func Test_BitManipulation_NegativeShift(t *testing.T) {
	assertEvaluation(t, nil, 0, "0x01 << -1") // underflow

	assertEvaluation(t, nil, 2, "0x01 >> -1")
	assertEvaluation(t, nil, 4, "0x01 >> -2")
	assertEvaluation(t, nil, 24, "0x03 >> -3")

	if BitSizeOfInt == 32 {
		assertEvaluation(t, nil, -2147483648, "0x01 >> -31") // 32bit, leading one (highest negative)
		assertEvaluation(t, nil, 0, "0x01 >> -32")           // 32bit, truncated
	} else if BitSizeOfInt == 64 {
		assertEvaluation(t, nil, -9223372036854775808, "0x01 >> -63") // 64bit, leading one (highest negative)
		assertEvaluation(t, nil, 0, "0x01 >> -64")                    // 64bit, truncated
	}

	if BitSizeOfInt == 32 {
		assertEvaluation(t, nil, 1, "0x40000000 << -30")
		assertEvaluation(t, nil, 2, "0x40000000 << -28")
		assertEvaluation(t, nil, 4, "0x40000000 << -28")
		assertEvaluation(t, nil, 12, "0x60000000 << -27")

		assertEvaluation(t, nil, 0, "0x40000000 << -31")  // underflow
		assertEvaluation(t, nil, -1, "0x80000000 << -31") // sign extension
	} else if BitSizeOfInt == 64 {
		assertEvaluation(t, nil, 1, "0x4000000000000000 << -62")
		assertEvaluation(t, nil, 2, "0x4000000000000000 << -61")
		assertEvaluation(t, nil, 4, "0x4000000000000000 << -60")
		assertEvaluation(t, nil, 12, "0x6000000000000000 << -59")

		assertEvaluation(t, nil, 0, "0x4000000000000000 << -63")  // underflow
		assertEvaluation(t, nil, -1, "0x8000000000000000 << -63") // sign extension
	}
}

func Test_BitManipulation_InvalidTypes(t *testing.T) {
	vars := getTestVars()
	allTypes := []string{"nil", "true", "false", "42", "4.0", `"text"`, `"0"`, "[0]", "[]", "arr", `{"a":0}`, "{}", "obj"}
	typeOfAllTypes := []string{"nil", "bool", "bool", "number", "number", "string", "string", "array", "array", "array", "object", "object", "object"}

	for idx1, t1 := range allTypes {
		for idx2, t2 := range allTypes {
			typ1 := typeOfAllTypes[idx1]
			typ2 := typeOfAllTypes[idx2]

			if typ1 == "number" && typ2 == "number" {
				continue
			}

			nonIntType := typ1
			if typ1 == "number" {
				nonIntType = typ2
			}
			// &
			expectedErr := fmt.Sprintf("type error: required number of type integer, but was %s", nonIntType)
			assertEvalError(t, vars, expectedErr, t1+"&"+t2)
			// |
			expectedErr = fmt.Sprintf("type error: required number of type integer, but was %s", nonIntType)
			assertEvalError(t, vars, expectedErr, t1+"|"+t2)
			// ^
			expectedErr = fmt.Sprintf("type error: required number of type integer, but was %s", nonIntType)
			assertEvalError(t, vars, expectedErr, t1+"^"+t2)
			// <<
			expectedErr = fmt.Sprintf("type error: required number of type integer, but was %s", nonIntType)
			assertEvalError(t, vars, expectedErr, t1+"<<"+t2)
			// >>
			expectedErr = fmt.Sprintf("type error: required number of type integer, but was %s", nonIntType)
			assertEvalError(t, vars, expectedErr, t1+">>"+t2)
		}

	}
}

func Test_BitManipulation_CannotCastFloat(t *testing.T) {
	expectedErr := "type error: cannot cast floating point number to integer without losing precision"

	// &
	assertEvalError(t, nil, expectedErr, "0 & 4.2")
	assertEvalError(t, nil, expectedErr, "4.2 & 0")
	assertEvalError(t, nil, expectedErr, "4.2 & 4.2")
	// |
	assertEvalError(t, nil, expectedErr, "0 | 4.2")
	assertEvalError(t, nil, expectedErr, "4.2 | 0")
	assertEvalError(t, nil, expectedErr, "4.2 | 4.2")
	// ^
	assertEvalError(t, nil, expectedErr, "0 ^ 4.2")
	assertEvalError(t, nil, expectedErr, "4.2 ^ 0")
	assertEvalError(t, nil, expectedErr, "4.2 ^ 4.2")
	// <<
	assertEvalError(t, nil, expectedErr, "0 << 4.2")
	assertEvalError(t, nil, expectedErr, "4.2 << 0")
	assertEvalError(t, nil, expectedErr, "4.2 << 4.2")
	// >>
	assertEvalError(t, nil, expectedErr, "0 >> 4.2")
	assertEvalError(t, nil, expectedErr, "4.2 >> 0")
	assertEvalError(t, nil, expectedErr, "4.2 >> 4.2")
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

	assertEvaluation(t, vars, 4.2, `{"a": 4.2}.a`)
	assertEvaluation(t, vars, 4.2, `({"a": 4.2}).a`)
	assertEvaluation(t, vars, 4.2, `{"a": 4.2}["a"]`)
	assertEvaluation(t, vars, 4.2, `({"a": 4.2})["a"]`)

	assertEvaluation(t, vars, 42, `{"a": {"b": 42}}.a.b`)
	assertEvaluation(t, vars, 42, `{"a": {"b": 42}}["a"]["b"]`)
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

	assertEvalError(t, vars, "syntax error: array index must be number, but was string", "arr.key")
	assertEvalError(t, vars, "syntax error: cannot access fields on type string", `"txt".key`)
	assertEvalError(t, vars, "syntax error: cannot access fields on type nil", `nil.key`)
	assertEvalError(t, vars, "syntax error: cannot access fields on type number", `4.2.key`)
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

	assertEvaluation(t, vars, 42, `[0, [1, 2, 42]][1][2]`)
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

	assertEvalError(t, vars, "syntax error: cannot access fields on type string", `"txt"[0]`)
	assertEvalError(t, vars, "syntax error: cannot access fields on type nil", `nil[0]`)
	assertEvalError(t, vars, "syntax error: cannot access fields on type number", `4.2[0]`)
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

func Test_VariableAccess_Structs(t *testing.T) {
	type NestedTestType struct {
		Name          string
		nonExportable string
	}

	type TestType struct {
		Title        string
		Nested       NestedTestType
		NestedPtr       *NestedTestType
		Interfaced   interface{}
		AnotherField interface{}
		AnotherField2 interface{}
		SliceField []NestedTestType
		SliceFieldInterface []interface{}
		SliceFieldPtr []*NestedTestType

		nonExportable string
	}

	nestedPointer := &NestedTestType{
		Name:          "a",
		nonExportable: "b",
	}

	vars := map[string]interface{}{
		"obj": TestType{
			Title:         "c",
			Nested:        NestedTestType{Name: "d", nonExportable: "e"},
			NestedPtr:        &NestedTestType{Name: "f", nonExportable: "g"},
			Interfaced:    NestedTestType{
				Name:          "h",
				nonExportable: "i",
			},
			AnotherField:  &NestedTestType{
				Name:          "j",
				nonExportable: "k",
			},
			AnotherField2:  &nestedPointer,
			SliceField: []NestedTestType{{Name:"l",nonExportable:"m"}},
			SliceFieldInterface: []interface{}{NestedTestType{Name:"n",nonExportable:"o"}},
			SliceFieldPtr: []*NestedTestType{{Name:"p",nonExportable:"q"}},
			nonExportable: "r",
		},
	}

	assertEvaluation(t, vars, "c", `obj.Title`)
	assertEvaluation(t, vars, "d", `obj.Nested.Name`)
	assertEvaluation(t, vars, "f", `obj.NestedPtr.Name`)
	assertEvaluation(t, vars, "h", `obj.Interfaced.Name`)
	assertEvaluation(t, vars, "j", `obj.AnotherField.Name`)
	assertEvaluation(t, vars, "a", `obj.AnotherField2.Name`)
	assertEvaluation(t, vars, "l", `obj.SliceField[0].Name`)
	assertEvaluation(t, vars, "n", `obj.SliceFieldInterface[0].Name`)
	assertEvaluation(t, vars, "p", `obj.SliceFieldPtr[0].Name`)

	assertEvalError(t, vars, `var error: object has no member "nonExistend"`, `obj.nonExistend`)
	assertEvalError(t, vars, `var error: object member "nonExportable" is inaccessible`, `obj.nonExportable`)
	assertEvalError(t, vars, `var error: object has no member "nonExistend"`, `obj.Nested.nonExistend`)
	assertEvalError(t, vars, `var error: object member "nonExportable" is inaccessible`, `obj.Nested.nonExportable`)
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

func Test_In(t *testing.T) {
	obj := map[string]interface{}{
		"a": 3,
		"b": 4.0,
		"c": 5.5,
	}
	arr := []interface{}{
		nil, true, false, 42, 4.2, 8.0, "", "abc", []interface{}{}, []interface{}{0, 1.0, 2.4}, obj,
	}
	vars := map[string]interface{}{
		"num":   42,
		"empty": []interface{}{},
		"obj":   obj,

		"arr": arr,
	}

	assertEvaluation(t, vars, false, `nil in []`)
	assertEvaluation(t, vars, false, `false in []`)
	assertEvaluation(t, vars, false, `true in [false]`)

	assertEvaluation(t, vars, true, `1 in [1]`)
	assertEvaluation(t, vars, true, `1 in [1.0]`)
	assertEvaluation(t, vars, true, `1.0 in [1]`)
	assertEvaluation(t, vars, true, `1.0 in [1.0]`)
	assertEvaluation(t, vars, true, `1 IN [1]`)

	assertEvaluation(t, vars, true, `false in [false, true]`)
	assertEvaluation(t, vars, true, `true in [false, true]`)

	assertEvaluation(t, vars, true, `nil in arr`)
	assertEvaluation(t, vars, true, `true in arr`)
	assertEvaluation(t, vars, true, `false in arr`)
	assertEvaluation(t, vars, true, `42 in arr`)
	assertEvaluation(t, vars, true, `42.0 in arr`)
	assertEvaluation(t, vars, true, `4.2 in arr`)
	assertEvaluation(t, vars, true, `8.0 in arr`)
	assertEvaluation(t, vars, true, `8 in arr`)
	assertEvaluation(t, vars, true, `"" in arr`)
	assertEvaluation(t, vars, true, `"abc" in arr`)
	assertEvaluation(t, vars, true, `[] in arr`)
	assertEvaluation(t, vars, true, `[0, 1.0, 2.4] in arr`)
	assertEvaluation(t, vars, true, `[0.0, 1, 2.4] in arr`)
	assertEvaluation(t, vars, true, `{"a":3, "b": 4.0, "c": 5.5} in arr`)
	assertEvaluation(t, vars, true, `{"a":3.0, "c": 5.5, "b": 4} in arr`)

	assertEvaluation(t, vars, false, `8.01 in arr`)
	assertEvaluation(t, vars, false, `[nil] in arr`)
	assertEvaluation(t, vars, false, `[0, 1.0] in arr`)
	assertEvaluation(t, vars, false, `1.0 in arr`)

	assertEvaluation(t, vars, false, `{"a":3, "b": 4.0} in arr`)
	assertEvaluation(t, vars, false, `{} in arr`)

	assertEvaluation(t, vars, true, `empty in arr`)
	assertEvaluation(t, vars, true, `num in arr`)
	assertEvaluation(t, vars, true, `obj in arr`)
	assertEvaluation(t, vars, false, `arr in arr`)
}

func Test_In_InvalidTypes(t *testing.T) {
	assertEvalError(t, nil, "syntax error: in-operator requires array, but was nil", "0 in nil")
	assertEvalError(t, nil, "syntax error: in-operator requires array, but was bool", "0 in true")
	assertEvalError(t, nil, "syntax error: in-operator requires array, but was bool", "0 in false")
	assertEvalError(t, nil, "syntax error: in-operator requires array, but was number", "0 in 42")
	assertEvalError(t, nil, "syntax error: in-operator requires array, but was number", "0 in 4.2")
	assertEvalError(t, nil, "syntax error: in-operator requires array, but was string", `0 in "text"`)
	assertEvalError(t, nil, "syntax error: in-operator requires array, but was object", "0 in {}")
}

func Test_String_Slice(t *testing.T) {
	assertEvaluation(t, nil, "abcdefg", `"abcdefg"[:]`)

	assertEvaluation(t, nil, "abcdefg", `"abcdefg"[0:]`)
	assertEvaluation(t, nil, "bcdefg", `"abcdefg"[1:]`)
	assertEvaluation(t, nil, "fg", `"abcdefg"[5:]`)
	assertEvaluation(t, nil, "g", `"abcdefg"[6:]`)
	assertEvaluation(t, nil, "", `"abcdefg"[7:]`)

	assertEvaluation(t, nil, "", `"abcdefg"[:0]`)
	assertEvaluation(t, nil, "a", `"abcdefg"[:1]`)
	assertEvaluation(t, nil, "abcde", `"abcdefg"[:5]`)
	assertEvaluation(t, nil, "abcdef", `"abcdefg"[:6]`)
	assertEvaluation(t, nil, "abcdefg", `"abcdefg"[:7]`)

	assertEvaluation(t, nil, "cde", `"abcdefg"[2:5]`)
	assertEvaluation(t, nil, "d", `"abcdefg"[3:4]`)
}

func Test_String_Slice_Unicode(t *testing.T) {
	// The characters 世 and 界 both require 3 bytes
	assertEvaluation(t, nil, "Hello, ", `"Hello, 世界"[:7]`)
	assertEvaluation(t, nil, "世界", `"Hello, 世界"[7:13]`)
	assertEvaluation(t, nil, "世", `"Hello, 世界"[7:10]`)
	assertEvaluation(t, nil, "界", `"Hello, 世界"[10:13]`)
}

func Test_String_Slice_OutOfRange(t *testing.T) {
	assertEvalError(t, nil, "range error: start-index -1 is negative", `"abcd"[-1:]`)
	assertEvalError(t, nil, "range error: start-index -42 is negative", `"abcd"[-42:]`)

	assertEvalError(t, nil, "range error: end-index -1 is out of range [0, 4]", `"abcd"[:-1]`)
	assertEvalError(t, nil, "range error: end-index 5 is out of range [0, 4]", `"abcd"[:5]`)
	assertEvalError(t, nil, "range error: end-index 42 is out of range [0, 4]", `"abcd"[:42]`)

	assertEvalError(t, nil, "range error: start-index 2 is greater than end-index 1", `"abcd"[2:1]`)
}

func Test_Array_Slice(t *testing.T) {
	arr := []interface{}{0, 1, 2, 3, 4, 5, 6}
	vars := map[string]interface{}{"arr": arr}

	assertEvaluation(t, vars, []interface{}{}, `[][:]`)
	assertEvaluation(t, vars, []interface{}{1}, `[1][:]`)

	assertEvaluation(t, vars, arr, `arr[:]`)

	assertEvaluation(t, vars, arr[0:], `arr[0:]`)
	assertEvaluation(t, vars, arr[1:], `arr[1:]`)
	assertEvaluation(t, vars, arr[5:], `arr[5:]`)
	assertEvaluation(t, vars, arr[6:], `arr[6:]`)
	assertEvaluation(t, vars, arr[7:], `arr[7:]`)

	assertEvaluation(t, vars, arr[:0], `arr[:0]`)
	assertEvaluation(t, vars, arr[:1], `arr[:1]`)
	assertEvaluation(t, vars, arr[:5], `arr[:5]`)
	assertEvaluation(t, vars, arr[:6], `arr[:6]`)
	assertEvaluation(t, vars, arr[:7], `arr[:7]`)

	assertEvaluation(t, vars, arr[2:5], `arr[2:5]`)
	assertEvaluation(t, vars, arr[3:4], `arr[3:4]`)
}

func Test_Array_Slice_OutOfRange(t *testing.T) {
	assertEvalError(t, nil, "range error: start-index -1 is negative", `[0,1,2,3][-1:]`)
	assertEvalError(t, nil, "range error: start-index -42 is negative", `[0,1,2,3][-42:]`)

	assertEvalError(t, nil, "range error: end-index -1 is out of range [0, 4]", `[0,1,2,3][:-1]`)
	assertEvalError(t, nil, "range error: end-index 5 is out of range [0, 4]", `[0,1,2,3][:5]`)
	assertEvalError(t, nil, "range error: end-index 42 is out of range [0, 4]", `[0,1,2,3][:42]`)

	assertEvalError(t, nil, "range error: start-index 2 is greater than end-index 1", `[0,1,2,3][2:1]`)
}

func Test_Slicing_InvalidTypes(t *testing.T) {
	vars := getTestVars()
	allTypes := []string{"nil", "true", "false", "42", "4.2", `{"a":0}`, "{}", "obj"}
	typeOfAllTypes := []string{"nil", "bool", "bool", "number", "number", "object", "object", "object"}

	for idx, e := range allTypes {
		typ := typeOfAllTypes[idx]

		expectedErr := fmt.Sprintf("syntax error: slicing requires an array or string, but was %s", typ)
		assertEvalError(t, vars, expectedErr, e+"[:]")
	}
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

	tests := map[string]interface{}{`nil`: nil, `true`: true, `false`: false, `42`: 42, `4.2`: 4.2, `"text"`: "text", `"0"`: "0"}

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

	assertEvalErrorFuncs(t, vars, functions, "syntax error: no such function \"noFunc\"", `noFunc()`)
	assertEvalErrorFuncs(t, vars, functions, "syntax error: unexpected $end", `func(`)
	assertEvalErrorFuncs(t, vars, functions, "syntax error: unexpected ')'", `func)`)
	assertEvalErrorFuncs(t, vars, functions, "syntax error: unexpected ','", `func((1, 2))`)
}

func assertEvaluation(t *testing.T, variables map[string]interface{}, expected interface{}, str string) {
	result, err := Evaluate(str, variables, nil)
	if assert.NoError(t, err) {
		assert.Equal(t, expected, result)
	}
}

func assertEvaluationFuncs(t *testing.T, variables map[string]interface{}, functions map[string]ExpressionFunction, expected interface{}, str string) {
	result, err := Evaluate(str, variables, functions)
	if assert.NoError(t, err) {
		assert.Equal(t, expected, result)
	}
}

func assertEvalError(t *testing.T, variables map[string]interface{}, expectedErr string, str string) {
	assertEvalErrorFuncs(t, variables, nil, expectedErr, str)
}

func assertEvalErrorFuncs(t *testing.T, variables map[string]interface{}, functions map[string]ExpressionFunction, expectedErr string, str string) {
	result, err := Evaluate(str, variables, functions)
	if assert.Error(t, err) {
		assert.Equal(t, expectedErr, err.Error())
	}
	assert.Nil(t, result)
}

func getTestVars() map[string]interface{} {
	return map[string]interface{}{
		"nl":    nil,
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

// func Test_TokenExperiment(t *testing.T) {
// 	tokenize("~2")
// }
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
