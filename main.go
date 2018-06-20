package main

import (
	"os"
	"fmt"
)

//go:generate goyacc.exe -o parser.go parser.go.y

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Invalid usage")
		os.Exit(1)
	}

	src := os.Args[1]

	// yyDebug = 5
	yyErrorVerbose = true

	lexer := NewLexer(src, nil, nil)
	yyParse(lexer)
	fmt.Printf("%#v\n", lexer.result)
}
