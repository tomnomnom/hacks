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

	mode := flag.Arg(0)
	search := flag.Arg(1)
	if mode == "" || search == "" {
		fmt.Fprintf(os.Stderr, "usage: ghsearch <mode> <search>\n")
		return
	}

	client, err := getClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error getting client: %s\n", err)
		return
	}

	var results []string

	switch mode {
	case "search":
		results, err = getResults(client, search)
	case "repos":
		results, err = getRepos(client, search)
	case "members":
		results, err = getMembers(client, search)
	default:
		fmt.Fprint(os.Stderr, "unknown mode. Try repos, members, or search")
		return
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "error getting results: %s\n", err)
		// There might still be results we can print, so continue...
	}

	for _, result := range results {
		fmt.Println(result)
	}

}

func getResults(client *github.Client, search string) ([]string, error) {
	var allResults []string

	opt := &github.SearchOptions{
		ListOptions: github.ListOptions{PerPage: 20},
	}

	for {

		ctx := context.Background()
		results, resp, err := client.Search.Code(ctx, search, opt)
		if err != nil {
			return allResults, err
		}

		for _, result := range results.CodeResults {
			allResults = append(allResults, fmt.Sprintf("https://github.com/%s.git", *result.GetRepository().FullName))
		}

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return allResults, nil

}

func getRepos(client *github.Client, user string) ([]string, error) {
	var allRepos []string

	opt := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 20},
	}

	for {

		ctx := context.Background()
		repos, resp, err := client.Repositories.List(ctx, user, opt)
		if err != nil {
			return allRepos, err
		}

		for _, repo := range repos {
			if *repo.Fork {
				continue
			}
			allRepos = append(allRepos, *repo.CloneURL)
		}

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return allRepos, nil

}

func getMembers(client *github.Client, org string) ([]string, error) {
	var allMembers []string

	opt := &github.ListMembersOptions{
		ListOptions: github.ListOptions{PerPage: 20},
	}

	for {

		ctx := context.Background()
		members, resp, err := client.Organizations.ListMembers(ctx, org, opt)
		if err != nil {
			return allMembers, err
		}

		for _, member := range members {

			allMembers = append(allMembers, *member.Login)
		}

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return allMembers, nil

}
