package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {

	var depths depthArgs
	flag.Var(&depths, "depth", "recursion depth; individual numbers, ranges (2-5), and lists (2,5) supported. Default: 1,2")

	var prefix string
	flag.StringVar(&prefix, "prefix", "", "prefix string")

	var suffix string
	flag.StringVar(&suffix, "suffix", "", "suffix string")

	var seps sepArgs
	flag.Var(&seps, "sep", "separator string (can be specified multiple times)")

	var noRepeats bool
	flag.BoolVar(&noRepeats, "no-repeats", false, "use each line of input only once per sequence")

	flag.Parse()

	if len(depths) == 0 {
		depths = append(depths, 1, 2)
	}

	if len(seps) == 0 {
		seps = append(seps, "")
	}

	alphabet := make([]string, 0)

	sc := bufio.NewScanner(os.Stdin)

	for sc.Scan() {
		alphabet = append(alphabet, sc.Text())
	}

	p := &permutator{
		depths:    depths,
		seps:      seps,
		prefix:    prefix,
		suffix:    suffix,
		noRepeats: noRepeats,
	}

	p.list("", alphabet)

}

type permutator struct {
	depth     int
	depths    depthArgs
	seps      sepArgs
	prefix    string
	suffix    string
	noRepeats bool
}

func (p *permutator) list(context string, alphabet []string) {

	if p.shouldStop() {
		p.depth--
		return
	}

	for i, a := range alphabet {

		// make a copy of the alphabet so that we can potentially
		// remove the 'letter' in use if --no-repeats is specified
		newAlpha := make([]string, len(alphabet))
		copy(newAlpha, alphabet)
		if p.noRepeats {
			newAlpha = append(newAlpha[:i], newAlpha[i+1:]...)
		}

		for _, sep := range p.seps {

			if context == "" {
				sep = ""
				context = p.prefix
			}

			newPerm := context + sep + a

			if p.shouldOutput() {
				fmt.Println(newPerm + p.suffix)
			}

			p.depth++
			p.list(newPerm, newAlpha)
		}
	}
	p.depth--
	return

}

func (p *permutator) shouldOutput() bool {
	// The user provides a one-indexed list of depths, but we're dealing
	// with a zero-indexed list so we have to add 1 to the current depth
	if p.depths.includes(p.depth + 1) {
		return true
	}

	return false
}

func (p *permutator) shouldStop() bool {
	if p.depth == p.depths.max() {
		return true
	}

	return false
}

type depthArgs []int

func (d *depthArgs) Set(val string) error {

	// list
	if strings.ContainsRune(val, ',') {
		vals := strings.Split(val, ",")
		for _, v := range vals {
			i, err := strconv.Atoi(v)
			if err != nil {
				return err
			}
			*d = append(*d, i)
		}
		return nil
	}

	// range
	if strings.ContainsRune(val, '-') {
		parts := strings.SplitN(val, "-", 2)

		min, err := strconv.Atoi(parts[0])
		if err != nil {
			return err
		}

		max, err := strconv.Atoi(parts[1])
		if err != nil {
			return err
		}

		if min >= max {
			return errors.New("range maximum must be more than range minimum")
		}

		for i := min; i <= max; i++ {
			*d = append(*d, i)
		}

		return nil
	}

	// default to an individual number
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

type sepArgs []string

func (s *sepArgs) Set(val string) error {
	*s = append(*s, val)
	return nil
}

func (s *sepArgs) String() string {
	return strings.Join(*s, ",")
}
