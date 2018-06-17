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
- Type-aware bit-operations (they only work with `int`)
- No special handling of strings that look like dates
- Array literals with `[]` as well as object literals with `{}`
- Useful error messages
- High performance and highly extensible due to the use of yacc


  
 


# Types

This library fully supports the following types: `int`, `float64`, `bool`, `string`, `[]interface{}` (=arrays) and `map[string]interface{}` (=objects). 

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

It will be possible to call custom-defined functions (not implemented yet).


# Literals

Any literal can be defined within expressions. 
String literals can be put in double-quotes " or back-ticks `.

Examples:

```
true
false
3
3.2
"Hello, 世界!\n"
"te\"xt"
`te"xt`
```

It will also be possible to add array- and object literals (not implemented yet).

# Operators

## Arithmetic

### Addition, concatenation `+`

If either the left or right side is a `string`, a string concatenation is performed.
Otherwise, both sides need to be numbers and will be added. If both sides are integers, the resulting value is
also an integer. Otherwise, the result will be a floating point number. 

Examples:

```
3 + 4           // 7
3 + 4.0         // 7.0
3.2 + 1.4       // 4.6

"text" + 42     // "text42"
42 + "text"     // "42text"
"text" + true   // "textrue"
```

### Arithmetic `-` `*` `/`

If both sides are integers, the resulting value is also an integer.
Otherwise, the result will be a floating point number.

Examples:

```
2 + 2 * 3           // 8
2 * 3 + 2.5         // 8.5
12 - 7 - 5          // 0
24 / 10             // 2
24.0 / 10           // 2.4
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

## Logic

### Not `!`

Inverts the boolean on the right.

Examples:

```
!true       // false
!false      // true
!!true      // true
!varName
```
