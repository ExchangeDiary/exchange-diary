package kakao

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Client ...
type Client struct {
	baseURL string
	client  *http.Client
}

// NewClient ...
func NewClient(baseURL string, client *http.Client) *Client {
	return &Client{
		baseURL: baseURL,
		client:  client,
	}
}

// GetKakaoUserInfo ...
func (c *Client) GetKakaoUserInfo() (*User, error) {
	uri := fmt.Sprintf("%s/v2/user/me", c.baseURL)
	res, err := c.client.Get(uri)
	if err != nil {
		return nil, err
	}
	kakaoUser := new(User)
	err = json.NewDecoder(res.Body).Decode(&kakaoUser)
	if err != nil {
		return nil, err
	}
	return kakaoUser, nil
}
