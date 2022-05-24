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
	token   string
}

// NewClient ...
func NewClient(baseURL string, token string) *Client {
	return &Client{
		baseURL: baseURL,
		client:  http.DefaultClient,
		token:   token,
	}
}

// GetKakaoUserInfo ...
func (c *Client) GetKakaoUserInfo() (*User, error) {
	uri := fmt.Sprintf("%s/v2/user/me", c.baseURL)
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)

	res, err := c.client.Do(req)
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
