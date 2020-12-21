package github_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/progfay/github-streaks/github"
)

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func Test_GetInfo(t *testing.T) {
	for _, testcase := range []struct {
		title string
		in    string
		want  *github.UserInfo
	}{
		{
			title: "exists user",
			in:    "progfay",
			want: &github.UserInfo{
				CreatedAt: "2016-05-25T09:24:46Z",
			},
		},
		{
			title: "not exists user",
			in:    "not-exists",
			want: &github.UserInfo{
				CreatedAt: "",
			},
		},
	} {
		t.Run(testcase.title, func(t *testing.T) {
			b, err := ioutil.ReadFile(fmt.Sprintf("data/users/%s.json", testcase.in))
			if err != nil {
				t.Error(err)
				return
			}
			user := github.NewUserWithroundTripper(testcase.in, RoundTripFunc(func(req *http.Request) *http.Response {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBuffer(b)),
					Header:     make(http.Header),
				}
			}))

			info, err := user.GetInfo()
			if err != nil {
				t.Error(err)
				return
			}
			if !cmp.Equal(*testcase.want, *info, nil) {
				t.Errorf("- expect\t+ actual\n%s", cmp.Diff(*testcase.want, *info, nil))
				return
			}
		})
	}
}
