package github

import "fmt"

var (
	// ErrUserNotFound throw on user not found
	ErrUserNotFound = fmt.Errorf("user not found")
)
