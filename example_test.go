package github_test

import (
	"context"
	"fmt"

	"github.com/gardenbed/go-github"
)

func ExampleClient_EnsureScopes() {
	c := github.NewClient("")

	err := c.EnsureScopes(context.Background(), github.ScopeRepo)
	if err != nil {
		panic(err)
	}
}

func ExampleUsersService_Get() {
	c := github.NewClient("")

	user, resp, err := c.Users.Get(context.Background(), "octocat")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Pages: %+v\n", resp.Pages)
	fmt.Printf("Rate: %+v\n\n", resp.Rate)
	fmt.Printf("Name: %s\n", user.Name)
}

func ExampleRepoService_Commits() {
	c := github.NewClient("")

	commits, resp, err := c.Repo("octocat", "Hello-World").Commits(context.Background(), 50, 1)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Pages: %+v\n", resp.Pages)
	fmt.Printf("Rate: %+v\n\n", resp.Rate)
	for _, c := range commits {
		fmt.Printf("%s\n", c.SHA)
	}
}
