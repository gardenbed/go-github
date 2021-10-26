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

func TestReleaseService_Latest(t *testing.T) {
	c := &Client{
		httpClient: &http.Client{},
		rates:      map[rateGroup]Rate{},
		apiURL:     publicAPIURL,
	}

	tests := []struct {
		name             string
		mockResponses    []MockResponse
		s                *ReleaseService
		ctx              context.Context
		expectedRelease  *Release
		expectedResponse *Response
		expectedError    string
	}{
		{
			name:          "NilContext",
			mockResponses: []MockResponse{},
			s: &ReleaseService{
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
			s: &ReleaseService{
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
			s: &ReleaseService{
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
			s: &ReleaseService{
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

			release, resp, err := tc.s.Latest(tc.ctx)

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

func TestReleaseService_Get(t *testing.T) {
	c := &Client{
		httpClient: &http.Client{},
		rates:      map[rateGroup]Rate{},
		apiURL:     publicAPIURL,
	}

	tests := []struct {
		name             string
		mockResponses    []MockResponse
		s                *ReleaseService
		ctx              context.Context
		id               int
		expectedRelease  *Release
		expectedResponse *Response
		expectedError    string
	}{
		{
			name:          "NilContext",
			mockResponses: []MockResponse{},
			s: &ReleaseService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:           nil,
			id:            1,
			expectedError: `net/http: nil Context`,
		},
		{
			name: "InvalidStatusCode",
			mockResponses: []MockResponse{
				{"GET", "/repos/octocat/Hello-World/releases/1", 401, http.Header{}, `{
					"message": "Bad credentials"
				}`},
			},
			s: &ReleaseService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:           context.Background(),
			id:            1,
			expectedError: `GET /repos/octocat/Hello-World/releases/1: 401 Bad credentials`,
		},
		{
			name: "ّInvalidResponse",
			mockResponses: []MockResponse{
				{"GET", "/repos/octocat/Hello-World/releases/1", 200, http.Header{}, `{`},
			},
			s: &ReleaseService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:           context.Background(),
			id:            1,
			expectedError: `unexpected EOF`,
		},
		{
			name: "Success",
			mockResponses: []MockResponse{
				{"GET", "/repos/octocat/Hello-World/releases/1", 200, header, releaseBody},
			},
			s: &ReleaseService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:             context.Background(),
			id:              1,
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

			release, resp, err := tc.s.Get(tc.ctx, tc.id)

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

func TestReleaseService_GetByTag(t *testing.T) {
	c := &Client{
		httpClient: &http.Client{},
		rates:      map[rateGroup]Rate{},
		apiURL:     publicAPIURL,
	}

	tests := []struct {
		name             string
		mockResponses    []MockResponse
		s                *ReleaseService
		ctx              context.Context
		tag              string
		expectedRelease  *Release
		expectedResponse *Response
		expectedError    string
	}{
		{
			name:          "NilContext",
			mockResponses: []MockResponse{},
			s: &ReleaseService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:           nil,
			tag:           "v1.0.0",
			expectedError: `net/http: nil Context`,
		},
		{
			name: "InvalidStatusCode",
			mockResponses: []MockResponse{
				{"GET", "/repos/octocat/Hello-World/releases/tags/v1.0.0", 401, http.Header{}, `{
					"message": "Bad credentials"
				}`},
			},
			s: &ReleaseService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:           context.Background(),
			tag:           "v1.0.0",
			expectedError: `GET /repos/octocat/Hello-World/releases/tags/v1.0.0: 401 Bad credentials`,
		},
		{
			name: "ّInvalidResponse",
			mockResponses: []MockResponse{
				{"GET", "/repos/octocat/Hello-World/releases/tags/v1.0.0", 200, http.Header{}, `{`},
			},
			s: &ReleaseService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:           context.Background(),
			tag:           "v1.0.0",
			expectedError: `unexpected EOF`,
		},
		{
			name: "Success",
			mockResponses: []MockResponse{
				{"GET", "/repos/octocat/Hello-World/releases/tags/v1.0.0", 200, header, releaseBody},
			},
			s: &ReleaseService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:             context.Background(),
			tag:             "v1.0.0",
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

			release, resp, err := tc.s.GetByTag(tc.ctx, tc.tag)

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

func TestReleaseService_Create(t *testing.T) {
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
		s                *ReleaseService
		ctx              context.Context
		params           ReleaseParams
		expectedRelease  *Release
		expectedResponse *Response
		expectedError    string
	}{
		{
			name:          "NilContext",
			mockResponses: []MockResponse{},
			s: &ReleaseService{
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
			s: &ReleaseService{
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
			s: &ReleaseService{
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
			s: &ReleaseService{
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

			release, resp, err := tc.s.Create(tc.ctx, tc.params)

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

func TestReleaseService_Update(t *testing.T) {
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
		s                *ReleaseService
		ctx              context.Context
		id               int
		params           ReleaseParams
		expectedRelease  *Release
		expectedResponse *Response
		expectedError    string
	}{
		{
			name:          "NilContext",
			mockResponses: []MockResponse{},
			s: &ReleaseService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:           nil,
			id:            1,
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
			s: &ReleaseService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:           context.Background(),
			id:            1,
			params:        params,
			expectedError: `PATCH /repos/octocat/Hello-World/releases/1: 401 Bad credentials`,
		},
		{
			name: "ّInvalidResponse",
			mockResponses: []MockResponse{
				{"PATCH", "/repos/octocat/Hello-World/releases/1", 200, http.Header{}, `{`},
			},
			s: &ReleaseService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:           context.Background(),
			id:            1,
			params:        params,
			expectedError: `unexpected EOF`,
		},
		{
			name: "Success",
			mockResponses: []MockResponse{
				{"PATCH", "/repos/octocat/Hello-World/releases/1", 200, header, releaseBody},
			},
			s: &ReleaseService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:             context.Background(),
			id:              1,
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

			release, resp, err := tc.s.Update(tc.ctx, tc.id, tc.params)

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

func TestReleaseService_Delete(t *testing.T) {
	c := &Client{
		httpClient: &http.Client{},
		rates:      map[rateGroup]Rate{},
		apiURL:     publicAPIURL,
	}

	tests := []struct {
		name             string
		mockResponses    []MockResponse
		s                *ReleaseService
		ctx              context.Context
		id               int
		expectedResponse *Response
		expectedError    string
	}{
		{
			name:          "NilContext",
			mockResponses: []MockResponse{},
			s: &ReleaseService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:           nil,
			id:            1,
			expectedError: `net/http: nil Context`,
		},
		{
			name: "InvalidStatusCode",
			mockResponses: []MockResponse{
				{"DELETE", "/repos/octocat/Hello-World/releases/1", 401, http.Header{}, `{
					"message": "Bad credentials"
				}`},
			},
			s: &ReleaseService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:           context.Background(),
			id:            1,
			expectedError: `DELETE /repos/octocat/Hello-World/releases/1: 401 Bad credentials`,
		},
		{
			name: "Success",
			mockResponses: []MockResponse{
				{"DELETE", "/repos/octocat/Hello-World/releases/1", 204, header, ``},
			},
			s: &ReleaseService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx: context.Background(),
			id:  1,
			expectedResponse: &Response{
				Rate: expectedRate,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ts := newHTTPTestServer(tc.mockResponses...)
			tc.s.client.apiURL, _ = url.Parse(ts.URL)

			resp, err := tc.s.Delete(tc.ctx, tc.id)

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

func TestReleaseService_UploadAsset(t *testing.T) {
	c := &Client{
		httpClient: &http.Client{},
		rates:      map[rateGroup]Rate{},
		uploadURL:  publicUploadURL,
	}

	tests := []struct {
		name                 string
		mockResponses        []MockResponse
		s                    *ReleaseService
		ctx                  context.Context
		id                   int
		assetFile            string
		assetLabel           string
		expectedReleaseAsset *ReleaseAsset
		expectedResponse     *Response
		expectedError        string
	}{
		{
			name:          "NilContext",
			mockResponses: []MockResponse{},
			s: &ReleaseService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:           nil,
			id:            1,
			assetFile:     "repo_test.go",
			assetLabel:    "test",
			expectedError: `net/http: nil Context`,
		},
		{
			name:          "NoFile",
			mockResponses: []MockResponse{},
			s: &ReleaseService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:           context.Background(),
			id:            1,
			assetFile:     "unknown",
			assetLabel:    "test",
			expectedError: `open unknown: no such file or directory`,
		},
		{
			name:          "BadFile",
			mockResponses: []MockResponse{},
			s: &ReleaseService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:           context.Background(),
			id:            1,
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
			s: &ReleaseService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:           context.Background(),
			id:            1,
			assetFile:     "repo_test.go",
			assetLabel:    "test",
			expectedError: `POST /repos/octocat/Hello-World/releases/1/assets: 401 Bad credentials`,
		},
		{
			name: "ّInvalidResponse",
			mockResponses: []MockResponse{
				{"POST", "/repos/octocat/Hello-World/releases/1/assets", 201, http.Header{}, `{`},
			},
			s: &ReleaseService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:           context.Background(),
			id:            1,
			assetFile:     "repo_test.go",
			assetLabel:    "test",
			expectedError: `unexpected EOF`,
		},
		{
			name: "Success",
			mockResponses: []MockResponse{
				{"POST", "/repos/octocat/Hello-World/releases/1/assets", 201, header, releaseAssetBody},
			},
			s: &ReleaseService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:                  context.Background(),
			id:                   1,
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

			asset, resp, err := tc.s.UploadAsset(tc.ctx, tc.id, tc.assetFile, tc.assetLabel)

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

func TestReleaseService_DownloadAsset(t *testing.T) {
	c := &Client{
		httpClient:  &http.Client{},
		rates:       map[rateGroup]Rate{},
		downloadURL: publicDownloadURL,
	}

	tests := []struct {
		name             string
		mockResponses    []MockResponse
		s                *ReleaseService
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
			s: &ReleaseService{
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
			s: &ReleaseService{
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
			s: &ReleaseService{
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

			resp, err := tc.s.DownloadAsset(tc.ctx, tc.releaseTag, tc.assetName, tc.w)

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

func TestReleaseService_DownloadTarArchive(t *testing.T) {
	c := &Client{
		httpClient: &http.Client{},
		rates:      map[rateGroup]Rate{},
		apiURL:     publicAPIURL,
	}

	tests := []struct {
		name             string
		mockResponses    []MockResponse
		s                *ReleaseService
		ctx              context.Context
		ref              string
		w                io.Writer
		expectedResponse *Response
		expectedError    string
	}{
		{
			name:          "NilContext",
			mockResponses: []MockResponse{},
			s: &ReleaseService{
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
			s: &ReleaseService{
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
			s: &ReleaseService{
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

func TestReleaseService_DownloadZipArchive(t *testing.T) {
	c := &Client{
		httpClient: &http.Client{},
		rates:      map[rateGroup]Rate{},
		apiURL:     publicAPIURL,
	}

	tests := []struct {
		name             string
		mockResponses    []MockResponse
		s                *ReleaseService
		ctx              context.Context
		ref              string
		w                io.Writer
		expectedResponse *Response
		expectedError    string
	}{
		{
			name:          "NilContext",
			mockResponses: []MockResponse{},
			s: &ReleaseService{
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
			s: &ReleaseService{
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
			s: &ReleaseService{
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
