package account_api

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type Client struct {
	token            token
	activeMembership *membership
	memberships      *[]membership
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

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf(string(data))
	}

	return data, nil
}

func (c Client) GetActiveCompanyId() int {
	return c.token.UserID
}

func (c Client) GetUserId() int {
	return c.token.UserAccountID
}

func (c Client) GetActiveMembership() *membership {
	return c.activeMembership
}

func (c Client) GetMemberships() *[]membership {
	return c.memberships
}
