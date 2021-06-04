package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ditashi/jsbeautifier-go/jsbeautifier"
)

func main() {
	sc := bufio.NewScanner(os.Stdin)

	for sc.Scan() {
		path := sc.Text()
		options := jsbeautifier.DefaultOptions()
		beautified := jsbeautifier.BeautifyFile(path, options)

		err := ioutil.WriteFile(path, []byte(*beautified), 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error writing file: %s\n", err)
		}
	}
}
