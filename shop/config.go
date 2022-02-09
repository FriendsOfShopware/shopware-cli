package shop

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

type Config struct {
	URL      string          `yaml:"url"`
	AdminApi *ConfigAdminApi `yaml:"admin_api"`
}

type ConfigAdminApi struct {
	ClientId     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
}

func ReadConfig(fileName string) (*Config, error) {
	config := &Config{}

	_, err := os.Stat(fileName)

	if os.IsNotExist(err) {
		return nil, fmt.Errorf("cannot find .shopware-project.yml")
	}

	if err != nil {
		return nil, err
	}

	fileHandle, err := ioutil.ReadFile(fileName)

	if err != nil {
		return nil, fmt.Errorf("ReadConfig: %v", err)
	}

	err = yaml.Unmarshal(fileHandle, &config)

	if err != nil {
		return nil, fmt.Errorf("ReadConfig: %v", err)
	}

	return config, nil
}

func (cfg ConfigAdminApi) GetTokenSource(ctx context.Context, shopURL string) (oauth2.TokenSource, error) {
	tokenURL := fmt.Sprintf("%s/api/oauth/token", shopURL)

	if cfg.Username != "" {
		oauthConf := &oauth2.Config{
			ClientID: "administration",
			Scopes:   []string{"write"},
			Endpoint: oauth2.Endpoint{
				TokenURL: tokenURL,
			},
		}

		token, err := oauthConf.PasswordCredentialsToken(ctx, cfg.Username, cfg.Password)
		if err != nil {
			return nil, err
		}
		return oauth2.StaticTokenSource(token), nil
	}

	oauthConf := &clientcredentials.Config{
		ClientID:     cfg.ClientId,
		ClientSecret: cfg.ClientSecret,
		TokenURL:     tokenURL,
	}

	return oauthConf.TokenSource(ctx), nil
}
