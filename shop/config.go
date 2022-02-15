package shop

import (
	"context"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"strings"

	"github.com/doutorfinancas/go-mad/core"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

type Config struct {
	URL        string          `yaml:"url"`
	AdminApi   *ConfigAdminApi `yaml:"admin_api,omitempty"`
	ConfigDump *ConfigDump     `yaml:"dump,omitempty"`
	Sync       *ConfigSync     `yaml:"sync,omitempty"`
}

type ConfigAdminApi struct {
	ClientId     string `yaml:"client_id,omitempty"`
	ClientSecret string `yaml:"client_secret,omitempty"`
	Username     string `yaml:"username,omitempty"`
	Password     string `yaml:"password,omitempty"`
}

type ConfigDump struct {
	Rewrite map[string]core.Rewrite `yaml:"rewrite,omitempty"`
	NoData  []string                `yaml:"nodata,omitempty"`
	Ignore  []string                `yaml:"ignore,omitempty"`
	Where   map[string]string       `yaml:"where,omitempty"`
}

type ConfigSync struct {
	Config []ConfigSyncConfig `yaml:"config"`
}

type ConfigSyncConfig struct {
	SalesChannel *string                `yaml:"sales_channel,omitempty"`
	Settings     map[string]interface{} `yaml:"settings"`
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

func NewUuid() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}
