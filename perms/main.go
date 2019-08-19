package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

func main() {

	var maxDepth int
	flag.IntVar(&maxDepth, "depth", 3, "maximum recursion depth")

	var prefix string
	flag.StringVar(&prefix, "prefix", "", "prefix string")

	var suffix string
	flag.StringVar(&suffix, "suffix", "", "suffix string")

	var sep string
	flag.StringVar(&sep, "sep", "", "separator string")

	flag.Parse()

	if maxDepth < 1 || maxDepth > 16 {
		fmt.Fprintln(os.Stderr, "depth can only be 1-16")
		return
	}

	alphabet := make([]string, 0)

	sc := bufio.NewScanner(os.Stdin)

	for sc.Scan() {
		alphabet = append(alphabet, sc.Text())
	}

	p := &permutator{
		maxDepth: maxDepth,
		alphabet: alphabet,
		sep:      sep,
		prefix:   prefix,
		suffix:   suffix,
	}

	perms := p.list("")

	for _, perm := range perms {
		fmt.Println(perm)
	}

}

type permutator struct {
	depth    int
	maxDepth int
	alphabet []string
	sep      string
	prefix   string
	suffix   string
}

func (p *permutator) list(context string) []string {
	out := make([]string, 0)

	if p.depth == p.maxDepth {
		p.depth--
		return out
	}

	sep := p.sep
	if context == "" {
		sep = ""
		context = p.prefix
	}

	for _, a := range p.alphabet {
		newPerm := context + sep + a
		out = append(out, newPerm+p.suffix)
		p.depth++
		out = append(out, p.list(newPerm)...)
	}
	p.depth--
	return out

}
