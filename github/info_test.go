package github_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/progfay/github-streaks/github"
)

func Test_GetInfo(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/users/progfay":
			http.ServeFile(w, r, "data/users/progfay.json")
			return

		case "/users/not-exists":
			w.WriteHeader(http.StatusNotFound)
			http.ServeFile(w, r, "data/users/not-exists.json")
			return
		}
	}))
	// http.FileServer(http.Dir("data")))
	defer ts.Close()

	endpoint, err := url.Parse(ts.URL)
	if err != nil {
		t.Error(err)
		return
	}

	for _, testcase := range []struct {
		title string
		in    string
		want  struct {
			info *github.UserInfo
			err  error
		}
	}{
		{
			title: "exists user",
			in:    "progfay",
			want: struct {
				info *github.UserInfo
				err  error
			}{
				info: &github.UserInfo{
					CreatedAt: "2016-05-25T09:24:46Z",
				},
				err: nil,
			},
		},
		{
			title: "not exists user",
			in:    "not-exists",
			want: struct {
				info *github.UserInfo
				err  error
			}{
				info: nil,
				err:  github.ErrUserNotFound,
			},
		},
	} {
		t.Run(testcase.title, func(t *testing.T) {
			user := github.NewUser(testcase.in)

			info, err := user.GetInfoWithCustomEndpoint(*endpoint)
			if err != testcase.want.err {
				t.Errorf("want error %v, got %v", testcase.want.err, err)
				return
			}
			if testcase.want.info == nil && info == nil {
				return
			}
			if testcase.want.info == nil || info == nil {
				t.Errorf("want info %v, got %v", testcase.want.info, info)
				return
			}
			if !cmp.Equal(*testcase.want.info, *info, nil) {
				t.Errorf("- expect\t+ actual\n%s", cmp.Diff(*testcase.want.info, *info, nil))
				return
			}
		})
	}
}
