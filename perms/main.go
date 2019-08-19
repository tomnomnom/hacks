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

	var noRepeats bool
	flag.BoolVar(&noRepeats, "no-repeats", false, "use each line of input only once per sequence")

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
		maxDepth:  maxDepth,
		sep:       sep,
		prefix:    prefix,
		suffix:    suffix,
		noRepeats: noRepeats,
	}

	perms := p.list("", alphabet)

	for _, perm := range perms {
		fmt.Println(perm)
	}

}

type permutator struct {
	depth     int
	maxDepth  int
	sep       string
	prefix    string
	suffix    string
	noRepeats bool
}

func (p *permutator) list(context string, alphabet []string) []string {
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

	for i, a := range alphabet {
		newPerm := context + sep + a
		out = append(out, newPerm+p.suffix)

		newAlpha := make([]string, len(alphabet))
		copy(newAlpha, alphabet)
		if p.noRepeats {
			newAlpha = append(newAlpha[:i], newAlpha[i+1:]...)
		}

		p.depth++
		out = append(out, p.list(newPerm, newAlpha)...)
	}
	p.depth--
	return out

}
