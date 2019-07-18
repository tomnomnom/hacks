package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func getClient() (*github.Client, error) {
	ghToken := os.Getenv("GITHUB_TOKEN")
	if ghToken == "" {
		return nil, errors.New("GITHUB_TOKEN not set")
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: ghToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	return client, nil
}

func main() {

	flag.Parse()

	search := flag.Arg(0)
	if search == "" {
		fmt.Fprintf(os.Stderr, "usage: ghsearch <search>\n")
		return
	}

	client, err := getClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error getting client: %s\n", err)
		return
	}

	results, err := getResults(client, search)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error getting results: %s\n", err)
		// There might still be members we can print here...
	}

	for _, result := range results {

		fmt.Printf("https://github.com/%s.git\n", *result.GetRepository().FullName)
	}

}

func getResults(client *github.Client, search string) ([]github.CodeResult, error) {
	var allResults []github.CodeResult

	opt := &github.SearchOptions{
		ListOptions: github.ListOptions{PerPage: 20},
	}

	for {

		ctx := context.Background()
		results, resp, err := client.Search.Code(ctx, search, opt)
		if err != nil {
			return allResults, err
		}

		allResults = append(allResults, results.CodeResults...)

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return allResults, nil

}
