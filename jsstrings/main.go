package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/robertkrimen/otto/ast"
	"github.com/robertkrimen/otto/parser"
)

type stringFinder struct {
}

func (w *stringFinder) Enter(n ast.Node) ast.Visitor {

	if id, ok := n.(*ast.StringLiteral); ok && id != nil {
		fmt.Println(id.Value)
	}

	return w
}

func (w *stringFinder) Exit(n ast.Node) {
}

func main() {
	flag.Parse()

	fn := flag.Arg(0)
	if fn == "" {
		fmt.Println("usage: jsstrings <filename>")
		return
	}

	program, err := parser.ParseFile(nil, fn, nil, 0)
	if err != nil {
		log.Fatal(err)
	}

	w := &stringFinder{}

	ast.Walk(w, program)
}
