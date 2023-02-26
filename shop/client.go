package shop

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"

	adminSdk "github.com/friendsofshopware/go-shopware-admin-api-sdk"
)

func newShopCredentials(config *Config) adminSdk.OAuthCredentials {
	var cred adminSdk.OAuthCredentials

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
			InsecureSkipVerify: config.AdminApi.DisableSSLCheck,
		},
	}
	client := &http.Client{Transport: tr}

	return adminSdk.NewApiClient(ctx, config.URL, newShopCredentials(config), client)
}
