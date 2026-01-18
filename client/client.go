package client

import (
	"github.com/unluckythoughts/go-microservice/tools/web"
)

// Client provides methods to interact with the Manga Reader API
type Client struct {
	client  web.Client
	baseURL string
}

// NewClient creates a new API client
func NewClient(baseURL string) *Client {
	return &Client{
		client:  web.NewClient(baseURL),
		baseURL: baseURL,
	}
}

// SetAuthToken sets the authentication token for the client
func (c *Client) SetAuthToken(token string) {
	c.client.SetBearerToken(token)
}
