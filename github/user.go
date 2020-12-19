package github

import (
	"strings"
)

// User represent GitHub User
type User struct {
	Name string
}

// NewUser generate initialized GitHub User struct with username
func NewUser(username string) *User {
	if strings.HasPrefix(username, "@") {
		username = username[1:]
	}

	return &User{
		Name: username,
	}
}
