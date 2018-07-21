goval
=====

[![Build Status](https://travis-ci.org/maja42/goval.svg?branch=master)](https://travis-ci.org/maja42/goval)
[![Go Report Card](https://goreportcard.com/badge/github.com/maja42/goval)](https://goreportcard.com/report/github.com/maja42/goval)
[![Coverage Status](https://coveralls.io/repos/github/maja42/goval/badge.svg?branch=master)](https://coveralls.io/github/maja42/goval?branch=master)
[![GoDoc](https://godoc.org/github.com/maja42/goval?status.svg)](https://godoc.org/github.com/maja42/goval)


This library allows programs to evaluate arbitrary arithmetic/string/logic expressions.
Custom extensions like variable access and function calls are supported.

Please use the issue-tracker for any questions, feedback, bug-reports and feature-requests. 

This project is licensed under the terms of the MIT license.

# Demo

A small demo that evaluates expressions given from stdin can be found in the example folder:

```
go get -u github.com/maja42/goval
cd $GOPATH/src/github.com/maja42/goval/
go run example/main.go
```

# Usage

Minimal example:

```go
eval := goval.NewEvaluator()
result, err := eval.Evaluate(`42 > 21`, nil, nil) // Returns <true, nil>
```

Accessing variables:
```go
eval := goval.NewEvaluator()
variables := map[string]interface{}{
    "os": runtime.GOOS,
    "arch": runtime.GOARCH,
}
result, err := eval.Evaluate(`os + "/" + arch`, variables, nil) // Returns <"linux/amd64", nil>
```

Calling functions:
```go
eval := goval.NewEvaluator()
variables := map[string]interface{}{
    "os":   runtime.GOOS,
    "arch": runtime.GOARCH,
}

functions := make(map[string]goval.ExpressionFunction)
functions["strlen"] = func(args ...interface{}) (interface{}, error) {
    str := args[0].(string)
    return len(str), nil
}

result, err := eval.Evaluate(`strlen(arch[:2]) + strlen("text")`, variables, functions) // Returns <6, nil>
```


# Alternative Libraries

If you are looking for a generic evaluation library, 
you can also take a look at [Knetic/govaluate](https://github.com/Knetic/govaluate).
I used that library myself, but due to a few shortcomings I decided to create goval. 
The main differences are:

- Full support for arrays and objects.
- Accessing variables (maps) via `.` and `[]` syntax
- Support for array- and object concatenation.
- Array literals with `[]` as well as object literals with `{}`
- Opaque differentiation between `int` and `float64`. \
  The underlying type is automatically converted as long as no precision is lost.
- Type-aware bit-operations (they only work with `int`-numbers).
- Hex-Literals (useful as soon as bit-operations are involved).
- No support for dates (strings are just strings, they don't have a special meaning, even if they look like dates).\
  Support for dates and structs *could* be added if needed.
- Useful error messages.
- Written with go/scanner and goyacc. \
    This vastly reduces code size (and therefore vulnerabilities to bugs),
    creates super-fast code and allows new features to be added in minutes, rather than days.
- High test coverage (including lots of special cases).\
  Also tested on 32 and 64bit architectures, where some (documented) operations like a bitwise-not can behave differently depending on the size of `int`. 

For a full list of features, please refer to the documentation below.

# Documentation

## Types

This library fully supports the following types: `nil`, `bool`, `int`, `float64`, `string`, `[]interface{}` (=arrays) and `map[string]interface{}` (=objects). 

Within expressions, `int` and `float64` both have the type `number` and are completely transparent.\
If necessary, numerical values will be automatically converted between `int` and `float64`, as long as no precision is lost.

Arrays and Objects are untyped. They can store any other value ("mixed arrays").

## Variables

It is possible to directly access custom-defined variables.
Variables are read-only and cannot be modified from within expressions.

Examples:

```
var
var.field
var[0]
var["field"]
var[anotherVar]

var["fie" + "ld"].field[42 - var2][0]
```

## Functions

It is possible to call custom-defined functions from within expressions.

Examples:

```
rand()
floor(42)
min(4, 3, 12, max(1, 3, 3))
len("te" + "xt")
```

## Literals

Any literal can be defined within expressions. 
String literals can be put in double-quotes `"` or back-ticks \`.
Hex-literals start with the prefix `0x`.

Examples:

```
nil
true
false
3
3.2
"Hello, 世界!\n"
"te\"xt"
`te"xt`
[0, 1, 2]
[]
[0, ["text", false], 4.2]
{}
{"a": 1, "b": {c: 3}}
{"key" + 42: "value"}
{"k" + "e" + "y": "value"}

0xA                 // 10
0x0A                // 10
0xFF                // 255 
0xFFFFFFFF          // 32bit appl.: -1  64bit appl.: 4294967295
0xFFFFFFFFFFFFFFFF  // 64bit appl.: -1  32bit appl.: error
```

It is possible to access elements of array and object literals:

Examples:

```
[1, 2, 3][1]                // 2
[1, [2, 3, 42][1][2]        // 42

{"a": 1}.a                  // 1
{"a": {"b": 42}}.a.b        // 42
{"a": {"b": 42}}["a"]["b"]  // 42
```

## Precedence

Operator precedence strictly follows [C/C++ rules](http://en.cppreference.com/w/cpp/language/operator_precedence).

Parenthesis `()` is used to control precedence.

Examples:

```
1 + 2 * 3    // 7
(1 + 2) * 3  // 9
```

## Operators

### Arithmetic

#### Arithmetic `+` `-` `*` `/`

If both sides are integers, the resulting value is also an integer.
Otherwise, the result will be a floating point number.

Examples:

```
3 + 4               // 7
2 + 2 * 3           // 8
2 * 3 + 2.5         // 8.5
12 - 7 - 5          // 0
24 / 10             // 2
24.0 / 10           // 2.4
```

#### Modulo `%`

If both sides are integers, the resulting value is also an integer.
Otherwise, the result will be a floating point number.

Examples:

```
4 % 3       // 1
144 % 85    // -55
5.5 % 2     // 1.5
10 % 3.5    // 3.0
```

#### Negation `-` (unary minus)

Negates the number on the right.

Examples:

```
-4       // -4
5 + -4   // 1
-5 - -4  // -1
1 + --1  // syntax error
-(4+3)   // -7
-varName
```


### Concatenation

#### String concatenation `+`

If either the left or right side of the `+` operator is a `string`, a string concatenation is performed.
Supports strings, numbers, booleans and nil.

Examples:

```
"text" + 42     // "text42"
"text" + 4.2    // "text4.2"
42 + "text"     // "42text"
"text" + nil    // "textnil"
"text" + true   // "texttrue"
```

#### Array concatenation `+`

If both sides of the `+` operator are arrays, they are concatenated

Examples:

```
[0, 1] + [2, 3]          // [0, 1, 2, 3]
[0] + [1] + [[2]] + []   // [0, 1, [2]]
```

#### Object concatenation `+`

If both sides of the `+` operator are objects, their fields are combined into a new object.
If both objects contain the same keys, the value of the right object will override those of the left.

Examples:

```
{"a": 1} + {"b": 2} + {"c": 3}         // {"a": 1, "b": 2, "c": 3}
{"a": 1, "b": 2} + {"b": 3, "c": 4}    // {"a": 1, "b": 3, "c": 4}
{"b": 3, "c": 4} + {"a": 1, "b": 2}    // {"a": 1, "b": 2, "c": 4}
```

### Logic

#### Equals `==`, NotEquals `!=`

Performs a deep-compare between the two operands.
When comparing `int` and `float64`, 
the integer will be cast to a floating point number.

#### Comparisons `<`, `>`, `<=`, `>=`

Compares two numbers. If one side of the operator is an integer and the other is a floating point number,
the integer number will be cast. This might lead to unexpected results for very big numbers which are rounded
during that process.

Examples:

```
3 <-4        // false
45 > 3.4     // false
-4 <= -1     // true
3.5 >= 3.5   // true
```

#### And `&&`, Or `||`

Examples:

```
true && true             // true
false || false           // false
true || false && false   // true
false && false || true   // true
```


#### Not `!`

Inverts the boolean on the right.

Examples:

```
!true       // false
!false      // true
!!true      // true
!varName
```

### Bit Manipulation

#### Logical Or `|`, Logical And `&`, Logical XOr `^`

If one side of the operator is a floating point number, the number is cast to an integer if possible. 
If decimal places would be lost during that process, it is considered a type error.
The resulting number is always an integer.

Examples:

```
8 | 2          // 10
9 | 5          // 13
8 | 2.0        // 10
8 | 2.1        // type error

13 & 10        // 8
10 & 15.0 & 2  // 2

13 ^ 10        // 7
10 ^ 15 ^ 1    // 4
```

#### Bitwise Not `~`

If performed on a floating point number, the number is cast to an integer if possible. 
If decimal places would be lost during that process, it is considered a type error.
The resulting number is always an integer.

The results can differ between 32bit and 64bit architectures.

Examples:

```
~-1                   // 0
(~0xA55A) & 0xFFFF    // 0x5AA5
(~0x5AA5) & 0xFFFF    // 0xA55A

~0xFFFFFFFF           // 64bit appl.: 0xFFFFFFFF 00000000; 32bit appl.: 0x00
~0xFFFFFFFF FFFFFFFF  // 64bit appl.: 0x00; 32bit: error
```

#### Bit-Shift `<<`, `>>`

If one side of the operator is a floating point number, the number is cast to an integer if possible. 
If decimal places would be lost during that process, it is considered a type error.
The resulting number is always an integer.

When shifting to the right, sign-extension is performed.
The results can differ between 32bit and 64bit architectures.

Examples:

```
1 << 0    // 1
1 << 1    // 2
1 << 2    // 4
8 << -1   // 4
8 >> -1   // 16

1 << 31   // 0x00000000 80000000   64bit appl.: 2147483648; 32bit appl.: -2147483648
1 << 32   // 0x00000001 00000000   32bit appl.: 0 (overflow)

1 << 63   // 0x80000000 00000000   32bit appl.: 0 (overflow); 64bit appl.: -9223372036854775808
1 << 64   // 0x00000000 00000000   0 (overflow)

0x80000000 00000000 >> 63     // 0xFFFFFFFF FFFFFFFF   64bit: -1 (sign extension); 32bit: error (cannot parse number literal)
0x80000000 >> 31              // 64bit: 0x00000000 0000001; 32bit: 0xFFFFFFFF (-1, sign extension)
```

### More

#### Array contains `in`

Returns true or false whether the array contains a specific element.

Examples:

```
"txt" in [nil, "hello", "txt", 42]   // true
true  in [nil, "hello", "txt", 42]   // false
nil   in [nil, "hello", "txt", 42]   // true
42.0  in [nil, "hello", "txt", 42]   // true
2         in [1, [2, 3], 4]          // false
[2, 3]    in [1, [2, 3], 4]          // true
[2, 3, 4] in [1, [2, 3], 4]          // false
```

#### Substrings `[a:b]`

Slices a string and returns the given substring.
Strings are indexed byte-wise. Multi-byte characters need to be treated carefully.

The start-index indicates the first byte to be present in the substring.\
The end-index indicates the last byte NOT to be present in the substring.\
Hence, valid indices are in the range `[0, len(str)]`.

Examples:

```
"abcdefg"[:]    // "abcdefg"
"abcdefg"[1:]   // "bcdefg"
"abcdefg"[:6]   // "abcdef"
"abcdefg"[2:5]  // "cde"
"abcdefg"[3:4]  // "d"

// The characters 世 and 界 both require 3 bytes:
"Hello, 世界"[7:13]    // "世界"
"Hello, 世界"[7:10]    // "世"
"Hello, 世界"[10:13]   // "界"
```


#### Array Slicing `[a:b]`

Slices an array and returns the given subarray.

The start-index indicates the first element to be present in the subarray.\
The end-index indicates the last element NOT to be present in the subarray.\
Hence, valid indices are in the range `[0, len(arr)]`.

Examples:

```
// Assuming `arr := [0, 1, 2, 3, 4, 5, 6]`:
arr[:]    // [0, 1, 2, 3, 4, 5, 6]
arr[1:]   // [1, 2, 3, 4, 5, 6]
arr[:6]   // [0, 1, 2, 3, 4, 5]
arr[2:5]  // [2, 3, 4]
arr[3:4]  // [3]
```



