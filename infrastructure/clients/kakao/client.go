package kakao

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	baseURL string
	client  *http.Client
}

func NewClient(baseURL string, client *http.Client) *Client {
	return &Client{
		baseURL: baseURL,
		client:  client,
	}
}

func (c *Client) GetKakaoUserInfo() (*KakaoUser, error) {
	uri := fmt.Sprintf("%s/v2/user/me", c.baseURL)
	res, err := c.client.Get(uri)
	if err != nil {
		return nil, err
	}
	kakaoUser := new(KakaoUser)
	err = json.NewDecoder(res.Body).Decode(&kakaoUser)
	if err != nil {
		return nil, err
	}
	return kakaoUser, nil
}
