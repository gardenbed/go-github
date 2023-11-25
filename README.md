[![Go Doc][godoc-image]][godoc-url]
[![Build Status][codeql-image]][codeql-url]
[![Build Status][workflow-image]][workflow-url]
[![Go Report Card][goreport-image]][goreport-url]
[![Test Coverage][codecov-image]][codecov-url]

# go-github

A simple Go client for [GitHub API v3](https://docs.github.com/rest).

## Quick Start

You can find more examples [here](./example).

```go
package main

import (
  "context"
  "fmt"

  "github.com/gardenbed/go-github"
)

func main() {
  client := github.NewClient("")
  commits, resp, err := client.Repo("octocat", "Hello-World").Commits(context.Background(), 50, 1)
  if err != nil {
    panic(err)
  }

  fmt.Printf("Pages: %+v\n", resp.Pages)
  fmt.Printf("Rate: %+v\n\n", resp.Rate)
  for _, commit := range commits {
    fmt.Printf("%s\n", commit.SHA)
  }
}
```


[godoc-url]: https://pkg.go.dev/github.com/gardenbed/go-github
[godoc-image]: https://pkg.go.dev/badge/github.com/gardenbed/go-github
[codeql-url]: https://github.com/gardenbed/basil-templates/actions/workflows/github-code-scanning/codeql
[codeql-image]: https://github.com/gardenbed/basil-templates/workflows/CodeQL/badge.svg
[workflow-url]: https://github.com/gardenbed/go-github/actions
[workflow-image]: https://github.com/gardenbed/go-github/workflows/Go/badge.svg
[goreport-url]: https://goreportcard.com/report/github.com/gardenbed/go-github
[goreport-image]: https://goreportcard.com/badge/github.com/gardenbed/go-github
[codecov-url]: https://codecov.io/gh/gardenbed/go-github
[codecov-image]: https://codecov.io/gh/gardenbed/go-github/branch/main/graph/badge.svg
