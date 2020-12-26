package github_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/progfay/github-streaks/github"
)

type snapshot struct {
	Contributions []github.Contribution `json:"contributions"`
}

var progfaySnapshot snapshot

func init() {
	raw, err := ioutil.ReadFile("./snapshot/contributions/progfay.json")
	if err != nil {
		panic(err)
	}

	json.Unmarshal(raw, &progfaySnapshot)
}

func Test_GetAnnualContributions(t *testing.T) {
	var wantQuery = url.Values{
		"from": []string{"2020-12-01"},
		"to":   []string{"2020-12-31"},
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if !cmp.Equal(q, wantQuery) {
			t.Errorf("want query %q, got %q", wantQuery.Encode(), q.Encode())
		}

		switch r.URL.Path {
		case "/users/progfay/contributions":
			http.ServeFile(w, r, "data/contributions/progfay.html")
			return

		case "/users/not-exists/contributions":
			w.WriteHeader(http.StatusNotFound)
			http.ServeFile(w, r, "data/contributions/not-exists.html")
			return

		default:
			t.Errorf("invalid http request: %s", r.URL.String())
		}
	}))
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
			contributions []github.Contribution
			err           error
		}
	}{
		{
			title: "exists user",
			in:    "progfay",
			want: struct {
				contributions []github.Contribution
				err           error
			}{
				contributions: progfaySnapshot.Contributions,
				err:           nil,
			},
		},
		{
			title: "not exists user",
			in:    "not-exists",
			want: struct {
				contributions []github.Contribution
				err           error
			}{
				contributions: nil,
				err:           github.ErrUserNotFound,
			},
		},
	} {
		t.Run(testcase.title, func(t *testing.T) {
			user := github.NewUser(testcase.in)

			got, err := user.GetAnnualContributionsWithCustomEndpoint(*endpoint, "2020")
			if err != testcase.want.err {
				t.Errorf("want error %v, got %v", testcase.want.err, err)
				return
			}
			if !cmp.Equal(testcase.want.contributions, got, nil) {
				t.Errorf("- expect\t+ actual\n%s", cmp.Diff(testcase.want.contributions, got, nil))
				return
			}
		})
	}
}
