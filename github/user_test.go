package github_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/progfay/github-streaks/github"
)

var (
	ignoreUserUnexported = cmpopts.IgnoreUnexported(github.User{})
)

func Test_NewUser(t *testing.T) {
	for _, testcase := range []struct {
		title string
		in    string
		want  *github.User
	}{
		{
			title: "without @ prefix",
			in:    "progfay",
			want: &github.User{
				Name: "progfay",
			},
		},
		{
			title: "with @ prefix",
			in:    "@progfay",
			want: &github.User{
				Name: "progfay",
			},
		},
	} {
		t.Run(testcase.title, func(t *testing.T) {
			user := github.NewUser(testcase.in)

			if !cmp.Equal(*testcase.want, *user, ignoreUserUnexported) {
				t.Errorf("- expect\t+ actual\n%s", cmp.Diff(*testcase.want, *user, ignoreUserUnexported))
				return
			}
		})
	}
}
