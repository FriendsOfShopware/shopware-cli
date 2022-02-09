package shop

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"io"
	"os"
)

type ClientCredentials struct {
	Id, Secret string
}

func (c ClientCredentials) getTokenSource(ctx context.Context, tokenURL string) (oauth2.TokenSource, error) {
	oauthConf := &clientcredentials.Config{
		ClientID: c.Id,
		ClientSecret: c.Secret,
	}

	return oauthConf.TokenSource(ctx), nil
}

func (c ClientCredentials) IsValid() bool {
	return len(c.Id) > 0 && len(c.Secret) > 0
}

type PasswordCredentials struct {
	Username, Password string
}

func (p PasswordCredentials) getTokenSource(ctx context.Context, tokenURL string) (oauth2.TokenSource, error) {
	oauthConf := &oauth2.Config{
		ClientID: "administration",
		Scopes:   []string{"write"},
		Endpoint: oauth2.Endpoint{
			TokenURL: tokenURL,
		},
	}

	token, err := oauthConf.PasswordCredentialsToken(ctx, p.Username, p.Password)
	if err != nil {
		return nil, err
	}
	return oauth2.StaticTokenSource(token), nil
}

func (p PasswordCredentials) IsValid() bool {
	return len(p.Password) > 0 && len(p.Username) > 0
}

func ReadCredentialsFromFile(file string) (Credentials, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	creds, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	var passwordGrant PasswordCredentials
	var clientCredentials ClientCredentials
	var grant Credentials
	if err := json.Unmarshal(creds, &passwordGrant); err == nil && passwordGrant.IsValid() {
		grant = passwordGrant
	}
	if err := json.Unmarshal(creds, &clientCredentials); err == nil && clientCredentials.IsValid() {
		grant = clientCredentials
	}
	if grant == nil {
		return nil,errors.New("file does not contain valid credentials")
	}
	return grant, nil
}