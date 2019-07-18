package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func main() {
	ghToken := os.Getenv("GITHUB_TOKEN")
	if ghToken == "" {
		fmt.Fprintf(os.Stderr, "GITHUB_TOKEN not set!\n")
		return
	}

	flag.Parse()

	user := flag.Arg(0)
	if user == "" {
		fmt.Fprintf(os.Stderr, "usage: ghrepos <username>\n")
		return
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: ghToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	opt := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 20},
	}

	var allRepos []*github.Repository

	for {

		repos, resp, err := client.Repositories.List(ctx, user, opt)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to list repos: %s", err)
			break
		}

		allRepos = append(allRepos, repos...)

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	for _, repo := range allRepos {
		if *repo.Fork {
			continue
		}

		fmt.Println(*repo.CloneURL)
	}

}
