package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"

	parser "github.com/Cgboal/DomainParser"
)

type node struct {
	val      string
	children map[string]*node
}

func newNode(val string) *node {
	return &node{
		val:      val,
		children: make(map[string]*node),
	}
}

func (n *node) add(c *node) {
	n.children[c.val] = c
}

func (n *node) addList(l []string) {
	leaf := n
	for _, part := range l {
		if _, exists := leaf.children[part]; !exists {
			leaf.children[part] = newNode(part)
		}
		leaf = leaf.children[part]
	}
}

func (n *node) tree(prefix string) string {
	out := strings.Builder{}

	if n.val != "." {
		out.WriteString(prefix)
		out.WriteString(n.val)
		out.WriteRune('\n')
		prefix = fmt.Sprintf("  %s", prefix)
	}

	for _, c := range n.children {
		out.WriteString(c.tree(prefix))
	}

	return out.String()
}

func main() {
	extractor := parser.NewDomainParser()

	sc := bufio.NewScanner(os.Stdin)

	root := newNode(".")

	for sc.Scan() {
		raw := strings.TrimSpace(sc.Text())

		tld := extractor.GetTld(raw)
		domain, _ := strings.CutSuffix(raw, fmt.Sprintf(".%s", tld))
		parts := strings.Split(domain, ".")
		slices.Reverse(parts)

		parts = append([]string{tld}, parts...)
		root.addList(parts)
	}

	fmt.Printf("%s", root.tree(""))
}
