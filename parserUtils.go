package main

import (
	"fmt"
	"reflect"
	"strconv"
)

func typeOf(val interface{}) string {
	kind := reflect.TypeOf(val).Kind()

	switch kind {
	case reflect.Bool:
		return "bool"
	case reflect.Int, reflect.Float64:
		return "number"
	case reflect.String:
		return "string"
	}

	if _, ok := val.([]interface{}); ok {
		return "array"
	}

	if _, ok := val.(map[string]interface{}); ok {
		return "object"
	}

	return "<unknown type>"
}

func asBool(val interface{}) bool {
	b, ok := val.(bool)
	if !ok {
		panic(fmt.Errorf("type error: required bool, but was %s", typeOf(val)))
	}
	return b
}

// func asInt(val interface{}) bool {
// 	i, ok := val.(int)
// 	if !ok {
// 		todo: cast float!
// 		panic(fmt.Errorf("required bool, but was %v", val))
// 	}
// 	return i
// }

func add(val1 interface{}, val2 interface{}) interface{} {
	str1, str1OK := val1.(string)
	str2, str2OK := val2.(string)

	if str1OK && str2OK { // string + string = string
		return str1 + str2
	}

	int1, int1OK := val1.(int)
	int2, int2OK := val2.(int)

	if int1OK && int2OK { // int + int = int
		return int1 + int2
	}

	float1, float1OK := val1.(float64)
	float2, float2OK := val2.(float64)

	if int1OK {
		float1 = float64(int1)
		float1OK = true
	}
	if int2OK {
		float2 = float64(int2)
		float2OK = true
	}

	if float1OK && float2OK {
		return float1 + float2
	}
	if str1OK && float2OK {
		return str1 + strconv.FormatFloat(float2, 'f', -1, 64)
	}
	if float1OK && str2OK {
		return strconv.FormatFloat(float1, 'f', -1, 64) + str2
	}

	bool1, bool1OK := val1.(bool)
	bool2, bool2OK := val2.(bool)

	if str1OK && bool2OK {
		return str1 + strconv.FormatBool(bool2)
	}
	if bool1OK && str2OK {
		return strconv.FormatBool(bool1) + str2
	}

	panic(fmt.Errorf("type error: cannot add or concatenate type %s and %s", typeOf(val1), typeOf(val2)))
}

func sub(val1 interface{}, val2 interface{}) interface{} {
	int1, int1OK := val1.(int)
	int2, int2OK := val2.(int)

	if int1OK && int2OK {
		return int1 - int2
	}

	float1, float1OK := val1.(float64)
	float2, float2OK := val2.(float64)

	if int1OK {
		float1 = float64(int1)
		float1OK = true
	}
	if int2OK {
		float2 = float64(int2)
		float2OK = true
	}

	if float1OK && float2OK {
		return float1 - float2
	}
	panic(fmt.Errorf("type error: cannot subtract type %s and %s", typeOf(val1), typeOf(val2)))
}

func mul(val1 interface{}, val2 interface{}) interface{} {
	int1, int1OK := val1.(int)
	int2, int2OK := val2.(int)

	if int1OK && int2OK {
		return int1 * int2
	}

	float1, float1OK := val1.(float64)
	float2, float2OK := val2.(float64)

	if int1OK {
		float1 = float64(int1)
		float1OK = true
	}
	if int2OK {
		float2 = float64(int2)
		float2OK = true
	}

	if float1OK && float2OK {
		return float1 * float2
	}
	panic(fmt.Errorf("type error: cannot multiply type %s and %s", typeOf(val1), typeOf(val2)))
}

func div(val1 interface{}, val2 interface{}) interface{} {
	int1, int1OK := val1.(int)
	int2, int2OK := val2.(int)

	if int1OK && int2OK {
		return int1 / int2
	}

	float1, float1OK := val1.(float64)
	float2, float2OK := val2.(float64)

	if int1OK {
		float1 = float64(int1)
		float1OK = true
	}
	if int2OK {
		float2 = float64(int2)
		float2OK = true
	}

	if float1OK && float2OK {
		return float1 / float2
	}
	panic(fmt.Errorf("type error: cannot divide type %s and %s", typeOf(val1), typeOf(val2)))
}

func accessVar(variables map[string]interface{}, varName string) interface{} {
	val, ok := variables[varName]
	if !ok {
		panic(fmt.Errorf("eval error: variable %q does not exist", varName))
	}
	return val
}

func accessField(s interface{}, field interface{}) interface{} {
	obj, ok := s.(map[string]interface{})
	if ok {
		key, ok := field.(string)
		if !ok {
			panic(fmt.Errorf("syntax error: object key must be string, but was %s", typeOf(field)))
		}
		val, ok := obj[key]
		if !ok {
			panic(fmt.Errorf("eval error: object has no member %q", field))
		}
		return val
	}

	arrVar, ok := s.([]interface{})
	if ok {
		intIdx, ok := field.(int)
		if !ok {
			floatIdx, ok := field.(float64)
			if !ok {
				panic(fmt.Errorf("syntax error: array index must be number, but was %s", typeOf(field)))
			}
			intIdx = int(floatIdx)
			if float64(intIdx) != floatIdx {
				panic(fmt.Errorf("eval error: array index must be whole number, but was %f", floatIdx))
			}
		}

		if intIdx < 0 || intIdx >= len(arrVar) {
			panic(fmt.Errorf("eval error: array index %d is out of range [%d, %d]", intIdx, 0, len(arrVar)))
		}
		return arrVar[intIdx]
	}

	return nil
}
