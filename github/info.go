package github

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

// UserInfo represent  information of GitHub User
type UserInfo struct {
	CreatedAt string `json:"created_at"`
}

// GetInfo fetch GitHub User Information
func (user *User) GetInfo() (*UserInfo, error) {
	req, err := user.newGetRequest(fmt.Sprintf("https://api.github.com/users/%s", url.PathEscape(user.Name)))
	if err != nil {
		return nil, err
	}

	res, err := user.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	defer io.Copy(ioutil.Discard, res.Body)

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
