package github

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"golang.org/x/net/html"
)

var (
	contributionsEndpointURL *url.URL
)

func init() {
	var err error
	contributionsEndpointURL, err = url.Parse("https://github.com")
	if err != nil {
		panic(err)
	}
}

// Contribution represent daily contribution information
type Contribution struct {
	Date  string
	Count int
}

// GetInfo fetch GitHub User Information
func (user *User) GetContributions(year string) ([]Contribution, error) {
	return user.GetContributionsWithCustomEndpoint(*contributionsEndpointURL, year)
}

// GetInfoWithCustomEndpoint fetch GitHub User Information from specific endpoint
func (user *User) GetContributionsWithCustomEndpoint(endpoint url.URL, year string) ([]Contribution, error) {
	endpoint.Path = fmt.Sprintf("/users/%s/contributions", url.PathEscape(user.Name))
	query := url.Values{
		"from": []string{fmt.Sprintf("%s-12-01", year)},
		"to":   []string{fmt.Sprintf("%s-12-31", year)},
	}
	endpoint.RawQuery = query.Encode()
	fmt.Println(endpoint.String())
	req, err := user.newGetRequest(endpoint.String())
	if err != nil {
		return nil, err
	}

	res, err := user.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	defer io.Copy(ioutil.Discard, res.Body)

	if res.StatusCode == http.StatusNotFound {
		return nil, ErrUserNotFound
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fail http request: %q", res.Status)
	}

	doc, err := html.Parse(res.Body)
	if err != nil {
		return nil, err
	}

	var contributions []Contribution
	visitNode := func(n *html.Node) error {
		if n.Type == html.ElementNode && n.Data == "rect" {
			contribution := Contribution{}
			for _, a := range n.Attr {
				switch a.Key {
				case "class":
					if a.Val != "day" {
						return nil
					}

				case "data-count":
					count, err := strconv.Atoi(a.Val)
					if err != nil {
						return err
					}
					contribution.Count = count

				case "data-date":
					contribution.Date = a.Val
				}
			}
			contributions = append(contributions, contribution)
		}
		return nil
	}
	err = forEachNode(doc, visitNode, nil)
	if err != nil {
		return nil, err
	}

	return contributions, nil
}

func forEachNode(n *html.Node, pre, post func(n *html.Node) error) error {
	if pre != nil {
		err := pre(n)
		if err != nil {
			return err
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		err := forEachNode(c, pre, post)
		if err != nil {
			return err
		}
	}

	if post != nil {
		err := post(n)
		if err != nil {
			return err
		}
	}

	return nil
}
