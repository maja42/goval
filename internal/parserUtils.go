package internal

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
)

func init() {
	yyDebug = 0           // can be increased for debugging the generated parser
	yyErrorVerbose = true // make sure to get better errors than just "syntax error"
}

// ExpressionFunction can be called from within expressions.
// The returned object needs to have one of the following types: `nil`, `bool`, `int`, `float64`, `string`, `[]interface{}` or `map[string]interface{}`.
type ExpressionFunction = func(args ...interface{}) (interface{}, error)

func typeOf(val interface{}) string {
	if val == nil {
		return "nil"
	}

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

func asInteger(val interface{}) int {
	i, ok := val.(int)
	if ok {
		return i
	}
	f, ok := val.(float64)
	if !ok {
		panic(fmt.Errorf("type error: required number of type integer, but was %s", typeOf(val)))
	}

	i = int(f)
	if float64(i) != f {
		panic(fmt.Errorf("type error: cannot cast floating point number to integer without losing precision"))
	}
	return i
}

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

	if str1OK && val2 == nil {
		return str1 + "nil"
	}
	if val1 == nil && str2OK {
		return "nil" + str2
	}

	bool1, bool1OK := val1.(bool)
	bool2, bool2OK := val2.(bool)

	if str1OK && bool2OK {
		return str1 + strconv.FormatBool(bool2)
	}
	if bool1OK && str2OK {
		return strconv.FormatBool(bool1) + str2
	}

	arr1, arr1OK := val1.([]interface{})
	arr2, arr2OK := val2.([]interface{})

	if arr1OK && arr2OK {
		return append(arr1, arr2...)
	}

	obj1, obj1OK := val1.(map[string]interface{})
	obj2, obj2OK := val2.(map[string]interface{})

	if obj1OK && obj2OK {
		sum := make(map[string]interface{})
		for k, v := range obj1 {
			sum[k] = v
		}
		for k, v := range obj2 {
			sum[k] = v
		}
		return sum
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

func mod(val1 interface{}, val2 interface{}) interface{} {
	int1, int1OK := val1.(int)
	int2, int2OK := val2.(int)

	if int1OK && int2OK {
		return int1 % int2
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
		return math.Mod(float1, float2)
	}
	panic(fmt.Errorf("type error: cannot perform modulo on type %s and %s", typeOf(val1), typeOf(val2)))
}

func unaryMinus(val interface{}) interface{} {
	intVal, ok := val.(int)
	if ok {
		return -intVal
	}
	floatVal, ok := val.(float64)
	if ok {
		return -floatVal
	}
	panic(fmt.Errorf("type error: unary minus requires number, but was %s", typeOf(val)))
}

func deepEqual(val1 interface{}, val2 interface{}) bool {
	switch typ1 := val1.(type) {

	case []interface{}:
		typ2, ok := val2.([]interface{})
		if !ok || len(typ1) != len(typ2) {
			return false
		}
		for idx := range typ1 {
			if !deepEqual(typ1[idx], typ2[idx]) {
				return false
			}
		}
		return true

	case map[string]interface{}:
		typ2, ok := val2.(map[string]interface{})
		if !ok || len(typ1) != len(typ2) {
			return false
		}
		for idx := range typ1 {
			if !deepEqual(typ1[idx], typ2[idx]) {
				return false
			}
		}
		return true

	case int:
		int2, ok := val2.(int)
		if ok {
			return typ1 == int2
		}
		float2, ok := val2.(float64)
		if ok {
			return float64(typ1) == float2
		}
		return false

	case float64:
		float2, ok := val2.(float64)
		if ok {
			return typ1 == float2
		}
		int2, ok := val2.(int)
		if ok {
			return typ1 == float64(int2)
		}
		return false
	}
	return val1 == val2
}

func compare(val1 interface{}, val2 interface{}, operation string) bool {
	int1, int1OK := val1.(int)
	int2, int2OK := val2.(int)

	if int1OK && int2OK {
		return compareInt(int1, int2, operation)
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
		return compareFloat(float1, float2, operation)
	}
	panic(fmt.Errorf("type error: cannot compare type %s and %s", typeOf(val1), typeOf(val2)))
}

func compareInt(val1 int, val2 int, operation string) bool {
	switch operation {
	case "<":
		return val1 < val2
	case "<=":
		return val1 <= val2
	case ">":
		return val1 > val2
	case ">=":
		return val1 >= val2
	}
	panic(fmt.Errorf("syntax error: unsupported operation %q", operation))
}

func compareFloat(val1 float64, val2 float64, operation string) bool {
	switch operation {
	case "<":
		return val1 < val2
	case "<=":
		return val1 <= val2
	case ">":
		return val1 > val2
	case ">=":
		return val1 >= val2
	}
	panic(fmt.Errorf("syntax error: unsupported operation %q", operation))
}

func asObjectKey(key interface{}) string {
	s, ok := key.(string)
	if !ok {
		panic(fmt.Errorf("syntax error: object key must be string, but was %s", typeOf(key)))
	}
	return s
}

func asObjectIdx(key interface{}) int {
	intIdx, ok := key.(int)
	if !ok {
		floatIdx, ok := key.(float64)
		if !ok {
			panic(fmt.Errorf("syntax error: array index must be number, but was %s", typeOf(key)))
		}
		intIdx = int(floatIdx)
		if float64(intIdx) != floatIdx {
			panic(fmt.Errorf("eval error: array index must be whole number, but was %f", floatIdx))
		}
	}

	return intIdx
}

func addObjectMember(obj map[string]interface{}, key, val interface{}) map[string]interface{} {
	s := asObjectKey(key)
	_, ok := obj[s]
	if ok {
		panic(fmt.Errorf("syntax error: duplicate object key %q", s))
	}
	obj[s] = val
	return obj
}

func accessVar(variables map[string]interface{}, varName string) interface{} {
	val, ok := variables[varName]
	if !ok {
		panic(fmt.Errorf("var error: variable %q does not exist", varName))
	}
	return val
}

func accessField(s interface{}, field interface{}) interface{} {
	obj, ok := s.(map[string]interface{})
	if ok {
		key := asObjectKey(field)
		val, ok := obj[key]
		if !ok {
			panic(fmt.Errorf("var error: object has no member %q", field))
		}
		return val
	}

	arrVar, ok := s.([]interface{})
	if ok {
		intIdx := asObjectIdx(field)
		if intIdx < 0 || intIdx >= len(arrVar) {
			panic(fmt.Errorf("var error: array index %d is out of range [%d, %d]", intIdx, 0, len(arrVar)))
		}
		return arrVar[intIdx]
	}

	v := reflect.ValueOf(s)
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}

	var fieldReflect *reflect.Value
	if v.Kind() == reflect.Struct {
		key := asObjectKey(field)

		if v.MethodByName(key).IsValid() {
			panic(fmt.Errorf("syntax error: object member %q is a method and currently unsupported", field))
		}

		name := v.FieldByName(key)
		fieldReflect = &name
	} else if v.Kind() == reflect.Slice {
		intIdx := asObjectIdx(field)

		if intIdx < 0 || intIdx >= v.Len() {
			panic(fmt.Errorf("var error: array index %d is out of range [%d, %d]", intIdx, 0, v.Len()))
		}

		idx := v.Index(intIdx)
		fieldReflect = &idx
	}

	if fieldReflect != nil {
		if !fieldReflect.IsValid() {
			panic(fmt.Errorf("var error: object has no member %q", field))
		}
		if !fieldReflect.CanInterface() {
			panic(fmt.Errorf("var error: object member %q is inaccessible", field))
		}

		return fieldReflect.Interface()
	}

	panic(fmt.Errorf("syntax error: cannot access fields on type %s", typeOf(s)))
}

func slice(v interface{}, from, to interface{}) interface{} {
	str, isStr := v.(string)
	arr, isArr := v.([]interface{})

	if !isStr && !isArr {
		panic(fmt.Errorf("syntax error: slicing requires an array or string, but was %s", typeOf(v)))
	}

	var fromInt, toInt int
	if from == nil {
		fromInt = 0
	} else {
		fromInt = asInteger(from)
	}

	if to == nil && isStr {
		toInt = len(str)
	} else if to == nil && isArr {
		toInt = len(arr)
	} else {
		toInt = asInteger(to)
	}

	if fromInt < 0 {
		panic(fmt.Errorf("range error: start-index %d is negative", fromInt))
	}

	if isStr {
		if toInt < 0 || toInt > len(str) {
			panic(fmt.Errorf("range error: end-index %d is out of range [0, %d]", toInt, len(str)))
		}
		if fromInt > toInt {
			panic(fmt.Errorf("range error: start-index %d is greater than end-index %d", fromInt, toInt))
		}
		return str[fromInt:toInt]
	}

	if toInt < 0 || toInt > len(arr) {
		panic(fmt.Errorf("range error: end-index %d is out of range [0, %d]", toInt, len(arr)))
	}
	if fromInt > toInt {
		panic(fmt.Errorf("range error: start-index %d is greater than end-index %d", fromInt, toInt))
	}
	return arr[fromInt:toInt]
}

func arrayContains(arr interface{}, val interface{}) bool {
	a, ok := arr.([]interface{})
	if !ok {
		panic(fmt.Errorf("syntax error: in-operator requires array, but was %s", typeOf(arr)))
	}

	for _, v := range a {
		if deepEqual(v, val) {
			return true
		}
	}
	return false
}

func callFunction(functions map[string]ExpressionFunction, name string, args []interface{}) interface{} {
	f, ok := functions[name]
	if !ok {
		panic(fmt.Errorf("syntax error: no such function %q", name))
	}

	res, err := f(args...)
	if err != nil {
		panic(fmt.Errorf("function error: %q - %s", name, err))
	}
	return res
}
