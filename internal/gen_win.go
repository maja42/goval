// +build windows

package internal

//go:generate echo "Generating parser 'using golang.org/x/tools/cmd/goyacc...'"
//go:generate goyacc.exe -o parser.go parser.go.y
