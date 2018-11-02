package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
)

type formater func(io.Reader) interface{}

func main() {
	var format string
	flag.StringVar(&format, "format", "array", "Output format to use (array, map, 2d-array)")
	flag.Parse()

	formats := map[string]formater{
		"array":    toArray,
		"map":      toMap,
		"2d-array": to2dArray,
	}

	f, ok := formats[format]
	if !ok {
		fmt.Fprintf(os.Stderr, "unknown format '%s'\n", format)
		return
	}

	out := f(os.Stdin)

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "    ")
	enc.Encode(out)
}

func toArray(r io.Reader) interface{} {
	sc := bufio.NewScanner(r)
	lines := make([]string, 0)
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	return lines
}

func to2dArray(r io.Reader) interface{} {
	sc := bufio.NewScanner(r)
	lines := make([][]string, 0)
	re := regexp.MustCompile("\\s+")
	for sc.Scan() {
		parts := re.Split(sc.Text(), -1)
		lines = append(lines, parts)
	}
	return lines
}

func toMap(r io.Reader) interface{} {
	sc := bufio.NewScanner(r)
	lines := make(map[string]string)
	re := regexp.MustCompile("\\s+")
	for sc.Scan() {
		parts := re.Split(sc.Text(), 2)
		key := parts[0]
		val := ""
		if len(parts) == 2 {
			val = parts[1]
		}
		lines[key] = val
	}
	return lines
}
