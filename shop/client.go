package shop

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"

	adminSdk "github.com/friendsofshopware/go-shopware-admin-api-sdk"
)

func newShopCredentials(config *Config) adminSdk.OAuthCredentials {
	var cred adminSdk.OAuthCredentials

	clientId, clientSecret := os.Getenv("SHOPWARE_CLI_API_CLIENT_ID"), os.Getenv("SHOPWARE_CLI_API_CLIENT_SECRET")

	if clientId != "" && clientSecret != "" {
		return adminSdk.NewIntegrationCredentials(clientId, clientSecret, []string{"write"})
	}

	username, password := os.Getenv("SHOPWARE_CLI_API_USERNAME"), os.Getenv("SHOPWARE_CLI_API_PASSWORD")

	if username != "" && password != "" {
		return adminSdk.NewPasswordCredentials(username, password, []string{"write"})
	}

	if config.AdminApi.Username != "" {
		cred = adminSdk.NewPasswordCredentials(config.AdminApi.Username, config.AdminApi.Password, []string{"write"})
	} else {
		cred = adminSdk.NewIntegrationCredentials(config.AdminApi.ClientId, config.AdminApi.ClientSecret, []string{"write"})
	}

	return cred
}

func NewShopClient(ctx context.Context, config *Config) (*adminSdk.Client, error) {
	if config.AdminApi == nil {
		return nil, fmt.Errorf("admin-api is not enabled in config")
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: config.AdminApi.DisableSSLCheck, // nolint:gosec
		},
	}
	client := &http.Client{Transport: tr}

	shopUrl := os.Getenv("SHOPWARE_CLI_API_URL")

	if shopUrl == "" {
		shopUrl = config.URL
	}

	return adminSdk.NewApiClient(ctx, shopUrl, newShopCredentials(config), client)
}
