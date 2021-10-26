package github_test

import (
	"context"
	"fmt"

	"github.com/gardenbed/go-github"
)

func ExampleClient_EnsureScopes() {
	client := github.NewClient("")
	if err := client.EnsureScopes(context.Background(), github.ScopeRepo); err != nil {
		panic(err)
	}
}

func ExampleUserService_Get() {
	client := github.NewClient("")
	user, resp, err := client.Users.Get(context.Background(), "octocat")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Rate: %+v\n\n", resp.Rate)
	fmt.Printf("Name: %s\n", user.Name)
}

func ExampleRepoService_Commits() {
	client := github.NewClient("")
	commits, resp, err := client.Repo("octocat", "Hello-World").Commits(context.Background(), 50, 1)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Pages: %+v\n", resp.Pages)
	fmt.Printf("Rate: %+v\n\n", resp.Rate)
	for _, c := range commits {
		fmt.Printf("SHA: %s\n", c.SHA)
	}
}

func ExampleIssueService_List() {
	client := github.NewClient("")
	issues, resp, err := client.Repo("octocat", "Hello-World").Issues.List(context.Background(), 50, 1, github.IssuesFilter{})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Pages: %+v\n", resp.Pages)
	fmt.Printf("Rate: %+v\n\n", resp.Rate)
	for _, i := range issues {
		fmt.Printf("Title: %s\n", i.Title)
	}
}

func ExamplePullService_List() {
	client := github.NewClient("")
	pull, resp, err := client.Repo("octocat", "Hello-World").Pulls.List(context.Background(), 50, 1, github.PullsFilter{})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Pages: %+v\n", resp.Pages)
	fmt.Printf("Rate: %+v\n\n", resp.Rate)
	for _, p := range pull {
		fmt.Printf("Title: %s\n", p.Title)
	}
}

func ExampleSearchService_SearchIssues() {
	client := github.NewClient("")

	query := github.SearchQuery{}
	query.IncludeKeywords("Fix")
	query.ExcludeKeywords("WIP")
	query.IncludeQualifiers(
		github.QualifierTypePR,
		github.QualifierInTitle,
		github.QualifierLabel("bug"),
	)

	result, resp, err := client.Search.SearchIssues(context.Background(), 20, 1, "", "", query)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Pages: %+v\n", resp.Pages)
	fmt.Printf("Rate: %+v\n\n", resp.Rate)
	for _, issue := range result.Items {
		fmt.Printf("%s\n", issue.HTMLURL)
	}
}
