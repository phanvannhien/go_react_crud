package llm

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type Client struct {
	Endpoint string
	APIKey   string
}

func NewClient(endpoint, key string) *Client {
	return &Client{Endpoint: endpoint, APIKey: key}
}

func (c *Client) call(prompt string, out interface{}) error {

	body, _ := json.Marshal(map[string]string{
		"prompt": prompt,
	})

	req, _ := http.NewRequest("POST", c.Endpoint, bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(out)
}
