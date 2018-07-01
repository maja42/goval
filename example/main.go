package main

import (
	"bufio"
	"fmt"
	"github.com/maja42/goval"
	"math/rand"
	"os"
	"os/signal"
	"reflect"
	"runtime"
)

func main() {
	// Ctrl+C should exit the application...
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		fmt.Printf("\nExiting\n")
		os.Exit(0)
	}()

	// Create some custom variables...
	variables := make(map[string]interface{})
	variables["os"] = runtime.GOOS
	variables["arch"] = runtime.GOARCH

	// Create sum custom functions...
	functions := make(map[string]goval.ExpressionFunction)

	functions["rand"] = func(...interface{}) (interface{}, error) {
		return rand.Float64(), nil
	}

	functions["len"] = func(args ...interface{}) (interface{}, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("expected exactly 1 argument")
		}
		str, ok := args[0].(string)
		if ok {
			return len(str), nil
		}
		arr, ok := args[0].([]interface{})
		if ok {
			return len(arr), nil
		}
		obj, ok := args[0].(map[string]interface{})
		if ok {
			return len(obj), nil
		}
		return nil, fmt.Errorf("expected string, array or object")
	}

	// Evaluate:
	fmt.Print("Enter expressions to evaluate them.\n" +
		"Variables:\n" +
		"\tos      runtime.GOOS\n" +
		"\tarch    runtime.ARCH\n" +
		"\tans     result of the last evaluation\n" +
		"Functions:\n" +
		"\trand()  returns a random number between [0, 1[\n" +
		"\tlen()   returns the length of a string, array or object\n" +
		"Press Ctrl+C to exit\n\n")

	eval := goval.NewEvaluator()

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		result, err := eval.Evaluate(input, variables, functions)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		} else {
			fmt.Printf("%v (%s)\n\n", result, reflect.TypeOf(result))
			variables["ans"] = result
		}
	}
}
