package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
)

func main() {

	var depths depthArgs
	flag.Var(&depths, "depth", "depth of recursion to output (can be used multiple times)")

	var maxDepth int
	flag.IntVar(&maxDepth, "max-depth", 3, "maximum recursion depth (cannot be used with --depth)")

	var minDepth int
	flag.IntVar(&minDepth, "min-depth", 1, "minimum recursion depth (cannot be used with --depth)")

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

	if minDepth < 1 || minDepth > maxDepth {
		fmt.Fprintln(os.Stderr, "min-depth can only be 1 to max-depth")
		return
	}

	alphabet := make([]string, 0)

	sc := bufio.NewScanner(os.Stdin)

	for sc.Scan() {
		alphabet = append(alphabet, sc.Text())
	}

	p := &permutator{
		depths:    depths,
		maxDepth:  maxDepth,
		minDepth:  minDepth,
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
	depths    depthArgs
	maxDepth  int
	minDepth  int
	sep       string
	prefix    string
	suffix    string
	noRepeats bool
}

func (p *permutator) list(context string, alphabet []string) []string {
	out := make([]string, 0)

	if p.shouldStop() {
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

		if p.shouldOutput() {
			out = append(out, newPerm+p.suffix)
		}

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

func (p *permutator) shouldOutput() bool {
	// The list of depths takes priority over the min and max depths.
	// The user provides a one-indexed list of depths, but we're dealing
	// with a zero-indexed list so we have to add 1 to the current depth
	if len(p.depths) > 0 && p.depths.includes(p.depth+1) {
		return true
	}

	// if we got to here and there's at least one depth option set, don't output
	if len(p.depths) > 0 {
		return false
	}

	// we only need one last check for being past the minimum depth as the
	// max depth is used to stop the recusion and is handled in p.shouldStop()
	return p.depth >= (p.minDepth - 1)
}

func (p *permutator) shouldStop() bool {
	// if there are depth flags and we're at max depth it's time to stop
	if len(p.depths) > 0 && p.depth == p.depths.max() {
		return true
	}

	// if there are depth flags and we didn't stop above, keep going :)
	if len(p.depths) > 0 {
		return false
	}

	// fall back to checking the max depth
	return p.depth == p.maxDepth
}

type depthArgs []int

func (d *depthArgs) Set(val string) error {
	i, err := strconv.Atoi(val)
	if err != nil {
		return err
	}
	*d = append(*d, i)
	return nil
}

func (d depthArgs) String() string {
	return "string"
}

func (d depthArgs) includes(search int) bool {
	for _, candidate := range d {
		if candidate == search {
			return true
		}
	}
	return false
}

func (d depthArgs) max() int {
	max := 0
	for _, candidate := range d {
		if candidate > max {
			max = candidate
		}
	}
	return max
}
