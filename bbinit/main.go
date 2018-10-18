package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/machinebox/graphql"
)

type scopesReponse struct {
	Query struct {
		Program program
	} `json:"query"`
}

type program struct {
	Name    string `json:"handle"`
	InScope struct {
		Assets []struct {
			Asset asset
		}
	}
	OutOfScope struct {
		Assets []struct {
			Asset asset
		}
	}
}

type asset struct {
	Identifier string `json:"asset_identifier"`
	Type       string `json:"asset_type"`
}

func (a asset) Domain() (string, error) {
	if a.Type != "URL" {
		return "", fmt.Errorf("asset with identifier %s is not a URL", a.Identifier)
	}

	i := strings.ToLower(a.Identifier)

	// If it has a scheme then parse it and return the hostname
	if a.hasScheme() {
		u, err := url.Parse(i)
		if err == nil {
			return u.Hostname(), nil
		}
	}

	if a.isWildcard() {
		return strings.TrimLeft(i, "*.%"), nil
	}

	return i, nil
}

func (a asset) isWildcard() bool {
	if len(a.Identifier) < 2 {
		return false
	}

	if a.Identifier[0] == '*' {
		return true
	}

	if a.Identifier[0] == '.' {
		return true
	}

	if a.Identifier[0] == '%' {
		return true
	}

	return false
}

func (a asset) hasScheme() bool {
	i := strings.ToLower(a.Identifier)
	if len(i) < 6 {
		return false
	}

	if i[:5] == "http:" {
		return true
	}

	if len(i) < 7 {
		return false
	}

	if i[:6] == "https:" {
		return true
	}

	return false
}

func main() {

	var risky bool
	flag.BoolVar(&risky, "risky", false, "treat all domains as wildcards")

	var appendScope bool
	flag.BoolVar(&appendScope, "append-scope", false, "append to the scope file instead of replacing it")

	flag.Parse()

	graphQLToken := os.Getenv("H1_GRAPHQL_TOKEN")
	if graphQLToken == "" {
		fmt.Println("H1_GRAPHQL_TOKEN not set. Go to https://hackerone.com/current_user/graphql_token.json to get one")
		return
	}

	scopesQuery := `
		query Team_assets($first_0:Int! $handle:String!) {
			query {
				id,
				...F0
			}
		}
		fragment F0 on Query {
			me {
				Membership:membership(team_handle:$handle) {
					permissions,
					id
				},
				id
			},
			Program:team(handle:$handle) {
				handle,
				_structured_scope_versions2ZWKHQ:structured_scope_versions(archived:false) {
					max_updated_at
				},
				InScope:structured_scopes(first:$first_0,archived:false,eligible_for_submission:true) {
					Assets:edges {
						Asset:node {
							id,
							asset_type,
							asset_identifier,
							rendered_instruction,
							max_severity,
							eligible_for_bounty
						},
						cursor
					},
					pageInfo {
						hasNextPage,
						hasPreviousPage
					}
				},
				OutOfScope:structured_scopes(first:$first_0,archived:false,eligible_for_submission:false) {
					Assets:edges {
						Asset:node {
							id,
							asset_type,
							asset_identifier,
							rendered_instruction
						},
						cursor
					},
					pageInfo {
						hasNextPage,
						hasPreviousPage
					}
				},
				id
			},
		id
		}
	`

	client := graphql.NewClient("https://hackerone.com/graphql")

	req := graphql.NewRequest(scopesQuery)

	req.Var("first_0", 250)
	req.Var("handle", flag.Arg(0))

	req.Header.Set("X-Auth-Token", graphQLToken)

	ctx := context.Background()

	var respData scopesReponse
	if err := client.Run(ctx, req, &respData); err != nil {
		log.Fatal(err)
	}

	// trucate the scopes file by default
	scopeFlags := os.O_TRUNC | os.O_CREATE | os.O_WRONLY
	if appendScope {
		scopeFlags = os.O_APPEND | os.O_CREATE | os.O_WRONLY
	}

	// scopes file
	sf, err := os.OpenFile(".scope", scopeFlags, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer sf.Close()

	// domains file
	df, err := os.OpenFile("domains", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer df.Close()

	// wildcards file
	wf, err := os.OpenFile("wildcards", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer wf.Close()

	for _, aa := range respData.Query.Program.InScope.Assets {
		a := aa.Asset
		d, err := a.Domain()
		if err != nil {
			continue
		}

		if a.isWildcard() || risky {
			// scope: subdomain
			fmt.Fprintf(sf, ".*\\.%s$\n", d)

			// wildcard file
			fmt.Fprintf(wf, "%s\n", d)
		}

		// scope: exact match
		fmt.Fprintf(sf, "^%s$\n", d)

		// domains file
		fmt.Fprintf(df, "%s\n", d)
	}

	for _, aa := range respData.Query.Program.OutOfScope.Assets {
		a := aa.Asset
		d, err := a.Domain()
		if err != nil {
			continue
		}

		if a.isWildcard() {
			// subdomain
			fmt.Fprintf(sf, "!.*\\.%s$\n", d)
		}

		// exact match
		fmt.Fprintf(sf, "!^%s$\n", d)
	}
}
