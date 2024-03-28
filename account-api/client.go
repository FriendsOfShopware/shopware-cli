package account_api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/FriendsOfShopware/shopware-cli/logging"
)

type Client struct {
	Token            token        `json:"token"`
	ActiveMembership Membership   `json:"active_membership"`
	Memberships      []Membership `json:"memberships"`
}

func (c *Client) NewAuthenticatedRequest(ctx context.Context, method, path string, body io.Reader) (*http.Request, error) {
	logging.FromContext(ctx).Debugf("%s: %s", method, path)
	r, err := http.NewRequestWithContext(ctx, method, path, body)
	if err != nil {
		return nil, err
	}

	r.Header.Set("content-type", "application/json")
	r.Header.Set("accept", "application/json")
	r.Header.Set("x-shopware-token", c.Token.Token)

	return r, nil
}

func (*Client) doRequest(request *http.Request) ([]byte, error) {
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		_ = resp.Body.Close()

		return nil, fmt.Errorf("doRequest: %v", err)
	}

	if err := resp.Body.Close(); err != nil {
		return nil, fmt.Errorf("doRequest: %v", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf(string(data))
	}

	return data, nil
}

func (c *Client) GetActiveCompanyID() int {
	return c.Token.UserID
}

func (c *Client) GetUserID() int {
	return c.Token.UserAccountID
}

func (c *Client) GetActiveMembership() Membership {
	return c.ActiveMembership
}

func (c *Client) GetMemberships() []Membership {
	return c.Memberships
}

func (c *Client) isTokenValid() bool {
	loc, err := time.LoadLocation(c.Token.Expire.Timezone)
	if err != nil {
		return false
	}

	expire, err := time.ParseInLocation("2006-01-02 15:04:05.000000", c.Token.Expire.Date, loc)
	if err != nil {
		return false
	}

	// When it will be expire in the next minute. Respond with false
	return expire.UTC().Sub(time.Now().UTC()).Seconds() > 60
}

const CacheFileName = "shopware-api-client-token.json"

func getApiTokenCacheFilePath() (string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", cacheDir, CacheFileName), nil
}

func createApiFromTokenCache(ctx context.Context) (*Client, error) {
	tokenFilePath, err := getApiTokenCacheFilePath()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(tokenFilePath); os.IsNotExist(err) {
		return nil, err
	}

	content, err := os.ReadFile(tokenFilePath)
	if err != nil {
		return nil, err
	}

	var client *Client
	err = json.Unmarshal(content, &client)
	if err != nil {
		return nil, err
	}

	logging.FromContext(ctx).Debugf("Using token cache from %s", tokenFilePath)
	logging.FromContext(ctx).Debugf("Impersonating currently as %s (%d)", client.ActiveMembership.Company.Name, client.ActiveMembership.Company.Id)

	if !client.isTokenValid() {
		return nil, fmt.Errorf("token is expired")
	}

	return client, nil
}

func saveApiTokenToTokenCache(client *Client) error {
	tokenFilePath, err := getApiTokenCacheFilePath()
	if err != nil {
		return err
	}

	content, err := json.Marshal(client)
	if err != nil {
		return err
	}

	tokenFileDirectory := filepath.Base(tokenFilePath)
	if _, err := os.Stat(tokenFileDirectory); os.IsNotExist(err) {
		err := os.MkdirAll(tokenFileDirectory, 0o750)
		if err != nil {
			return err
		}
	}

	err = os.WriteFile(tokenFilePath, content, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func InvalidateTokenCache() error {
	tokenFilePath, err := getApiTokenCacheFilePath()
	if err != nil {
		return err
	}

	if _, err := os.Stat(tokenFilePath); os.IsNotExist(err) {
		return nil
	}

	return os.Remove(tokenFilePath)
}
