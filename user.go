package github

import (
	"context"
	"fmt"
	"time"
)

// UsersService provides GitHub APIs for users.
// See https://docs.github.com/en/rest/reference/users
type UsersService struct {
	client *Client
}

// User is a GitHub user object.
type User struct {
	ID         int       `json:"id"`
	Login      string    `json:"login"`
	Type       string    `json:"type"`
	Email      string    `json:"email"`
	Name       string    `json:"name"`
	URL        string    `json:"url"`
	HTMLURL    string    `json:"html_url"`
	OrgsURL    string    `json:"organizations_url"`
	AvatarURL  string    `json:"avatar_url"`
	GravatarID string    `json:"gravatar_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// User returns the authenticated user.
// If the access token does not have the user scope, then the response includes only the public information.
// If the access token has the user scope, then the response includes the public and private information.
// See https://docs.github.com/rest/reference/users#get-the-authenticated-user
func (s *UsersService) User(ctx context.Context) (*User, *Response, error) {
	req, err := s.client.NewRequest(ctx, "GET", "/user", nil)
	if err != nil {
		return nil, nil, err
	}

	user := new(User)

	resp, err := s.client.Do(req, user)
	if err != nil {
		return nil, nil, err
	}

	return user, resp, nil
}

// Get retrieves a user by its username (login).
// See https://docs.github.com/rest/reference/users#get-a-user
func (s *UsersService) Get(ctx context.Context, username string) (*User, *Response, error) {
	url := fmt.Sprintf("/users/%s", username)
	req, err := s.client.NewRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	user := new(User)

	resp, err := s.client.Do(req, user)
	if err != nil {
		return nil, nil, err
	}

	return user, resp, nil
}
