package github

import (
	"net/http"
	"strings"
)

const (
	userAgent = "@progfay/github-streaks"
	endpoint  = "https://github.com"
)

// User represent GitHub User
type User struct {
	Name   string
	client *http.Client
}

// NewUser generate initialized GitHub User struct with username
func NewUser(username string) *User {
	if strings.HasPrefix(username, "@") {
		username = username[1:]
	}

	return &User{
		Name:   username,
		client: new(http.Client),
	}
}

func (user *User) newGetRequest(url string) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", userAgent)

	return req, nil
}
