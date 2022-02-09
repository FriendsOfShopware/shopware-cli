package shop

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/doutorfinancas/go-mad/core"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"gopkg.in/yaml.v3"
)

type Config struct {
	URL        string          `yaml:"url"`
	AdminApi   *ConfigAdminApi `yaml:"admin_api"`
	ConfigDump *ConfigDump     `yaml:"dump"`
}

type ConfigAdminApi struct {
	ClientId     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
}

type ConfigDump struct {
	Rewrite map[string]core.Rewrite `yaml:"rewrite"`
	NoData  []string                `yaml:"nodata"  json:"nodata"`
	Ignore  []string                `yaml:"ignore"  json:"ignore"`
	Where   map[string]string       `yaml:"where"   json:"where"`
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
