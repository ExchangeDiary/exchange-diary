package apple

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	baseURL string
	client  *http.Client
	token   string
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		client:  http.DefaultClient,
	}
}

func (c *Client) getApplePublicKey(kid string, alg string) *AppleKey {

	keys, err := c.getApplePublicKeys()
	if err != nil || keys == nil {
		return nil
	}

	for _, key := range keys {
		if key.Kid == kid && key.Alg == alg {
			return &key
		}
	}

	return nil
}

func (c *Client) getApplePublicKeys() ([]AppleKey, error) {
	uri := fmt.Sprintf("%s/auth/keys", c.baseURL)
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	var keys struct {
		Keys []AppleKey `json:"keys"`
	}
	err = json.NewDecoder(res.Body).Decode(&keys)
	if err != nil {
		return nil, err
	}
	return keys.Keys, nil
}
