// +build linux darwin

package internal

//go:generate echo "Generating parser 'using golang.org/x/tools/cmd/goyacc...'"
//go:generate goyacc -o parser.go parser.go.y
