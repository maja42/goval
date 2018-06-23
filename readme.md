goval
=====

This library is currently under development and not ready for production.

Many features are still missing and will be added in the future. 

**What's the difference to Knetic/govaluate?**

If you are looking for a generic evaluation library, 
you probably know about [Knetic/govaluate](https://github.com/Knetic/govaluate).

However, this library has several shortcomings and limitations, which cannot be worked around. 
That's how this library was born. The main differences are:

- Full support for arrays and objects
- Accessing variables (maps) via `.` and `[]` syntax
- Differentiation between `int` and `float64`
- Support for array- and object concatenation
- Array literals with `[]` as well as object literals with `{}`
- Type-aware bit-operations (they only work with `int`)
- No special handling of strings that look like dates
- Useful error messages
- High performance and highly extensible due to the use of yacc
- High test coverage


# Types

This library fully supports the following types: `bool`, `int`, `float64`, `string`, `[]interface{}` (=arrays) and `map[string]interface{}` (=objects). 

If necessary, numerical values will be automatically converted between `int` and `float64`, as long as no precision is lost.

Arrays and Objects are untyped. They can store any other value ("mixed arrays").

# Variables

It is possible directly access custom-defined variables.
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

# Functions

It is possible to call custom-defined functions from within expressions.

Examples:

```
rand()
floor(42)
min(4, 3, 12, max(1, 3, 3))
len("te" + "xt")
```

# Literals

Any literal can be defined within expressions. 
String literals can be put in double-quotes `"` or back-ticks \`.
Hex-literals start with the prefix `0x`.

Examples:

```
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

# Precedence

Operator precedence strictly follows [C/C++ rules](http://en.cppreference.com/w/cpp/language/operator_precedence).

Parenthesis `()` is used to control precedence.

Examples:

```
1 + 2 * 3    // 7
(1 + 2) * 3  // 9
```

# Operators


## Arithmetic

### Arithmetic `+` `-` `*` `/`

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

### Modulo `%`

If both sides are integers, the resulting value is also an integer.
Otherwise, the result will be a floating point number.

Examples:

```
4 % 3       // 1
144 % 85    // -55
5.5 % 2     // 1.5
10 % 3.5    // 3.0
```

### Negation `-` (unary minus)

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


## Concatenation

### String concatenation `+`

If either the left or right side of the `+` operator is a `string`, a string concatenation is performed.
Supports strings, numbers and booleans.

Examples:

```
"text" + 42     // "text42"
42 + "text"     // "42text"
"text" + true   // "texttrue"
```

### Array concatenation `+`

If both sides of the `+` operator are arrays, they are concatenated

Examples:

```
[0, 1] + [2, 3]          // [0, 1, 2, 3]
[0] + [1] + [[2]] + []   // [0, 1, [2]]
```

### Object concatenation `+`

If both sides of the `+` operator are objects, their fields are combined into a new object.
If both objects contain the same keys, the value of the right object will override those of the left.

Examples:

```
{"a": 1} + {"b": 2} + {"c": 3}         // {"a": 1, "b": 2, "c": 3}
{"a": 1, "b": 2} + {"b": 3, "c": 4}    // {"a": 1, "b": 3, "c": 4}
{"b": 3, "c": 4} + {"a": 1, "b": 2}    // {"a": 1, "b": 2, "c": 4}
```

## Logic

### Equals `==`, NotEquals `!=`

Performs a deep-compare between the two operands.
When comparing `int` and `float64`, 
the integer will be casted to a floating point number.

### Comparisons `<`, `>`, `<=`, `>=`

Compares two numbers. If one side of the operator is an integer and the other is a floating point number,
the integer number will be casted. This might lead to unexpected results for very big numbers which are rounded
during that process.

Examples:

```
3 <-4        // false
45 > 3.4     // false
-4 <= -1     // true
3.5 >= 3.5   // true
```

### And `&&`, Or `||`

Examples:

```
true && true             // true
false || false           // false
true || false && false   // true
false && false || true   // true
```


### Not `!`

Inverts the boolean on the right.

Examples:

```
!true       // false
!false      // true
!!true      // true
!varName
```

## Bit Manipulation

### Logical Or `|`, Logical And `&`, Logical XOr `^`

If one side of the operator is a floating point number, the number is casted to an integer if possible. 
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
