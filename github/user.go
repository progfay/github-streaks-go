package github

import (
	"net/http"
	"strings"
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
