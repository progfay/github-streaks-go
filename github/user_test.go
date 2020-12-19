package github_test

import (
	"reflect"
	"testing"

	"github.com/progfay/github-streaks/github"
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
			if !reflect.DeepEqual(*user, *testcase.want) {
				t.Errorf("want user %#v, got %#v", *testcase.want, *user)
				return
			}
		})
	}
}
