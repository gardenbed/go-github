package github

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	repositoryBody = `{
		"id": 1296269,
		"name": "Hello-World",
		"full_name": "octocat/Hello-World",
		"owner": {
			"login": "octocat",
			"id": 1,
			"type": "User"
		},
		"private": false,
		"description": "This your first repo!",
		"fork": false,
		"default_branch": "main",
		"topics": [
			"octocat",
			"api"
		],
		"archived": false,
		"disabled": false,
		"visibility": "public",
		"pushed_at": "2020-10-31T14:00:00Z",
		"created_at": "2020-01-20T09:00:00Z",
		"updated_at": "2020-10-31T14:00:00Z"
	}`

	permissionBody = `{
		"permission": "admin",
		"user": {
			"login": "octocat",
			"id": 1,
			"type": "User"
		}
	}`

	commitBody1 = `{
		"sha": "6dcb09b5b57875f334f61aebed695e2e4193db5e",
		"commit": {
			"author": {
				"name": "The Octocat",
				"email": "octocat@github.com",
				"date": "2020-10-20T19:59:59Z"
			},
			"committer": {
				"name": "The Octocat",
				"email": "octocat@github.com",
				"date": "2020-10-20T19:59:59Z"
			},
			"message": "Fix all the bugs"
		},
		"author": {
			"login": "octocat",
			"id": 1,
			"type": "User"
		},
		"committer": {
			"login": "octocat",
			"id": 1,
			"type": "User"
		},
		"parents": [
			{
				"url": "https://api.github.com/repos/octocat/Hello-World/commits/c3d0be41ecbe669545ee3e94d31ed9a4bc91ee3c",
				"sha": "c3d0be41ecbe669545ee3e94d31ed9a4bc91ee3c"
			}
  	]
	}`

	commitsBody = `[
		{
			"sha": "c3d0be41ecbe669545ee3e94d31ed9a4bc91ee3c",
			"commit": {
				"author": {
					"name": "The Octocat",
					"email": "octocat@github.com",
					"date": "2020-10-27T23:59:59Z"
				},
				"committer": {
					"name": "The Octocat",
					"email": "octocat@github.com",
					"date": "2020-10-27T23:59:59Z"
				},
				"message": "Release v0.1.0"
			},
			"author": {
				"login": "octocat",
				"id": 1,
				"type": "User"
			},
			"committer": {
				"login": "octocat",
				"id": 1,
				"type": "User"
			}
		},
		{
			"sha": "6dcb09b5b57875f334f61aebed695e2e4193db5e",
			"commit": {
				"author": {
					"name": "The Octocat",
					"email": "octocat@github.com",
					"date": "2020-10-20T19:59:59Z"
				},
				"committer": {
					"name": "The Octocat",
					"email": "octocat@github.com",
					"date": "2020-10-20T19:59:59Z"
				},
				"message": "Fix all the bugs"
			},
			"author": {
				"login": "octocat",
				"id": 1,
				"type": "User"
			},
			"committer": {
				"login": "octocat",
				"id": 1,
				"type": "User"
			},
			"parents": [
				{
					"url": "https://api.github.com/repos/octocat/Hello-World/commits/c3d0be41ecbe669545ee3e94d31ed9a4bc91ee3c",
					"sha": "c3d0be41ecbe669545ee3e94d31ed9a4bc91ee3c"
				}
			]
		}
	]`

	branchBody = `{
		"name": "main",
		"commit": {
			"sha": "c3d0be41ecbe669545ee3e94d31ed9a4bc91ee3c",
			"commit": {
				"author": {
					"name": "The Octocat",
					"email": "octocat@github.com",
					"date": "2020-10-27T23:59:59Z"
				},
				"committer": {
					"name": "The Octocat",
					"email": "octocat@github.com",
					"date": "2020-10-27T23:59:59Z"
				},
				"message": "Release v0.1.0"
			},
			"author": {
				"login": "octocat",
				"id": 1,
				"type": "User"
			},
			"committer": {
				"login": "octocat",
				"id": 1,
				"type": "User"
			}
		},
		"protected": true
	}`

	tagsBody = `[
		{
			"name": "v0.1.0",
			"commit": {
				"sha": "c3d0be41ecbe669545ee3e94d31ed9a4bc91ee3c",
				"url": "https://api.github.com/repos/octocat/Hello-World/commits/c3d0be41ecbe669545ee3e94d31ed9a4bc91ee3c"
			}
		}
	]`

	releaseBody = `{
		"id": 1,
		"tag_name": "v1.0.0",
		"target_commitish": "main",
		"name": "v1.0.0",
		"body": "Description of the release",
		"draft": false,
		"prerelease": false,
		"author": {
			"login": "octocat",
			"id": 1,
			"type": "User"
		},
		"assets": [
			{
				"id": 1,
				"name": "example.zip",
				"label": "short description",
				"state": "uploaded",
				"content_type": "application/zip",
				"size": 1024,
				"uploader": {
					"login": "octocat",
					"id": 1,
					"type": "User"
				}
			}
		]
	}`

	releaseAssetBody = `{
		"id": 1,
		"name": "example.zip",
		"label": "short description",
		"state": "uploaded",
		"content_type": "application/zip",
		"size": 1024,
		"uploader": {
			"id": 1,
			"login": "octocat",
			"type": "User"
		}
	}`
)

var (
	permission = PermissionAdmin

	repository = Repository{
		ID:            1296269,
		Name:          "Hello-World",
		FullName:      "octocat/Hello-World",
		Description:   "This your first repo!",
		Topics:        []string{"octocat", "api"},
		Private:       false,
		Fork:          false,
		Archived:      false,
		Disabled:      false,
		DefaultBranch: "main",
		Owner: User{
			ID:    1,
			Login: "octocat",
			Type:  "User",
		},
		CreatedAt: parseGitHubTime("2020-01-20T09:00:00Z"),
		UpdatedAt: parseGitHubTime("2020-10-31T14:00:00Z"),
		PushedAt:  parseGitHubTime("2020-10-31T14:00:00Z"),
	}

	commit1 = Commit{
		SHA: "6dcb09b5b57875f334f61aebed695e2e4193db5e",
		Commit: RawCommit{
			Message: "Fix all the bugs",
			Author: Signature{
				Name:  "The Octocat",
				Email: "octocat@github.com",
				Time:  parseGitHubTime("2020-10-20T19:59:59Z"),
			},
			Committer: Signature{
				Name:  "The Octocat",
				Email: "octocat@github.com",
				Time:  parseGitHubTime("2020-10-20T19:59:59Z"),
			},
		},
		Author: User{
			ID:    1,
			Login: "octocat",
			Type:  "User",
		},
		Committer: User{
			ID:    1,
			Login: "octocat",
			Type:  "User",
		},
		Parents: []Hash{
			{
				SHA: "c3d0be41ecbe669545ee3e94d31ed9a4bc91ee3c",
				URL: "https://api.github.com/repos/octocat/Hello-World/commits/c3d0be41ecbe669545ee3e94d31ed9a4bc91ee3c",
			},
		},
	}

	commit2 = Commit{
		SHA: "c3d0be41ecbe669545ee3e94d31ed9a4bc91ee3c",
		Commit: RawCommit{
			Message: "Release v0.1.0",
			Author: Signature{
				Name:  "The Octocat",
				Email: "octocat@github.com",
				Time:  parseGitHubTime("2020-10-27T23:59:59Z"),
			},
			Committer: Signature{
				Name:  "The Octocat",
				Email: "octocat@github.com",
				Time:  parseGitHubTime("2020-10-27T23:59:59Z"),
			},
		},
		Author: User{
			ID:    1,
			Login: "octocat",
			Type:  "User",
		},
		Committer: User{
			ID:    1,
			Login: "octocat",
			Type:  "User",
		},
	}

	branch = Branch{
		Name:      "main",
		Protected: true,
		Commit:    commit2,
	}

	tag = Tag{
		Name: "v0.1.0",
		Commit: Hash{
			SHA: "c3d0be41ecbe669545ee3e94d31ed9a4bc91ee3c",
			URL: "https://api.github.com/repos/octocat/Hello-World/commits/c3d0be41ecbe669545ee3e94d31ed9a4bc91ee3c",
		},
	}

	release = Release{
		ID:         1,
		Name:       "v1.0.0",
		TagName:    "v1.0.0",
		Target:     "main",
		Draft:      false,
		Prerelease: false,
		Body:       "Description of the release",
		Author: User{
			ID:    1,
			Login: "octocat",
			Type:  "User",
		},
		Assets: []ReleaseAsset{
			{
				ID:          1,
				Name:        "example.zip",
				Label:       "short description",
				State:       "uploaded",
				ContentType: "application/zip",
				Size:        1024,
				Uploader: User{
					ID:    1,
					Login: "octocat",
					Type:  "User",
				},
			},
		},
	}

	releaseAsset = ReleaseAsset{
		ID:          1,
		Name:        "example.zip",
		Label:       "short description",
		State:       "uploaded",
		ContentType: "application/zip",
		Size:        1024,
		Uploader: User{
			ID:    1,
			Login: "octocat",
			Type:  "User",
		},
	}
)

func TestRepoService_Get(t *testing.T) {
	c := &Client{
		httpClient: &http.Client{},
		rates:      map[rateGroup]Rate{},
		apiURL:     publicAPIURL,
	}

	tests := []struct {
		name               string
		mockResponses      []MockResponse
		s                  *RepoService
		ctx                context.Context
		expectedRepository *Repository
		expectedResponse   *Response
		expectedError      string
	}{
		{
			name:          "NilContext",
			mockResponses: []MockResponse{},
			s: &RepoService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:           nil,
			expectedError: `net/http: nil Context`,
		},
		{
			name: "InvalidStatusCode",
			mockResponses: []MockResponse{
				{"GET", "/repos/octocat/Hello-World", 401, http.Header{}, `{
					"message": "Bad credentials"
				}`},
			},
			s: &RepoService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:           context.Background(),
			expectedError: `GET /repos/octocat/Hello-World: 401 Bad credentials`,
		},
		{
			name: "ّInvalidResponse",
			mockResponses: []MockResponse{
				{"GET", "/repos/octocat/Hello-World", 200, http.Header{}, `[`},
			},
			s: &RepoService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:           context.Background(),
			expectedError: `unexpected EOF`,
		},
		{
			name: "Success",
			mockResponses: []MockResponse{
				{"GET", "/repos/octocat/Hello-World", 200, header, repositoryBody},
			},
			s: &RepoService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:                context.Background(),
			expectedRepository: &repository,
			expectedResponse: &Response{
				Rate: expectedRate,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ts := newHTTPTestServer(tc.mockResponses...)
			tc.s.client.apiURL, _ = url.Parse(ts.URL)

			repository, resp, err := tc.s.Get(tc.ctx)

			if tc.expectedError != "" {
				assert.Nil(t, repository)
				assert.Nil(t, resp)
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedRepository, repository)
				assert.NotNil(t, resp)
				assert.NotNil(t, resp.Response)
				assert.Equal(t, tc.expectedResponse.Rate, resp.Rate)
			}
		})
	}
}

func TestRepoService_Permission(t *testing.T) {
	c := &Client{
		httpClient: &http.Client{},
		rates:      map[rateGroup]Rate{},
		apiURL:     publicAPIURL,
	}

	tests := []struct {
		name               string
		mockResponses      []MockResponse
		s                  *RepoService
		ctx                context.Context
		username           string
		expectedPermission Permission
		expectedResponse   *Response
		expectedError      string
	}{
		{
			name:          "NilContext",
			mockResponses: []MockResponse{},
			s: &RepoService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:           nil,
			username:      "octocat",
			expectedError: `net/http: nil Context`,
		},
		{
			name: "InvalidStatusCode",
			mockResponses: []MockResponse{
				{"GET", "/repos/octocat/Hello-World/collaborators/octocat/permission", 401, http.Header{}, `{
					"message": "Bad credentials"
				}`},
			},
			s: &RepoService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:           context.Background(),
			username:      "octocat",
			expectedError: `GET /repos/octocat/Hello-World/collaborators/octocat/permission: 401 Bad credentials`,
		},
		{
			name: "ّInvalidResponse",
			mockResponses: []MockResponse{
				{"GET", "/repos/octocat/Hello-World/collaborators/octocat/permission", 200, http.Header{}, `[`},
			},
			s: &RepoService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:           context.Background(),
			username:      "octocat",
			expectedError: `unexpected EOF`,
		},
		{
			name: "Success",
			mockResponses: []MockResponse{
				{"GET", "/repos/octocat/Hello-World/collaborators/octocat/permission", 200, header, permissionBody},
			},
			s: &RepoService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:                context.Background(),
			username:           "octocat",
			expectedPermission: permission,
			expectedResponse: &Response{
				Rate: expectedRate,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ts := newHTTPTestServer(tc.mockResponses...)
			tc.s.client.apiURL, _ = url.Parse(ts.URL)

			permission, resp, err := tc.s.Permission(tc.ctx, tc.username)

			if tc.expectedError != "" {
				assert.Empty(t, permission)
				assert.Nil(t, resp)
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedPermission, permission)
				assert.NotNil(t, resp)
				assert.NotNil(t, resp.Response)
				assert.Equal(t, tc.expectedResponse.Rate, resp.Rate)
			}
		})
	}
}

func TestRepoService_Commit(t *testing.T) {
	c := &Client{
		httpClient: &http.Client{},
		rates:      map[rateGroup]Rate{},
		apiURL:     publicAPIURL,
	}

	tests := []struct {
		name             string
		mockResponses    []MockResponse
		s                *RepoService
		ctx              context.Context
		ref              string
		expectedCommit   *Commit
		expectedResponse *Response
		expectedError    string
	}{
		{
			name:          "NilContext",
			mockResponses: []MockResponse{},
			s: &RepoService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:           nil,
			ref:           "6dcb09b5b57875f334f61aebed695e2e4193db5e",
			expectedError: `net/http: nil Context`,
		},
		{
			name: "InvalidStatusCode",
			mockResponses: []MockResponse{
				{"GET", "/repos/octocat/Hello-World/commits/6dcb09b5b57875f334f61aebed695e2e4193db5e", 401, http.Header{}, `{
					"message": "Bad credentials"
				}`},
			},
			s: &RepoService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:           context.Background(),
			ref:           "6dcb09b5b57875f334f61aebed695e2e4193db5e",
			expectedError: `GET /repos/octocat/Hello-World/commits/6dcb09b5b57875f334f61aebed695e2e4193db5e: 401 Bad credentials`,
		},
		{
			name: "ّInvalidResponse",
			mockResponses: []MockResponse{
				{"GET", "/repos/octocat/Hello-World/commits/6dcb09b5b57875f334f61aebed695e2e4193db5e", 200, http.Header{}, `{`},
			},
			s: &RepoService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:           context.Background(),
			ref:           "6dcb09b5b57875f334f61aebed695e2e4193db5e",
			expectedError: `unexpected EOF`,
		},
		{
			name: "Success",
			mockResponses: []MockResponse{
				{"GET", "/repos/octocat/Hello-World/commits/6dcb09b5b57875f334f61aebed695e2e4193db5e", 200, header, commitBody1},
			},
			s: &RepoService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:            context.Background(),
			ref:            "6dcb09b5b57875f334f61aebed695e2e4193db5e",
			expectedCommit: &commit1,
			expectedResponse: &Response{
				Rate: expectedRate,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ts := newHTTPTestServer(tc.mockResponses...)
			tc.s.client.apiURL, _ = url.Parse(ts.URL)

			commit, resp, err := tc.s.Commit(tc.ctx, tc.ref)

			if tc.expectedError != "" {
				assert.Nil(t, commit)
				assert.Nil(t, resp)
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedCommit, commit)
				assert.NotNil(t, resp)
				assert.NotNil(t, resp.Response)
				assert.Equal(t, tc.expectedResponse.Rate, resp.Rate)
			}
		})
	}
}

func TestRepoService_Commits(t *testing.T) {
	c := &Client{
		httpClient: &http.Client{},
		rates:      map[rateGroup]Rate{},
		apiURL:     publicAPIURL,
	}

	tests := []struct {
		name             string
		mockResponses    []MockResponse
		s                *RepoService
		ctx              context.Context
		pageSize         int
		pageNo           int
		expectedCommits  []Commit
		expectedResponse *Response
		expectedError    string
	}{
		{
			name:          "NilContext",
			mockResponses: []MockResponse{},
			s: &RepoService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:           nil,
			pageSize:      10,
			pageNo:        1,
			expectedError: `net/http: nil Context`,
		},
		{
			name: "InvalidStatusCode",
			mockResponses: []MockResponse{
				{"GET", "/repos/octocat/Hello-World/commits", 401, http.Header{}, `{
					"message": "Bad credentials"
				}`},
			},
			s: &RepoService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:           context.Background(),
			pageSize:      10,
			pageNo:        1,
			expectedError: `GET /repos/octocat/Hello-World/commits: 401 Bad credentials`,
		},
		{
			name: "ّInvalidResponse",
			mockResponses: []MockResponse{
				{"GET", "/repos/octocat/Hello-World/commits", 200, http.Header{}, `[`},
			},
			s: &RepoService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:           context.Background(),
			pageSize:      10,
			pageNo:        1,
			expectedError: `unexpected EOF`,
		},
		{
			name: "Success",
			mockResponses: []MockResponse{
				{"GET", "/repos/octocat/Hello-World/commits", 200, header, commitsBody},
			},
			s: &RepoService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:             context.Background(),
			pageSize:        10,
			pageNo:          1,
			expectedCommits: []Commit{commit2, commit1},
			expectedResponse: &Response{
				Pages: expectedPages,
				Rate:  expectedRate,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ts := newHTTPTestServer(tc.mockResponses...)
			tc.s.client.apiURL, _ = url.Parse(ts.URL)

			commits, resp, err := tc.s.Commits(tc.ctx, tc.pageSize, tc.pageNo)

			if tc.expectedError != "" {
				assert.Nil(t, commits)
				assert.Nil(t, resp)
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedCommits, commits)
				assert.NotNil(t, resp)
				assert.NotNil(t, resp.Response)
				assert.Equal(t, tc.expectedResponse.Pages, resp.Pages)
				assert.Equal(t, tc.expectedResponse.Rate, resp.Rate)
			}
		})
	}
}

func TestRepoService_Branch(t *testing.T) {
	c := &Client{
		httpClient: &http.Client{},
		rates:      map[rateGroup]Rate{},
		apiURL:     publicAPIURL,
	}

	tests := []struct {
		name             string
		mockResponses    []MockResponse
		s                *RepoService
		ctx              context.Context
		branchName       string
		expectedBranch   *Branch
		expectedResponse *Response
		expectedError    string
	}{
		{
			name:          "NilContext",
			mockResponses: []MockResponse{},
			s: &RepoService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:           nil,
			branchName:    "main",
			expectedError: `net/http: nil Context`,
		},
		{
			name: "InvalidStatusCode",
			mockResponses: []MockResponse{
				{"GET", "/repos/octocat/Hello-World/branches/main", 401, http.Header{}, `{
					"message": "Bad credentials"
				}`},
			},
			s: &RepoService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:           context.Background(),
			branchName:    "main",
			expectedError: `GET /repos/octocat/Hello-World/branches/main: 401 Bad credentials`,
		},
		{
			name: "ّInvalidResponse",
			mockResponses: []MockResponse{
				{"GET", "/repos/octocat/Hello-World/branches/main", 200, http.Header{}, `{`},
			},
			s: &RepoService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:           context.Background(),
			branchName:    "main",
			expectedError: `unexpected EOF`,
		},
		{
			name: "Success",
			mockResponses: []MockResponse{
				{"GET", "/repos/octocat/Hello-World/branches/main", 200, header, branchBody},
			},
			s: &RepoService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:            context.Background(),
			branchName:     "main",
			expectedBranch: &branch,
			expectedResponse: &Response{
				Rate: expectedRate,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ts := newHTTPTestServer(tc.mockResponses...)
			tc.s.client.apiURL, _ = url.Parse(ts.URL)

			branch, resp, err := tc.s.Branch(tc.ctx, tc.branchName)

			if tc.expectedError != "" {
				assert.Nil(t, branch)
				assert.Nil(t, resp)
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedBranch, branch)
				assert.NotNil(t, resp)
				assert.NotNil(t, resp.Response)
				assert.Equal(t, tc.expectedResponse.Rate, resp.Rate)
			}
		})
	}
}

func TestRepoService_BranchProtection(t *testing.T) {
	c := &Client{
		httpClient: &http.Client{},
		rates:      map[rateGroup]Rate{},
		apiURL:     publicAPIURL,
	}

	tests := []struct {
		name             string
		mockResponses    []MockResponse
		s                *RepoService
		ctx              context.Context
		branch           string
		enabled          bool
		expectedResponse *Response
		expectedError    string
	}{
		{
			name:          "NilContext",
			mockResponses: []MockResponse{},
			s: &RepoService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:           nil,
			branch:        "main",
			enabled:       true,
			expectedError: `net/http: nil Context`,
		},
		{
			name: "InvalidStatusCode",
			mockResponses: []MockResponse{
				{"POST", "/repos/octocat/Hello-World/branches/main/protection/enforce_admins", 401, http.Header{}, `{
					"message": "Bad credentials"
				}`},
			},
			s: &RepoService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:           context.Background(),
			branch:        "main",
			enabled:       true,
			expectedError: `POST /repos/octocat/Hello-World/branches/main/protection/enforce_admins: 401 Bad credentials`,
		},
		{
			name: "Success_Enable",
			mockResponses: []MockResponse{
				{"POST", "/repos/octocat/Hello-World/branches/main/protection/enforce_admins", 200, header, ``},
			},
			s: &RepoService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:     context.Background(),
			branch:  "main",
			enabled: true,
			expectedResponse: &Response{
				Rate: expectedRate,
			},
		},
		{
			name: "Success_Disable",
			mockResponses: []MockResponse{
				{"DELETE", "/repos/octocat/Hello-World/branches/main/protection/enforce_admins", 204, header, ``},
			},
			s: &RepoService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:     context.Background(),
			branch:  "main",
			enabled: false,
			expectedResponse: &Response{
				Rate: expectedRate,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ts := newHTTPTestServer(tc.mockResponses...)
			tc.s.client.apiURL, _ = url.Parse(ts.URL)

			resp, err := tc.s.BranchProtection(tc.ctx, tc.branch, tc.enabled)

			if tc.expectedError != "" {
				assert.Nil(t, resp)
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.NotNil(t, resp.Response)
				assert.Equal(t, tc.expectedResponse.Rate, resp.Rate)
			}
		})
	}
}

func TestRepoService_Tags(t *testing.T) {
	c := &Client{
		httpClient: &http.Client{},
		rates:      map[rateGroup]Rate{},
		apiURL:     publicAPIURL,
	}

	tests := []struct {
		name             string
		mockResponses    []MockResponse
		s                *RepoService
		ctx              context.Context
		pageSize         int
		pageNo           int
		expectedTags     []Tag
		expectedResponse *Response
		expectedError    string
	}{
		{
			name:          "NilContext",
			mockResponses: []MockResponse{},
			s: &RepoService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:           nil,
			pageSize:      10,
			pageNo:        1,
			expectedError: `net/http: nil Context`,
		},
		{
			name: "InvalidStatusCode",
			mockResponses: []MockResponse{
				{"GET", "/repos/octocat/Hello-World/tags", 401, http.Header{}, `{
					"message": "Bad credentials"
				}`},
			},
			s: &RepoService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:           context.Background(),
			pageSize:      10,
			pageNo:        1,
			expectedError: `GET /repos/octocat/Hello-World/tags: 401 Bad credentials`,
		},
		{
			name: "ّInvalidResponse",
			mockResponses: []MockResponse{
				{"GET", "/repos/octocat/Hello-World/tags", 200, http.Header{}, `[`},
			},
			s: &RepoService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:           context.Background(),
			pageSize:      10,
			pageNo:        1,
			expectedError: `unexpected EOF`,
		},
		{
			name: "Success",
			mockResponses: []MockResponse{
				{"GET", "/repos/octocat/Hello-World/tags", 200, header, tagsBody},
			},
			s: &RepoService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:          context.Background(),
			pageSize:     10,
			pageNo:       1,
			expectedTags: []Tag{tag},
			expectedResponse: &Response{
				Pages: expectedPages,
				Rate:  expectedRate,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ts := newHTTPTestServer(tc.mockResponses...)
			tc.s.client.apiURL, _ = url.Parse(ts.URL)

			tags, resp, err := tc.s.Tags(tc.ctx, tc.pageSize, tc.pageNo)

			if tc.expectedError != "" {
				assert.Nil(t, tags)
				assert.Nil(t, resp)
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedTags, tags)
				assert.NotNil(t, resp)
				assert.NotNil(t, resp.Response)
				assert.Equal(t, tc.expectedResponse.Pages, resp.Pages)
				assert.Equal(t, tc.expectedResponse.Rate, resp.Rate)
			}
		})
	}
}

func TestRepoService_LatestRelease(t *testing.T) {
	c := &Client{
		httpClient: &http.Client{},
		rates:      map[rateGroup]Rate{},
		apiURL:     publicAPIURL,
	}

	tests := []struct {
		name             string
		mockResponses    []MockResponse
		s                *RepoService
		ctx              context.Context
		expectedRelease  *Release
		expectedResponse *Response
		expectedError    string
	}{
		{
			name:          "NilContext",
			mockResponses: []MockResponse{},
			s: &RepoService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:           nil,
			expectedError: `net/http: nil Context`,
		},
		{
			name: "InvalidStatusCode",
			mockResponses: []MockResponse{
				{"GET", "/repos/octocat/Hello-World/releases/latest", 401, http.Header{}, `{
					"message": "Bad credentials"
				}`},
			},
			s: &RepoService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:           context.Background(),
			expectedError: `GET /repos/octocat/Hello-World/releases/latest: 401 Bad credentials`,
		},
		{
			name: "ّInvalidResponse",
			mockResponses: []MockResponse{
				{"GET", "/repos/octocat/Hello-World/releases/latest", 200, http.Header{}, `{`},
			},
			s: &RepoService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:           context.Background(),
			expectedError: `unexpected EOF`,
		},
		{
			name: "Success",
			mockResponses: []MockResponse{
				{"GET", "/repos/octocat/Hello-World/releases/latest", 200, header, releaseBody},
			},
			s: &RepoService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:             context.Background(),
			expectedRelease: &release,
			expectedResponse: &Response{
				Rate: expectedRate,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ts := newHTTPTestServer(tc.mockResponses...)
			tc.s.client.apiURL, _ = url.Parse(ts.URL)

			release, resp, err := tc.s.LatestRelease(tc.ctx)

			if tc.expectedError != "" {
				assert.Nil(t, release)
				assert.Nil(t, resp)
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedRelease, release)
				assert.NotNil(t, resp)
				assert.NotNil(t, resp.Response)
				assert.Equal(t, tc.expectedResponse.Rate, resp.Rate)
			}
		})
	}
}

func TestRepoService_CreateRelease(t *testing.T) {
	c := &Client{
		httpClient: &http.Client{},
		rates:      map[rateGroup]Rate{},
		apiURL:     publicAPIURL,
	}

	params := ReleaseParams{
		Name:       "v1.0.0",
		TagName:    "v1.0.0",
		Target:     "main",
		Draft:      false,
		Prerelease: false,
		Body:       "Description of the release",
	}

	tests := []struct {
		name             string
		mockResponses    []MockResponse
		s                *RepoService
		ctx              context.Context
		params           ReleaseParams
		expectedRelease  *Release
		expectedResponse *Response
		expectedError    string
	}{
		{
			name:          "NilContext",
			mockResponses: []MockResponse{},
			s: &RepoService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:           nil,
			params:        params,
			expectedError: `net/http: nil Context`,
		},
		{
			name: "InvalidStatusCode",
			mockResponses: []MockResponse{
				{"POST", "/repos/octocat/Hello-World/releases", 401, http.Header{}, `{
					"message": "Bad credentials"
				}`},
			},
			s: &RepoService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:           context.Background(),
			params:        params,
			expectedError: `POST /repos/octocat/Hello-World/releases: 401 Bad credentials`,
		},
		{
			name: "ّInvalidResponse",
			mockResponses: []MockResponse{
				{"POST", "/repos/octocat/Hello-World/releases", 201, http.Header{}, `{`},
			},
			s: &RepoService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:           context.Background(),
			params:        params,
			expectedError: `unexpected EOF`,
		},
		{
			name: "Success",
			mockResponses: []MockResponse{
				{"POST", "/repos/octocat/Hello-World/releases", 201, header, releaseBody},
			},
			s: &RepoService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:             context.Background(),
			params:          params,
			expectedRelease: &release,
			expectedResponse: &Response{
				Rate: expectedRate,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ts := newHTTPTestServer(tc.mockResponses...)
			tc.s.client.apiURL, _ = url.Parse(ts.URL)

			release, resp, err := tc.s.CreateRelease(tc.ctx, tc.params)

			if tc.expectedError != "" {
				assert.Nil(t, release)
				assert.Nil(t, resp)
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedRelease, release)
				assert.NotNil(t, resp)
				assert.NotNil(t, resp.Response)
				assert.Equal(t, tc.expectedResponse.Rate, resp.Rate)
			}
		})
	}
}

func TestRepoService_UpdateRelease(t *testing.T) {
	c := &Client{
		httpClient: &http.Client{},
		rates:      map[rateGroup]Rate{},
		apiURL:     publicAPIURL,
	}

	params := ReleaseParams{
		Name:       "v1.0.0",
		TagName:    "v1.0.0",
		Target:     "main",
		Draft:      false,
		Prerelease: false,
		Body:       "Description of the release",
	}

	tests := []struct {
		name             string
		mockResponses    []MockResponse
		s                *RepoService
		ctx              context.Context
		releaseID        int
		params           ReleaseParams
		expectedRelease  *Release
		expectedResponse *Response
		expectedError    string
	}{
		{
			name:          "NilContext",
			mockResponses: []MockResponse{},
			s: &RepoService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:           nil,
			releaseID:     1,
			params:        params,
			expectedError: `net/http: nil Context`,
		},
		{
			name: "InvalidStatusCode",
			mockResponses: []MockResponse{
				{"PATCH", "/repos/octocat/Hello-World/releases/1", 401, http.Header{}, `{
					"message": "Bad credentials"
				}`},
			},
			s: &RepoService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:           context.Background(),
			releaseID:     1,
			params:        params,
			expectedError: `PATCH /repos/octocat/Hello-World/releases/1: 401 Bad credentials`,
		},
		{
			name: "ّInvalidResponse",
			mockResponses: []MockResponse{
				{"PATCH", "/repos/octocat/Hello-World/releases/1", 200, http.Header{}, `{`},
			},
			s: &RepoService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:           context.Background(),
			releaseID:     1,
			params:        params,
			expectedError: `unexpected EOF`,
		},
		{
			name: "Success",
			mockResponses: []MockResponse{
				{"PATCH", "/repos/octocat/Hello-World/releases/1", 200, header, releaseBody},
			},
			s: &RepoService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:             context.Background(),
			releaseID:       1,
			params:          params,
			expectedRelease: &release,
			expectedResponse: &Response{
				Rate: expectedRate,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ts := newHTTPTestServer(tc.mockResponses...)
			tc.s.client.apiURL, _ = url.Parse(ts.URL)

			release, resp, err := tc.s.UpdateRelease(tc.ctx, tc.releaseID, tc.params)

			if tc.expectedError != "" {
				assert.Nil(t, release)
				assert.Nil(t, resp)
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedRelease, release)
				assert.NotNil(t, resp)
				assert.NotNil(t, resp.Response)
				assert.Equal(t, tc.expectedResponse.Rate, resp.Rate)
			}
		})
	}
}

func TestRepoService_UploadReleaseAsset(t *testing.T) {
	c := &Client{
		httpClient: &http.Client{},
		rates:      map[rateGroup]Rate{},
		uploadURL:  publicUploadURL,
	}

	tests := []struct {
		name                 string
		mockResponses        []MockResponse
		s                    *RepoService
		ctx                  context.Context
		releaseID            int
		assetFile            string
		assetLabel           string
		expectedReleaseAsset *ReleaseAsset
		expectedResponse     *Response
		expectedError        string
	}{
		{
			name:          "NilContext",
			mockResponses: []MockResponse{},
			s: &RepoService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:           nil,
			releaseID:     1,
			assetFile:     "repo_test.go",
			assetLabel:    "test",
			expectedError: `net/http: nil Context`,
		},
		{
			name:          "NoFile",
			mockResponses: []MockResponse{},
			s: &RepoService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:           context.Background(),
			releaseID:     1,
			assetFile:     "unknown",
			assetLabel:    "test",
			expectedError: `open unknown: no such file or directory`,
		},
		{
			name:          "BadFile",
			mockResponses: []MockResponse{},
			s: &RepoService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:           context.Background(),
			releaseID:     1,
			assetFile:     "/dev/null",
			assetLabel:    "test",
			expectedError: `EOF`,
		},
		{
			name: "InvalidStatusCode",
			mockResponses: []MockResponse{
				{"POST", "/repos/octocat/Hello-World/releases/1/assets", 401, http.Header{}, `{
					"message": "Bad credentials"
				}`},
			},
			s: &RepoService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:           context.Background(),
			releaseID:     1,
			assetFile:     "repo_test.go",
			assetLabel:    "test",
			expectedError: `POST /repos/octocat/Hello-World/releases/1/assets: 401 Bad credentials`,
		},
		{
			name: "ّInvalidResponse",
			mockResponses: []MockResponse{
				{"POST", "/repos/octocat/Hello-World/releases/1/assets", 201, http.Header{}, `{`},
			},
			s: &RepoService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:           context.Background(),
			releaseID:     1,
			assetFile:     "repo_test.go",
			assetLabel:    "test",
			expectedError: `unexpected EOF`,
		},
		{
			name: "Success",
			mockResponses: []MockResponse{
				{"POST", "/repos/octocat/Hello-World/releases/1/assets", 201, header, releaseAssetBody},
			},
			s: &RepoService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:                  context.Background(),
			releaseID:            1,
			assetFile:            "repo_test.go",
			assetLabel:           "test",
			expectedReleaseAsset: &releaseAsset,
			expectedResponse: &Response{
				Rate: expectedRate,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ts := newHTTPTestServer(tc.mockResponses...)
			tc.s.client.uploadURL, _ = url.Parse(ts.URL)

			asset, resp, err := tc.s.UploadReleaseAsset(tc.ctx, tc.releaseID, tc.assetFile, tc.assetLabel)

			if tc.expectedError != "" {
				assert.Nil(t, asset)
				assert.Nil(t, resp)
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedReleaseAsset, asset)
				assert.NotNil(t, resp)
				assert.NotNil(t, resp.Response)
				assert.Equal(t, tc.expectedResponse.Rate, resp.Rate)
			}
		})
	}
}

func TestRepoService_DownloadReleaseAsset(t *testing.T) {
	c := &Client{
		httpClient:  &http.Client{},
		rates:       map[rateGroup]Rate{},
		downloadURL: publicDownloadURL,
	}

	tests := []struct {
		name             string
		mockResponses    []MockResponse
		s                *RepoService
		ctx              context.Context
		releaseTag       string
		assetName        string
		w                io.Writer
		expectedResponse *Response
		expectedError    string
	}{
		{
			name:          "NilContext",
			mockResponses: []MockResponse{},
			s: &RepoService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:           nil,
			releaseTag:    "v1.0.0",
			assetName:     "example.zip",
			w:             nil,
			expectedError: `net/http: nil Context`,
		},
		{
			name: "InvalidStatusCode",
			mockResponses: []MockResponse{
				{"GET", "/octocat/Hello-World/releases/download/v1.0.0/example.zip", 401, http.Header{}, ``},
			},
			s: &RepoService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:           context.Background(),
			releaseTag:    "v1.0.0",
			assetName:     "example.zip",
			w:             ioutil.Discard,
			expectedError: `GET /octocat/Hello-World/releases/download/v1.0.0/example.zip: 401 `,
		},
		{
			name: "Success",
			mockResponses: []MockResponse{
				{"GET", "/octocat/Hello-World/releases/download/v1.0.0/example.zip", 200, header, `content`},
			},
			s: &RepoService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:        context.Background(),
			releaseTag: "v1.0.0",
			assetName:  "example.zip",
			w:          ioutil.Discard,
			expectedResponse: &Response{
				Rate: expectedRate,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ts := newHTTPTestServer(tc.mockResponses...)
			tc.s.client.downloadURL, _ = url.Parse(ts.URL)

			resp, err := tc.s.DownloadReleaseAsset(tc.ctx, tc.releaseTag, tc.assetName, tc.w)

			if tc.expectedError != "" {
				assert.Nil(t, resp)
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.NotNil(t, resp.Response)
				assert.Equal(t, tc.expectedResponse.Rate, resp.Rate)
			}
		})
	}
}

func TestRepoService_DownloadTarArchive(t *testing.T) {
	c := &Client{
		httpClient: &http.Client{},
		rates:      map[rateGroup]Rate{},
		apiURL:     publicAPIURL,
	}

	tests := []struct {
		name             string
		mockResponses    []MockResponse
		s                *RepoService
		ctx              context.Context
		ref              string
		w                io.Writer
		expectedResponse *Response
		expectedError    string
	}{
		{
			name:          "NilContext",
			mockResponses: []MockResponse{},
			s: &RepoService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:           nil,
			ref:           "main",
			w:             nil,
			expectedError: `net/http: nil Context`,
		},
		{
			name: "InvalidStatusCode",
			mockResponses: []MockResponse{
				{"GET", "/repos/octocat/Hello-World/tarball/main", 401, http.Header{}, ``},
			},
			s: &RepoService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:           context.Background(),
			ref:           "main",
			w:             ioutil.Discard,
			expectedError: `GET /repos/octocat/Hello-World/tarball/main: 401 `,
		},
		{
			name: "Success",
			mockResponses: []MockResponse{
				{"GET", "/repos/octocat/Hello-World/tarball/main", 200, header, `content`},
			},
			s: &RepoService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx: context.Background(),
			ref: "main",
			w:   ioutil.Discard,
			expectedResponse: &Response{
				Rate: expectedRate,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ts := newHTTPTestServer(tc.mockResponses...)
			tc.s.client.apiURL, _ = url.Parse(ts.URL)

			resp, err := tc.s.DownloadTarArchive(tc.ctx, tc.ref, tc.w)

			if tc.expectedError != "" {
				assert.Nil(t, resp)
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.NotNil(t, resp.Response)
				assert.Equal(t, tc.expectedResponse.Rate, resp.Rate)
			}
		})
	}
}

func TestRepoService_DownloadZipArchive(t *testing.T) {
	c := &Client{
		httpClient: &http.Client{},
		rates:      map[rateGroup]Rate{},
		apiURL:     publicAPIURL,
	}

	tests := []struct {
		name             string
		mockResponses    []MockResponse
		s                *RepoService
		ctx              context.Context
		ref              string
		w                io.Writer
		expectedResponse *Response
		expectedError    string
	}{
		{
			name:          "NilContext",
			mockResponses: []MockResponse{},
			s: &RepoService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:           nil,
			ref:           "main",
			w:             nil,
			expectedError: `net/http: nil Context`,
		},
		{
			name: "InvalidStatusCode",
			mockResponses: []MockResponse{
				{"GET", "/repos/octocat/Hello-World/zipball/main", 401, http.Header{}, ``},
			},
			s: &RepoService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:           context.Background(),
			ref:           "main",
			w:             ioutil.Discard,
			expectedError: `GET /repos/octocat/Hello-World/zipball/main: 401 `,
		},
		{
			name: "Success",
			mockResponses: []MockResponse{
				{"GET", "/repos/octocat/Hello-World/zipball/main", 200, header, `content`},
			},
			s: &RepoService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx: context.Background(),
			ref: "main",
			w:   ioutil.Discard,
			expectedResponse: &Response{
				Rate: expectedRate,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ts := newHTTPTestServer(tc.mockResponses...)
			tc.s.client.apiURL, _ = url.Parse(ts.URL)

			resp, err := tc.s.DownloadZipArchive(tc.ctx, tc.ref, tc.w)

			if tc.expectedError != "" {
				assert.Nil(t, resp)
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.NotNil(t, resp.Response)
				assert.Equal(t, tc.expectedResponse.Rate, resp.Rate)
			}
		})
	}
}
