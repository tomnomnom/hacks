package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"path"
	"strconv"
)

func main() {

	flag.Parse()

	dir := flag.Arg(0)
	if dir == "" {
		log.Fatal("No dir specified")
	}

	threshold := 1.0
	if t, err := strconv.ParseFloat(flag.Arg(1), 64); err == nil {
		threshold = t
	}

	contents, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatalf("failed to read dir: %s", err)
	}

	sizes := make([]int64, 0)

	for _, f := range contents {
		// don't bother with directories
		if f.IsDir() {
			continue
		}

		// see how different the file is to the sizes we've seen so far
		// start of assuming the file is different, but if it's within, say,
		// 1% of a filesize we've already seen then skip it
		isDifferent := true
		for _, s := range sizes {
			diff := math.Abs((float64(s-f.Size()) / float64(s)) * 100)
			if diff < threshold {
				isDifferent = false
			}
		}
		if isDifferent {
			sizes = append(sizes, f.Size())
			fmt.Println(path.Join(dir, f.Name()))
		}
	}
}
