package vk

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

const (
	postsMethod = "wall.get"
)

type Client struct {
	host	string
	basePath string
	version string
	token string
	client http.Client
}

func NewClient(host string, basePath, version string, token string) *Client {
	return &Client{
		host: host,
		basePath: basePath,
		version: version,
		client: http.Client{},
	}
}

func (c *Client) Posts(domain string, offset int, count int) ([]Post, error) {
	q := url.Values{}
	q.Add("domain", domain)
	q.Add("offset", strconv.Itoa(offset))
	q.Add("count", strconv.Itoa(count))
	q.Add("v", c.version)
	q.Add("access_token", c.token)
	
	data, err := c.makeRequest(postsMethod, q)
	if err != nil {
		return nil, err
	}

	var resp PostsResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	return resp.Posts, nil
}

func (c *Client) makeRequest(method string, query url.Values) ([]byte, error) {
	u := url.URL{
		Scheme: "https",
		Host: c.host,
		Path: path.Join(c.basePath, method),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to build request: %w", err)
	}

	req.URL.RawQuery = query.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to do request: %w", err)
	}
	// WARN: ignored error
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read response body: %w", err)
	}

	return body, nil
}