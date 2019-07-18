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

	user := flag.Arg(0)
	if user == "" {
		fmt.Fprintf(os.Stderr, "usage: ghmembers <org>\n")
		return
	}

	client, err := getClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error getting client: %s\n")
		return
	}

	members, err := getMembers(client, user)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error getting members: %s\n")
		// There might still be members we can print here...
	}

	for _, member := range members {

		fmt.Println(*member.Login)
	}

}

func getMembers(client *github.Client, org string) ([]*github.User, error) {
	var allUsers []*github.User

	opt := &github.ListMembersOptions{
		ListOptions: github.ListOptions{PerPage: 20},
	}

	for {

		ctx := context.Background()
		repos, resp, err := client.Organizations.ListMembers(ctx, org, opt)
		if err != nil {
			return allUsers, fmt.Errorf("failed to list repos: %s", err)
		}

		allUsers = append(allUsers, repos...)

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return allUsers, nil

}
