package shop

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"io"
	"net/http"
)

type Client struct {
	url        string
	httpClient *http.Client
}

func NewShopClient(ctx context.Context,config *Config, httpClient *http.Client) (*Client, error) {
	shopClient := &Client{config.URL, httpClient}

	if err := shopClient.authorize(ctx, config); err != nil {
		return nil, err
	}

	return shopClient, nil
}

func (c *Client) authorize(ctx context.Context, config *Config) error {
	if c.httpClient != nil {
		ctx = context.WithValue(ctx, oauth2.HTTPClient, c.httpClient)
	}

	if config.AdminApi == nil {
		return fmt.Errorf("admin-api is not enabled in config")
	}

	tokenSrc, err := config.AdminApi.GetTokenSource(ctx, config.URL)
	if err != nil {
		return err
	}
	c.httpClient = oauth2.NewClient(ctx, tokenSrc)
	return nil
}

func (c *Client) newRequest(ctx context.Context, method, path string, body io.Reader) (*http.Request, error) {
	return http.NewRequestWithContext(ctx, method, c.url+path, body)
}