package main

import (
	"flag"
	"fmt"
	"os"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

func main() {
	flag.Parse()

	repoPath := "."
	if flag.Arg(0) != "" {
		repoPath = flag.Arg(0)
	}

	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open repo: %s", err)
		return
	}
	_ = repo

	tress, err := repo.TreeObjects()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get tree objects: %s", err)
		return
	}

	seen := make(map[string]bool)
	tress.ForEach(func(t *object.Tree) error {
		t.Files().ForEach(func(f *object.File) error {
			id := f.ID().String()

			if seen[id] {
				return nil
			}
			seen[id] = true

			fmt.Println(id, f.Name)

			return nil
		})

		return nil
	})
}
