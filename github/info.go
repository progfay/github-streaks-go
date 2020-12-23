package github

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

var (
	infoEndpointURL *url.URL
)

func init() {
	var err error
	infoEndpointURL, err = url.Parse("https://api.github.com")
	if err != nil {
		panic(err)
	}
}

// UserInfo represent  information of GitHub User
type UserInfo struct {
	CreatedAt string `json:"created_at"`
}

// GetInfo fetch GitHub User Information
func (user *User) GetInfo() (*UserInfo, error) {
	return user.GetInfoWithCustomEndpoint(*infoEndpointURL)
}

// GetInfoWithCustomEndpoint fetch GitHub User Information from specific endpoint
func (user *User) GetInfoWithCustomEndpoint(endpoint url.URL) (*UserInfo, error) {
	endpoint.Path = fmt.Sprintf("/users/%s", url.PathEscape(user.Name))
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

	var info UserInfo
	err = json.NewDecoder(res.Body).Decode(&info)
	if err != nil {
		return nil, err
	}

	return &info, nil
}
