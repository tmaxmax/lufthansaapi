package ratelimithttp

import (
	"context"
	"io"
	"net/http"
	"net/url"

	"golang.org/x/time/rate"
)

type Client struct {
	client       *http.Client
	rateLimiters []*rate.Limiter
}

func (c *Client) wait(ctx context.Context) error {
	for _, rl := range c.rateLimiters {
		err := rl.Wait(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) Get(url string) (resp *http.Response, err error) {
	err = c.wait(context.Background())
	if err != nil {
		return
	}

	return c.client.Get(url)
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	err := c.wait(req.Context())
	if err != nil {
		return nil, err
	}

	return c.client.Do(req)
}

func (c *Client) Post(url string, contentType string, body io.Reader) (resp *http.Response, err error) {
	err = c.wait(context.Background())
	if err != nil {
		return
	}

	return c.client.Post(url, contentType, body)
}

func (c *Client) PostForm(url string, data url.Values) (resp *http.Response, err error) {
	err = c.wait(context.Background())
	if err != nil {
		return
	}

	return c.client.PostForm(url, data)
}

func (c *Client) Head(url string) (resp *http.Response, err error) {
	err = c.wait(context.Background())
	if err != nil {
		return
	}

	return c.client.Head(url)
}

func (c *Client) CloseIdleConnections() {
	c.client.CloseIdleConnections()
}

func NewClient(client *http.Client, limiters ...*rate.Limiter) *Client {
	return &Client{
		client:       client,
		rateLimiters: limiters,
	}
}
