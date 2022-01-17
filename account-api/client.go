package account_api

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type Client struct {
	token            token
	activeMembership Membership
	memberships      []Membership
}

func (c Client) NewAuthenticatedRequest(method, path string, body io.Reader) (*http.Request, error) {
	r, err := http.NewRequestWithContext(context.TODO(), method, path, body) // TODO: pass real context
	if err != nil {
		return nil, err
	}

	r.Header.Set("content-type", "application/json")
	r.Header.Set("accept", "application/json")
	r.Header.Set("x-shopware-token", c.token.Token)

	return r, nil
}

func (c Client) doRequest(request *http.Request) ([]byte, error) {
	resp, err := http.DefaultClient.Do(request)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("doRequest: %v", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf(string(data))
	}

	return data, nil
}

func (c Client) GetActiveCompanyID() int {
	return c.token.UserID
}

func (c Client) GetUserID() int {
	return c.token.UserAccountID
}

func (c Client) GetActiveMembership() Membership {
	return c.activeMembership
}

func (c Client) GetMemberships() []Membership {
	return c.memberships
}
