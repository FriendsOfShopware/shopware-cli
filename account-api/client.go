package account_api

import (
	"io"
	"net/http"
)

type Client struct {
	token token
}

func (c Client) NewAuthenticatedRequest(method, path string, body io.Reader) (*http.Request, error) {
	r, err := http.NewRequest(method, path, body)

	if err != nil {
		return nil, err
	}

	r.Header.Set("content-type", "application/json")
	r.Header.Set("accept", "application/json")
	r.Header.Set("x-shopware-token", c.token.Token)

	return r, nil

}
